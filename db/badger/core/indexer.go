package core

import (
	"aurora-relayer-go-common/db/badger/core/dbkey"
	"aurora-relayer-go-common/types/primitives"
)

func (txn *ViewTxn) ReadIndexerState(chainId uint64) ([]byte, error) {
	data, err := read[primitives.VarData](txn, dbkey.IndexerState.Get(chainId))
	if err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}
