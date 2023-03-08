package main

import (
	"math/rand"
	"testing"

	"github.com/aurora-is-near/relayer2-base/db/badger/core"
	"github.com/aurora-is-near/relayer2-base/db/badger/core/dbkey"
	"github.com/aurora-is-near/relayer2-base/db/codec"
	"github.com/aurora-is-near/relayer2-base/types/indexer"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/aurora-is-near/relayer2-base/utils"
	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestExportImport(t *testing.T) {
	chainID := uint64(42)
	blocks := generateBlocks(50, chainID)
	codec := codec.NewTinypackCodec()
	dbConf := core.Config{
		GcIntervalSeconds: 10,
		BadgerConfig: badger.DefaultOptions("").
			WithInMemory(true).
			WithLoggingLevel(badger.ERROR),
	}
	exportDB, err := core.NewDB(dbConf, codec)
	require.Nil(t, err, err)
	defer exportDB.Close()

	for _, b := range blocks {
		err := insertBlock(exportDB, b)
		require.NoError(t, err)
	}

	fs := afero.NewMemMapFs()
	archiver, err := NewArchiver(fs, codec)
	require.NoError(t, err)

	e := Exporter{
		DB:       exportDB.BadgerDB(),
		Archiver: archiver,
		ChainID:  chainID,
		Decoder:  codec,
		Height:   0,
	}
	err = e.Export()
	require.NoError(t, err)

	exportDBData, err := allDBData(exportDB.BadgerDB(), chainID)
	require.NoError(t, err)
	require.NotEqual(t, 0, len(exportDBData))
	exportDB.Close()

	unarchiver, err := NewUnarchiver(fs, codec)
	require.NoError(t, err)

	importDB, err := core.NewDB(dbConf, codec)
	require.NoError(t, err)
	defer importDB.Close()

	i := Importer{
		DB:         importDB,
		Unarchiver: unarchiver,
		ChainID:    chainID,
	}
	err = i.Import()
	require.NoError(t, err)

	importDBData, err := allDBData(importDB.BadgerDB(), chainID)
	require.NoError(t, err)
	importDB.Close()

	require.Equal(t, len(exportDBData), len(importDBData), "database item count doesn't match")
	for i := range exportDBData {
		require.Equal(t, exportDBData[i], importDBData[i])
	}
}

func insertBlock(db *core.DB, block *indexer.Block) error {
	writer := db.NewWriter()
	defer writer.Cancel()

	chainId := block.ChainId
	height := block.Height
	hash := block.Hash
	err := writer.InsertBlock(chainId, height, hash, utils.IndexerBlockToDbBlock(block))
	if err != nil {
		return err
	}

	for i, t := range block.Transactions {
		txnIndex := uint64(i)
		err = writer.InsertTransaction(
			chainId,
			height,
			txnIndex,
			t.Hash,
			utils.IndexerTxnToDbTxn(t),
		)
		if err != nil {
			return err
		}
		for j, l := range t.Logs {
			err = writer.InsertLog(chainId, height, txnIndex, uint64(j), utils.IndexerLogToDbLog(l))
			if err != nil {
				return err
			}
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func allDBData(db *badger.DB, chainID uint64) ([][2][]byte, error) {
	data := make([][2][]byte, 0)
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = dbkey.Chain.Get(chainID)
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().Key()
			val, err := it.Item().ValueCopy(nil)
			if err != nil {
				return err
			}
			data = append(data, [2][]byte{key, val})
		}
		return nil
	})
	return data, err
}

func generateBlocks(amount int, chainID uint64) []*indexer.Block {
	rand.Seed(0)
	var counter uint64
	nextHex := func() string {
		counter++
		return hexutil.EncodeUint64(counter)
	}
	blocks := make([]*indexer.Block, 0, amount)
	for i := 0; i < amount; i++ {
		block := new(indexer.Block)

		block.ChainId = chainID
		block.Height = uint64(i)
		block.Sequence = uint64(i)
		block.Timestamp = indexer.Timestamp(i)

		block.GasLimit = primitives.QuantityFromHex(nextHex())
		block.GasUsed = primitives.QuantityFromHex(nextHex())
		block.Hash = primitives.Data32FromHex(nextHex())
		block.ParentHash = primitives.Data32FromHex(nextHex())
		block.TransactionsRoot = primitives.Data32FromHex(nextHex())
		block.ReceiptsRoot = primitives.Data32FromHex(nextHex())
		block.StateRoot = primitives.Data32FromHex(nextHex())
		block.Miner = primitives.Data20FromHex(nextHex())
		block.LogsBloom = primitives.Data256FromHex(nextHex())
		block.Transactions = make([]*indexer.Transaction, 0)

		for j := 0; j < 3; j++ {
			if 0.5 < rand.Float32() {
				break
			}
			tx := new(indexer.Transaction)

			tx.Hash = primitives.Data32FromHex(nextHex())
			tx.BlockHash = primitives.Data32FromHex(nextHex())
			tx.BlockHeight = uint64(i)
			tx.ChainId = chainID
			tx.TransactionIndex = uint64(j)
			tx.From = primitives.Data20FromHex(nextHex())
			// tx.To = func() *primitives.Data20 {
			// 	h := primitives.Data20FromHex(nextHex())
			// 	return &h
			// }()
			tx.Nonce = primitives.QuantityFromHex(nextHex())
			tx.GasPrice = primitives.QuantityFromHex(nextHex())
			tx.GasLimit = primitives.QuantityFromHex(nextHex())
			tx.GasUsed = 1
			tx.MaxPriorityFeePerGas = primitives.QuantityFromHex(nextHex())
			tx.MaxFeePerGas = primitives.QuantityFromHex(nextHex())
			tx.Value = primitives.QuantityFromHex(nextHex())
			tx.Input = indexer.InputOutputData(
				primitives.DataFromHex[primitives.VarLen](nextHex()),
			)
			tx.Output = indexer.InputOutputData(
				primitives.DataFromHex[primitives.VarLen](nextHex()),
			)
			tx.AccessList = []indexer.AccessList{}
			tx.TxType = 0
			tx.LogsBloom = primitives.Data256FromHex(nextHex())
			tx.ContractAddress = func() *primitives.Data20 {
				h := primitives.Data20FromHex(nextHex())
				return &h
			}()
			tx.V = 0
			tx.R = primitives.QuantityFromHex(nextHex())
			tx.S = primitives.QuantityFromHex(nextHex())

			nearHash := indexer.NearHash(primitives.Data32FromHex(nextHex()))
			tx.NearTransaction.Hash = &nearHash
			tx.NearTransaction.ReceiptHash = indexer.NearHash(primitives.Data32FromHex(nextHex()))

			tx.Logs = make([]*indexer.Log, 0)
			for k := 0; k < 5; k++ {
				log := new(indexer.Log)

				log.Address = primitives.Data20FromHex(nextHex())
				log.Topics = []indexer.Topic{
					indexer.Topic(primitives.Data32FromHex(nextHex())),
					indexer.Topic(primitives.Data32FromHex(nextHex())),
					indexer.Topic(primitives.Data32FromHex(nextHex())),
					indexer.Topic(primitives.Data32FromHex(nextHex())),
				}
				log.Data = indexer.InputOutputData(
					primitives.DataFromHex[primitives.VarLen](nextHex()),
				)

				tx.Logs = append(tx.Logs, log)

				if 0.5 < rand.Float32() {
					break
				}
			}

			block.Transactions = append(block.Transactions, tx)
		}

		blocks = append(blocks, block)
	}
	return blocks
}
