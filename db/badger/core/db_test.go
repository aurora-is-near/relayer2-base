package core

import (
	"aurora-relayer-go-common/db/codec"
	dbt "aurora-relayer-go-common/types/db"
	"aurora-relayer-go-common/types/primitives"
	"aurora-relayer-go-common/types/response"
	"context"
	"encoding/binary"
	"github.com/dgraph-io/badger/v3"
	"log"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

const testChainId = 88005553535

func genBytes(length int, seeds ...uint64) []byte {
	hash := crypto.NewKeccakState()
	for _, seed := range seeds {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, seed)
		hash.Write(buf)
	}
	result := make([]byte, length)
	hash.Read(result)
	return result
}

func genUint64(seeds ...uint64) uint64 {
	return binary.BigEndian.Uint64(genBytes(8, seeds...))
}

func genBool(seeds ...uint64) bool {
	return genUint64(seeds...)%2 > 0
}

func genAddress(seeds ...uint64) primitives.Data20 {
	return primitives.Data20FromBytes(genBytes(20, seeds...))
}

func genHash(seeds ...uint64) primitives.Data32 {
	return primitives.Data32FromBytes(genBytes(32, seeds...))
}

func genQuantity(seeds ...uint64) primitives.Quantity {
	return primitives.QuantityFromBytes(genBytes(32, seeds...))
}

func genLogsBloom(seeds ...uint64) primitives.Data256 {
	return primitives.Data256FromBytes(genBytes(256, seeds...))
}

func genVarData(minLength int, maxLength int, seeds ...uint64) primitives.VarData {
	length := minLength
	if maxLength > minLength {
		length += int(genUint64(append(seeds, 1)...) % uint64(maxLength-minLength))
	}
	return primitives.VarDataFromBytes(genBytes(length, append(seeds, 2)...))
}

func genBlock(seed uint64) *dbt.Block {
	return &dbt.Block{
		ParentHash:       genHash(seed, 1),
		Miner:            genAddress(seed, 2),
		Timestamp:        genUint64(seed, 3),
		GasLimit:         genQuantity(seed, 4),
		GasUsed:          genQuantity(seed, 5),
		LogsBloom:        genLogsBloom(seed, 6),
		TransactionsRoot: genHash(seed, 7),
		StateRoot:        genHash(seed, 8),
		ReceiptsRoot:     genHash(seed, 9),
		Size:             genUint64(seed, 10),
	}
}

func genTx(seed uint64) *dbt.Transaction {
	tx := &dbt.Transaction{
		Type:                 genUint64(seed, 1) % 3,
		From:                 genAddress(seed, 2),
		IsContractDeployment: genBool(seed, 3),
		ToOrContract:         genAddress(seed, 4),
		Nonce:                genQuantity(seed, 5),
		GasPrice:             genQuantity(seed, 6),
		GasLimit:             genQuantity(seed, 7),
		GasUsed:              genUint64(seed, 8),
		Value:                genQuantity(seed, 9),
		Input:                genVarData(0, 512, seed, 10),
		NearHash:             genHash(seed, 11),
		NearReceiptHash:      genHash(seed, 12),
		Status:               genBool(seed, 13),
		V:                    genUint64(seed, 14),
		R:                    genQuantity(seed, 15),
		S:                    genQuantity(seed, 16),
		LogsBloom:            genLogsBloom(seed, 17),
	}
	if tx.Type >= 1 {
		tx.AccessList.Content = []dbt.AccessListEntry{}
		n := int(genUint64(seed, 18) % 10)
		for i := 0; i < n; i++ {
			ale := dbt.AccessListEntry{
				Address: genAddress(seed, 19, uint64(i)),
			}
			m := int(genUint64(seed, 20, uint64(i)) % 10)
			ale.StorageKeys.Content = []primitives.Data32{}
			for j := 0; j < m; j++ {
				sk := genHash(seed, 21, uint64(i), uint64(j))
				ale.StorageKeys.Content = append(ale.StorageKeys.Content, sk)
			}
			tx.AccessList.Content = append(tx.AccessList.Content, ale)
		}
	}
	if tx.Type >= 2 {
		tx.MaxPriorityFeePerGas = genQuantity(seed, 22)
		tx.MaxFeePerGas = genQuantity(seed, 23)
	}
	return tx
}

