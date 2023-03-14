package main

import (
	"io"
	"log"

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
	var (
		writer       = i.DB.NewWriter()
		pendingCount uint64
		total        uint64
		block        *dbt.Block
		hash         *primitives.Data32
		height       uint64
		err          error
	)

	commit := func() error {
		err := writer.Flush()
		if err != nil {
			return err
		}
		log.Println("committed", pendingCount, "blocks to the DB, total", total)
		pendingCount = 0
		return nil
	}

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

		total++
		pendingCount++
		if pendingCount >= i.PendingLimit {
			err = commit()
			if err != nil {
				break
			}
			writer = i.DB.NewWriter()
		}
	}
	if err != nil && err != io.EOF {
		writer.Cancel()
		return err
	}
	return commit()
}

func (i *Importer) importTxs() error {
	var (
		writer       = i.DB.NewWriter()
		pendingCount uint64
		total        uint64
		tx           *dbt.Transaction
		hash         *primitives.Data32
		index        uint64
		height       uint64
		err          error
	)

	commit := func() error {
		err := writer.Flush()
		if err != nil {
			return err
		}
		log.Println("committed", pendingCount, "transactions to the DB, total", total)
		pendingCount = 0
		return nil
	}

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

		total++
		pendingCount++
		if pendingCount >= i.PendingLimit {
			err = commit()
			if err != nil {
				break
			}
			writer = i.DB.NewWriter()
		}
	}
	if err != nil && err != io.EOF {
		writer.Cancel()
		return err
	}
	return commit()
}

func (i *Importer) importLogs() error {
	var (
		writer       = i.DB.NewWriter()
		pendingCount uint64
		total        uint64
		lg           *dbt.Log
		index        uint64
		txIndex      uint64
		height       uint64
		err          error
	)

	commit := func() error {
		err := writer.Flush()
		if err != nil {
			return err
		}
		log.Println("committed", pendingCount, "logs to the DB, total", total)
		pendingCount = 0
		return nil
	}

	for {
		lg, err = i.Unarchiver.ReadLog()
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

		err = writer.InsertLog(i.ChainID, height, txIndex, index, lg)
		if err != nil {
			break
		}

		total++
		pendingCount++
		if pendingCount >= i.PendingLimit {
			err = commit()
			if err != nil {
				break
			}
			writer = i.DB.NewWriter()
		}
	}
	if err != nil && err != io.EOF {
		writer.Cancel()
		return err
	}
	return commit()
}
