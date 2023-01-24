package endpoint

import (
	"relayer2-base/types/common"
	"relayer2-base/types/response"

	"golang.org/x/net/context"
)

type DebugProcessorAware struct {
	*Debug
}

func NewDebugProcessorAware(d *Debug) *DebugProcessorAware {
	return &DebugProcessorAware{d}
}

func (d *DebugProcessorAware) TraceTransaction(ctx context.Context, hash common.H256) (*response.CallFrame, error) {
	return Process(ctx, "debug_traceTransaction", d.Endpoint, func(ctx context.Context) (*response.CallFrame, error) {
		return d.Debug.TraceTransaction(ctx, hash)
	}, hash)
}