func genLog(addressSeed uint64, dataSeed uint64, topicSeeds ...uint64) *dbt.Log {
	log := &dbt.Log{
		Address: genAddress(addressSeed),
		Data:    genVarData(0, 100, dataSeed),
	}
	log.Topics.Content = []primitives.Data32{}
	for _, topicSeed := range topicSeeds {
		log.Topics.Content = append(log.Topics.Content, genHash(topicSeed))
	}
	return log
}

type blockSeed struct {
	height uint64
}

func (bs *blockSeed) getBlockKey() *dbt.BlockKey {
	return &dbt.BlockKey{
		Height: bs.height,
	}
}

func (bs *blockSeed) getBlockHash() primitives.Data32 {
	return genHash(bs.height * 10)
}

func (bs *blockSeed) getBlockData() *dbt.Block {
	return genBlock(bs.height*10 + 1)
}

func (bs *blockSeed) getBlockResponse(txSeeds []txSeed, full bool) *response.Block {
	txs := []any{}
	for _, txSeed := range txSeeds {
		if full {
			txs = append(txs, txSeed.getTxResponse())
		} else {
			txs = append(txs, txSeed.getTxHash())
		}
	}

	return makeBlockResponse(
		bs.height,
		bs.getBlockHash(),
		*bs.getBlockData(),
		txs,
	)
}

type txSeed struct {
	height uint64
	index  uint64
}

func (ts *txSeed) getTxKey() *dbt.TransactionKey {
	return &dbt.TransactionKey{
		BlockHeight:      ts.height,
		TransactionIndex: ts.index,
	}
}

func (ts *txSeed) getTxHash() primitives.Data32 {
	return genHash((ts.height*1000+ts.index)*10 + 2)
}

func (ts *txSeed) getTxData() *dbt.Transaction {
	return genTx((ts.height*1000+ts.index)*10 + 3)
}

func (ts *txSeed) getTxResponse() *response.Transaction {
	bs := blockSeed{height: ts.height}
	return makeTransactionResponse(
		testChainId,
		ts.height,
		ts.index,
		bs.getBlockHash(),
		ts.getTxHash(),
		ts.getTxData(),
	)
}

func (ts *txSeed) getTxReceiptResponse(logSeeds []logSeed) *response.TransactionReceipt {
	logResponses := []*response.Log{}
	for _, logSeed := range logSeeds {
		logResponses = append(logResponses, logSeed.getLogResponse())
	}

	bs := blockSeed{height: ts.height}
	return makeTransactionReceiptResponse(
		ts.height,
		ts.index,
		bs.getBlockHash(),
		ts.getTxHash(),
		ts.getTxData(),
		logResponses,
	)
}

type logSeed struct {
	height      uint64
	txIndex     uint64
	logIndex    uint64
	addressSeed uint64
	topicSeeds  []uint64
}

func (ls *logSeed) getLogKey() *dbt.LogKey {
	return &dbt.LogKey{
		BlockHeight:      ls.height,
		TransactionIndex: ls.txIndex,
		LogIndex:         ls.logIndex,
	}
}

func (ls *logSeed) getLogData() *dbt.Log {
	dataSeed := ((ls.height*1000+ls.txIndex)*1000+ls.logIndex)*10 + 4
	return genLog(ls.addressSeed, dataSeed, ls.topicSeeds...)
}

func (ls *logSeed) getLogResponse() *response.Log {
	bs := blockSeed{height: ls.height}
	ts := txSeed{height: ls.height, index: ls.txIndex}
	return makeLogResponse(
		ls.height,
		ls.txIndex,
		ls.logIndex,
		bs.getBlockHash(),
		ts.getTxHash(),
		ls.getLogData(),
	)
}

var blockSeeds = []blockSeed{
	{101},
	{103}, {104}, {105},
	{120}, {121},
	{1_000_001},
	{9_000_001}, {9_000_002}, {9_000_003}, {9_000_004}, {9_000_005},
	{9_000_007}, {9_000_008},
}

