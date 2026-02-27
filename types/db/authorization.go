package db

import (
	"github.com/aurora-is-near/relayer2-base/types/primitives"
)

type Authorization struct {
	ChainId uint64
	Address primitives.Data20
	Nonce   uint64
	YParity uint8
	R       primitives.Quantity
	S       primitives.Quantity
}
