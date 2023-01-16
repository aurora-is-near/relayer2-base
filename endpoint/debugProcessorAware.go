package endpoint

import (
	"aurora-relayer-go-common/types/common"
	"aurora-relayer-go-common/types/response"

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
