package core

import (
	"bytes"
	"context"
	"fmt"
	dbt "relayer2-base/types/db"
	"relayer2-base/types/primitives"
	"relayer2-base/types/response"

	"github.com/dgraph-io/badger/v3"
	"relayer2-base/db/badger/core/dbkey"
)

func (txn *ViewTxn) ReadTxKey(chainId uint64, hash primitives.Data32) (*dbt.TransactionKey, error) {
	return read[dbt.TransactionKey](txn, dbkey.TxKeyByHash.Get(chainId, hash.Bytes()))
}

func (txn *ViewTxn) ReadEarliestTxKey(chainId uint64) (*dbt.TransactionKey, error) {
	it := txn.txn.NewIterator(badger.IteratorOptions{
		Prefix: dbkey.TxHashes.Get(chainId),
	})
	defer it.Close()
	it.Rewind()
	key, err := readTxKeyFromTxHashIterator(it)
	if err != nil {
		txn.db.logger.Errorf("DB: can't read earliest tx key: %v", err)
	}
	return key, err
}

func (txn *ViewTxn) ReadLatestTxKey(chainId uint64) (*dbt.TransactionKey, error) {
	it := txn.txn.NewIterator(badger.IteratorOptions{
		Reverse: true,
		Prefix:  dbkey.TxHashes.Get(chainId),
	})
	defer it.Close()
	it.Seek(dbkey.TxHash.Get(chainId, dbkey.MaxBlockHeight, dbkey.MaxTxIndex))
	key, err := readTxKeyFromTxHashIterator(it)
	if err != nil {
		txn.db.logger.Errorf("DB: can't read latest tx key: %v", err)
	}
	return key, err
}

func readTxKeyFromTxHashIterator(it *badger.Iterator) (*dbt.TransactionKey, error) {
	if !it.Valid() {
		return nil, nil
	}
	if !dbkey.TxHash.Matches(it.Item().Key()) {
		err := fmt.Errorf("found unexpected key format (expected to match dbkey.TxHash)")
		return nil, err
	}
	return &dbt.TransactionKey{
		BlockHeight:      dbkey.TxHash.ReadUintVar(it.Item().Key(), 1),
		TransactionIndex: dbkey.TxHash.ReadUintVar(it.Item().Key(), 2),
	}, nil
}

func (txn *ViewTxn) ReadTx(chainId uint64, key dbt.TransactionKey) (*response.Transaction, error) {
	blockHash, err := read[primitives.Data32](txn, dbkey.BlockHash.Get(chainId, key.BlockHeight))
	if err != nil || blockHash == nil {
		return nil, err
	}
	txHash, err := read[primitives.Data32](txn, dbkey.TxHash.Get(chainId, key.BlockHeight, key.TransactionIndex))
	if err != nil || txHash == nil {
		return nil, err
	}
	txData, err := read[dbt.Transaction](txn, dbkey.TxData.Get(chainId, key.BlockHeight, key.TransactionIndex))
	if err != nil || txData == nil {
		return nil, err
	}
	return makeTransactionResponse(
		chainId,
		key.BlockHeight,
		key.TransactionIndex,
		*blockHash,
		*txHash,
		txData,
	), nil
}

func (txn *ViewTxn) ReadTxReceipt(chainId uint64, key dbt.TransactionKey) (*response.TransactionReceipt, error) {
	blockHash, err := readCached[primitives.Data32](txn, dbkey.BlockHash.Get(chainId, key.BlockHeight))
	if err != nil || blockHash == nil {
		return nil, err
	}
	txHash, err := readCached[primitives.Data32](txn, dbkey.TxHash.Get(chainId, key.BlockHeight, key.TransactionIndex))
	if err != nil || txHash == nil {
		return nil, err
	}
	txData, err := read[dbt.Transaction](txn, dbkey.TxData.Get(chainId, key.BlockHeight, key.TransactionIndex))
	if err != nil || txData == nil {
		return nil, err
	}

	logs, _, err := txn.ReadLogs(
		context.Background(),
		chainId,
		&dbt.LogKey{
			BlockHeight:      key.BlockHeight,
			TransactionIndex: key.TransactionIndex,
			LogIndex:         0,
		},
		&dbt.LogKey{
			BlockHeight:      key.BlockHeight,
			TransactionIndex: key.TransactionIndex,
			LogIndex:         dbkey.MaxLogIndex,
		},
		nil,
		nil,
		int(dbkey.MaxLogIndex)+1,
	)
	if err != nil {
		errCtx := fmt.Sprintf("chainId=%v, block=%v, tx=%v", chainId, key.BlockHeight, key.TransactionIndex)
		txn.db.logger.Errorf("DB: errors reading logs for %v: %v", errCtx, err)
		return nil, err
	}

	return makeTransactionReceiptResponse(
		key.BlockHeight,
		key.TransactionIndex,
		*blockHash,
		*txHash,
		txData,
		logs,
	), nil
}

