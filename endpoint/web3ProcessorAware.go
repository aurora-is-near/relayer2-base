package endpoint

import (
	"context"
)

type Web3ProcessorAware struct {
	*Web3
}

func NewWeb3ProcessorAware(web3 *Web3) *Web3ProcessorAware {
	return &Web3ProcessorAware{web3}
}

func (e *Web3ProcessorAware) ClientVersion(ctx context.Context) (string, error) {
	return Process(ctx, "web3_clientVersion", e.Endpoint, func(ctx context.Context) (string, error) {
		return e.Web3.ClientVersion(ctx)
	})
}

func (e *Web3ProcessorAware) Sha3(ctx context.Context, in string) (string, error) {
	return Process(ctx, "web3_sha3", e.Endpoint, func(ctx context.Context) (string, error) {
		return e.Web3.Sha3(ctx, in)
	}, in)
}
