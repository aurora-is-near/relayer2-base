package endpoint

import (
	"relayer2-base/types/common"
	"relayer2-base/utils"

	"golang.org/x/net/context"
)

type Parity struct {
	*Endpoint
}

func NewParity(endpoint *Endpoint) *Parity {
	return &Parity{endpoint}
}

// PendingTransactions returns a list of txs currently in the queue.
// As of now method always returns empty array since the relayer has no txs queue support.
//
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
func (p *Parity) PendingTransactions(_ context.Context, _ *common.Uint64, _ *interface{}) (*[]string, error) {
	return utils.Constants.EmptyArray(), nil
}
