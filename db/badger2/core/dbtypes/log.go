package dbtypes

import (
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
	tp "aurora-relayer-go-common/tinypack"
)

type Log struct {
	Address dbp.Data20
	Data    dbp.VarData
	Topics  tp.VarList[dbp.Data32]
}

func (l *Log) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&l.Address,
		&l.Data,
		&l.Topics,
	}, nil
}
