package github_neonxp_jsonrpc2

import (
	"context"

	"go.neonxp.dev/jsonrpc2/rpc"
)

func DefaultParams() rpc.Middleware {
	defaultParams := [2]byte{'[', ']'}
	return func(handler rpc.RpcHandler) rpc.RpcHandler {
		return func(ctx context.Context, req *rpc.RpcRequest) *rpc.RpcResponse {
			if req.Params == nil {
				req.Params = defaultParams[:]
			}
			return handler(ctx, req)
		}
	}
}
