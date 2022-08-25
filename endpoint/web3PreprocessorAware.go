package endpoint

import (
	"context"
)

type Web3PreprocessorAware struct {
	*Web3
}

func NewWeb3PreprocessorAware(web3 *Web3) *Web3PreprocessorAware {
	return &Web3PreprocessorAware{web3}
}

func (e *Web3PreprocessorAware) ClientVersion(ctx context.Context) (string, error) {
	return Preprocess("web3_clientVersion", e.Endpoint, func() (string, error) {
		return e.Web3.ClientVersion(ctx)
	}, ctx)
}

func (e Web3PreprocessorAware) Sha3(ctx context.Context, in *string) (string, error) {
	return Preprocess("web3_sha3", e.Endpoint, func() (string, error) {
		return e.Web3.Sha3(ctx, in)
	}, ctx, in)
}
