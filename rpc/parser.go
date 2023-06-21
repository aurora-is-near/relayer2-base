package rpc

import (
	"errors"
	"fmt"

	"github.com/buger/jsonparser"

	"github.com/aurora-is-near/relayer2-base/rpc/types"
	errs "github.com/aurora-is-near/relayer2-base/types/errors"
)

const (
	maxBatchLevel = 1
)

// prepareRpcContext generates parsed rpc context using the receieved data
func (r *RpcServer) prepareRpcContext(body []byte, wsCtx *WebSocketContext) *RpcContext {
	_, dataType, _, err := jsonparser.Get(body)

	rpcCtx, _, err := r.parseRpcHierarchy(body, dataType, err, 0, r.maxBatchRequests, wsCtx)
	if err != nil {
		rpcCtx.SetErrorObject(&errs.GenericError{Err: err})
	} else {
		r.parseRpcContext(rpcCtx)
	}

	return rpcCtx
}

// parseRpcHierarchy generates unparsed rpc context using the receieved data
// returns the created RpcContext, number of requests each level and error, respectively
func (r *RpcServer) parseRpcHierarchy(
	body []byte,
	dataType jsonparser.ValueType,
	parseErr error,
	batchLevel uint,
	maxRequests uint,
	wsCtx *WebSocketContext,
) (*RpcContext, uint, error) {

	if batchLevel > maxBatchLevel {
		return nil, 0, fmt.Errorf("batch level exceeded")
	}
	if maxRequests == 0 {
		return nil, 0, fmt.Errorf("batch requests exceeded")
	}

	rpcCtx := r.newRpcContext(wsCtx, body, dataType, parseErr)
	if rpcCtx.bodyParseError != nil || rpcCtx.bodyDataType != jsonparser.Array {
		return rpcCtx, 1, nil
	}

	var anyErr error
	requestsCnt := uint(1)
	_, rpcCtx.bodyParseError = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if anyErr != nil {
			return
		}
		ch, rc, cerr := r.parseRpcHierarchy(value, dataType, err, batchLevel+1, maxRequests-requestsCnt, wsCtx)
		if cerr != nil {
			anyErr = cerr
			return
		}
		rpcCtx.batchChildren = append(rpcCtx.batchChildren, ch)
		requestsCnt += rc
	})
	return rpcCtx, requestsCnt, anyErr
}

// parseRpcContext checks the rpc context and calls the related batch or single request parser
func (r *RpcServer) parseRpcContext(rpcCtx *RpcContext) *RpcContext {
	if rpcCtx.bodyParseError != nil {
		return rpcCtx.setError(errs.ParseError, fmt.Sprintf("can't parse jsonrpc object: %v", rpcCtx.bodyParseError))
	}

	switch rpcCtx.bodyDataType {
	case jsonparser.Array:
		return r.parseBatchRequest(rpcCtx)
	case jsonparser.Object:
		return r.parseSingleRequest(rpcCtx)
	default:
		return rpcCtx.SetErrorObject(&errs.InvalidRequestError{Message: "unknown request format"})
	}
}

// parseBatchRequest parses the batch rpc requests and updates the batchChiledren in RpcContext object
func (r *RpcServer) parseBatchRequest(rpcCtx *RpcContext) *RpcContext {
	for _, child := range rpcCtx.batchChildren {
		r.parseRpcContext(child)
	}

	return rpcCtx
}

// parseSingleRequest parses the single rpc request and updates the RpcContext object accordingly
func (r *RpcServer) parseSingleRequest(rpcCtx *RpcContext) *RpcContext {
	var err error
	rpcCtx.parsedBody, err = types.ParseRPCRequestBody(rpcCtx.body)
	if err != nil {
		return rpcCtx.setError(errs.ParseError, fmt.Sprintf("can't parse json rpc body %v", err))
	}
	if len(rpcCtx.parsedBody.Method.Value) == 0 {
		return rpcCtx.SetErrorObject(&errs.GenericError{Err: errors.New("method name can't be empty")})
	}

	rpcCtx.setMethod(rpcCtx.parsedBody.Method.Str())

	return rpcCtx
}
