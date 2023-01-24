package core

import (
	"relayer2-base/db/badger/core/dbkey"
	"relayer2-base/types/primitives"
)

func (txn *ViewTxn) ReadIndexerState(chainId uint64) ([]byte, error) {
	data, err := read[primitives.VarData](txn, dbkey.IndexerState.Get(chainId))
	if err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}
