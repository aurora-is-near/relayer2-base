package processor

import (
	"context"
	"github.com/aurora-is-near/relayer2-base/endpoint"
	errs "github.com/aurora-is-near/relayer2-base/types/errors"
)

type EnableDisable struct{}

func NewEnableDisable() endpoint.Processor {
	return &EnableDisable{}
}

func (p *EnableDisable) Pre(ctx context.Context, name string, endpoint *endpoint.Endpoint, _ *any, _ ...any) (context.Context, bool, error) {
	if endpoint.Config.DisabledEndpoints[name] {
		return ctx, true, &errs.MethodNotFoundError{Method: name}
	}
	return ctx, false, nil
}

func (p *EnableDisable) Post(ctx context.Context, _ string, _ *any, _ *error) context.Context {
	return ctx
}
