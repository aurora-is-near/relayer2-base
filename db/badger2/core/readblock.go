package db

import (
	"bytes"
	"context"
	"fmt"

	"aurora-relayer-go-common/db/badger2/core/dbkey"
	"aurora-relayer-go-common/db/badger2/core/dbprimitives"
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
	"aurora-relayer-go-common/db/badger2/core/dbresponses"
	"aurora-relayer-go-common/db/badger2/core/dbtypes"

	"github.com/dgraph-io/badger/v3"
)

func (txn *ViewTxn) ReadBlockKey(chainId uint64, hash dbp.Data32) (*dbtypes.BlockKey, error) {
	return read[dbtypes.BlockKey](txn, dbkey.BlockKeyByHash.Get(chainId, hash.Content))
}

func (txn *ViewTxn) ReadEarliestBlockKey(chainId uint64) (*dbtypes.BlockKey, error) {
	it := txn.txn.NewIterator(badger.IteratorOptions{
		Prefix: dbkey.BlockHashes.Get(chainId),
	})
	defer it.Close()
	it.Rewind()
	key, err := readBlockKeyFromBlockHashIterator(it)
	if err != nil {
		txn.db.logger.Errorf("DB: can't read earliest block key: %v", err)
	}
	return key, err
}

func (txn *ViewTxn) ReadLatestBlockKey(chainId uint64) (*dbtypes.BlockKey, error) {
	it := txn.txn.NewIterator(badger.IteratorOptions{
		Reverse: true,
		Prefix:  dbkey.BlockHashes.Get(chainId),
	})
	defer it.Close()
	it.Seek(dbkey.BlockHash.Get(chainId, dbkey.MaxBlockHeight))
	key, err := readBlockKeyFromBlockHashIterator(it)
	if err != nil {
		txn.db.logger.Errorf("DB: can't read latest block key: %v", err)
	}
	return key, err
}

func readBlockKeyFromBlockHashIterator(it *badger.Iterator) (*dbtypes.BlockKey, error) {
	if !it.Valid() {
		return nil, nil
	}
	if !dbkey.BlockHash.Matches(it.Item().Key()) {
		return nil, fmt.Errorf("found unexpected key format (expected to match dbkey.BlockHash)")
	}
	return &dbtypes.BlockKey{Height: dbkey.BlockHash.ReadUintVar(it.Item().Key(), 1)}, nil
}

func (txn *ViewTxn) ReadBlock(chainId uint64, key dbtypes.BlockKey, fullTransactions bool) (*dbresponses.Block, error) {
	height := key.Height
	hash, err := readCached[dbp.Data32](txn, dbkey.BlockHash.Get(chainId, height))
	if err != nil || hash == nil {
		return nil, err
	}
	data, err := read[dbtypes.Block](txn, dbkey.BlockData.Get(chainId, height))
	if err != nil || data == nil {
		return nil, err
	}

	txs, err := txn.ReadTransactions(
		context.Background(),
		chainId,
		&dbtypes.TransactionKey{
			BlockHeight:      key.Height,
			TransactionIndex: 0,
		},
		&dbtypes.TransactionKey{
			BlockHeight:      key.Height,
			TransactionIndex: dbkey.MaxTxIndex,
		},
		fullTransactions,
		1000,
	)
	if err != nil {
		errCtx := fmt.Sprintf("chainId=%v, block=%v", chainId, key.Height)
		txn.db.logger.Errorf("DB: error reading txs for %v: %v", errCtx, err)
		return nil, err
	}

	return makeBlockResponse(height, *hash, *data, txs), nil
}

func (txn *ViewTxn) ReadBlockTxCount(chainId uint64, key dbtypes.BlockKey) (dbp.HexUint, error) {
	it := txn.txn.NewIterator(badger.IteratorOptions{
		Reverse: true,
		Prefix:  dbkey.TxsForBlock.Get(chainId, key.Height),
	})
	defer it.Close()
	it.Seek(dbkey.Tx.Get(chainId, key.Height, dbkey.MaxTxIndex))
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

func (txn *ViewTxn) ReadBlockHashes(
	ctx context.Context,
	chainId uint64,
	from *dbtypes.BlockKey,
	to *dbtypes.BlockKey,
	limit int,
) ([]dbprimitives.Data32, error) {

	if limit <= 0 {
		limit = 100_000
	}

	if from.CompareTo(to) > 0 {
		return nil, fmt.Errorf("from > to")
	}

	fromKey := dbkey.BlockHash.Get(chainId, from.Height)
	toKey := dbkey.BlockHash.Get(chainId, to.Height)
	it := txn.txn.NewIterator(badger.IteratorOptions{
		PrefetchValues: true,
		PrefetchSize:   1000,
		Prefix:         getCommonPrefix(fromKey, toKey),
	})
	defer it.Close()

	response := []dbprimitives.Data32{}
	for it.Seek(fromKey); it.Valid(); it.Next() {
		select {
		case <-ctx.Done():
			return response, ctx.Err()
		default:
		}
		if bytes.Compare(it.Item().Key(), toKey) > 0 {
			break
		}

		if !dbkey.BlockHash.Matches(it.Item().Key()) {
			txn.db.logger.Errorf("DB: detected corrupted BlockHash key, will ignore: %v", it.Item().Key())
			continue
		}

		blockHash, err := readItem[dbp.Data32](txn.db, it.Item())
		if err != nil || blockHash == nil {
			txn.db.logger.Errorf("DB: can't read BlockHash, will ignore: %v", err)
			continue
		}

		if len(response) == limit {
			return response, ErrLimited
		}
		response = append(response, *blockHash)
	}
	return response, nil
}
