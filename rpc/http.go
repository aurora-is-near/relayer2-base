package rpc

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	cors "github.com/adhityaramadhanus/fasthttpcors"
	"github.com/aurora-is-near/relayer2-base/log"
	"github.com/aurora-is-near/relayer2-base/utils"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/valyala/fasthttp/reuseport"
)

const (
	maxRequestContentLength = 1024 * 1024 * 5
	DefaultContentType      = "application/json"

	wsReadBuffer       = 1024
	wsWriteBuffer      = 1024
	wsMessageSizeLimit = 15 * 1024 * 1024
)

// https://www.jsonrpc.org/historical/json-rpc-over-http.html#id13
var acceptedContentTypes = []string{DefaultContentType, "application/json-rpc", "application/jsonrequest"}

type HttpServer struct {
	Logger     *log.Logger
	Config     HttpConfig
	resolver   Resolver
	listener   net.Listener
	wsUpgrader *websocket.FastHTTPUpgrader
}

// HttpConfig holds both http and websocket configuration elements
type HttpConfig struct {
	HttpEndpoint       string
	HttpCors           []string
	HttpCompress       bool
	HttpTimeout        time.Duration
	WsEndpoint         string
	WsHandshakeTimeout time.Duration
	WsOnly             bool
}

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  wsReadBuffer,
	WriteBufferSize: wsWriteBuffer,
	WriteBufferPool: new(sync.Pool),
}

// Run enables CORS and Compression handlers if enabled on configuration and
// starts the fasthttp server
func (h *HttpServer) Run(ctx context.Context, resolver Resolver) error {
	h.Logger.Info().Msgf("starting HTTP server on %s", h.Config.HttpEndpoint)
	var err error
	h.listener, err = reuseport.Listen("tcp4", h.Config.HttpEndpoint)
	if err != nil {
		return err
	}

	h.resolver = resolver
	h.wsUpgrader = &websocket.FastHTTPUpgrader{
		HandshakeTimeout: h.Config.WsHandshakeTimeout * time.Second,
		CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
			return true
		},
	}

	var reqHandler fasthttp.RequestHandler
	// Use CorsHandler if CORS settings are applied
	withCors := h.newCorsHandler(h.Config.HttpCors)
	if withCors != nil {
		reqHandler = withCors.CorsMiddleware(h.mainHandler)
	} else {
		reqHandler = h.mainHandler
	}
	// Transparently compress the response body if the request contains 'gzip' or 'deflate' 'Accept-Encoding' header
	// and configuration is enabled
	if h.Config.HttpCompress {
		reqHandler = fasthttp.CompressHandler(reqHandler)
	}
	reqHandler = fasthttp.TimeoutHandler(reqHandler, h.Config.HttpTimeout*time.Second, "request timeout")

	go func() {
		if err := fasthttp.Serve(h.listener, reqHandler); err != nil {
			h.Logger.Fatal().Err(err).Msg("error while serving http")
		}
	}()
	return nil
}

// mainHandler is the handler called by the fasthttp server when a proper http
// request is received. It calls the appropriate handler based on transport type
func (h *HttpServer) mainHandler(ctx *fasthttp.RequestCtx) {
	r := &http.Request{}
	if err := fasthttpadaptor.ConvertRequest(ctx, r, false); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	// HttpServer can handle both http and ws requests. If http and ws configurations use the same port, then single
	// HttpServer is running and handling all requests. If they have different ports, then 2 separate servers are running.
	// The following code differentiates the cases based on the configurations and Config.WsOnly variable
	if strings.EqualFold(h.Config.HttpEndpoint, h.Config.WsEndpoint) {
		if isWebsocket(r) {
			h.fastWsHandler(ctx)
		} else if !h.Config.WsOnly {
			h.fastHTTPHandler(ctx, r)
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		}
	} else {
		h.fastHTTPHandler(ctx, r)
	}
}

