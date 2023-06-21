package rpc

import (
	"context"
)

type TransportOption func(s *RpcServer)

func WithTransport(transport Transport) TransportOption {
	return func(s *RpcServer) {
		s.transports = append(s.transports, transport)
	}
}

type Transport interface {
	Run(ctx context.Context, resolver Resolver) error
	Stop() error
}

type Resolver interface {
	ResolveHttp(ctx *context.Context, rpcMessage []byte) []byte
	ResolveWs(ctx *context.Context, wsCtx *WebSocketContext, rpcMessage []byte) []byte
	CloseWsConn(wsCtx *WebSocketContext)
}
