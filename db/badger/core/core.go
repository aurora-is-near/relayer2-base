package core

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aurora-is-near/relayer2-base/log"
	"github.com/aurora-is-near/relayer2-base/syncutils"

	"github.com/dgraph-io/badger/v3"
)

var bdbPtr syncutils.LockablePtr[badger.DB]
var gcStop chan bool

func Open(options badger.Options, gcIntervalSeconds int, recreateOnCorruption bool) (*badger.DB, error) {
	var err error
	bdb, unlock := bdbPtr.LockIfNil()
	if unlock != nil {
		bdb, err = tryOpen(options, gcIntervalSeconds, recreateOnCorruption)
		unlock(bdb)
	}
	return bdb, err
}

func Close() error {
	bdb, unlock := bdbPtr.LockIfNotNil()
	if unlock != nil {
		gcStop <- true
		log.Log().Info().Msg("closing database")
		if err := bdb.Close(); err != nil {
			unlock(bdb)
			return err
		}
		unlock(nil)
	}
	return nil
}

func tryOpen(options badger.Options, gcIntervalSeconds int, recreateOnCorruption bool) (*badger.DB, error) {
	var err error
	logger := log.Log()

	if !options.InMemory {
		logger.Info().Msgf("opening database with path [%s]", options.Dir)
	} else {
		options.Dir = ""
		options.ValueDir = ""
		logger.Info().Msg("opening database as in-memory")
	}
	bdb, err := badger.Open(options)
	if err != nil {
		logger.Error().Err(err).Msg("failed to tryOpen database")

		// No proper error class for that:
		// https://github.com/dgraph-io/badger/blob/c08da0b80769f86aa44366e2da77b5c662810eca/dir_unix.go#L67
		if strings.Contains(err.Error(), "Cannot acquire directory lock") {
			return nil, fmt.Errorf("db is busy: %w", err)
		}

		if !recreateOnCorruption {
			return nil, err
		}

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
		gcStop = make(chan bool)
		go runGC(bdb, gcIntervalSeconds)
	}

	return bdb, err
}

func runGC(bdb *badger.DB, gcIntervalSeconds int) {
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