// fastWSHandler is the handler to serve WS requests
func (h *HttpServer) fastWsHandler(ctx *fasthttp.RequestCtx) {
	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		defer conn.Close()
		wsCtx := &WebSocketContext{ws: conn, subscriptions: make(map[ID]*Subscription), subscriptionsMtx: sync.Mutex{}}
		wsCtx.output = make(chan []byte, 100)
		wsCtx.closed.Store(false)

		wsCtx.outputWg.Add(1)
		go h.handleWebSocketOutput(wsCtx)

		for {
			messageType, message, err := wsCtx.ws.ReadMessage()
			if err != nil {
				break
			}
			if messageType != websocket.TextMessage {
				h.Logger.Warn().Msgf("websocket: got message with unknown type #%v", messageType)
				continue
			}
			// get the clientIp and add it to context so rpcserver can use it when needed
			clientIp := ctx.RemoteIP()
			cCtx := utils.PutClientIpKey(ctx, clientIp)
			resp := h.resolver.ResolveWs(&cCtx, wsCtx, message)
			if resp != nil {
				wsCtx.output <- resp
			}
		}

		wsCtx.closed.Store(true)
		h.resolver.CloseWsConn(wsCtx)
		wsCtx.outputWg.Wait()
	})
	if err != nil {
		h.Logger.Error().Err(err).Msg("error upgrading to websocket")
	}
}

// handleWebSocketOutput listens the websocket output channel and writes incoming data to connection
func (h *HttpServer) handleWebSocketOutput(wsCtx *WebSocketContext) {
	defer wsCtx.outputWg.Done()

	for data := range wsCtx.output {
		if !wsCtx.closed.Load() {
			err := wsCtx.ws.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				h.Logger.Error().Err(err).Msg("error while writing to websocket")
			}
		}
	}
}

// fastHTTPHandler is the handler to serve HTTP requests
func (h *HttpServer) fastHTTPHandler(ctx *fasthttp.RequestCtx, r *http.Request) {
	ctx.Response.Header.SetServer("Relayer " + *utils.Constants.RelayerVersion())

	if code, err := validateRequest(r); err != nil {
		ctx.Error(err.Error(), code)
		return
	}

	ctx.SetContentType(DefaultContentType)
	ctx.SetStatusCode(http.StatusOK)

	// get the clientIp and add it to context so rpcserver can use it when needed
	clientIp := ctx.RemoteIP()
	cCtx := utils.PutClientIpKey(ctx, clientIp)
	resp := h.resolver.ResolveHttp(&cCtx, ctx.Request.Body())
	ctx.SetBody(resp)
}

// newCorsHandler creates and returns a fasthttp compliant CORS handler if CORS configuration is enabled
func (h *HttpServer) newCorsHandler(headers []string) *cors.CorsHandler {
	// disable CORS support if user has not specified a custom CORS configuration
	if len(h.Config.HttpCors) == 0 {
		return nil
	}
	c := cors.NewCorsHandler(cors.Options{
		AllowedOrigins: h.Config.HttpCors,
		AllowedMethods: []string{http.MethodPost, http.MethodGet},
		AllowedHeaders: headers,
		AllowMaxAge:    int(h.Config.HttpTimeout.Seconds()),
	})
	return c
}

// Stop stops http server
func (h *HttpServer) Stop() error {
	h.Logger.Info().Msg("stopping http listener...")
	if err := h.listener.Close(); err != nil {
		h.Logger.Error().Msgf("error while closing listener: %v", err)
		return err
	}
	return nil
}

// validateRequest checks and validates incoming http request
func validateRequest(r *http.Request) (int, error) {
	if r.Method != http.MethodPost && r.Method != http.MethodHead && r.Method != http.MethodOptions {
		return http.StatusMethodNotAllowed, errors.New("method not allowed")
	}
	if r.ContentLength > maxRequestContentLength {
		err := fmt.Errorf("content length too large (%d>%d)", r.ContentLength, maxRequestContentLength)
		return http.StatusRequestEntityTooLarge, err
	}
	// Allow OPTIONS (regardless of content-type)
	if r.Method == http.MethodOptions {
		return 0, nil
	}
	// Check content-type
	if mt, _, err := mime.ParseMediaType(r.Header.Get("content-type")); err == nil {
		for _, accepted := range acceptedContentTypes {
			if accepted == mt {
				return 0, nil
			}
		}
	}
	// Invalid content-type
	err := fmt.Errorf("invalid content type, only %s are supported", strings.Join(acceptedContentTypes, ","))
	return http.StatusUnsupportedMediaType, err
}

// isWebsocket checks the header of an http request for a websocket upgrade request
func isWebsocket(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("Upgrade"), "websocket") &&
		strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade")
}