var txSeeds = []txSeed{
	{103, 0}, {103, 1}, {103, 2}, {103, 3}, {103, 4},
	{104, 0}, {104, 1},
	{105, 0}, {105, 1},
	{120, 0}, {120, 1},
	{121, 0}, {121, 1}, {121, 2},
	{1_000_001, 0},
	{9_000_002, 0},
	{9_000_003, 0},
	{9_000_007, 0}, {9_000_007, 1}, {9_000_007, 2},
	{9_000_008, 0}, {9_000_008, 1}, {9_000_008, 2},
}

var logSeeds = []logSeed{
	{103, 0, 0, 555_0_0, []uint64{555_1_0, 555_2_1, 555_3_0}},
	{103, 1, 0, 555_0_2, []uint64{555_1_1, 555_2_1, 555_3_0, 555_4_1}},
	{103, 1, 1, 555_0_0, []uint64{555_1_2, 555_2_0}},
	{103, 1, 2, 555_0_1, []uint64{555_1_1, 555_2_1, 555_3_1}},
	{103, 2, 0, 555_0_0, []uint64{555_1_2, 555_2_2}},
	{103, 2, 1, 555_0_2, []uint64{555_1_1, 555_2_2, 555_3_0, 555_4_1}},
	{103, 3, 0, 555_0_1, []uint64{555_1_1, 555_2_2}},
	{103, 3, 1, 555_0_2, []uint64{555_1_1, 555_2_2, 555_3_1, 555_4_2}},
	{103, 4, 0, 555_0_1, []uint64{555_1_2, 555_2_0, 555_3_1}},
	{103, 4, 1, 555_0_1, []uint64{555_1_1, 555_2_2, 555_3_2}},
	{104, 0, 0, 555_0_0, []uint64{}},
	{104, 0, 1, 555_0_2, []uint64{555_1_1, 555_2_2, 555_3_2}},
	{104, 1, 0, 555_0_2, []uint64{555_1_2}},
	{105, 0, 0, 555_0_2, []uint64{555_1_0, 555_2_2, 555_3_0}},
	{105, 0, 1, 555_0_1, []uint64{555_1_0, 555_2_1, 555_3_2}},
	{105, 0, 2, 555_0_1, []uint64{}},
	{120, 0, 0, 555_0_2, []uint64{}},
	{120, 0, 1, 555_0_0, []uint64{555_1_1, 555_2_1, 555_3_2, 555_4_2}},
	{121, 0, 0, 555_0_0, []uint64{555_1_1, 555_2_1, 555_3_0}},
	{121, 1, 0, 555_0_2, []uint64{}},
	{121, 2, 0, 555_0_1, []uint64{555_1_0, 555_2_1, 555_3_1}},
	{121, 2, 1, 555_0_2, []uint64{555_1_0, 555_2_1, 555_3_0}},
	{1000001, 0, 0, 555_0_0, []uint64{555_1_2, 555_2_0, 555_3_1, 555_4_2}},
	{1000001, 0, 1, 555_0_2, []uint64{555_1_2, 555_2_0, 555_3_2, 555_4_1}},
	{1000001, 0, 2, 555_0_0, []uint64{555_1_1, 555_2_1}},
	{9000003, 0, 0, 555_0_2, []uint64{555_1_2, 555_2_0, 555_3_0}},
	{9000003, 0, 1, 555_0_0, []uint64{555_1_0, 555_2_2, 555_3_2}},
	{9000007, 0, 0, 555_0_2, []uint64{555_1_0}},
	{9000007, 0, 1, 555_0_2, []uint64{555_1_0}},
	{9000007, 0, 2, 555_0_2, []uint64{}},
	{9000007, 1, 0, 555_0_0, []uint64{}},
	{9000007, 1, 1, 555_0_1, []uint64{555_1_2}},
	{9000008, 0, 0, 555_0_1, []uint64{555_1_0, 555_2_1, 555_3_1}},
	{9000008, 0, 1, 555_0_1, []uint64{555_1_0}},
	{9000008, 1, 0, 555_0_2, []uint64{555_1_2, 555_2_1}},
	{9000008, 1, 1, 555_0_2, []uint64{555_1_1, 555_2_2, 555_3_0}},
}

