package core

import (
	"github.com/aurora-is-near/relayer2-base/db/badger/core/dbkey"
	dbp "github.com/aurora-is-near/relayer2-base/types/primitives"
)

func (txn *ViewTxn) DeleteFilter(chainId uint64, filterId dbp.Data32) error {
	return txn.txn.Delete(dbkey.LogFilter.Get(chainId, filterId.Bytes()))
}

func (txn *ViewTxn) DeleteBlockFilter(chainId uint64, filterId dbp.Data32) error {
	return txn.txn.Delete(dbkey.BlockFilter.Get(chainId, filterId.Bytes()))
}

func (txn *ViewTxn) DeleteTransactionFilter(chainId uint64, filterId dbp.Data32) error {
	return txn.txn.Delete(dbkey.TxFilter.Get(chainId, filterId.Bytes()))
}

func (txn *ViewTxn) DeleteLogFilter(chainId uint64, filterId dbp.Data32) error {
	return txn.txn.Delete(dbkey.LogFilter.Get(chainId, filterId.Bytes()))
}
