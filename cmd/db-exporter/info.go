package main

import (
	"fmt"
	"io"

	"github.com/aurora-is-near/relayer2-base/db/badger/core"
	"github.com/aurora-is-near/relayer2-base/db/badger/core/dbkey"
	"github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
)

func PrintDBInfo(db *core.DB, w io.Writer) error {
	return db.BadgerDB().View(func(txn *badger.Txn) error {
		chainIDs, err := getChainIDs(txn)
		if err != nil {
			return errors.Wrap(err, "failed to get chainIDs from database")
		}

		fmt.Fprintln(w, "Found ChainIDs:")
		for _, chainID := range chainIDs {
			fmt.Fprintf(w, " - %d / 0x%x\n", chainID, chainID)
		}

		for _, chainID := range chainIDs {
			fmt.Fprintf(w, "\nFound for ChainID %d / 0x%x:\n", chainID, chainID)
			it := txn.NewIterator(badger.IteratorOptions{
				PrefetchValues: false,
				Prefix:         dbkey.BlocksData.Get(chainID),
			})
			blockCount := uint64(0)
			for it.Rewind(); it.Valid(); it.Next() {
				blockCount++
			}
			it.Close()
			fmt.Fprintf(w, " - %d blocks\n", blockCount)

			it = txn.NewIterator(badger.IteratorOptions{
				PrefetchValues: false,
				Prefix:         dbkey.TxsData.Get(chainID),
			})
			txCount := uint64(0)
			for it.Rewind(); it.Valid(); it.Next() {
				txCount++
			}
			it.Close()
			fmt.Fprintf(w, " - %d transactions\n", txCount)

			it = txn.NewIterator(badger.IteratorOptions{
				PrefetchValues: false,
				Prefix:         dbkey.Logs.Get(chainID),
			})
			logCount := uint64(0)
			for it.Rewind(); it.Valid(); it.Next() {
				logCount++
			}
			it.Close()
			fmt.Fprintf(w, " - %d logs\n", logCount)
		}

		return nil
	})
}

func getChainIDs(txn *badger.Txn) ([]uint64, error) {
	it := txn.NewIterator(badger.IteratorOptions{
		PrefetchValues: false,
		Prefix:         dbkey.Chains.Get(),
	})
	defer it.Close()

	chainIDs := make([]uint64, 0)
	nextChainID := uint64(0)

	for it.Rewind(); it.Valid(); it.Seek(dbkey.Chain.Get(nextChainID)) {
		currentChainId := dbkey.Chain.ReadUintVar(it.Item().Key(), 0)
		if len(chainIDs) > 0 && currentChainId == chainIDs[len(chainIDs)-1] {
			break
		}
		chainIDs = append(chainIDs, currentChainId)
		nextChainID = currentChainId + 1
	}

	return chainIDs, nil
}
