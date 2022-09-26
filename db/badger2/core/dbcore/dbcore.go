package dbcore

import (
	"os"
	"path"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type DBCoreOpts struct {
	Dir               string
	GCIntervalSeconds uint
	InMemory          bool
}

type DBCore struct {
	opts     *DBCoreOpts
	logger   badger.Logger
	badgerDB *badger.DB
	gcStop   chan bool
}

func (opts DBCoreOpts) FillMissingFields() *DBCoreOpts {
	if opts.GCIntervalSeconds == 0 {
		opts.GCIntervalSeconds = 10
	}
	return &opts
}

func Open(opts *DBCoreOpts, logger badger.Logger) (*DBCore, error) {
	opts = opts.FillMissingFields()
	dbc := &DBCore{
		opts:   opts,
		logger: logger,
		gcStop: make(chan bool),
	}

	logger.Infof("DBCore: opening badger database")
	if err := dbc.tryOpen(); err != nil {
		logger.Warningf("DBCore: will save old snapshot and create new db: %v", err)
		oldBaseName := path.Base(opts.Dir) + "_" + time.Now().Format("2006-01-02T15-04-05.000000000")
		oldPath := path.Join(path.Dir(opts.Dir), oldBaseName)
		if err := os.Rename(opts.Dir, oldPath); err != nil {
			return nil, err
		}
		logger.Infof("DBCore: creating new badger database")
		if err := dbc.tryOpen(); err != nil {
			return nil, err
		}
	}

	logger.Infof("DBCore: starting garbage collector")
	go dbc.runGC()

	return dbc, nil
}

func (dbc *DBCore) BadgerDB() *badger.DB {
	return dbc.badgerDB
}

func (dbc *DBCore) Close() error {
	dbc.logger.Infof("DBCore: stopping garbage collector")
	dbc.gcStop <- true

	dbc.logger.Infof("DBCore: closing badger database")
	if err := dbc.BadgerDB().Close(); err != nil {
		dbc.logger.Errorf("DBCore: unable to close badger database: %v", err)
		return err
	}

	return nil
}

func (dbc *DBCore) tryOpen() error {
	opts := badger.DefaultOptions(dbc.opts.Dir)
	opts = opts.WithInMemory(dbc.opts.InMemory)

	var err error
	if dbc.badgerDB, err = badger.Open(opts); err != nil {
		dbc.logger.Errorf("DBCore: unable to open badger database: %v", err)
		return err
	}

	return nil
}

func (dbc *DBCore) runGC() {
	ticker := time.NewTicker(time.Duration(dbc.opts.GCIntervalSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-dbc.gcStop:
			return
		case <-ticker.C:
			for {
				select {
				case <-dbc.gcStop:
					return
				default:
				}
				if dbc.badgerDB.RunValueLogGC(0.5) != nil {
					break
				}
			}
		}
	}
}
