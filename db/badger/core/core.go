package core

import (
	"aurora-relayer-go-common/log"
	"github.com/dgraph-io/badger/v3"
	"os"
	"path"
	"sync"
	"time"
)

var lock = &sync.Mutex{}
var bdb *badger.DB
var gcStop chan bool

func Open(options badger.Options, gcIntervalSeconds int) (*badger.DB, error) {
	var err error
	if bdb == nil {
		lock.Lock()
		defer lock.Unlock()
		if bdb == nil {
			bdb, err = tryOpen(options, gcIntervalSeconds)
		}
	}
	return bdb, err
}

func Close() error {
	if bdb != nil {
		gcStop <- true
		log.Log().Info().Msg("closing database")
		err := bdb.Close()
		if err != nil {
			return err
		}
		bdb = nil
	}
	return nil
}

func tryOpen(options badger.Options, gcIntervalSeconds int) (*badger.DB, error) {

	var err error
	logger := log.Log()

	if !options.InMemory {
		logger.Info().Msgf("opening database with path [%s]", options.Dir)
	} else {
		options.Dir = ""
		options.ValueDir = ""
		logger.Info().Msg("opening database as in-memory")
	}
	bdb, err = badger.Open(options)
	if err != nil {
		logger.Error().Err(err).Msg("failed to tryOpen database")
		snapshotBaseName := path.Base(options.Dir) + "_" + time.Now().Format("2006-01-02T15-04-05.000000000")
		snapshotPath := path.Join(path.Dir(options.Dir), snapshotBaseName)
		logger.Warn().Err(err).Msgf("saving old database snapshot at [%s]", snapshotPath)
		if err := os.Rename(options.Dir, snapshotPath); err != nil {
			logger.Error().Err(err).Msg("failed to save old snapshot")
			return nil, err
		}
		logger.Info().Err(err).Msg("creating new database")
		bdb, err = badger.Open(options)
	}
	if bdb != nil {
		go runGC(gcIntervalSeconds)
	}

	return bdb, err
}

func runGC(gcIntervalSeconds int) {
	gcStop = make(chan bool)
	ticker := time.NewTicker(time.Duration(gcIntervalSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-gcStop:
			return
		case <-ticker.C:
			for {
				select {
				case <-gcStop:
					return
				default:
				}
				if bdb.RunValueLogGC(0.5) != nil {
					break
				}
			}
		}
	}
}