func (txn *ViewTxn) ReadTransactions(
	ctx context.Context,
	chainId uint64,
	from *dbt.TransactionKey,
	to *dbt.TransactionKey,
	full bool,
	limit int,
) ([]any, *dbt.TransactionKey, error) {

	if limit <= 0 {
		limit = 100_000
	}

	if from.CompareTo(to) > 0 {
		return nil, nil, fmt.Errorf("from > to")
	}

	fromHashKey := dbkey.TxHash.Get(chainId, from.BlockHeight, from.TransactionIndex)
	toHashKey := dbkey.TxHash.Get(chainId, to.BlockHeight, to.TransactionIndex)
	hashIt := txn.txn.NewIterator(badger.IteratorOptions{
		PrefetchValues: true,
		PrefetchSize:   1000,
		Prefix:         getCommonPrefix(fromHashKey, toHashKey),
	})
	defer hashIt.Close()
	hashIt.Seek(fromHashKey)

	fromDataKey := dbkey.TxData.Get(chainId, from.BlockHeight, from.TransactionIndex)
	toDataKey := dbkey.TxData.Get(chainId, to.BlockHeight, to.TransactionIndex)
	var dataIt *badger.Iterator
	if full {
		dataIt = txn.txn.NewIterator(badger.IteratorOptions{
			PrefetchValues: true,
			PrefetchSize:   1000,
			Prefix:         getCommonPrefix(fromDataKey, toDataKey),
		})
		defer dataIt.Close()
		dataIt.Seek(fromDataKey)
	}

	response := []any{}
	lastKey := from.Prev()
	for {
		select {
		case <-ctx.Done():
			return response, lastKey, ctx.Err()
		default:
		}

		if !hashIt.Valid() || bytes.Compare(hashIt.Item().Key(), toHashKey) > 0 {
			break
		}
		if !dbkey.TxHash.Matches(hashIt.Item().Key()) {
			txn.db.logger.Errorf("DB: key was expected to match dbkey.TxHash, found %v, will ignore", hashIt.Item().Key())
			hashIt.Next()
			continue
		}
		txHash, err := readItem[primitives.Data32](txn.db, hashIt.Item())
		if err != nil || txHash == nil {
			txn.db.logger.Errorf("DB: can't read TxHash, will ignore [key=%v]: %v", hashIt.Item().Key(), err)
			hashIt.Next()
			continue
		}
		hashKey := &dbt.TransactionKey{
			BlockHeight:      dbkey.TxHash.ReadUintVar(hashIt.Item().Key(), 1),
			TransactionIndex: dbkey.TxHash.ReadUintVar(hashIt.Item().Key(), 2),
		}
		if !full {
			if len(response) == limit {
				return response, lastKey, ErrLimited
			}
			response = append(response, *txHash)
			lastKey = hashKey
			hashIt.Next()
			continue
		}

		if !dataIt.Valid() || bytes.Compare(dataIt.Item().Key(), toDataKey) > 0 {
			txn.db.logger.Errorf("DB: found dangling TxHash, will ignore [key=%v]", hashIt.Item().Key())
			break
		}
		if !dbkey.TxData.Matches(dataIt.Item().Key()) {
			txn.db.logger.Errorf("DB: key was expected to match dbkey.TxData, found %v", dataIt.Item().Key())
			dataIt.Next()
			continue
		}
		txData, err := readItem[dbt.Transaction](txn.db, dataIt.Item())
		if err != nil || txData == nil {
			txn.db.logger.Errorf("DB: can't read TxData, will ignore [key=%v]: %v", dataIt.Item().Key(), err)
			dataIt.Next()
			continue
		}
		dataKey := &dbt.TransactionKey{
			BlockHeight:      dbkey.TxHash.ReadUintVar(dataIt.Item().Key(), 1),
			TransactionIndex: dbkey.TxHash.ReadUintVar(dataIt.Item().Key(), 2),
		}

		keysCompare := hashKey.CompareTo(dataKey)
		if keysCompare < 0 {
			txn.db.logger.Errorf("DB: found dangling TxHash, will ignore [key=%v]", hashIt.Item().Key())
			hashIt.Next()
			continue
		}
		if keysCompare > 0 {
			txn.db.logger.Errorf("DB: found dangling TxData, will ignore [key=%v]", dataIt.Item().Key())
			dataIt.Next()
			continue
		}

		blockHash, err := readCached[primitives.Data32](txn, dbkey.BlockHash.Get(chainId, hashKey.BlockHeight))
		if err != nil || blockHash == nil {
			txn.db.logger.Errorf("DB: can't read BlockHash for tx, will ignore [key=%v]: %v", hashIt.Item().Key(), err)
			hashIt.Next()
			dataIt.Next()
			continue
		}

		if len(response) == limit {
			return response, lastKey, ErrLimited
		}
		response = append(response, makeTransactionResponse(
			chainId,
			hashKey.BlockHeight,
			hashKey.TransactionIndex,
			*blockHash,
			*txHash,
			txData,
		))
		lastKey = hashKey

		hashIt.Next()
		dataIt.Next()
	}
	return response, to, nil
}
