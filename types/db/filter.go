package db

import (
	tp "relayer2-base/tinypack"
	"relayer2-base/types/primitives"
)

type BlockFilter struct {
	CreatedAt uint64
	Metadata  primitives.VarData
	From      BlockKey
	Next      BlockKey
	To        BlockKey
}

func (f *BlockFilter) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&f.CreatedAt,
		&f.Metadata,
		&f.From,
		&f.Next,
		&f.To,
	}, nil
}

type TransactionFilter struct {
	CreatedAt uint64
	Metadata  primitives.VarData
	From      TransactionKey
	Next      TransactionKey
	To        TransactionKey
}

func (f *TransactionFilter) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&f.CreatedAt,
		&f.Metadata,
		&f.From,
		&f.Next,
		&f.To,
	}, nil
}

type LogFilter struct {
	CreatedAt uint64
	Metadata  primitives.VarData
	From      LogKey
	Next      LogKey
	To        LogKey
	Addresses tp.VarList[primitives.Data20]
	Topics    tp.VarList[tp.VarList[primitives.Data32]]
}

func (f *LogFilter) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&f.CreatedAt,
		&f.Metadata,
		&f.From,
		&f.Next,
		&f.To,
		&f.Addresses,
		&f.Topics,
	}, nil
}
