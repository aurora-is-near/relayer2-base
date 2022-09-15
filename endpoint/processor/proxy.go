package processor

import (
	"aurora-relayer-go-common/endpoint"
	"context"
)

type Proxy struct{}

func NewProxy() endpoint.Processor {
	return &Proxy{}
}

func (p *Proxy) Pre(ctx context.Context, name string, endpoint *endpoint.Endpoint, _ ...any) (context.Context, bool, *any, error) {
	if endpoint.Config.ProxyEndpoints[name] {
		// TODO send request to proxy server
		return ctx, true, nil, nil
	}
	return ctx, false, nil, nil
}

func (p *Proxy) Post(ctx context.Context, _ string, r *any, err *error) (context.Context, *any, *error) {
	return ctx, r, err
}
