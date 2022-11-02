package db

import (
	"aurora-relayer-go-common/db/badger2/core/dbkey"
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
	"aurora-relayer-go-common/db/badger2/core/dbtypes"
)

func (txn *ViewTxn) ReadBlockFilter(chainId uint64, filterId dbp.Data32) (*dbtypes.BlockFilter, error) {
	return read[dbtypes.BlockFilter](txn, dbkey.BlockFilter.Get(chainId, filterId.Bytes()))
}

func (txn *ViewTxn) ReadTransactionFilter(chainId uint64, filterId dbp.Data32) (*dbtypes.TransactionFilter, error) {
	return read[dbtypes.TransactionFilter](txn, dbkey.TxFilter.Get(chainId, filterId.Bytes()))
}

func (txn *ViewTxn) ReadLogFilter(chainId uint64, filterId dbp.Data32) (*dbtypes.LogFilter, error) {
	return read[dbtypes.LogFilter](txn, dbkey.LogFilter.Get(chainId, filterId.Bytes()))
}

func (txn *ViewTxn) ReadFilter(chainId uint64, filterId dbp.Data32) (any, error) {
	if blockFilter, err := txn.ReadBlockFilter(chainId, filterId); blockFilter != nil || err != nil {
		return blockFilter, err
	}
	if txFilter, err := txn.ReadTransactionFilter(chainId, filterId); txFilter != nil || err != nil {
		return txFilter, err
	}
	return txn.ReadLogFilter(chainId, filterId)
}
