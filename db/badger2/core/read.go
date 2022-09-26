package db

import (
	"fmt"

	"aurora-relayer-go-common/db/badger2/core/dbkey"
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
	"aurora-relayer-go-common/db/badger2/core/dbresponses"
	"aurora-relayer-go-common/db/badger2/core/dbtypes"
	"github.com/dgraph-io/badger/v3"
)

// Wrapper that hides internals from outer packages
type ViewTxn struct {
	db  *DB
	txn *badger.Txn
}

func readItem[T any](db *DB, item *badger.Item) (*T, error) {
	res := new(T)
	err := item.Value(func(val []byte) error {
		if err := db.Decoder.Unmarshal(val, res); err != nil {
			return fmt.Errorf("can't unmarshal value of type %T: %v", res, err)
		}
		return nil
	})
	if err != nil {
		db.logger.Errorf("DB: can't read item: %v", err)
		return nil, err
	}
	return res, nil
}

func read[T any](txn *ViewTxn, key []byte) (*T, error) {
	item, err := txn.txn.Get(key)
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	if err != nil {
		txn.db.logger.Errorf("DB: Can't fetch item: %v", err)
	}
	return readItem[T](txn.db, item)
}

func (txn *ViewTxn) ReadBlockKey(chainId uint64, hash dbp.Data32) (*dbtypes.BlockKey, error) {
	return read[dbtypes.BlockKey](txn, dbkey.BlockKeyByHash.Get(chainId, hash.Content))
}

func (txn *ViewTxn) ReadEarliestBlockKey(chainId uint64) (*dbtypes.BlockKey, error) {
	it := txn.txn.NewIterator(badger.IteratorOptions{
		Prefix: dbkey.Blocks.Get(chainId),
	})
	defer it.Close()
	it.Rewind()
	if !it.Valid() {
		return nil, nil
	}
	if !dbkey.Block.Matches(it.Item().Key()) {
		err := fmt.Errorf("found unexpected key format (expected to match dbkey.Block)")
		txn.db.logger.Errorf("DB: can't read earliest block number: %v", err)
		return nil, err
	}
	return &dbtypes.BlockKey{Height: dbkey.Block.ReadUintVar(it.Item().Key(), 1)}, nil
}

func (txn *ViewTxn) ReadLatestBlockKey(chainId uint64) (*dbtypes.BlockKey, error) {
	it := txn.txn.NewIterator(badger.IteratorOptions{
		Reverse: true,
		Prefix:  dbkey.Blocks.Get(chainId),
	})
	defer it.Close()
	it.Seek(dbkey.Block.Get(chainId, 0xffffffff))
	if !it.Valid() {
		return nil, nil
	}
	if !dbkey.Block.Matches(it.Item().Key()) {
		err := fmt.Errorf("found unexpected key format (expected to match dbkey.Block)")
		txn.db.logger.Errorf("DB: can't read latest block number: %v", err)
		return nil, err
	}
	return &dbtypes.BlockKey{Height: dbkey.Block.ReadUintVar(it.Item().Key(), 1)}, nil
}

