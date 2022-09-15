package processor

import (
	"aurora-relayer-go-common/endpoint"
	"aurora-relayer-go-common/utils"
	"context"
)

type EnableDisable struct{}

func NewEnableDisable() endpoint.Processor {
	return &EnableDisable{}
}

func (p *EnableDisable) Pre(ctx context.Context, name string, endpoint *endpoint.Endpoint, _ ...any) (context.Context, bool, *any, error) {
	if endpoint.Config.DisabledEndpoints[name] {
		return ctx, true, nil, &utils.MethodNotFoundError{Method: name}
	}
	return ctx, false, nil, nil
}

func (p *EnableDisable) Post(ctx context.Context, _ string, r *any, err *error) (context.Context, *any, *error) {
	return ctx, r, err
}
