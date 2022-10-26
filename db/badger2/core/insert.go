package db

import (
	"aurora-relayer-go-common/db/badger2/core/dbkey"
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
	"aurora-relayer-go-common/db/badger2/core/dbtypes"
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
	// TODO: populate scan-index
	return nil
}
