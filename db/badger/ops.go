package badger

import (
	"bytes"
	"github.com/dgraph-io/badger/v3"
	"github.com/fxamacker/cbor/v2"
)

func fetchFromDB[T any](db *badger.DB, key []byte) (*T, error) {
	res := new(T)
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		valCopy, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return cbor.Unmarshal(valCopy, res)
	})
	return res, err
}

func fetchPrefixedFromDB[T any](db *badger.DB, prefix []byte) ([]T, error) {
	results := make([]T, 0)
	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			valCopy, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			res := new(T)
			if err := cbor.Unmarshal(valCopy, res); err != nil {
				return err
			}
			results = append(results, *res)
		}
		return nil
	})
	return results, err
}

func insertToDB[T any](db *badger.DB, key []byte, val T) error {
	b, err := cbor.Marshal(val)
	if err != nil {
		return err
	}
	return db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry(key, b)
		return txn.SetEntry(e)
	})
}

func concatBytes(pieces ...[]byte) []byte {
	var buf bytes.Buffer
	for _, piece := range pieces {
		buf.Write(piece)
	}
	return buf.Bytes()
}
