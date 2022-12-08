package endpoint

import (
	"aurora-relayer-go-common/types/common"

	"golang.org/x/net/context"
)

type ParityProcessorAware struct {
	*Parity
}

func NewParityProcessorAware(p *Parity) *ParityProcessorAware {
	return &ParityProcessorAware{p}
}

func (e *ParityProcessorAware) PendingTransactions(ctx context.Context, limit *common.Uint64, filter *interface{}) (*[]string, error) {
	return Process(ctx, "parity_pendingTransactions", e.Endpoint, func(ctx context.Context) (*[]string, error) {
		return e.Parity.PendingTransactions(ctx, limit, filter)
	})
}
