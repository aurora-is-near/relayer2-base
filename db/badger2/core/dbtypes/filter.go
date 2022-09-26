package dbtypes

import (
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
	tp "aurora-relayer-go-common/tinypack"
)

type Filter struct {
	Type      dbp.VarData
	CreatedBy dbp.VarData
	PollBlock uint64
	FromBlock tp.Nullable[dbp.Data32]
	ToBlock   tp.Nullable[dbp.Data32]
	Addresses tp.VarList[dbp.Data20]
	Topics    tp.VarList[tp.VarList[dbp.Data32]]
}

func (f *Filter) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&f.Type,
		&f.CreatedBy,
		&f.PollBlock,
		&f.FromBlock,
		&f.ToBlock,
		&f.Addresses,
		&f.Topics,
	}, nil
}
