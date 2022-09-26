package dbtypes

import dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"

type Block struct {
	ParentHash       dbp.Data32
	Miner            dbp.Data20
	Timestamp        uint64
	GasLimit         dbp.Quantity
	GasUsed          dbp.Quantity
	LogsBloom        dbp.Data256
	TransactionsRoot dbp.Data32
	StateRoot        dbp.Data32
	ReceiptsRoot     dbp.Data32
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