// func TestGenerateLogSeeds(t *testing.T) {
// 	rand.Seed(654342)
// 	for _, txSeed := range txSeeds {
// 		n := rand.Intn(4)
// 		for i := 0; i < n; i++ {
// 			fmt.Printf("{%v, %v, %v, 555_0_%v, []uint64{", txSeed.height, txSeed.index, i, rand.Intn(3))
// 			t := rand.Intn(5)
// 			for j := 0; j < t; j++ {
// 				fmt.Printf("555_%v_%v", j+1, rand.Intn(3))
// 				if j < t-1 {
// 					fmt.Printf(", ")
// 				}
// 			}
// 			fmt.Println("}},")
// 		}
// 	}
// }

const suppressSecondaryLogging = true

type testLogger struct {
	errCnt  int32
	warnCnt int32
}

func (l *testLogger) getErrCnt() int32 {
	return atomic.LoadInt32(&l.errCnt)
}

func (l *testLogger) getWarnCnt() int32 {
	return atomic.LoadInt32(&l.warnCnt)
}

func (l testLogger) Errorf(f string, v ...interface{}) {
	atomic.AddInt32(&l.errCnt, 1)
	log.Printf("ERROR: "+f, v...)
}

func (l testLogger) Warningf(f string, v ...interface{}) {
	atomic.AddInt32(&l.warnCnt, 1)
	log.Printf("WARNING: "+f, v...)
}

func (l testLogger) Infof(f string, v ...interface{}) {
	if !suppressSecondaryLogging {
		log.Printf("INFO: "+f, v...)
	}
}

func (l testLogger) Debugf(f string, v ...interface{}) {
	if !suppressSecondaryLogging {
		log.Printf("DEBUG: "+f, v...)
	}
}

func initTestDb(t *testing.T) (*DB, *testLogger) {

	opts := badger.DefaultOptions("")
	opts.InMemory = true
	core, err := Open(opts, 10)
	require.NoError(t, err, "DB must tryOpen")

	logger := testLogger{}
	testDb := &DB{
		codec:                 codec.NewTinypackCodec(),
		maxLogScanIterators:   1000,
		logScanRangeThreshold: 1000,
		filterTtlMinutes:      15,
		logger:                logger,
		core:                  core,
	}

	// testDb := &DB{
	// 	CoreOpts: &dbcore.DBCoreOpts{
	// 		GCIntervalSeconds: 10,
	// 		InMemory:          true,
	// 	},
	// 	Encoder:               tinypack.DefaultEncoder(),
	// 	Decoder:               tinypack.DefaultDecoder(),
	// 	maxLogScanIterators:   1000,
	// 	logScanRangeThreshold: 1000,
	// }
	//
	// logger := &testLogger{}
	// err := testDb.Open(logger)

	writer := testDb.NewWriter()

	for _, blockSeed := range blockSeeds {
		require.NoError(t, writer.InsertBlock(
			testChainId,
			blockSeed.height,
			blockSeed.getBlockHash(),
			blockSeed.getBlockData(),
		), "InsertBlock must work")
	}
	for _, txSeed := range txSeeds {
		require.NoError(t, writer.InsertTransaction(
			testChainId,
			txSeed.height,
			txSeed.index,
			txSeed.getTxHash(),
			txSeed.getTxData(),
		), "InsertTransaction must work")
	}
	for _, logSeed := range logSeeds {
		require.NoError(t, writer.InsertLog(
			testChainId,
			logSeed.height,
			logSeed.txIndex,
			logSeed.logIndex,
			logSeed.getLogData(),
		), "InsertLog must work")
	}
	require.NoError(t, writer.Flush(), "Flush must work")

	require.EqualValues(t, 0, logger.getErrCnt(), "There should be no errors")
	require.EqualValues(t, 0, logger.getWarnCnt(), "There should be no warnings")

	return testDb, &logger
}

func testView(t *testing.T, ensureNoErrors bool, ensureNoWarnings bool, fn func(txn *ViewTxn) error) {
	testDb, logger := initTestDb(t)
	viewErr := testDb.View(fn)
	require.NoError(t, viewErr, "db.View must work")
	require.NoError(t, testDb.Close(), "close must work")
	if ensureNoErrors {
		require.EqualValues(t, 0, logger.getErrCnt(), "There should be no errors")
	}
	if ensureNoWarnings {
		require.EqualValues(t, 0, logger.getWarnCnt(), "There should be no warnings")
	}
}

