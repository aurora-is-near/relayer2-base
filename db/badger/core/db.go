package core

import (
	"aurora-relayer-go-common/db/codec"
	"github.com/dgraph-io/badger/v3"
)

type DB struct {
	codec                 codec.Codec
	maxLogScanIterators   uint // Should be somewhere between 1k and 100k
	logScanRangeThreshold uint // Minimum block range size for using index instead of simple iteration
	filterTtlMinutes      int
	logger                badger.Logger
	core                  *badger.DB
}

func NewDB(config Config, codec codec.Codec) (*DB, error) {
	core, err := Open(config.BadgerConfig, config.GcIntervalSeconds)
	if err != nil {
		return nil, err
	}
	db := &DB{
		codec:                 codec,
		maxLogScanIterators:   config.MaxScanIterators,
		logScanRangeThreshold: config.ScanRangeThreshold,
		filterTtlMinutes:      config.FilterTtlMinutes,
		logger:                config.BadgerConfig.Logger,
		core:                  core,
	}
	return db, nil
}

func (db *DB) NewWriter() *Writer {
	return &Writer{
		db:     db,
		writer: db.core.NewWriteBatch(),
	}
}

func (db *DB) View(fn func(txn *ViewTxn) error) error {
	return db.core.View(func(txn *badger.Txn) error {
		return fn(&ViewTxn{
			db:  db,
			txn: txn,
		})
	})
}

func (db *DB) Update(fn func(txn *ViewTxn) error) error {
	return db.core.Update(func(txn *badger.Txn) error {
		return fn(&ViewTxn{
			db:  db,
			txn: txn,
		})
	})
}

func (db *DB) Close() error {
	return Close()
}

// For debugging/testing/etc purposes
// For production purposes use db.View(...) and db.InsertSomething functions
func (db *DB) BadgerDB() *badger.DB {
	return db.core
}
