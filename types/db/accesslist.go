package db

import (
	tp "relayer2-base/tinypack"
	"relayer2-base/types/primitives"
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
