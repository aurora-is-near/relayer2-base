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
	logger := log.Log()
	if bdb == nil {
		lock.Lock()
		defer lock.Unlock()
		if bdb == nil {
			if !options.InMemory {
				logger.Info().Msgf("opening database with path [%s]", options.Dir)
			} else {
				options.Dir = ""
				options.ValueDir = ""
				logger.Info().Msg("opening database as in-memory")
			}
			bdb, err = open(options)
			if err != nil {
				snapshotBaseName := path.Base(options.Dir) + "_" + time.Now().Format("2006-01-02T15-04-05.000000000")
				snapshotPath := path.Join(path.Dir(options.Dir), snapshotBaseName)
				logger.Warn().Err(err).Msgf("saving old database snapshot at [%s]", snapshotPath)
				if err := os.Rename(options.Dir, snapshotPath); err != nil {
					logger.Error().Err(err).Msg("failed to save old snapshot")
					return nil, err
				}
				bdb, err = open(options)
			}
			if bdb != nil {
				go runGC(gcIntervalSeconds)
			}
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

func Fetch(key []byte) (*[]byte, error) {
	res := new([]byte)
	err := bdb.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		valueCopy, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		res = &valueCopy
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func FetchPrefixed(prefix []byte) ([][]byte, error) {
	res := make([][]byte, 0)
	err := bdb.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			valueCopy, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			res = append(res, valueCopy)
		}
		return nil
	})
	return res, err
}

func FetchPrefixedWithLimitAndTimeout(limit uint, timeout uint, prefix []byte) ([][]byte, error) {
	res := make([][]byte, 0)
	to := time.NewTimer(time.Duration(time.Second * time.Duration(timeout)))
	defer to.Stop()
	err := bdb.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			valueCopy, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			res = append(res, valueCopy)
			if uint(len(res)) >= limit {
				return nil
			}
			select {
			case <-to.C:
				return nil
			default:
			}
		}
		return nil
	})
	return res, err
}

func Insert(key []byte, value []byte) error {
	return bdb.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key, value)
		return txn.SetEntry(e)
	})
}

func Delete(key []byte) error {
	return bdb.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func open(options badger.Options) (*badger.DB, error) {
	bdb, err := badger.Open(options)
	if err != nil {
		return nil, err
	}
	return bdb, nil
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