func TestReadBlock(t *testing.T) {
	testView(t, true, true, func(txn *ViewTxn) error {
		earliestBlockKey, err := txn.ReadEarliestBlockKey(testChainId)
		require.NoError(t, err, "ReadEarliestBlockKey must work")
		require.Equal(t, blockSeeds[0].getBlockKey(), earliestBlockKey, "Earliest block key must be right")

		latestBlockKey, err := txn.ReadLatestBlockKey(testChainId)
		require.NoError(t, err, "ReadLatestBlockKey must work")
		require.Equal(t, blockSeeds[len(blockSeeds)-1].getBlockKey(), latestBlockKey, "Latest block key must be right")

		for _, blockSeed := range blockSeeds {
			blockKey, err := txn.ReadBlockKey(testChainId, blockSeed.getBlockHash())
			require.NoError(t, err, "ReadBlockKey must work")
			require.Equal(t, blockSeed.getBlockKey(), blockKey, "ReadBlockKey must return right value")

			txs := []txSeed{}
			for _, txSeed := range txSeeds {
				if txSeed.height == blockSeed.height {
					txs = append(txs, txSeed)
				}
			}

			txCount, err := txn.ReadBlockTxCount(testChainId, *blockSeed.getBlockKey())
			require.NoError(t, err, "ReadBlockTxCount must work")
			require.EqualValues(t, len(txs), txCount, "ReadBlockTxCount must return right value")

			block, err := txn.ReadBlock(testChainId, *blockSeed.getBlockKey(), false)
			require.NoError(t, err, "ReadBlock (full=false) must work")
			require.Equal(t, blockSeed.getBlockResponse(txs, false), block, "ReadBlock (full=false) must return right value")

			block, err = txn.ReadBlock(testChainId, *blockSeed.getBlockKey(), true)
			require.NoError(t, err, "ReadBlock (full=true) must work")
			require.Equal(t, blockSeed.getBlockResponse(txs, true), block, "ReadBlock (full=true) must return right value")
		}

		bounds := []*dbt.BlockKey{
			{Height: 0},
			{Height: 20},
			{Height: 50},
			{Height: 102},
			{Height: 110},
			{Height: 500},
			{Height: 555_555},
			{Height: 1_000_000},
			{Height: 5_000_000},
			{Height: 9_000_000},
			{Height: 9_000_006},
			{Height: 10_000_000},
			{Height: 1_000_000_000},
		}
		for _, blockSeed := range blockSeeds {
			bounds = append(bounds, blockSeed.getBlockKey())
		}

		for _, from := range bounds {
			for _, to := range bounds {
				if from.CompareTo(to) > 0 {
					continue
				}
				for _, limit := range []int{1, 2, 5, 10, 100} {
					expectedLastKey := from.Prev()
					limited := false
					expectedResult := []primitives.Data32{}
					for _, blockSeed := range blockSeeds {
						if blockSeed.getBlockKey().CompareTo(from) < 0 {
							continue
						}
						if blockSeed.getBlockKey().CompareTo(to) > 0 {
							continue
						}
						if len(expectedResult) == limit {
							limited = true
							break
						}
						expectedResult = append(expectedResult, blockSeed.getBlockHash())
						expectedLastKey = blockSeed.getBlockKey()
					}
					if !limited {
						expectedLastKey = to
					}

					blockHashes, lastKey, err := txn.ReadBlockHashes(context.Background(), testChainId, from, to, limit)
					if limited {
						require.Equal(t, ErrLimited, err, "ReadBlockHashes must return ErrLimited")
					} else {
						require.NoError(t, err, "ReadBlockHashes must work without errors")
					}
					require.Equal(t, expectedLastKey, lastKey, "ReadBlockHashes must return right lastKey")
					require.Equal(t, expectedResult, blockHashes, "ReadBlockHashes must return right result")
				}
			}
		}

		return nil
	})
}

