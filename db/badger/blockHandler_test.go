package badger

import (
	"aurora-relayer-go-common/db"
	"aurora-relayer-go-common/utils"
	"bytes"
	"context"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"math"
	"strings"
	"testing"
)

const blockHandlerTestYaml = `
db:
  badger:
      gcIntervalSeconds: 1
      iterationTimeoutSeconds: 5
      iterationMaxItems: 10000
      logFilterTtlMinutes: 15
      index:
        maxJumps: 1000
        maxRangeScanners: 2
        maxValueFetchers: 2
        keysOnly: false
      options:
        InMemory: true
        DetectConflicts: true
`

func TestGetFunctions(t *testing.T) {
	blocks := [...]*utils.Block{
		{Height: 1, Sequence: 0, Hash: utils.HexStringToHash("a"), Transactions: []*utils.Transaction{{}}},
		{Height: 2, Sequence: 1, Hash: utils.HexStringToHash("b"), Transactions: []*utils.Transaction{{}, {}}},
		{Height: 3, Sequence: 2, Hash: utils.HexStringToHash("c"), Transactions: []*utils.Transaction{{}, {}, {}}},
		{Height: 4, Sequence: 3, Hash: utils.HexStringToHash("d"), Transactions: []*utils.Transaction{{}, {}, {}, {}}},
	}

	var bh db.BlockHandler
	ctx := context.Background()

	viper.SetConfigType("yml")
	err := viper.ReadConfig(strings.NewReader(blockHandlerTestYaml))
	if err != nil {
		panic(err)
	}

	ttable := []struct {
		name           string
		getter         func(key interface{}) (interface{}, error)
		key            interface{}
		want           interface{}
		wantMissing    interface{}
		wantMissingErr string
		testGot        func(t *testing.T, got interface{})
	}{
		{
			name: "BlockNumberToHash",
			getter: func(key interface{}) (interface{}, error) {
				return bh.BlockNumberToHash(ctx, utils.UintToUint256(key.(uint64)))
			},
			key:            blocks[2].Height,
			want:           &blocks[2].Hash,
			wantMissingErr: "Key not found",
		},
		{
			name: "BlockHashToNumber",
			getter: func(key interface{}) (interface{}, error) {
				return bh.BlockHashToNumber(ctx, key.(utils.H256))
			},
			testGot: func(t *testing.T, got interface{}) {
				assert.Equal(t, blocks[2].Height, got.(*utils.Uint256).Uint64())
			},
			key:            blocks[2].Hash,
			wantMissingErr: "Key not found",
		},
		{
			name: "CurrentBlockNumber",
			getter: func(key interface{}) (interface{}, error) {
				return bh.BlockNumber(ctx)
			},
			testGot: func(t *testing.T, got interface{}) {
				assert.Equal(t, blocks[3].Height, got.(*utils.Uint256).Uint64())
			},
			wantMissingErr: "Key not found",
		},
		{
			name: "CurrentBlockSequence",
			getter: func(key interface{}) (interface{}, error) {
				return bh.CurrentBlockSequence(ctx), nil
			},
			want:        blocks[3].Sequence,
			wantMissing: uint64(0),
		},
		{
			name: "GetBlockByHash",
			getter: func(key interface{}) (interface{}, error) {
				return bh.GetBlockByHash(ctx, key.(utils.H256))
			},
			key: blocks[2].Hash,
			testGot: func(t *testing.T, got interface{}) {
				assert.Equal(t, blocks[2].Height, got.(*utils.Block).Height)
			},
			wantMissingErr: "Key not found",
		},
		{
			name: "GetBlockByNumber",
			getter: func(key interface{}) (interface{}, error) {
				return bh.GetBlockByNumber(ctx, key.(utils.Uint256))
			},
			key: utils.UintToUint256(blocks[2].Height),
			testGot: func(t *testing.T, got interface{}) {
				assert.Equal(t, blocks[2].Hash, got.(*utils.Block).Hash)
			},
			wantMissingErr: "Key not found",
		},
		{
			name: "GetBlockHashesSinceNumber",
			getter: func(key interface{}) (interface{}, error) {
				return bh.GetBlockHashesSinceNumber(ctx, key.(utils.Uint256))
			},
			key:         utils.UintToUint256(blocks[1].Height),
			want:        []utils.H256{blocks[2].Hash, blocks[3].Hash},
			wantMissing: []utils.H256{},
		},
		{
			name: "GetBlockTransactionCountByNumber",
			getter: func(key interface{}) (interface{}, error) {
				return bh.GetBlockTransactionCountByNumber(ctx, key.(utils.Uint256))
			},
			key:            utils.UintToUint256(blocks[2].Height),
			want:           blocks[2].TxCount(),
			wantMissing:    int64(0),
			wantMissingErr: "Key not found",
		},
		{
			name: "GetBlockTransactionCountByHash",
			getter: func(key interface{}) (interface{}, error) {
				return bh.GetBlockTransactionCountByHash(ctx, key.(utils.H256))
			},
			key:            blocks[2].Hash,
			want:           blocks[2].TxCount(),
			wantMissing:    int64(0),
			wantMissingErr: "Key not found",
		},
	}
	for _, tc := range ttable {
		t.Run(tc.name, func(t *testing.T) {

			// bh.SaveBlock = true TODO saveBlock, saveTxn, saveLog config to be ported

			bh, err = NewBlockHandler()
			if err != nil {
				panic(err)
			}
			defer bh.Close()

			got, err := tc.getter(tc.key)
			if tc.wantMissingErr != "" {
				assert.ErrorContains(t, err, tc.wantMissingErr)
			} else {
				assert.Nil(t, err)
			}

			if tc.wantMissing != nil {
				assert.Equal(t, tc.wantMissing, got, "value should be missing while data hasn't been inserted")
			} else {
				assert.Nil(t, got)
			}

			for _, b := range blocks {
				err := bh.InsertBlock(b)
				assert.Nil(t, err)
			}

			got, err = tc.getter(tc.key)
			assert.Nil(t, err)
			if tc.testGot != nil {
				tc.testGot(t, got)
			} else {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestGetBlockByHashFetchesTransactions(t *testing.T) {
	block := &utils.Block{
		ChainId:  5,
		Hash:     utils.HexStringToHash("abc"),
		Sequence: 0,
		Transactions: []*utils.Transaction{
			{
				Hash:             utils.HexStringToHash("a"),
				TransactionIndex: 0,
				Nonce:            utils.IntToUint256(1),
				R:                utils.IntToUint256(1),
				S:                utils.IntToUint256(1),
				GasLimit:         utils.IntToUint256(1),
				GasPrice:         utils.IntToUint256(1),
				GasUsed:          utils.IntToUint256(1),
				Value:            utils.IntToUint256(1),
			},
			{
				Hash:             utils.HexStringToHash("b"),
				TransactionIndex: 1,
				Nonce:            utils.IntToUint256(2),
				R:                utils.IntToUint256(2),
				S:                utils.IntToUint256(2),
				GasLimit:         utils.IntToUint256(2),
				GasPrice:         utils.IntToUint256(2),
				GasUsed:          utils.IntToUint256(2),
				Value:            utils.IntToUint256(2),
			},
			{
				Hash:             utils.HexStringToHash("c"),
				TransactionIndex: 2,
				Nonce:            utils.IntToUint256(3),
				R:                utils.IntToUint256(3),
				S:                utils.IntToUint256(3),
				GasLimit:         utils.IntToUint256(3),
				GasPrice:         utils.IntToUint256(3),
				GasUsed:          utils.IntToUint256(3),
				Value:            utils.IntToUint256(3),
			},
		},
	}

	viper.SetConfigType("yml")
	err := viper.ReadConfig(strings.NewReader(blockHandlerTestYaml))
	if err != nil {
		panic(err)
	}
	bh, err := NewBlockHandler()
	if err != nil {
		panic(err)
	}
	defer bh.Close()
	// bh.SaveBlock = true TODO saveBlock, saveTxn, saveLog config to be ported
	// bh.SaveTx = true TODO saveBlock, saveTxn, saveLog config to be ported

	got, err := bh.GetBlockByHash(context.Background(), block.Hash)
	assert.ErrorContains(t, err, "Key not found")
	assert.Nil(t, got)

	err = bh.InsertBlock(block)
	assert.Nil(t, err)

	got, err = bh.GetBlockByHash(context.Background(), block.Hash)
	assert.Nil(t, err)
	assert.Equal(t, block.Transactions, got.Transactions)
}

func TestGetTransaction(t *testing.T) {

	viper.SetConfigType("yml")
	err := viper.ReadConfig(strings.NewReader(blockHandlerTestYaml))
	if err != nil {
		panic(err)
	}
	bh, err := NewBlockHandler()
	if err != nil {
		panic(err)
	}
	defer bh.Close()

	ctx := context.Background()

	block := &utils.Block{
		Hash:   utils.HexStringToHash("abc"),
		Height: 1,
		Transactions: []*utils.Transaction{
			{Hash: utils.HexStringToHash("a")},
			{Hash: utils.HexStringToHash("b")},
			{Hash: utils.HexStringToHash("c")},
		},
	}

	// bh.SaveTx = true TODO saveBlock, saveTxn, saveLog config to be ported

	one := utils.IntToUint256(1)

	tx, err := bh.GetTransactionByBlockHashAndIndex(ctx, block.Hash, one)
	assert.ErrorContains(t, err, "Key not found")
	assert.Nil(t, tx)

	tx, err = bh.GetTransactionByBlockNumberAndIndex(ctx, utils.UintToUint256(block.Height), one)
	assert.ErrorContains(t, err, "Key not found")
	assert.Nil(t, tx)

	tx, err = bh.GetTransactionByHash(ctx, block.Transactions[1].Hash)
	assert.ErrorContains(t, err, "Key not found")
	assert.Nil(t, tx)

	err = bh.InsertBlock(block)
	assert.Nil(t, err)

	tx, err = bh.GetTransactionByBlockHashAndIndex(ctx, block.Hash, one)
	assert.Nil(t, err)
	assert.Equal(t, block.Transactions[1].Hash, tx.Hash)

	tx, err = bh.GetTransactionByBlockNumberAndIndex(ctx, utils.UintToUint256(block.Height), one)
	assert.Nil(t, err)
	assert.Equal(t, block.Transactions[1].Hash, tx.Hash)

	tx, err = bh.GetTransactionByHash(ctx, block.Transactions[1].Hash)
	assert.Nil(t, err)
	assert.Equal(t, block.Transactions[1].Hash, tx.Hash)
}

func TestGetLogs(t *testing.T) {
	ctx := context.Background()
	blocks := []*utils.Block{
		{
			Hash:   utils.HexStringToHash("a"),
			Height: 0,
			Transactions: []*utils.Transaction{
				{
					Hash: utils.HexStringToHash("aa"),
					Logs: []*utils.Log{
						{Address: utils.HexStringToAddress("a1"), Topics: []utils.Bytea{utils.Bytea("1"), utils.Bytea("2"), utils.Bytea("3"), utils.Bytea("4")}},
						{Address: utils.HexStringToAddress("a2"), Topics: []utils.Bytea{utils.Bytea("1"), utils.Bytea("2"), utils.Bytea("3"), utils.Bytea("4")}},
					},
				},
				{
					Hash: utils.HexStringToHash("ab"),
					Logs: []*utils.Log{
						{Address: utils.HexStringToAddress("a2"), Topics: []utils.Bytea{utils.Bytea("1"), utils.Bytea("2"), utils.Bytea("3"), utils.Bytea("4")}},
					},
				},
			},
		},
		{
			Hash:   utils.HexStringToHash("b"),
			Height: 1,
			Transactions: []*utils.Transaction{
				{
					Hash: utils.HexStringToHash("ba"),
					Logs: []*utils.Log{
						{Address: utils.HexStringToAddress("a2"), Topics: []utils.Bytea{utils.Bytea("4"), utils.Bytea("3"), utils.Bytea("22"), utils.Bytea("1")}},
					},
				},
				{
					Hash: utils.HexStringToHash("bb"),
					Logs: []*utils.Log{
						{Address: utils.HexStringToAddress("a2"), Topics: []utils.Bytea{utils.Bytea("1"), utils.Bytea("2"), utils.Bytea("3"), utils.Bytea("4")}},
					},
				},
			},
		},
		{
			Hash:   utils.HexStringToHash("c"),
			Height: 2,
			Transactions: []*utils.Transaction{
				{
					Hash: utils.HexStringToHash("ca"),
					Logs: []*utils.Log{
						{Address: utils.HexStringToAddress("a1"), Topics: []utils.Bytea{utils.Bytea("4"), utils.Bytea("3"), utils.Bytea("2"), utils.Bytea("1")}},
					},
				},
			},
		},
	}

	viper.SetConfigType("yml")
	err := viper.ReadConfig(strings.NewReader(blockHandlerTestYaml))
	if err != nil {
		panic(err)
	}
	bh, err := NewBlockHandler()
	if err != nil {
		panic(err)
	}
	defer bh.Close()

	// no items have been inserted yet
	filter := newFilter()
	addFromAndTo(filter, 0, 2)
	l, err := bh.GetLogs(ctx, *filter)
	logs := *l
	assert.Nil(t, err)
	assert.Len(t, logs, 0)

	for _, b := range blocks {
		err := bh.InsertBlock(b)
		assert.Nil(t, err)
	}

	// returns all items
	filter = newFilter()
	addFromAndTo(filter, 0, 2)
	l, err = bh.GetLogs(ctx, *filter)
	logs = *l
	assert.Nil(t, err)
	assert.Len(t, logs, 6)

	// there are 3 logs in the first block
	filter = newFilter()
	addFromAndTo(filter, 0, 0)
	l, err = bh.GetLogs(ctx, *filter)
	logs = *l
	assert.Nil(t, err)
	assert.Len(t, logs, 3)

	// there are only three blocks
	filter = newFilter()
	addFromAndTo(filter, 3, 4)
	l, err = bh.GetLogs(ctx, *filter)
	logs = *l
	assert.Nil(t, err)
	assert.Len(t, logs, 0)

	// two items have 3 as second topic
	filter = newFilter()
	addTopic(filter, 1, "3")
	l, err = bh.GetLogs(ctx, *filter)
	logs = *l
	assert.Nil(t, err)
	assert.Len(t, logs, 2)
	// TODO set 1 worker in DB index scan options for tests.
	// Currently tests fails occasionally due to ordering:
	// assert.Equal(t, blocks[1].Transactions[0].Hash, logs[0].TransactionHash)
	// assert.Equal(t, blocks[1].Hash, logs[0].BlockHash)
	// assert.Equal(t, blocks[2].Transactions[0].Hash, logs[1].TransactionHash)
	// assert.Equal(t, blocks[2].Hash, logs[1].BlockHash)

	// four items have 1 as first topic
	filter = newFilter()
	addTopic(filter, 0, "1")
	l, err = bh.GetLogs(ctx, *filter)
	logs = *l
	assert.Nil(t, err)
	assert.Len(t, logs, 4)

	// all items have either 1 or 4 as first topic
	filter = newFilter()
	addTopic(filter, 0, "1")
	addTopic(filter, 0, "4")
	l, err = bh.GetLogs(ctx, *filter)
	logs = *l
	assert.Nil(t, err)
	assert.Len(t, logs, 6)

	// two items have a1 as address
	filter = newFilter()
	addAddress(filter, "a1")
	l, err = bh.GetLogs(ctx, *filter)
	logs = *l
	assert.Nil(t, err)
	assert.Len(t, logs, 2)
	// assert.Equal(t, blocks[0].Transactions[0].Hash, logs[0].TransactionHash)
	// assert.Equal(t, blocks[0].Hash, logs[0].BlockHash)
	// assert.Equal(t, blocks[2].Transactions[0].Hash, logs[1].TransactionHash)
	// assert.Equal(t, blocks[2].Hash, logs[1].BlockHash)

	// all items have either a1 or a2 as address
	filter = newFilter()
	addAddress(filter, "a1")
	addAddress(filter, "a2")
	l, err = bh.GetLogs(ctx, *filter)
	logs = *l
	assert.Nil(t, err)
	assert.Len(t, logs, 6)

	// no items have 1 as the first and last topic
	filter = newFilter()
	addTopic(filter, 0, "1")
	addTopic(filter, 3, "1")
	l, err = bh.GetLogs(ctx, *filter)
	logs = *l
	assert.Nil(t, err)
	assert.Len(t, logs, 0)

	// only matching log is in the second block's second transaction
	filter = newFilter()
	addTopic(filter, 0, "1")
	addFromAndTo(filter, 1, 2)
	addAddress(filter, "a2")
	l, err = bh.GetLogs(ctx, *filter)
	logs = *l
	assert.Nil(t, err)
	assert.Len(t, logs, 1)
	assert.Equal(t, blocks[1].Hash, logs[0].BlockHash)
	assert.Equal(t, blocks[1].Transactions[1].Hash, logs[0].TransactionHash)
}

func TestGetLogsForTransaction(t *testing.T) {
	ctx := context.Background()
	blocks := []*utils.Block{
		{
			Hash:   utils.HexStringToHash("a"),
			Height: 0,
			Transactions: []*utils.Transaction{
				{
					Hash:             utils.HexStringToHash("aa"),
					BlockHash:        utils.HexStringToHash("a"),
					TransactionIndex: 0,
					BlockHeight:      0,
					Logs: []*utils.Log{
						{Address: utils.HexStringToAddress("a1")},
					},
				},
				{
					Hash:             utils.HexStringToHash("ab"),
					BlockHash:        utils.HexStringToHash("a"),
					TransactionIndex: 1,
					BlockHeight:      0,
					Logs: []*utils.Log{
						{Address: utils.HexStringToAddress("a2")},
						{Address: utils.HexStringToAddress("a2")},
					},
				},
			},
		},
		{
			Hash:   utils.HexStringToHash("b"),
			Height: 1,
			Transactions: []*utils.Transaction{
				{
					Hash:             utils.HexStringToHash("ba"),
					BlockHash:        utils.HexStringToHash("b"),
					TransactionIndex: 0,
					BlockHeight:      1,
					Logs: []*utils.Log{
						{Address: utils.HexStringToAddress("a2")},
					},
				},
				{
					Hash:             utils.HexStringToHash("bb"),
					BlockHash:        utils.HexStringToHash("b"),
					TransactionIndex: 1,
					BlockHeight:      1,
					Logs: []*utils.Log{
						{Address: utils.HexStringToAddress("a2")},
					},
				},
			},
		},
	}

	viper.SetConfigType("yml")
	err := viper.ReadConfig(strings.NewReader(blockHandlerTestYaml))
	if err != nil {
		panic(err)
	}
	bh, err := NewBlockHandler()
	if err != nil {
		panic(err)
	}
	defer bh.Close()
	// bh.SaveLog = true TODO saveBlock, saveTxn, saveLog config to be ported

	tx := blocks[0].Transactions[1]

	// no items have been inserted yet
	logs, err := bh.GetLogsForTransaction(ctx, tx)
	assert.Nil(t, err)
	assert.Len(t, logs, 0)

	for _, b := range blocks {
		err := bh.InsertBlock(b)
		assert.Nil(t, err)
	}

	logs, err = bh.GetLogsForTransaction(ctx, tx)
	assert.Nil(t, err)
	assert.Len(t, logs, 2)
	assert.Equal(t, logs[0].BlockHash, blocks[0].Hash)
	assert.Equal(t, logs[0].TransactionHash, blocks[0].Transactions[1].Hash)
	assert.Equal(t, logs[0].LogIndex, utils.IntToUint256(0))
	assert.Equal(t, logs[1].BlockHash, blocks[0].Hash)
	assert.Equal(t, logs[1].TransactionHash, blocks[0].Transactions[1].Hash)
	assert.Equal(t, logs[1].LogIndex, utils.IntToUint256(1))
}

func newFilter() *utils.LogFilter {
	from, to := utils.IntToUint256(0), utils.IntToUint256(math.MaxInt64)
	return &utils.LogFilter{FromBlock: &from, ToBlock: &to}
}

func addFromAndTo(filter *utils.LogFilter, from, to int) {
	*filter.FromBlock = utils.IntToUint256(from)
	*filter.ToBlock = utils.IntToUint256(to)
}

func addTopic(filter *utils.LogFilter, idx int, topic string) {
	if filter.Topics == nil {
		filter.Topics = [][][]byte{{}, {}, {}, {}}
	}
	filter.Topics[idx] = append(filter.Topics[idx], []byte(topic))
}

func addAddress(filter *utils.LogFilter, address string) {
	toAdd := utils.HexStringToAddress(address).Bytes()
	for _, a := range filter.Address {
		if bytes.Compare(a, toAdd) == 0 {
			return
		}
	}
	filter.Address = append(filter.Address, toAdd)
}
