package db

import (
	"aurora-relayer-go-common/db/badger2/core/dbcore"
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
	"aurora-relayer-go-common/db/badger2/core/dbresponses"
	"aurora-relayer-go-common/db/badger2/core/dbtypes"
	"aurora-relayer-go-common/tinypack"
	"context"
	"encoding/binary"
	"log"
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

func genAddress(seeds ...uint64) dbp.Data20 {
	return dbp.Data20FromBytes(genBytes(20, seeds...))
}

func genHash(seeds ...uint64) dbp.Data32 {
	return dbp.Data32FromBytes(genBytes(32, seeds...))
}

func genQuantity(seeds ...uint64) dbp.Quantity {
	return dbp.QuantityFromBytes(genBytes(32, seeds...))
}

func genLogsBloom(seeds ...uint64) dbp.Data256 {
	return dbp.Data256FromBytes(genBytes(256, seeds...))
}

func genVarData(minLength int, maxLength int, seeds ...uint64) dbp.VarData {
	length := minLength
	if maxLength > minLength {
		length += int(genUint64(append(seeds, 1)...) % uint64(maxLength-minLength))
	}
	return dbp.VarDataFromBytes(genBytes(length, append(seeds, 2)...))
}

func genBlock(seed uint64) *dbtypes.Block {
	return &dbtypes.Block{
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

func genTx(seed uint64) *dbtypes.Transaction {
	tx := &dbtypes.Transaction{
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
		tx.AccessList.Content = []dbtypes.AccessListEntry{}
		n := int(genUint64(seed, 18) % 10)
		for i := 0; i < n; i++ {
			ale := dbtypes.AccessListEntry{
				Address: genAddress(seed, 19, uint64(i)),
			}
			m := int(genUint64(seed, 20, uint64(i)) % 10)
			ale.StorageKeys.Content = []dbp.Data32{}
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

func genLog(addressSeed uint64, dataSeed uint64, topicSeeds ...uint64) *dbtypes.Log {
	log := &dbtypes.Log{
		Address: genAddress(addressSeed),
		Data:    genVarData(0, 100, dataSeed),
	}
	log.Topics.Content = []dbp.Data32{}
	for _, topicSeed := range topicSeeds {
		log.Topics.Content = append(log.Topics.Content, genHash(topicSeed))
	}
	return log
}

type blockSeed struct {
	height uint64
}

func (bs *blockSeed) getBlockKey() *dbtypes.BlockKey {
	return &dbtypes.BlockKey{
		Height: bs.height,
	}
}

func (bs *blockSeed) getBlockHash() dbp.Data32 {
	return genHash(bs.height * 10)
}

func (bs *blockSeed) getBlockData() *dbtypes.Block {
	return genBlock(bs.height*10 + 1)
}

func (bs *blockSeed) getBlockResponse(txSeeds []txSeed, full bool) *dbresponses.Block {
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

func (ts *txSeed) getTxKey() *dbtypes.TransactionKey {
	return &dbtypes.TransactionKey{
		BlockHeight:      ts.height,
		TransactionIndex: ts.index,
	}
}

func (ts *txSeed) getTxHash() dbp.Data32 {
	return genHash((ts.height*1000+ts.index)*10 + 2)
}

func (ts *txSeed) getTxData() *dbtypes.Transaction {
	return genTx((ts.height*1000+ts.index)*10 + 3)
}

func (ts *txSeed) getTxResponse() *dbresponses.Transaction {
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

func (ts *txSeed) getTxReceiptResponse(logSeeds []*logSeed) *dbresponses.TransactionReceipt {
	logResponses := []*dbresponses.Log{}
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

func (ls *logSeed) getLogKey() *dbtypes.LogKey {
	return &dbtypes.LogKey{
		BlockHeight:      ls.height,
		TransactionIndex: ls.txIndex,
		LogIndex:         ls.logIndex,
	}
}

func (ls *logSeed) getLogData() *dbtypes.Log {
	dataSeed := ((ls.height*1000+ls.txIndex)*1000+ls.logIndex)*10 + 4
	return genLog(ls.addressSeed, dataSeed, ls.topicSeeds...)
}

func (ls *logSeed) getLogResponse() *dbresponses.Log {
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
	{103, 0, 0, 555_0_0, []uint64{555_1_0}},
	{103, 1, 0, 555_0_0, []uint64{555_1_2, 555_2_2, 555_3_1}},
	{103, 2, 0, 555_0_0, []uint64{555_1_0}},
	{104, 0, 0, 555_0_1, []uint64{555_1_1, 555_2_1}},
	{104, 1, 0, 555_0_1, []uint64{555_1_0, 555_2_2, 555_3_2}},
	{104, 1, 1, 555_0_2, []uint64{555_1_1}},
	{105, 0, 0, 555_0_0, []uint64{555_1_2}},
	{120, 1, 0, 555_0_2, []uint64{}},
	{120, 1, 1, 555_0_1, []uint64{555_1_2}},
	{120, 1, 2, 555_0_1, []uint64{555_1_1}},
	{121, 0, 0, 555_0_0, []uint64{555_1_0, 555_2_1}},
	{121, 1, 0, 555_0_0, []uint64{555_1_2, 555_2_2, 555_3_1}},
	{121, 1, 1, 555_0_0, []uint64{555_1_2, 555_2_2, 555_3_1}},
	{121, 1, 2, 555_0_2, []uint64{555_1_0, 555_2_2, 555_3_0}},
	{121, 2, 0, 555_0_0, []uint64{}},
	{1000001, 0, 0, 555_0_0, []uint64{555_1_0, 555_2_1, 555_3_1}},
	{9000003, 0, 0, 555_0_2, []uint64{555_1_1, 555_2_0, 555_3_2}},
	{9000007, 2, 0, 555_0_2, []uint64{}},
	{9000007, 2, 1, 555_0_1, []uint64{}},
	{9000008, 0, 0, 555_0_0, []uint64{555_1_0}},
	{9000008, 2, 0, 555_0_1, []uint64{}},
	{9000008, 2, 1, 555_0_1, []uint64{555_1_1, 555_2_1}},
	{9000008, 2, 2, 555_0_1, []uint64{555_1_0, 555_2_1}},
}

// func TestGenerateLogSeeds(t *testing.T) {
// 	rand.Seed(654342)
// 	for _, txSeed := range txSeeds {
// 		n := rand.Intn(4)
// 		for i := 0; i < n; i++ {
// 			fmt.Printf("{%v, %v, %v, 555_0_%v, []uint64{", txSeed.height, txSeed.index, i, rand.Intn(3))
// 			t := rand.Intn(4)
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

func (l *testLogger) Errorf(f string, v ...interface{}) {
	atomic.AddInt32(&l.errCnt, 1)
	log.Printf("ERROR: "+f, v...)
}

func (l *testLogger) Warningf(f string, v ...interface{}) {
	atomic.AddInt32(&l.warnCnt, 1)
	log.Printf("WARNING: "+f, v...)
}

func (l *testLogger) Infof(f string, v ...interface{}) {
	if !suppressSecondaryLogging {
		log.Printf("INFO: "+f, v...)
	}
}

func (l *testLogger) Debugf(f string, v ...interface{}) {
	if !suppressSecondaryLogging {
		log.Printf("DEBUG: "+f, v...)
	}
}

func initTestDb(t *testing.T) (*DB, *testLogger) {
	testDb := &DB{
		CoreOpts: &dbcore.DBCoreOpts{
			GCIntervalSeconds: 10,
			InMemory:          true,
		},
		Encoder:               tinypack.DefaultEncoder(),
		Decoder:               tinypack.DefaultDecoder(),
		MaxLogScanIterators:   1000,
		LogScanRangeThreshold: 1000,
	}

	logger := &testLogger{}
	err := testDb.Open(logger)
	require.NoError(t, err, "DB must open")

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

	return testDb, logger
}

func testView(t *testing.T, ensureNoErrors bool, ensureNoWarnings bool, fn func(txn *ViewTxn) error) {
	testDb, logger := initTestDb(t)
	viewErr := testDb.View(fn)
	require.NoError(t, viewErr, "db.View must work")
	require.NoError(t, testDb.Close(), "Close must work")
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
		require.Equal(t, &dbtypes.BlockKey{Height: 101}, earliestBlockKey, "Earliest block key must be right")

		latestBlockKey, err := txn.ReadLatestBlockKey(testChainId)
		require.NoError(t, err, "ReadLatestBlockKey must work")
		require.Equal(t, &dbtypes.BlockKey{Height: 9_000_008}, latestBlockKey, "Latest block key must be right")

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

		bounds := []*dbtypes.BlockKey{
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
					expectedResult := []dbp.Data32{}
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
