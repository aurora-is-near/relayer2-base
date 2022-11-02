package core

import (
	"errors"
	"fmt"
	"sync"

	"github.com/dgraph-io/badger/v3"
)

var ErrLimited = errors.New("limited")

// Wrapper that hides internals from outer packages
type ViewTxn struct {
	db    *DB
	txn   *badger.Txn
	cache sync.Map
}

type cachedRead[T any] struct {
	value *T
	err   error
	ready chan struct{}
}

func readItem[T any](db *DB, item *badger.Item) (*T, error) {
	res := new(T)
	err := item.Value(func(val []byte) error {
		if err := db.codec.Unmarshal(val, res); err != nil {
			return fmt.Errorf("can't unmarshal value of type %T: %v", res, err)
		}
		return nil
	})
	if err != nil {
		db.logger.Errorf("DB: can't read item: %v", err)
		return nil, err
	}
	return res, nil
}

func read[T any](txn *ViewTxn, key []byte) (*T, error) {
	item, err := txn.txn.Get(key)
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	if err != nil {
		txn.db.logger.Errorf("DB: Can't fetch item for key %v: %v", key, err)
	}
	return readItem[T](txn.db, item)
}

func readCached[T any](txn *ViewTxn, key []byte) (*T, error) {
	r, loaded := txn.cache.LoadOrStore(string(key), &cachedRead[T]{ready: make(chan struct{})})
	cr := r.(*cachedRead[T])
	if loaded {
		<-cr.ready
	} else {
		cr.value, cr.err = read[T](txn, key)
		close(cr.ready)
	}
	return cr.value, cr.err
}
