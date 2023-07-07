package rpc

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/buger/jsonparser"
	jsoniter "github.com/json-iterator/go"

	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/aurora-is-near/relayer2-base/log"
	errs "github.com/aurora-is-near/relayer2-base/types/errors"
	"golang.org/x/sync/errgroup"

	"github.com/aurora-is-near/relayer2-base/rpc/types"
)

const (
	subscriptionPrefix = "eth_"
)

type RpcServer struct {
	serviceMap       ServiceMap
	logger           *log.Logger
	middlewares      []Middleware
	transports       []Transport
	mu               sync.RWMutex
	maxBatchRequests uint
}

func New(l *log.Logger, maxBatchReq uint, transports ...TransportOption) *RpcServer {
	s := &RpcServer{
		serviceMap:       ServiceMap{},
		logger:           l,
		transports:       []Transport{},
		mu:               sync.RWMutex{},
		maxBatchRequests: maxBatchReq,
	}
	for _, transport := range transports {
		transport(s)
	}
	return s
}

// RegisterEndpoints creates a map of service handlers for the given receiver by adding its exposed methods and their arguments
func (r *RpcServer) RegisterEndpoints(nameSpace string, sh any) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	err := r.serviceMap.register(nameSpace, sh, false)
	if err != nil {
		r.logger.Error().Err(err).Msgf("can't register service with %s namespace", nameSpace)
		return err
	}
	return nil
}

// RegisterEvents creates a map of subscription handlers for the given receiver by adding its exposed methods and their arguments
func (r *RpcServer) RegisterEvents(nameSpace string, sh any) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	err := r.serviceMap.register(nameSpace, sh, true)
	if err != nil {
		r.logger.Error().Err(err).Msgf("can't register events with %s namespace", nameSpace)
		return err
	}
	return nil
}

// WithMiddleware places the given handler function to the middlewares chain.
func (r *RpcServer) WithMiddleware(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}

// Run starts the configured transports (HTTP and WS) for the json rpc server
func (r *RpcServer) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, t := range r.transports {
		eg.Go(func(t Transport) func() error {
			return func() error { return t.Run(ctx, r) }
		}(t))
	}
	return eg.Wait()
}

func (r *RpcServer) ResolveHttp(ctx *context.Context, rpcMessage []byte) []byte {
	rpcCtx := r.prepareRpcContext(rpcMessage, nil)

	if rpcCtx.hasParseError() {
		return rpcCtx.response
	} else if rpcCtx.isBatch() {
		return r.executeBatchRequest(ctx, rpcCtx, false).response
	} else {
		return r.executeSingleRequest(ctx, rpcCtx, false).response
	}
}

func (r *RpcServer) ResolveWs(ctx *context.Context, wsCtx *WebSocketContext, rpcMessage []byte) []byte {
	rpcCtx := r.prepareRpcContext(rpcMessage, wsCtx)

	if rpcCtx.hasParseError() {
		return rpcCtx.response
	} else if rpcCtx.isBatch() {
		return r.executeBatchRequest(ctx, rpcCtx, true).response
	} else {
		return r.executeSingleRequest(ctx, rpcCtx, true).response
	}
}

func (r *RpcServer) CloseWsConn(wsCtx *WebSocketContext) {
	wsCtx.subscriptionsMtx.Lock()
	defer wsCtx.subscriptionsMtx.Unlock()
	// WS connection is closing, if all subscriptions are prevously unsubscribed subsciptions list is empty
	// if it is not empty, notify err channel of each subscription, so that the records can be cleaned
	//  from broker notify list
	for k := range wsCtx.subscriptions {
		close(wsCtx.subscriptions[k].err)
	}
	// clear the subscriptions map
	wsCtx.subscriptions = make(map[ID]*Subscription)
	wsCtx.ws = nil
}

// executeBatchRequest runs each json-rpc request in parallel and formats the responses and returns updated RpcContext object with the total response
func (r *RpcServer) executeBatchRequest(ctx *context.Context, rpcCtx *RpcContext, isWs bool) *RpcContext {
	childResponsesChan := make(chan *RpcContext)

	for i, child := range rpcCtx.batchChildren {
		capturedChild := child.setIndexInBatch(i)
		go func() {
			childResponsesChan <- r.executeSingleRequest(ctx, capturedChild, isWs)
		}()
	}

	childResponses := make([][]byte, len(rpcCtx.batchChildren))
	for range rpcCtx.batchChildren {
		childResponse := <-childResponsesChan
		childResponses[childResponse.indexInBatch] = childResponse.response
	}

	rpcCtx.response = []byte("[")
	cnt := 0
	for _, childResponse := range childResponses {
		// response of the notification requests are empty so handle them carefully
		if len(childResponse) > 0 {
			if cnt > 0 {
				rpcCtx.response = append(rpcCtx.response, []byte(",")...)
			}
			rpcCtx.response = append(rpcCtx.response, childResponse...)
			cnt++
		}
	}
	rpcCtx.response = append(rpcCtx.response, []byte("]")...)

	return rpcCtx
}

