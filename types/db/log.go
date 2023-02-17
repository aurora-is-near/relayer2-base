package db

import (
	tp "github.com/aurora-is-near/relayer2-base/tinypack"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
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
