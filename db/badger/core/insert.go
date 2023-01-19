package core

import (
	"aurora-relayer-go-common/db/badger/core/dbkey"
	"aurora-relayer-go-common/db/badger/core/logscan"
	"aurora-relayer-go-common/tinypack"
	dbt "aurora-relayer-go-common/types/db"
	"aurora-relayer-go-common/types/primitives"

	"github.com/dgraph-io/badger/v3"
)

// Wrapper that hides internals from outer packages
type Writer struct {
	db     *DB
	writer *badger.WriteBatch
}

func (w *Writer) Flush() error {
	return w.writer.Flush()
}

func (w *Writer) Cancel() {
	w.writer.Cancel()
}

func insert[T any](w *Writer, key []byte, value *T) error {
	b, err := w.db.codec.Marshal(value)
	if err != nil {
		w.db.logger.Errorf("DB: Can't marshal value of type %T: %v", value, err)
		return err
	}
	if err := w.writer.Set(key, b); err != nil {
		w.db.logger.Errorf("DB: Can't write value: %v", err)
		return err
	}
	return nil
}

func insertInstantly[T any](db *DB, key []byte, value *T) error {
	b, err := db.codec.Marshal(value)
	if err != nil {
		db.logger.Errorf("DB: Can't marshal value of type %T: %v", value, err)
		return err
	}
	err = db.core.Update(func(txn *badger.Txn) error {
		return txn.Set(key, b)
	})
	if err != nil {
		db.logger.Errorf("DB: Can't write value: %v", err)
		return err
	}
	return nil
}

func (w *Writer) InsertBlock(chainId, height uint64, hash primitives.Data32, data *dbt.Block) error {
	if err := insert(w, dbkey.BlockHash.Get(chainId, height), &hash); err != nil {
		w.db.logger.Errorf("DB: Can't insert block hash: %v", err)
		return err
	}
	if err := insert(w, dbkey.BlockData.Get(chainId, height), data); err != nil {
		w.db.logger.Errorf("DB: Can't insert block data: %v", err)
		return err
	}
	blockKey := &dbt.BlockKey{Height: height}
	if err := insert(w, dbkey.BlockKeyByHash.Get(chainId, hash.Bytes()), blockKey); err != nil {
		w.db.logger.Errorf("DB: Can't insert block key: %v", err)
		return err
	}
	return nil
}

func (w *Writer) InsertTransaction(chainId, height, index uint64, hash primitives.Data32, data *dbt.Transaction) error {
	if err := insert(w, dbkey.TxHash.Get(chainId, height, index), &hash); err != nil {
		w.db.logger.Errorf("DB: Can't insert transaction hash: %v", err)
		return err
	}
	if err := insert(w, dbkey.TxData.Get(chainId, height, index), data); err != nil {
		w.db.logger.Errorf("DB: Can't insert transaction data: %v", err)
		return err
	}
	txKey := &dbt.TransactionKey{BlockHeight: height, TransactionIndex: index}
	if err := insert(w, dbkey.TxKeyByHash.Get(chainId, hash.Bytes()), txKey); err != nil {
		w.db.logger.Errorf("DB: Can't insert transaction key: %v", err)
		return err
	}
	return nil
}

func (w *Writer) InsertLog(chainId, height, txIndex, logIndex uint64, data *dbt.Log) error {
	if err := insert(w, dbkey.Log.Get(chainId, height, txIndex, logIndex), data); err != nil {
		w.db.logger.Errorf("DB: Can't insert log: %v", err)
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
		err := w.writer.Set(dbkey.LogScanEntry.Get(chainId, uint64(scanBitmask), hash, height, txIndex, logIndex), nil)
		if err != nil {
			w.db.logger.Errorf("DB: Can't insert LogScanEntry: %v", err)
			return err
		}
	}
	return nil
}

func (db *DB) InsertIndexerState(chainId uint64, data []byte) error {
	d := tinypack.CreateVarData(data...)
	if err := insertInstantly(db, dbkey.IndexerState.Get(chainId), &d); err != nil {
		return err
	}
	return nil
}

func (db *DB) InsertBlockFilter(chainId uint64, filterId primitives.Data32, filter *dbt.BlockFilter) error {
	if err := insertInstantly(db, dbkey.BlockFilter.Get(chainId, filterId.Bytes()), filter); err != nil {
		db.logger.Errorf("DB: Can't insert BlockFilter: %v", err)
		return err
	}
	return nil
}

func (db *DB) InsertTransactionFilter(chainId uint64, filterId primitives.Data32, filter *dbt.TransactionFilter) error {
	if err := insertInstantly(db, dbkey.TxFilter.Get(chainId, filterId.Bytes()), filter); err != nil {
		db.logger.Errorf("DB: Can't insert TransactionFilter: %v", err)
		return err
	}
	return nil
}

func (db *DB) InsertLogFilter(chainId uint64, filterId primitives.Data32, filter *dbt.LogFilter) error {
	if err := insertInstantly(db, dbkey.LogFilter.Get(chainId, filterId.Bytes()), filter); err != nil {
		db.logger.Errorf("DB: Can't insert LogFilter: %v", err)
		return err
	}
	return nil
}