// executeSingleRequest runs the retrieved json-rpc request and returns the updated RpcContext object with the response
func (r *RpcServer) executeSingleRequest(ctx *context.Context, rpcCtx *RpcContext, isWs bool) *RpcContext {
	switch {
	case rpcCtx.parseFailed:
		//Do nothing, parse error had already added to response while parsing
	case rpcCtx.parsedBody.IsNotification():
		rpcCtx.setNotificationResponse()
	case rpcCtx.parsedBody.IsMethodCall():
		var h func(ctx *context.Context, rpcCtx *RpcContext) *RpcContext
		if rpcCtx.parsedBody.IsSubscribe() {
			if isWs {
				h = r.callSubscribe
			} else {
				return rpcCtx.SetErrorObject(&errs.InvalidRequestError{Message: "subscribe method is for websocket only"})
			}
		} else if rpcCtx.parsedBody.IsUnsubscribe() {
			if isWs {
				h = r.callUnsubscribe
			} else {
				return rpcCtx.SetErrorObject(&errs.InvalidRequestError{Message: "unsubscribe method is for websocket only"})
			}
		} else { // standard method call
			h = r.callMethod
		}
		for _, m := range r.middlewares {
			h = m(h)
		}
		rpcCtx = h(ctx, rpcCtx)
	default:
		rpcCtx = rpcCtx.SetErrorObject(&errs.InvalidRequestError{Message: "invalid request"})
	}
	return rpcCtx
}

// processSubscriptionRequest handles the incoming subscription request
func (r *RpcServer) processSubscriptionRequest(req *types.RPCRequestBody) (*handler, error) {
	// subscription method is the first argument
	method, err := parseSubscriptionMethod(req.Params.Value)
	if err != nil {
		return nil, &errs.InvalidParamsError{Message: err.Error()}
	}

	r.mu.RLock()
	s, ok := r.serviceMap.subscriptions[subscriptionPrefix+strings.ToLower(method)]
	r.mu.RUnlock()
	if !ok {
		return nil, &errs.SubscriptionNotFoundError{Subscription: method}
	}
	return s.handler, nil
}

// callSubscribe processes *_subscribe methods
func (r *RpcServer) callSubscribe(ctx *context.Context, rpcCtx *RpcContext) *RpcContext {
	req := rpcCtx.parsedBody
	handler, err := r.processSubscriptionRequest(req)
	if err != nil {
		return rpcCtx.SetErrorObject(&errs.GenericError{Err: err})
	}

	// Parse subscription name arg too, but remove it before calling the callback.
	argTypes := append([]reflect.Type{stringType}, handler.argTypes...)
	args, err := prepareArguments(req.Params.Value, argTypes)
	if err != nil {
		return rpcCtx.SetErrorObject(&errs.InvalidParamsError{Message: err.Error()})
	}
	args = args[1:]

	// Add notifier to context so that subscription handler can use it
	n := &Notifier{h: handler, wsCtx: rpcCtx.wsCtx}
	*ctx = PutNotifierKey(*ctx, n)
	resp, err := handler.call(ctx, args)
	if err != nil {
		return rpcCtx.SetErrorObject(&errs.GenericError{Err: err})
	}

	respJson, err := jsoniter.Marshal(resp)
	if err != nil {
		return rpcCtx.SetErrorObject(&errs.GenericError{Err: err})
	}

	return rpcCtx.setResult(respJson)
}

