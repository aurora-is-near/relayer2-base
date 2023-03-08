package main

import (
	"fmt"
	"io"

	"github.com/aurora-is-near/relayer2-base/db/badger/core"
	"github.com/aurora-is-near/relayer2-base/db/badger/core/dbkey"
	"github.com/dgraph-io/badger/v3"
)

func PrintDBInfo(db *core.DB, w io.Writer) error {
	return db.BadgerDB().View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchValues: false,
			Prefix:         dbkey.Chains.Get(),
		})

		chainIDs := make(map[uint64]bool, 0)
		for it.Rewind(); it.Valid(); it.Next() {
			chainID := dbkey.Chain.ReadUintVar(it.Item().Key(), 0)
			chainIDs[chainID] = true
		}
		it.Close()

		fmt.Fprintln(w, "Found ChainIDs:")
		for cid := range chainIDs {
			fmt.Fprintf(w, " - %d / 0x%x\n", cid, cid)
		}

		for chainID := range chainIDs {
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
