package rpc

import (
	"context"
)

type RpcHandler func(ctx *context.Context, rpcCtx *RpcContext) *RpcContext

type Middleware func(handler RpcHandler) RpcHandler