// callUnsubscribe processes *_unsubscribe methods
func (r *RpcServer) callUnsubscribe(ctx *context.Context, rpcCtx *RpcContext) *RpcContext {
	subscriptionId, subscriptionIdDataType, _, err := jsonparser.Get(rpcCtx.parsedBody.Params.Value, "[0]")
	if err != nil {
		return rpcCtx.SetErrorObject(&errs.InvalidParamsError{Message: fmt.Sprintf("can't parse eth_unsibscribe params: %v", err)})
	}
	if subscriptionIdDataType != jsonparser.String {
		return rpcCtx.SetErrorObject(&errs.InvalidParamsError{Message: "subscription ID must be string"})
	}

	rpcCtx.wsCtx.subscriptionsMtx.Lock()
	defer rpcCtx.wsCtx.subscriptionsMtx.Unlock()

	sub, ok := rpcCtx.wsCtx.subscriptions[ID(subscriptionId)]
	if !ok {
		// return false if the subsciptionId couldn't be found
		return rpcCtx.setResult([]byte("false"))
	}
	// unsubscribe call received, notify err channel of each subscription
	// so that records can be cleaned from broker notify list
	close(sub.err)
	// delete the related subscription from subscriptions list
	delete(rpcCtx.wsCtx.subscriptions, ID(subscriptionId))

	return rpcCtx.setResult([]byte("true"))
}

// callMethod processes incoming service methods
func (r *RpcServer) callMethod(ctx *context.Context, rpcCtx *RpcContext) *RpcContext {
	r.mu.RLock()
	s, ok := r.serviceMap.services[strings.ToLower(rpcCtx.parsedBody.Method.Str())]
	r.mu.RUnlock()
	if !ok {
		return rpcCtx.SetErrorObject(&errs.MethodNotFoundError{Method: rpcCtx.parsedBody.Method.Str()})
	}

	args, err := prepareArguments(rpcCtx.parsedBody.Params.Value, s.handler.argTypes)
	if err != nil {
		return rpcCtx.SetErrorObject(&errs.InvalidParamsError{Message: err.Error()})
	}

	resp, err := s.handler.call(ctx, args)
	if err != nil {
		e, ok := err.(errs.Error)
		if ok {
			return rpcCtx.SetErrorObject(e)
		} else {
			return rpcCtx.setError(errs.Generic, err.Error())
		}
	}

	respJson, err := jsoniter.Marshal(resp)
	if err != nil {
		return rpcCtx.SetErrorObject(&errs.GenericError{Err: err})
	}

	return rpcCtx.setResult(respJson)
}

// Close gracefully stops RpcServer
func (r *RpcServer) Close() {
	for _, t := range r.transports {
		t.Stop()
	}
}

// prepareArguments tries to parse the given args to an array of values with the
// given types. It returns the parsed values or an error when the args could not be
// parsed.
func prepareArguments(rawArgs jsoniter.RawMessage, types []reflect.Type) ([]reflect.Value, error) {
	dec := json.NewDecoder(bytes.NewReader(rawArgs))
	var args []reflect.Value
	tok, err := dec.Token()
	switch {
	case err == io.EOF || tok == nil && err == nil:
		// "params" is optional and may be empty. Also allow "params":null even though it's not in the spec
	case err != nil:
		return nil, err
	case tok == json.Delim('['):
		if args, err = parseArgumentArray(dec, types); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("non-array args")
	}
	// Set any missing args to nil.
	for i := len(args); i < len(types); i++ {
		if types[i].Kind() != reflect.Ptr {
			return nil, fmt.Errorf("missing value for required argument %d", i)
		}
		args = append(args, reflect.Zero(types[i]))
	}
	return args, nil
}

// parseArgumentArray parses the arguments given in the provided types array and returns the parsed arguments
func parseArgumentArray(dec *json.Decoder, types []reflect.Type) ([]reflect.Value, error) {
	args := make([]reflect.Value, 0, len(types))
	for i := 0; dec.More(); i++ {
		if i >= len(types) {
			return args, fmt.Errorf("too many arguments, want at most %d", len(types))
		}
		argval := reflect.New(types[i])
		if err := dec.Decode(argval.Interface()); err != nil {
			return args, fmt.Errorf("invalid argument %d: %v", i, err)
		}
		if argval.IsNil() && types[i].Kind() != reflect.Ptr {
			return args, fmt.Errorf("missing value for required argument %d", i)
		}
		args = append(args, argval.Elem())
	}
	// Read end of args array.
	_, err := dec.Token()
	return args, err
}

// parseSubscriptionMethod extracts the subscription method name
func parseSubscriptionMethod(rawArgs jsoniter.RawMessage) (string, error) {
	dec := json.NewDecoder(bytes.NewReader(rawArgs))
	if tok, _ := dec.Token(); tok != json.Delim('[') {
		return "", errors.New("non-array args")
	}
	v, _ := dec.Token()
	method, ok := v.(string)
	if !ok {
		return "", errors.New("expected subscription name as first argument")
	}
	return method, nil
}
