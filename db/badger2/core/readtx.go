package db

import (
	"bytes"
	"context"
	"fmt"

	"aurora-relayer-go-common/db/badger2/core/dbkey"
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
	"aurora-relayer-go-common/db/badger2/core/dbresponses"
	"aurora-relayer-go-common/db/badger2/core/dbtypes"

	"github.com/dgraph-io/badger/v3"
)

func (txn *ViewTxn) ReadTxKey(chainId uint64, hash dbp.Data32) (*dbtypes.TransactionKey, error) {
	return read[dbtypes.TransactionKey](txn, dbkey.TxKeyByHash.Get(chainId, hash.Content))
}

func (txn *ViewTxn) ReadEarliestTxKey(chainId uint64) (*dbtypes.TransactionKey, error) {
	it := txn.txn.NewIterator(badger.IteratorOptions{
		Prefix: dbkey.Txs.Get(chainId),
	})
	defer it.Close()
	it.Rewind()
	key, err := readTxKeyFromTxIterator(it)
	if err != nil {
		txn.db.logger.Errorf("DB: can't read earliest tx key: %v", err)
	}
	return key, err
}

func (txn *ViewTxn) ReadLatestTxKey(chainId uint64) (*dbtypes.TransactionKey, error) {
	it := txn.txn.NewIterator(badger.IteratorOptions{
		Reverse: true,
		Prefix:  dbkey.Txs.Get(chainId),
	})
	defer it.Close()
	it.Seek(dbkey.Tx.Get(chainId, dbkey.MaxBlockHeight, dbkey.MaxLogIndex))
	key, err := readTxKeyFromTxIterator(it)
	if err != nil {
		txn.db.logger.Errorf("DB: can't read latest tx key: %v", err)
	}
	return key, err
}

func readTxKeyFromTxIterator(it *badger.Iterator) (*dbtypes.TransactionKey, error) {
	if !it.Valid() {
		return nil, nil
	}
	if !dbkey.Tx.Matches(it.Item().Key()) {
		err := fmt.Errorf("found unexpected key format (expected to match dbkey.Tx)")
		return nil, err
	}
	return &dbtypes.TransactionKey{
		BlockHeight:      dbkey.Tx.ReadUintVar(it.Item().Key(), 1),
		TransactionIndex: dbkey.Tx.ReadUintVar(it.Item().Key(), 2),
	}, nil
}

