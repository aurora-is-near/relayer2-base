package github_neonxp_jsonrpc2

import (
	"context"
	"fmt"

	"go.neonxp.dev/jsonrpc2/rpc"
	"go.neonxp.dev/jsonrpc2/transport"

	"github.com/aurora-is-near/relayer2-base/log"
)

type JsonRpc2 struct {
	rpc.RpcServer
}

func New(config *Config) (*JsonRpc2, error) {
	n := rpc.New(
		rpc.WithLogger(NewNeonxpJsonRpc2Logger(log.Log())),
		rpc.WithTransport(&transport.HTTP{Bind: fmt.Sprintf("%s:%d", config.HttpHost, config.HttpPort)}),
		rpc.WithMiddleware(DefaultParams()),
	)
	return &JsonRpc2{*n}, nil
}

func (n *JsonRpc2) Start() error {
	return n.Run(context.Background())
}