func (txn *ViewTxn) ReadBlock(chainId uint64, key dbtypes.BlockKey, fullTransactions bool) (*dbresponses.Block, error) {
	height := key.Height
	hash, err := read[dbp.Data32](txn, dbkey.BlockHash.Get(chainId, height))
	if err != nil || hash == nil {
		return nil, err
	}
	data, err := read[dbtypes.Block](txn, dbkey.BlockData.Get(chainId, height))
	if err != nil || data == nil {
		return nil, err
	}

	it := txn.txn.NewIterator(badger.IteratorOptions{
		PrefetchValues: true,
		PrefetchSize:   20,
		Prefix:         dbkey.TxsForBlock.Get(chainId, height),
	})
	defer it.Close()
	expectedIndex := uint64(0)
	txs := []any{}
	for it.Rewind(); it.Valid(); it.Next() {
		if !dbkey.TxHash.Matches(it.Item().Key()) {
			err := fmt.Errorf("transaction doesn't start with hash")
			txn.db.logger.Errorf("DB: error reading transactions for chainId=%v, block=%v: %v", chainId, height, err)
			return nil, err
		}
		txIndex := dbkey.TxHash.ReadUintVar(it.Item().Key(), 2)
		if txIndex != expectedIndex {
			err := fmt.Errorf("unexpected tx index, expected=%v, found=%v", expectedIndex, txIndex)
			txn.db.logger.Errorf("DB: error reading transactions for chainId=%v, block=%v: %v", chainId, height, err)
			return nil, err
		}
		expectedIndex++
		txHash, err := readItem[dbp.Data32](txn.db, it.Item())
		if err != nil {
			return nil, err
		}
		it.Next()
		if !it.Valid() || !dbkey.TxData.Matches(it.Item().Key()) || dbkey.TxData.ReadUintVar(it.Item().Key(), 2) != txIndex {
			err := fmt.Errorf("transaction hash isn't followed by it's data")
			txn.db.logger.Errorf("DB: error reading transactions for chainId=%v, block=%v: %v", chainId, height, err)
			return nil, err
		}
		if !fullTransactions {
			txs = append(txs, *txHash)
			continue
		}
		txData, err := readItem[dbtypes.Transaction](txn.db, it.Item())
		if err != nil {
			return nil, err
		}
		tx := makeTransactionResponse(chainId, height, txIndex, *hash, *txHash, txData)
		txs = append(txs, tx)
	}

	return makeBlockResponse(height, *hash, *data, txs), nil
}

func (txn *ViewTxn) ReadBlockTxCount(chainId uint64, key dbtypes.BlockKey) (dbp.HexUint, error) {
	it := txn.txn.NewIterator(badger.IteratorOptions{
		Reverse: true,
		Prefix:  dbkey.TxsForBlock.Get(chainId, key.Height),
	})
	defer it.Close()
	it.Seek(dbkey.Tx.Get(chainId, key.Height, 0xffff))
	if !it.Valid() {
		return 0, nil
	}
	if !dbkey.Tx.Matches(it.Item().Key()) {
		err := fmt.Errorf("found unexpected key format (expected to match dbkey.Tx)")
		txn.db.logger.Errorf("DB: can't read block tx count: %v", err)
		return 0, err
	}
	return dbp.HexUint(dbkey.Tx.ReadUintVar(it.Item().Key(), 2)) + 1, nil
}

func (txn *ViewTxn) ReadTxKey(chainId uint64, hash dbp.Data32) (*dbtypes.TransactionKey, error) {
	return read[dbtypes.TransactionKey](txn, dbkey.TxKeyByHash.Get(chainId, hash.Content))
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

	it := txn.txn.NewIterator(badger.IteratorOptions{
		PrefetchValues: true,
		PrefetchSize:   20,
		Prefix:         dbkey.LogsForTx.Get(chainId, key.BlockHeight, key.TransactionIndex),
	})
	defer it.Close()
	expectedIndex := uint64(0)
	logs := []*dbresponses.Log{}
	var anyErr error
	for it.Rewind(); it.Valid(); it.Next() {
		if !dbkey.Log.Matches(it.Item().Key()) {
			anyErr = fmt.Errorf("log key doesn't match the pattern of dbkey.Log")
			break
		}
		logIndex := dbkey.Log.ReadUintVar(it.Item().Key(), 3)
		if logIndex != expectedIndex {
			anyErr = fmt.Errorf("unexpected log index, expected=%v, found=%v", expectedIndex, logIndex)
			break
		}
		expectedIndex++
		var logData *dbtypes.Log
		if logData, anyErr = readItem[dbtypes.Log](txn.db, it.Item()); anyErr != nil {
			break
		}
		logs = append(logs, makeLogResponse(
			key.BlockHeight,
			key.TransactionIndex,
			logIndex,
			*blockHash,
			*txHash,
			logData,
		))
	}
	if anyErr != nil {
		errCtx := fmt.Sprintf("chainId=%v, block=%v, tx=%v", chainId, key.BlockHeight, key.TransactionIndex)
		txn.db.logger.Errorf("DB: error reading logs for %v: %v", errCtx, anyErr)
		return nil, anyErr
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

func (txn *ViewTxn) ReadLogs(chainId uint64, filter *dbtypes.Filter) ([]*dbresponses.Log, error) {
	// TODO
	return []*dbresponses.Log{}, nil
}