func TestReadTransaction(t *testing.T) {
	testView(t, true, true, func(txn *ViewTxn) error {
		earliestTxKey, err := txn.ReadEarliestTxKey(testChainId)
		require.NoError(t, err, "ReadEarliestTxKey must work")
		require.Equal(t, txSeeds[0].getTxKey(), earliestTxKey, "Earliest tx key must be right")

		latestTxKey, err := txn.ReadLatestTxKey(testChainId)
		require.NoError(t, err, "ReadLatestBlockKey must work")
		require.Equal(t, txSeeds[len(txSeeds)-1].getTxKey(), latestTxKey, "Latest tx key must be right")

		for _, txSeed := range txSeeds {
			txKey, err := txn.ReadTxKey(testChainId, txSeed.getTxHash())
			require.NoError(t, err, "ReadTxKey must work")
			require.Equal(t, txSeed.getTxKey(), txKey, "ReadTxKey must return right value")

			tx, err := txn.ReadTx(testChainId, *txSeed.getTxKey())
			require.NoError(t, err, "ReadTx must work")
			require.Equal(t, txSeed.getTxResponse(), tx, "ReadTx must return right value")

			logs := []logSeed{}
			for _, logSeed := range logSeeds {
				if logSeed.height == txSeed.height && logSeed.txIndex == txSeed.index {
					logs = append(logs, logSeed)
				}
			}

			txReceipt, err := txn.ReadTxReceipt(testChainId, *txSeed.getTxKey())
			require.NoError(t, err, "ReadTxReceipt must work")
			require.Equal(t, txSeed.getTxReceiptResponse(logs), txReceipt, "ReadTx must return right value")
		}

		bounds := []*dbt.TransactionKey{
			{BlockHeight: 0, TransactionIndex: 0},
			{BlockHeight: 0, TransactionIndex: 1},
			{BlockHeight: 0, TransactionIndex: 5},
			{BlockHeight: 0, TransactionIndex: 1000},
			{BlockHeight: 50, TransactionIndex: 0},
			{BlockHeight: 50, TransactionIndex: 7},
			{BlockHeight: 102, TransactionIndex: 0},
			{BlockHeight: 102, TransactionIndex: 3},
			{BlockHeight: 103, TransactionIndex: 5},
			{BlockHeight: 104, TransactionIndex: 6},
			{BlockHeight: 110, TransactionIndex: 2},
			{BlockHeight: 500, TransactionIndex: 500},
			{BlockHeight: 1_000_000, TransactionIndex: 0},
			{BlockHeight: 1_000_000, TransactionIndex: 1000},
			{BlockHeight: 9_000_001, TransactionIndex: 0},
			{BlockHeight: 9_000_001, TransactionIndex: 1},
			{BlockHeight: 9_000_001, TransactionIndex: 1000},
			{BlockHeight: 9_000_003, TransactionIndex: 1},
			{BlockHeight: 9_000_004, TransactionIndex: 0},
			{BlockHeight: 9_000_004, TransactionIndex: 500},
			{BlockHeight: 9_000_006, TransactionIndex: 1000},
			{BlockHeight: 9_000_007, TransactionIndex: 3},
			{BlockHeight: 9_000_007, TransactionIndex: 999},
			{BlockHeight: 9_000_008, TransactionIndex: 3},
			{BlockHeight: 9_000_008, TransactionIndex: 500},
			{BlockHeight: 1_000_000_000, TransactionIndex: 0},
			{BlockHeight: 1_000_000_000, TransactionIndex: 1000},
		}
		for _, txSeed := range txSeeds {
			bounds = append(bounds, txSeed.getTxKey())
		}

		for _, limit := range []int{1, 2, 5, 10, 100} {
			log.Printf("Testing ReadTransactions with limit=%v", limit)
			for _, from := range bounds {
				for _, to := range bounds {
					if from.CompareTo(to) > 0 {
						continue
					}
					expectedLastKey := from.Prev()
					limited := false
					expectedHashes := []any{}
					expectedTxes := []any{}
					for _, txSeed := range txSeeds {
						if txSeed.getTxKey().CompareTo(from) < 0 {
							continue
						}
						if txSeed.getTxKey().CompareTo(to) > 0 {
							continue
						}
						if len(expectedHashes) == limit {
							limited = true
							break
						}
						expectedHashes = append(expectedHashes, txSeed.getTxHash())
						expectedTxes = append(expectedTxes, txSeed.getTxResponse())
						expectedLastKey = txSeed.getTxKey()
					}
					if !limited {
						expectedLastKey = to
					}

					hashes, lastKey, err := txn.ReadTransactions(context.Background(), testChainId, from, to, false, limit)
					if limited {
						require.Equal(t, ErrLimited, err, "ReadTransactions (full=false) must return ErrLimited")
					} else {
						require.NoError(t, err, "ReadTransactions (full=false) must work without errors")
					}
					require.Equal(t, expectedLastKey, lastKey, "ReadTransactions (full=false) must return right lastKey")
					require.Equal(t, expectedHashes, hashes, "ReadTransactions (full=false) must return right result")

					txes, lastKey, err := txn.ReadTransactions(context.Background(), testChainId, from, to, true, limit)
					if limited {
						require.Equal(t, ErrLimited, err, "ReadTransactions (full=true) must return ErrLimited")
					} else {
						require.NoError(t, err, "ReadTransactions (full=true) must work without errors")
					}
					require.Equal(t, expectedLastKey, lastKey, "ReadTransactions (full=true) must return right lastKey")
					require.Equal(t, expectedTxes, txes, "ReadTransactions (full=true) must return right result")
				}
			}
		}

		return nil
	})
}

