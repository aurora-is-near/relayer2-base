package core

import (
	"github.com/aurora-is-near/relayer2-base/db/badger/core/dbkey"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
)

func (txn *ViewTxn) ReadIndexerState(chainId uint64) ([]byte, error) {
	data, err := read[primitives.VarData](txn, dbkey.IndexerState.Get(chainId))
	if err != nil {
		return nil, err
	}
	if data != nil {
		return data.Bytes(), nil
	}
	return nil, nil
}
