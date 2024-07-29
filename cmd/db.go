package cmd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/aurora-is-near/relayer2-base/db/badger/core"
	"github.com/aurora-is-near/relayer2-base/db/codec"
	dbt "github.com/aurora-is-near/relayer2-base/types/db"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/spf13/cobra"
)

var chainId uint64
var blockType string

func GetLastBlockCmd() *cobra.Command {
	getLastBlockCmd := &cobra.Command{
		Use:   "get-last-block <dbPath>",
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

func GetBlockCmd() *cobra.Command {
	getLastBlockCmd := &cobra.Command{
		Use:   "get-block <dbPath> <height>",
		Short: "Command to retrieve block by height from db",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dbPath := args[0]
			height := args[1]
			heightUint64, err := strconv.ParseUint(height, 10, 64)
			if err != nil {
				return err
			}

			return dbView(dbPath, func(txn *core.ViewTxn) error {
				block, err := txn.ReadBlock(chainId, dbt.BlockKey{Height: heightUint64}, true)
				if err != nil {
					return err
				}

				if block == nil {
					return fmt.Errorf("no blocks found at height %d, check chain ID and DB path", heightUint64)
				}

				jsonBlock, err := json.MarshalIndent(block, "", "  ")
				if err != nil {
					return err
				}

				fmt.Println(string(jsonBlock))
				return nil
			})
		},
	}
	getLastBlockCmd.PersistentFlags().Uint64VarP(&chainId, "chain-id", "c", 1313161554, "Chain ID")
	return getLastBlockCmd
}

func FlattenDB() *cobra.Command {
	return &cobra.Command{
		Use:   "flatten-db <dbPath>",
		Short: "Command to flatten db (merge all LSM levels) to improve performance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dbPath := args[0]

			log.Printf("Opening database at %s...", dbPath)
			config := core.Config{
				BadgerConfig:      badger.DefaultOptions(dbPath),
				GcIntervalSeconds: 60 * 60 * 24 * 365, // We don't want GC to be running during the flattening
			}
			db, err := core.NewDB(config, codec.NewTinypackCodec())
			if err != nil {
				return fmt.Errorf("unable to open database: %w", err)
			}

			defer func() {
				log.Printf("Closing database")
				if err := db.Close(); err != nil {
					log.Printf("Error: unable to close database normally: %v", err)
				}
			}()

			log.Printf("Flattening database, that might take some time...")
			threads := runtime.GOMAXPROCS(0)
			if err := db.BadgerDB().Flatten(threads); err != nil {
				return fmt.Errorf("can't flatten database: %w", err)
			}

			return nil
		},
	}
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
