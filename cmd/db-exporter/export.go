package main

import (
	"github.com/aurora-is-near/relayer2-base/db/badger/core/dbkey"
	"github.com/aurora-is-near/relayer2-base/db/codec"
	dbt "github.com/aurora-is-near/relayer2-base/types/db"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
)

type Archiver interface {
	WriteBlock(*dbt.Block) error
	WriteBlockHash(*primitives.Data32) error
	WriteBlockHeight(uint64) error

	WriteTx(*dbt.Transaction) error
	WriteTxHash(*primitives.Data32) error
	WriteTxIndex(uint64) error
	WriteTxHeight(uint64) error

	WriteLog(*dbt.Log) error
	WriteLogIndex(uint64) error
	WriteLogTxIndex(uint64) error
	WriteLogHeight(uint64) error

	Close() error
}

type Exporter struct {
	DB       *badger.DB
	Archiver Archiver
	ChainID  uint64
	Height   uint64
	Decoder  codec.Decoder
}

func (e *Exporter) Export() error {
	err := e.DB.View(func(txn *badger.Txn) error {
		if err := e.exportBlocks(txn); err != nil {
			return err
		}
		if err := e.exportTxs(txn); err != nil {
			return err
		}
		if err := e.exportLogs(txn); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return e.Archiver.Close()
}

func (e *Exporter) exportBlocks(txn *badger.Txn) error {
	if err := exportHelper(txn, e.Decoder, e.Archiver.WriteBlock, dbkey.BlocksData.Get(e.ChainID), dbkey.BlockData.Get(e.ChainID, e.Height)); err != nil {
		return err
	}

	it := txn.NewIterator(badger.IteratorOptions{
		PrefetchValues: true,
		PrefetchSize:   1000,
		Prefix:         dbkey.BlockHashes.Get(e.ChainID),
	})
	defer it.Close()

	hash := new(primitives.Data32)
	for it.Seek(dbkey.BlockHash.Get(e.ChainID, e.Height)); it.Valid(); it.Next() {
		blockHeight := dbkey.BlockHash.ReadUintVar(it.Item().Key(), 1)
		if err := e.Archiver.WriteBlockHeight(blockHeight); err != nil {
			return err
		}
		if err := readItem(e.Decoder, it.Item(), hash); err != nil {
			return err
		}
		if err := e.Archiver.WriteBlockHash(hash); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) exportTxs(txn *badger.Txn) error {
	if err := exportHelper(txn, e.Decoder, e.Archiver.WriteTx, dbkey.TxsData.Get(e.ChainID), dbkey.TxsDataForBlock.Get(e.ChainID, e.Height)); err != nil {
		return err
	}

	it := txn.NewIterator(badger.IteratorOptions{
		PrefetchValues: true,
		PrefetchSize:   1000,
		Prefix:         dbkey.TxHashes.Get(e.ChainID),
	})
	defer it.Close()

	hash := new(primitives.Data32)
	for it.Seek(dbkey.TxHashesForBlock.Get(e.ChainID, e.Height)); it.Valid(); it.Next() {
		key := it.Item().Key()
		blockHeight := dbkey.TxHash.ReadUintVar(key, 1)
		if err := e.Archiver.WriteTxHeight(blockHeight); err != nil {
			return err
		}

		txIndex := dbkey.TxHash.ReadUintVar(key, 2)
		if err := e.Archiver.WriteTxIndex(txIndex); err != nil {
			return err
		}

		if err := readItem(e.Decoder, it.Item(), hash); err != nil {
			return err
		}
		if err := e.Archiver.WriteTxHash(hash); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) exportLogs(txn *badger.Txn) error {
	it := txn.NewIterator(badger.IteratorOptions{
		PrefetchValues: true,
		PrefetchSize:   1000,
		Prefix:         dbkey.Logs.Get(e.ChainID),
	})
	defer it.Close()

	log := new(dbt.Log)
	for it.Seek(dbkey.LogsForBlock.Get(e.ChainID, e.Height)); it.Valid(); it.Next() {
		key := it.Item().Key()
		blockHeight := dbkey.Log.ReadUintVar(key, 1)
		if err := e.Archiver.WriteLogHeight(blockHeight); err != nil {
			return err
		}

		txIndex := dbkey.Log.ReadUintVar(key, 2)
		if err := e.Archiver.WriteLogTxIndex(txIndex); err != nil {
			return err
		}

		logIndex := dbkey.Log.ReadUintVar(key, 3)
		if err := e.Archiver.WriteLogIndex(logIndex); err != nil {
			return err
		}

		if err := readItem(e.Decoder, it.Item(), log); err != nil {
			return err
		}
		if err := e.Archiver.WriteLog(log); err != nil {
			return err
		}
	}
	return nil
}

func exportHelper[T any](
	txn *badger.Txn,
	dec codec.Decoder,
	write func(*T) error,
	prefix, start []byte,
) error {
	it := txn.NewIterator(badger.IteratorOptions{
		PrefetchValues: true,
		PrefetchSize:   1000,
		Prefix:         prefix,
	})
	defer it.Close()

	data := new(T)
	for it.Seek(start); it.Valid(); it.Next() {
		if err := readItem(dec, it.Item(), data); err != nil {
			return err
		}
		if err := write(data); err != nil {
			return err
		}
	}
	return nil
}

func readItem[T any](dec codec.Decoder, item *badger.Item, res *T) error {
	return item.Value(func(val []byte) error {
		if err := dec.Unmarshal(val, res); err != nil {
			return errors.Wrapf(err, "can't unmarshal value of type %T", res)
		}
		return nil
	})
}