func (txn *ViewTxn) ReadTx(chainId uint64, key dbtypes.TransactionKey) (*dbresponses.Transaction, error) {
	blockHash, err := read[dbp.Data32](txn, dbkey.BlockHash.Get(chainId, key.BlockHeight))
	if err != nil || blockHash == nil {
		return nil, err
	}
	txHash, err := read[dbp.Data32](txn, dbkey.TxHash.Get(chainId, key.BlockHeight, key.TransactionIndex))
	if err != nil || txHash == nil {
		return nil, err
	}
	txData, err := read[dbtypes.Transaction](txn, dbkey.TxData.Get(chainId, key.BlockHeight, key.TransactionIndex))
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

func (txn *ViewTxn) ReadTxReceipt(chainId uint64, key dbtypes.TransactionKey) (*dbresponses.TransactionReceipt, error) {
	blockHash, err := readCached[dbp.Data32](txn, dbkey.BlockHash.Get(chainId, key.BlockHeight))
	if err != nil || blockHash == nil {
		return nil, err
	}
	txHash, err := readCached[dbp.Data32](txn, dbkey.TxHash.Get(chainId, key.BlockHeight, key.TransactionIndex))
	if err != nil || txHash == nil {
		return nil, err
	}
	txData, err := read[dbtypes.Transaction](txn, dbkey.TxData.Get(chainId, key.BlockHeight, key.TransactionIndex))
	if err != nil || txData == nil {
		return nil, err
	}

	logs, err := txn.ReadLogs(
		context.Background(),
		chainId,
		&dbtypes.LogKey{
			BlockHeight:      key.BlockHeight,
			TransactionIndex: key.TransactionIndex,
			LogIndex:         0,
		},
		&dbtypes.LogKey{
			BlockHeight:      key.BlockHeight,
			TransactionIndex: key.TransactionIndex,
			LogIndex:         dbkey.MaxLogIndex,
		},
		nil,
		nil,
		1000,
	)
	if err != nil {
		errCtx := fmt.Sprintf("chainId=%v, block=%v, tx=%v", chainId, key.BlockHeight, key.TransactionIndex)
		txn.db.logger.Errorf("DB: error reading logs for %v: %v", errCtx, err)
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
	from *dbtypes.TransactionKey,
	to *dbtypes.TransactionKey,
	full bool,
	limit int,
) ([]any, error) {

	if limit <= 0 {
		limit = 100_000
	}

	if from.CompareTo(to) > 0 {
		return nil, fmt.Errorf("from > to")
	}

	fromKey := dbkey.Txs.Get(chainId, from.BlockHeight, from.TransactionIndex)
	toKey := dbkey.Txs.Get(chainId, to.BlockHeight, to.TransactionIndex)
	it := txn.txn.NewIterator(badger.IteratorOptions{
		PrefetchValues: true,
		PrefetchSize:   1000,
		Prefix:         getCommonPrefix(fromKey, toKey),
	})
	defer it.Close()

	response := []any{}
	var pendingHashKey *dbtypes.TransactionKey
	var pendingHash dbp.Data32

	for it.Seek(fromKey); it.Valid(); it.Next() {
		select {
		case <-ctx.Done():
			return response, ctx.Err()
		default:
		}
		if bytes.Compare(it.Item().Key(), toKey) > 0 {
			break
		}

		if dbkey.TxHash.Matches(it.Item().Key()) {
			if pendingHashKey != nil {
				txn.db.logger.Errorf("DB: TxHash isn't followed by TxData, will ignore [key=%v]", pendingHashKey)
				pendingHashKey = nil
			}
			txHash, err := readItem[dbp.Data32](txn.db, it.Item())
			if err != nil || txHash == nil {
				txn.db.logger.Errorf("DB: can't read TxHash, will ignore [key=%v]: %v", it.Item().Key(), err)
				continue
			}
			if full {
				pendingHashKey = &dbtypes.TransactionKey{
					BlockHeight:      dbkey.TxHash.ReadUintVar(it.Item().Key(), 1),
					TransactionIndex: dbkey.TxHash.ReadUintVar(it.Item().Key(), 2),
				}
				pendingHash = *txHash
			} else {
				if len(response) == limit {
					return response, ErrLimited
				}
				response = append(response, txHash)
			}
			continue
		}

		if dbkey.TxData.Matches(it.Item().Key()) {
			if full {
				continue
			}
			if pendingHashKey == nil {
				txn.db.logger.Errorf("DB: transaction doesn't start with hash, will ignore [key=%v]", it.Item().Key())
				continue
			}
			key := &dbtypes.TransactionKey{
				BlockHeight:      dbkey.TxData.ReadUintVar(it.Item().Key(), 1),
				TransactionIndex: dbkey.TxData.ReadUintVar(it.Item().Key(), 2),
			}
			if key.BlockHeight != pendingHashKey.BlockHeight || key.TransactionIndex != pendingHashKey.TransactionIndex {
				txn.db.logger.Errorf("DB: TxHash isn't followed by TxData, will ignore [key=%v]", pendingHashKey)
				pendingHashKey = nil
				txn.db.logger.Errorf("DB: transaction doesn't start with hash, will ignore [key=%v]", it.Item().Key())
				continue
			}
			pendingHashKey = nil
			txData, err := readItem[dbtypes.Transaction](txn.db, it.Item())
			if err != nil || txData == nil {
				txn.db.logger.Errorf("DB: can't read TxData, will ignore [key=%v]: %v", it.Item().Key(), err)
				continue
			}
			blockHash, err := readCached[dbp.Data32](txn, dbkey.BlockHash.Get(chainId, key.BlockHeight))
			if err != nil || blockHash == nil {
				txn.db.logger.Errorf("DB: can't read BlockHash for tx, will ignore [key=%v]: %v", it.Item().Key(), err)
				continue
			}
			if len(response) == limit {
				return response, ErrLimited
			}
			response = append(response, makeTransactionResponse(
				chainId,
				key.BlockHeight,
				key.TransactionIndex,
				*blockHash,
				pendingHash,
				txData,
			))
			continue
		}

		txn.db.logger.Errorf("DB: found unknown key while iterating through txs, will ignore: [key=%v]", it.Item().Key())
	}
	if pendingHashKey != nil {
		txn.db.logger.Errorf("DB: TxHash isn't followed by TxData, will ignore [key=%v]", pendingHashKey)
	}
	return response, nil
}
