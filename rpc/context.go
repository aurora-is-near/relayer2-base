package rpc

import (
	"sync"

	"github.com/buger/jsonparser"
	"github.com/fasthttp/websocket"
	"go.uber.org/atomic"

	"github.com/aurora-is-near/relayer2-base/rpc/types"
	errs "github.com/aurora-is-near/relayer2-base/types/errors"
)

type RpcContext struct {
	wsCtx          *WebSocketContext
	body           []byte
	bodyDataType   jsonparser.ValueType
	bodyParseError error
	indexInBatch   int
	parsedBody     *types.RPCRequestBody
	method         string
	parseFailed    bool
	response       []byte
	batchChildren  []*RpcContext
}

func (r *RpcServer) newRpcContext(wsCtx *WebSocketContext, body []byte, dataType jsonparser.ValueType, parseError error) *RpcContext {
	ctx := &RpcContext{
		wsCtx:          wsCtx,
		body:           body,
		bodyDataType:   dataType,
		bodyParseError: parseError,
	}
	return ctx
}

type WebSocketContext struct {
	ws               *websocket.Conn
	output           chan []byte
	subscriptions    map[ID]*Subscription
	subscriptionsMtx sync.Mutex
	outputWg         sync.WaitGroup
	closed           atomic.Bool
}

// hasParseError returns if there is parse error or not
func (ctx *RpcContext) hasParseError() bool {
	return ctx.parseFailed
}

// isBatch returns true if the request is a batch request
func (ctx *RpcContext) isBatch() bool {
	return len(ctx.batchChildren) > 0
}

// setIndexInBatch sets the provided argument as the batch index
func (ctx *RpcContext) setIndexInBatch(index int) *RpcContext {
	ctx.indexInBatch = index
	return ctx
}

// getRpcIdRepr returns the Id of the rpc request
func (ctx *RpcContext) getRpcIdRepr() []byte {
	if ctx.parsedBody == nil {
		return []byte("1")
	}
	return ctx.parsedBody.ID.Repr()
}

// setNotificationResponse sets empty response
func (ctx *RpcContext) setNotificationResponse() *RpcContext {
	ctx.response = []byte{}
	return ctx
}

// setError sets response as error using the error code and message
func (ctx *RpcContext) setError(code int64, message string) *RpcContext {
	ctx.parseFailed = true
	ctx.response = createErrorResponse(ctx.getRpcIdRepr(), code, message)
	return ctx
}

// SetErrorObject sets response as error using the provided error object
func (ctx *RpcContext) SetErrorObject(e errs.Error) *RpcContext {
	ctx.parseFailed = true
	de, ok := e.(errs.DataError)
	if ok {
		ctx.response = createDataErrorResponse(ctx.getRpcIdRepr(), int64(e.ErrorCode()), e.Error(), de.ErrorData())
	} else {
		ctx.response = createErrorResponse(ctx.getRpcIdRepr(), int64(e.ErrorCode()), e.Error())
	}
	return ctx
}

// setResult sets result field of the response using the provided byte slice
func (ctx *RpcContext) setResult(resultRepr []byte) *RpcContext {
	ctx.response = createResponse(ctx.getRpcIdRepr(), resultRepr)
	return ctx
}

// setMethod sets the rpc method
func (ctx *RpcContext) setMethod(method string) *RpcContext {
	ctx.method = method
	return ctx
}

// SetResponse sets whole response using the provided byte slice
func (ctx *RpcContext) SetResponse(response []byte) *RpcContext {
	ctx.response = response
	return ctx
}

// GetMethod returns the method name
func (ctx *RpcContext) GetMethod() string {
	return ctx.method
}

// GetBody returns the body as a byte slice
func (ctx *RpcContext) GetBody() []byte {
	return ctx.body
}
