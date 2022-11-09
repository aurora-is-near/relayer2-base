package db

import (
	"aurora-relayer-go-common/db/badger2/core/dbkey"
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
	"aurora-relayer-go-common/db/badger2/core/dbtypes"
	"aurora-relayer-go-common/db/badger2/core/logscan"

	badger "github.com/dgraph-io/badger/v3"
)

func insert[T any](db *DB, key []byte, value *T) error {
	b, err := db.Encoder.Marshal(value)
	if err != nil {
		db.logger.Errorf("DB: Can't marshal value of type %T: %v", value, err)
		return err
	}
	if err := db.writer.Set(key, b); err != nil {
		db.logger.Errorf("DB: Can't write value: %v", err)
		return err
	}
	return nil
}

func insertInstantly[T any](db *DB, key []byte, value *T) error {
	b, err := db.Encoder.Marshal(value)
	if err != nil {
		db.logger.Errorf("DB: Can't marshal value of type %T: %v", value, err)
		return err
	}
	err = db.BadgerDB().Update(func(txn *badger.Txn) error {
		return txn.Set(key, b)
	})
	if err != nil {
		db.logger.Errorf("DB: Can't write value: %v", err)
		return err
	}
	return nil
}

func (db *DB) InsertBlock(chainId, height uint64, hash dbp.Data32, data *dbtypes.Block) error {
	if err := insert(db, dbkey.BlockHash.Get(chainId, height), &hash); err != nil {
		db.logger.Errorf("DB: Can't insert block hash: %v", err)
		return err
	}
	if err := insert(db, dbkey.BlockData.Get(chainId, height), data); err != nil {
		db.logger.Errorf("DB: Can't insert block data: %v", err)
		return err
	}
	blockKey := &dbtypes.BlockKey{Height: height}
	if err := insert(db, dbkey.BlockKeyByHash.Get(chainId, hash.Bytes()), blockKey); err != nil {
		db.logger.Errorf("DB: Can't insert block key: %v", err)
		return err
	}
	return nil
}

func (db *DB) InsertTransaction(chainId, height, index uint64, hash dbp.Data32, data *dbtypes.Transaction) error {
	if err := insert(db, dbkey.TxHash.Get(chainId, height, index), &hash); err != nil {
		db.logger.Errorf("DB: Can't insert transaction hash: %v", err)
		return err
	}
	if err := insert(db, dbkey.TxData.Get(chainId, height, index), data); err != nil {
		db.logger.Errorf("DB: Can't insert transaction data: %v", err)
		return err
	}
	txKey := &dbtypes.TransactionKey{BlockHeight: height, TransactionIndex: index}
	if err := insert(db, dbkey.TxKeyByHash.Get(chainId, hash.Bytes()), txKey); err != nil {
		db.logger.Errorf("DB: Can't insert transaction key: %v", err)
		return err
	}
	return nil
}

func (db *DB) InsertLog(chainId, height, txIndex, logIndex uint64, data *dbtypes.Log) error {
	if err := insert(db, dbkey.Log.Get(chainId, height, txIndex, logIndex), data); err != nil {
		db.logger.Errorf("DB: Can't insert log: %v", err)
		return err
	}

	scanFeatures := make([][]byte, len(data.Topics.Content)+1)
	scanFeatures[0] = data.Address.Bytes()
	for i, t := range data.Topics.Content {
		scanFeatures[i+1] = t.Bytes()
	}

	maxScanBitmask := 1<<len(scanFeatures) - 1
	for scanBitmask := 1; scanBitmask <= maxScanBitmask; scanBitmask++ {
		hash := logscan.CalcHash(scanFeatures, scanBitmask)
		err := db.writer.Set(dbkey.LogScanEntry.Get(chainId, uint64(scanBitmask), hash, height, txIndex, logIndex), nil)
		if err != nil {
			db.logger.Errorf("DB: Can't insert LogScanEntry: %v", err)
			return err
		}
	}
	return nil
}

func (db *DB) InsertBlockFilter(chainId uint64, filterId dbp.Data32, filter *dbtypes.BlockFilter) error {
	if err := insertInstantly(db, dbkey.BlockFilter.Get(chainId, filterId.Bytes()), filter); err != nil {
		db.logger.Errorf("DB: Can't insert BlockFilter: %v", err)
		return err
	}
	return nil
}

func (db *DB) InsertTransactionFilter(chainId uint64, filterId dbp.Data32, filter *dbtypes.TransactionFilter) error {
	if err := insertInstantly(db, dbkey.TxFilter.Get(chainId, filterId.Bytes()), filter); err != nil {
		db.logger.Errorf("DB: Can't insert TransactionFilter: %v", err)
		return err
	}
	return nil
}

func (db *DB) InsertLogFilter(chainId uint64, filterId dbp.Data32, filter *dbtypes.LogFilter) error {
	if err := insertInstantly(db, dbkey.LogFilter.Get(chainId, filterId.Bytes()), filter); err != nil {
		db.logger.Errorf("DB: Can't insert LogFilter: %v", err)
		return err
	}
	return nil
}