func assembleLogFilter(filterDesc []string) (
	addressMap map[uint64]struct{},
	addressFilter []primitives.Data20,
	topicMaps []map[uint64]struct{},
	topicFilters [][]primitives.Data32,
) {
	addressMap = map[uint64]struct{}{}
	addressFilter = []primitives.Data20{}
	topicMaps = []map[uint64]struct{}{}
	topicFilters = [][]primitives.Data32{}

	for i, filter := range filterDesc {
		if i > 0 {
			topicMaps = append(topicMaps, map[uint64]struct{}{})
			topicFilters = append(topicFilters, []primitives.Data32{})
		}
		for j := 0; j <= 9; j++ {
			if strings.Contains(filter, strconv.FormatUint(uint64(j), 10)) {
				seed := uint64((555*10+i)*10 + j)
				if i == 0 {
					addressMap[seed] = struct{}{}
					addressFilter = append(addressFilter, genAddress(seed))
				} else {
					topicMaps[i-1][seed] = struct{}{}
					topicFilters[i-1] = append(topicFilters[i-1], genHash(seed))
				}
			}
		}
	}
	return
}

func getLogSearchExpectedResult(
	from *dbt.LogKey,
	to *dbt.LogKey,
	addressMap map[uint64]struct{},
	topicMaps []map[uint64]struct{},
	limit int,
) ([]*response.Log, *dbt.LogKey, error) {

	lastKey := from.Prev()
	result := []*response.Log{}

	for _, logSeed := range logSeeds {
		if logSeed.getLogKey().CompareTo(from) < 0 {
			continue
		}
		if logSeed.getLogKey().CompareTo(to) > 0 {
			continue
		}
		if len(addressMap) > 0 {
			if _, ok := addressMap[logSeed.addressSeed]; !ok {
				continue
			}
		}
		valid := true
		for i, topicMap := range topicMaps {
			if len(topicMap) == 0 {
				continue
			}
			if i >= len(logSeed.topicSeeds) {
				valid = false
				break
			}
			if _, ok := topicMap[logSeed.topicSeeds[i]]; !ok {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}
		if len(result) == limit {
			return result, lastKey, ErrLimited
		}
		result = append(result, logSeed.getLogResponse())
		lastKey = logSeed.getLogKey()
	}
	return result, to, nil
}

func TestReadLog(t *testing.T) {
	testView(t, true, true, func(txn *ViewTxn) error {
		earliestLogKey, err := txn.ReadEarliestLogKey(testChainId)
		require.NoError(t, err, "ReadEarliestLogKey must work")
		require.Equal(t, logSeeds[0].getLogKey(), earliestLogKey, "Earliest log key must be right")

		latestLogKey, err := txn.ReadLatestLogKey(testChainId)
		require.NoError(t, err, "ReadLatestLogKey must work")
		require.Equal(t, logSeeds[len(logSeeds)-1].getLogKey(), latestLogKey, "Latest log key must be right")

		bounds := []*dbt.LogKey{
			{BlockHeight: 0, TransactionIndex: 0, LogIndex: 0},
			{BlockHeight: 0, TransactionIndex: 1000, LogIndex: 0},
			{BlockHeight: 0, TransactionIndex: 0, LogIndex: 1000},
			{BlockHeight: 0, TransactionIndex: 1000, LogIndex: 1000},
			{BlockHeight: 50, TransactionIndex: 50, LogIndex: 50},
			{BlockHeight: 103, TransactionIndex: 500, LogIndex: 0},
			{BlockHeight: 103, TransactionIndex: 500, LogIndex: 500},
			{BlockHeight: 110, TransactionIndex: 0, LogIndex: 0},
			{BlockHeight: 120, TransactionIndex: 2, LogIndex: 0},
			{BlockHeight: 121, TransactionIndex: 1000, LogIndex: 1000},
			{BlockHeight: 500_000, TransactionIndex: 0, LogIndex: 0},
			{BlockHeight: 500_000, TransactionIndex: 500, LogIndex: 500},
			{BlockHeight: 500_000, TransactionIndex: 1000, LogIndex: 1000},
			{BlockHeight: 5_000_000, TransactionIndex: 1, LogIndex: 1},
			{BlockHeight: 1_000_000_000, TransactionIndex: 0, LogIndex: 0},
			{BlockHeight: 1_000_000_000, TransactionIndex: 1000, LogIndex: 1000},
		}
		addBound := func(bound *dbt.LogKey) {
			if bound == nil {
				return
			}
			for _, existingBound := range bounds {
				if existingBound.CompareTo(bound) == 0 {
					return
				}
			}
			bounds = append(bounds, bound)
		}
		for _, logSeed := range logSeeds {
			addBound(logSeed.getLogKey().Prev())
			addBound(logSeed.getLogKey())
			addBound(logSeed.getLogKey().Next())
		}

		filterDescs := [][]string{
			{},
			{"", "", "", "", ""},
			{"012", "012", "012", "012", "012"},
			{"01", "", "12"},
			{"023456", "1789", "", "056", ""},
			{"0123456789", "0123456789", "0123456789", "0123456789", "0123456789"},
			{"0123456789", "123456789", "0123456789", "013456789", "0123456789"},
			{"02", "02"},
			{"023456789", "023456789"},
			{"", "", "", "", "0"},
			{"", "", "", "", "1"},
			{"02", "02", "02", "02", "01"},
		}

		for filterNumber, filterDesc := range filterDescs {
			log.Printf("Testing ReadLogs with filter [%v/%v]", filterNumber+1, len(filterDescs))

			addressMap, addressFilter, topicMaps, topicFilters := assembleLogFilter(filterDesc)
			// log.Printf("Addresses: %v", addressMap)
			// log.Printf("Topics: %v", topicMaps)

			for _, limit := range []int{1, 2, 5, 10, 100} {
				log.Printf("Testing ReadLogs with limit=%v", limit)
				for _, from := range bounds {
					for _, to := range bounds {
						if from.CompareTo(to) > 0 {
							continue
						}
						expectedLogs, expectedLastKey, expectedErr := getLogSearchExpectedResult(
							from, to, addressMap, topicMaps, limit,
						)

						logs, lastKey, err := txn.ReadLogs(
							context.Background(),
							testChainId,
							from,
							to,
							addressFilter,
							topicFilters,
							limit,
						)
						require.Equal(t, expectedErr, err, "ReadLogs must return right errors")
						require.Equal(t, expectedLastKey, lastKey, "ReadLogs must return right lastKey")
						require.Equal(t, expectedLogs, logs, "ReadLogs must return right result")
					}
				}
			}
		}

		return nil
	})
}

// func TestReadLogManual(t *testing.T) {
// 	testView(t, true, true, func(txn *ViewTxn) errors {
// 		from := &db.LogKey{0, 0, 0}
// 		to := &db.LogKey{1e9, 0, 0}
// 		addressMap, addressFilter, topicMaps, topicFilters := assembleLogFilter([]string{"02", "02", "02", "02", "01"})
// 		limit := 100

// 		expectedResult, _, _ := getLogSearchExpectedResult(from, to, addressMap, topicMaps, limit)
// 		expectedJson, _ := json.MarshalIndent(expectedResult, "", "  ")
// 		fmt.Printf("Expected: %s\n", expectedJson)

// 		result, _, _ := txn.ReadLogs(context.Background(), testChainId, from, to, addressFilter, topicFilters, limit)
// 		js, _ := json.MarshalIndent(result, "", "  ")
// 		fmt.Printf("Actual: %s\n", js)
// 		_ = js

// 		return nil
// 	})
// }
