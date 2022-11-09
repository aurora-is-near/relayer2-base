package dbtypes

import (
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
	tp "aurora-relayer-go-common/tinypack"
)

type BlockFilter struct {
	CreatedAt uint64
	Metadata  dbp.VarData
	Next      BlockKey
	Last      BlockKey
}

func (f *BlockFilter) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&f.CreatedAt,
		&f.Metadata,
		&f.Next,
		&f.Last,
	}, nil
}

type TransactionFilter struct {
	CreatedAt uint64
	Metadata  dbp.VarData
	Next      TransactionKey
	Last      TransactionKey
}

func (f *TransactionFilter) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&f.CreatedAt,
		&f.Metadata,
		&f.Next,
		&f.Last,
	}, nil
}

type LogFilter struct {
	CreatedAt uint64
	Metadata  dbp.VarData
	Next      LogKey
	Last      LogKey
	Addresses tp.VarList[dbp.Data20]
	Topics    tp.VarList[tp.VarList[dbp.Data32]]
}

func (f *LogFilter) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&f.CreatedAt,
		&f.Metadata,
		&f.Next,
		&f.Last,
		&f.Addresses,
		&f.Topics,
	}, nil
}
