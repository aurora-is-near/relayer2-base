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

func (txn *ViewTxn) ReadBlockKey(chainId uint64, hash primitives.Data32) (*dbt.BlockKey, error) {
	return read[dbt.BlockKey](txn, dbkey.BlockKeyByHash.Get(chainId, hash.Bytes()))
}

func (txn *ViewTxn) ReadEarliestBlockKey(chainId uint64) (*dbt.BlockKey, error) {
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

func (txn *ViewTxn) ReadLatestBlockKey(chainId uint64) (*dbt.BlockKey, error) {
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

func readBlockKeyFromBlockHashIterator(it *badger.Iterator) (*dbt.BlockKey, error) {
	if !it.Valid() {
		return nil, nil
	}
	if !dbkey.BlockHash.Matches(it.Item().Key()) {
		return nil, fmt.Errorf("found unexpected key format (expected to match dbkey.BlockHash)")
	}
	return &dbt.BlockKey{Height: dbkey.BlockHash.ReadUintVar(it.Item().Key(), 1)}, nil
}

func (txn *ViewTxn) ReadBlock(chainId uint64, key dbt.BlockKey, fullTransactions bool) (*response.Block, error) {
	height := key.Height
	hash, err := readCached[primitives.Data32](txn, dbkey.BlockHash.Get(chainId, height))
	if err != nil || hash == nil {
		return nil, err
	}
	data, err := read[dbt.Block](txn, dbkey.BlockData.Get(chainId, height))
	if err != nil || data == nil {
		return nil, err
	}

	txs, _, err := txn.ReadTransactions(
		context.Background(),
		chainId,
		&dbt.TransactionKey{
			BlockHeight:      key.Height,
			TransactionIndex: 0,
		},
		&dbt.TransactionKey{
			BlockHeight:      key.Height,
			TransactionIndex: dbkey.MaxTxIndex,
		},
		fullTransactions,
		int(dbkey.MaxTxIndex)+1,
	)
	if err != nil {
		errCtx := fmt.Sprintf("chainId=%v, block=%v", chainId, key.Height)
		txn.db.logger.Errorf("DB: errors reading txs for %v: %v", errCtx, err)
		return nil, err
	}

	return makeBlockResponse(height, *hash, *data, txs), nil
}

func (txn *ViewTxn) ReadBlockTxCount(chainId uint64, key dbt.BlockKey) (primitives.HexUint, error) {
	it := txn.txn.NewIterator(badger.IteratorOptions{
		Reverse: true,
		Prefix:  dbkey.TxHashesForBlock.Get(chainId, key.Height),
	})
	defer it.Close()
	it.Seek(dbkey.TxHash.Get(chainId, key.Height, dbkey.MaxTxIndex))
	if !it.Valid() {
		return 0, nil
	}
	if !dbkey.TxHash.Matches(it.Item().Key()) {
		err := fmt.Errorf("found unexpected key format (expected to match dbkey.Tx)")
		txn.db.logger.Errorf("DB: can't read block tx count: %v", err)
		return 0, err
	}
	return primitives.HexUint(dbkey.TxHash.ReadUintVar(it.Item().Key(), 2)) + 1, nil
}

func (txn *ViewTxn) ReadBlockHashes(
	ctx context.Context,
	chainId uint64,
	from *dbt.BlockKey,
	to *dbt.BlockKey,
	limit int,
) ([]primitives.Data32, *dbt.BlockKey, error) {

	if limit <= 0 {
		limit = 100_000
	}

	if from.CompareTo(to) > 0 {
		return nil, nil, fmt.Errorf("from > to")
	}

	fromKey := dbkey.BlockHash.Get(chainId, from.Height)
	toKey := dbkey.BlockHash.Get(chainId, to.Height)
	it := txn.txn.NewIterator(badger.IteratorOptions{
		PrefetchValues: true,
		PrefetchSize:   1000,
		Prefix:         getCommonPrefix(fromKey, toKey),
	})
	defer it.Close()

	response := []primitives.Data32{}
	lastKey := from.Prev()
	for it.Seek(fromKey); it.Valid(); it.Next() {
		select {
		case <-ctx.Done():
			return response, lastKey, ctx.Err()
		default:
		}
		if bytes.Compare(it.Item().Key(), toKey) > 0 {
			break
		}

		if !dbkey.BlockHash.Matches(it.Item().Key()) {
			txn.db.logger.Errorf("DB: detected corrupted BlockHash key, will ignore: %v", it.Item().Key())
			continue
		}

		blockHash, err := readItem[primitives.Data32](txn.db, it.Item())
		if err != nil || blockHash == nil {
			txn.db.logger.Errorf("DB: can't read BlockHash, will ignore: %v", err)
			continue
		}

		if len(response) == limit {
			return response, lastKey, ErrLimited
		}
		response = append(response, *blockHash)
		lastKey = &dbt.BlockKey{Height: dbkey.BlockHash.ReadUintVar(it.Item().Key(), 1)}
	}
	return response, to, nil
}
