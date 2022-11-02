package db

import (
	tp "aurora-relayer-go-common/tinypack"
	"aurora-relayer-go-common/types/primitives"
)

type AccessListEntry struct {
	Address     primitives.Data20
	StorageKeys tp.VarList[primitives.Data32]
}

func (e *AccessListEntry) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&e.Address,
		&e.StorageKeys,
	}, nil
}
