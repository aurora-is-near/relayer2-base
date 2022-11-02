package db

import (
	tp "aurora-relayer-go-common/tinypack"
	"aurora-relayer-go-common/types/primitives"
)

type Log struct {
	Address primitives.Data20
	Data    primitives.VarData
	Topics  tp.VarList[primitives.Data32]
}

func (l *Log) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&l.Address,
		&l.Data,
		&l.Topics,
	}, nil
}
