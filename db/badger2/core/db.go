package db

import (
	"aurora-relayer-go-common/db/badger/core"
	"aurora-relayer-go-common/db/badger2/core/dbcore"
	"aurora-relayer-go-common/tinypack"
	badger "github.com/dgraph-io/badger/v3"
)

type DB struct {
	CoreOpts    *dbcore.DBCoreOpts
	Encoder     *tinypack.Encoder
	Decoder     *tinypack.Decoder
	LogScanOpts *core.ScanOpts

	logger badger.Logger
	core   *dbcore.DBCore
	writer *badger.WriteBatch
}

func (db *DB) Open(logger badger.Logger) error {
	db.logger = logger
	db.logger.Infof("DB: Opening")
	var err error
	db.core, err = dbcore.Open(db.CoreOpts, logger)
	if err != nil {
		db.logger.Errorf("DB: Unable to open database: %v", err)
		return err
	}
	db.writer = db.core.BadgerDB().NewWriteBatch()
	return nil
}

func (db *DB) View(fn func(txn *ViewTxn) error) error {
	return db.BadgerDB().View(func(txn *badger.Txn) error {
		return fn(&ViewTxn{
			db:  db,
			txn: txn,
		})
	})
}

func (db *DB) OpenWriter() {
	db.writer = db.core.BadgerDB().NewWriteBatch()
}

func (db *DB) CloseWriter() {
	db.writer.Cancel()
}

func (db *DB) FlushWriter() error {
	if err := db.writer.Flush(); err != nil {
		db.logger.Errorf("DB: unable to flush writer: %v", err)
		return err
	}
	return nil
}

func (db *DB) Close() error {
	db.logger.Infof("DB: Closing")
	_ = db.FlushWriter()
	if err := db.core.Close(); err != nil {
		db.logger.Errorf("DB: unable to close database: %v", err)
		return err
	}
	return nil
}

// For debugging/testing/etc purposes
// For production purposes use db.View(...) and db.InsertSomething functions
func (db *DB) BadgerDB() *badger.DB {
	return db.core.BadgerDB()
}
