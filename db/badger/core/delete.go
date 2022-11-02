package core

import (
	"aurora-relayer-go-common/db/badger/core/dbkey"
	dbp "aurora-relayer-go-common/types/primitives"
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
