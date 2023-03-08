package main

import (
	"io"
	L "log"

	"github.com/aurora-is-near/relayer2-base/db/badger/core"
	dbt "github.com/aurora-is-near/relayer2-base/types/db"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
)

type Unarchiver interface {
	ReadBlock() (*dbt.Block, error)
	ReadBlockHash() (*primitives.Data32, error)
	ReadBlockHeight() (uint64, error)

	ReadTx() (*dbt.Transaction, error)
	ReadTxHash() (*primitives.Data32, error)
	ReadTxIndex() (uint64, error)
	ReadTxHeight() (uint64, error)

	ReadLog() (*dbt.Log, error)
	ReadLogIndex() (uint64, error)
	ReadLogTxIndex() (uint64, error)
	ReadLogHeight() (uint64, error)

	Close() error
}

type Writer interface {
	InsertBlock(chainId, height uint64, hash primitives.Data32, data *dbt.Block) error
	InsertTransaction(
		chainId, height, index uint64,
		hash primitives.Data32,
		data *dbt.Transaction,
	) error
	InsertLog(chainId, height, txIndex, logIndex uint64, data *dbt.Log) error
}

type Importer struct {
	DB           *core.DB
	Unarchiver   Unarchiver
	ChainID      uint64
	PendingLimit uint64
}

func (i *Importer) Import() error {
	err := i.importBlocks()
	if err != nil {
		return err
	}

	err = i.importTxs()
	if err != nil {
		return err
	}

	err = i.importLogs()
	if err != nil {
		return err
	}
	return i.Unarchiver.Close()
}

func (i *Importer) importBlocks() error {
	writer := i.DB.NewWriter()
	defer writer.Cancel()
	pendingCount := uint64(0)

	var (
		block  *dbt.Block
		hash   *primitives.Data32
		height uint64
		err    error
	)

	for {
		block, err = i.Unarchiver.ReadBlock()
		if err != nil {
			break
		}

		hash, err = i.Unarchiver.ReadBlockHash()
		if err != nil {
			break
		}

		height, err = i.Unarchiver.ReadBlockHeight()
		if err != nil {
			break
		}

		err = writer.InsertBlock(i.ChainID, height, *hash, block)
		if err != nil {
			break
		}

		pendingCount++
		if pendingCount >= i.PendingLimit {
			err = writer.Flush()
			if err != nil {
				break
			}
			L.Println("committed", pendingCount, "blocks to the DB")
			writer = i.DB.NewWriter()
			pendingCount = 0
		}
	}
	if err != nil && err != io.EOF {
		return err
	}
	return writer.Flush()
}

func (i *Importer) importTxs() error {
	writer := i.DB.NewWriter()
	defer writer.Cancel()
	pendingCount := uint64(0)

	var (
		tx     *dbt.Transaction
		hash   *primitives.Data32
		index  uint64
		height uint64
		err    error
	)

	for {
		tx, err = i.Unarchiver.ReadTx()
		if err != nil {
			break
		}

		hash, err = i.Unarchiver.ReadTxHash()
		if err != nil {
			break
		}

		index, err = i.Unarchiver.ReadTxIndex()
		if err != nil {
			break
		}

		height, err = i.Unarchiver.ReadTxHeight()
		if err != nil {
			break
		}

		err = writer.InsertTransaction(i.ChainID, height, index, *hash, tx)
		if err != nil {
			break
		}

		pendingCount++
		if pendingCount >= i.PendingLimit {
			err = writer.Flush()
			if err != nil {
				break
			}
			L.Println("committed", pendingCount, "transactions to the DB")
			writer = i.DB.NewWriter()
			pendingCount = 0
		}
	}
	if err != nil && err != io.EOF {
		return err
	}
	return writer.Flush()
}

func (i *Importer) importLogs() error {
	writer := i.DB.NewWriter()
	defer writer.Cancel()
	pendingCount := uint64(0)

	var (
		log     *dbt.Log
		index   uint64
		txIndex uint64
		height  uint64
		err     error
	)

	for {
		log, err = i.Unarchiver.ReadLog()
		if err != nil {
			break
		}

		index, err = i.Unarchiver.ReadLogIndex()
		if err != nil {
			break
		}

		txIndex, err = i.Unarchiver.ReadLogTxIndex()
		if err != nil {
			break
		}

		height, err = i.Unarchiver.ReadLogHeight()
		if err != nil {
			break
		}

		err = writer.InsertLog(i.ChainID, height, txIndex, index, log)
		if err != nil {
			break
		}

		pendingCount++
		if pendingCount >= i.PendingLimit {
			err = writer.Flush()
			if err != nil {
				break
			}
			L.Println("committed", pendingCount, "logs to the DB")
			writer = i.DB.NewWriter()
			pendingCount = 0
		}
	}
	if err != nil && err != io.EOF {
		return err
	}
	return writer.Flush()
}
