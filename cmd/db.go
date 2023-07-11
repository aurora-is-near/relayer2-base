package cmd

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/aurora-is-near/relayer2-base/db/badger/core"
	"github.com/aurora-is-near/relayer2-base/db/codec"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/spf13/cobra"
)

var chainId uint64
var blockType string

func GetLastBlockCmd() *cobra.Command {
	getLastBlockCmd := &cobra.Command{
		Use:   "get-last-block",
		Short: "Command to retieve last block height or sequence number from db",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please provide db path")
				return
			}

			dbPath := args[0]
			err := dbView(dbPath, func(txn *core.ViewTxn) error {
				switch blockType {
				case "sequence":
					data, err := txn.ReadIndexerState(chainId)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					if data == nil {
						fmt.Println("No blocks in db, perhaps you need to provide --chain-id flag?")
						os.Exit(1)
					}
					fmt.Println(binary.BigEndian.Uint64(data))
				case "height":
					key, err := txn.ReadLatestBlockKey(chainId)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					if key == nil {
						fmt.Println("No blocks in db, perhaps you need to provide --chain-id flag?")
						os.Exit(1)
					}
					fmt.Println(key.Height)
				default:
					fmt.Println("Unknown block type")
					os.Exit(1)
				}

				return nil
			})

			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}
	getLastBlockCmd.PersistentFlags().Uint64VarP(&chainId, "chain-id", "c", 1313161554, "Chain ID")
	getLastBlockCmd.PersistentFlags().StringVarP(&blockType, "type", "t", "height", "Type of block head: height or sequence")
	return getLastBlockCmd
}

// dbView opens db in read-only mode and calls fn with ViewTxn
func dbView(dbPath string, fn func(txn *core.ViewTxn) error) error {
	config := core.Config{
		BadgerConfig:      badger.DefaultOptions(dbPath).WithLogger(nil).WithReadOnly(true),
		GcIntervalSeconds: 10,
	}
	db, err := core.NewDB(config, codec.NewTinypackCodec())
	if err != nil {
		return err
	}
	return db.View(fn)
}
