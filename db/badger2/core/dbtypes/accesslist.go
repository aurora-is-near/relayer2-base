package dbtypes

import (
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
	tp "aurora-relayer-go-common/tinypack"
)

type AccessListEntry struct {
	Address     dbp.Data20
	StorageKeys tp.VarList[dbp.Data32]
}

func (e *AccessListEntry) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&e.Address,
		&e.StorageKeys,
	}, nil
}
