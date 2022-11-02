package db

import (
	"aurora-relayer-go-common/types/primitives"
)

type Block struct {
	ParentHash       primitives.Data32
	Miner            primitives.Data20
	Timestamp        uint64
	GasLimit         primitives.Quantity
	GasUsed          primitives.Quantity
	LogsBloom        primitives.Data256
	TransactionsRoot primitives.Data32
	StateRoot        primitives.Data32
	ReceiptsRoot     primitives.Data32
	Size             uint64
}

func (b *Block) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&b.ParentHash,
		&b.Miner,
		&b.Timestamp,
		&b.GasLimit,
		&b.GasUsed,
		&b.LogsBloom,
		&b.TransactionsRoot,
		&b.StateRoot,
		&b.ReceiptsRoot,
		&b.Size,
	}, nil
}
