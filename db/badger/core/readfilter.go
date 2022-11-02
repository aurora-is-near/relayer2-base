package core

import (
	"aurora-relayer-go-common/db/badger/core/dbkey"
	"aurora-relayer-go-common/types/db"
	dbp "aurora-relayer-go-common/types/primitives"
)

func (txn *ViewTxn) ReadBlockFilter(chainId uint64, filterId dbp.Data32) (*db.BlockFilter, error) {
	return read[db.BlockFilter](txn, dbkey.BlockFilter.Get(chainId, filterId.Bytes()))
}

func (txn *ViewTxn) ReadTransactionFilter(chainId uint64, filterId dbp.Data32) (*db.TransactionFilter, error) {
	return read[db.TransactionFilter](txn, dbkey.TxFilter.Get(chainId, filterId.Bytes()))
}

func (txn *ViewTxn) ReadLogFilter(chainId uint64, filterId dbp.Data32) (*db.LogFilter, error) {
	return read[db.LogFilter](txn, dbkey.LogFilter.Get(chainId, filterId.Bytes()))
}

func (txn *ViewTxn) ReadFilter(chainId uint64, filterId dbp.Data32) (any, error) {
	var err error
	if blockFilter, err := txn.ReadBlockFilter(chainId, filterId); blockFilter != nil || err != nil {
		return blockFilter, err
	}
	if txFilter, err := txn.ReadTransactionFilter(chainId, filterId); txFilter != nil || err != nil {
		return txFilter, err
	}
	if logFilter, err := txn.ReadLogFilter(chainId, filterId); logFilter != nil || err != nil {
		return logFilter, err
	}
	return nil, err
}
