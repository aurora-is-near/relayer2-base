package endpoint

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aurora-is-near/relayer2-base/db"
	"github.com/aurora-is-near/relayer2-base/db/badger"
	"github.com/aurora-is-near/relayer2-base/types"
	"github.com/aurora-is-near/relayer2-base/types/common"
	"github.com/aurora-is-near/relayer2-base/types/indexer"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/aurora-is-near/relayer2-base/types/request"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const ethTestYaml = `
db:
  badger:
    core:
      gcIntervalSeconds: 10
      scanRangeThreshold: 3000
      maxScanIterators: 10000
      filterTtlMinutes: 15
      options:
        Dir: /tmp/relayer/data
        InMemory: true
        DetectConflicts: false
`

func TestLogFilterUnmarshalJSON(t *testing.T) {
	data := `{"address":["` +
		fmt.Sprintf("0x%040x", 0x1) +
		`","` +
		fmt.Sprintf("0x%040x", 0x2) +
		`"],"fromBlock":"0xb1","toBlock":"0xb2","blockHash":"` +
		fmt.Sprintf("0x%064x", 0xabcdf) +
		`"}`
	wantHash := common.MustHexStringToHash("0xabcdf")
	wantFromBlock := common.IntToBN64(0xb1)
	wantToBlock := common.IntToBN64(0xb2)
	want := request.Filter{
		Addresses: []common.Address{
			common.MustHexStringToAddress("1"),
			common.MustHexStringToAddress("2"),
		},
		FromBlock: &wantFromBlock,
		ToBlock:   &wantToBlock,
		BlockHash: &wantHash,
	}
	var result request.Filter
	err := jsoniter.Unmarshal([]byte(data), &result)
	assert.Nil(t, err)
	assert.Equal(t, want, result)
}

func TestFormatFilterOptions(t *testing.T) {

	var eth *Eth
	ctx := context.Background()

	viper.SetConfigType("yml")
	err := viper.ReadConfig(strings.NewReader(ethTestYaml))
	if err != nil {
		panic(err)
	}

	ca1 := primitives.MustData20FromHex(fmt.Sprintf("0x%040x", 0x2))
	ca2 := primitives.MustData20FromHex(fmt.Sprintf("0x%040x", 0x1))
	ca3 := primitives.MustData20FromHex(fmt.Sprintf("0x%040x", 0x2))
	ca4 := primitives.MustData20FromHex(fmt.Sprintf("0x%040x", 0x1))
	ca5 := primitives.MustData20FromHex(fmt.Sprintf("0x%040x", 0x3))

	blockHash := primitives.MustData32FromHex(fmt.Sprintf("0x%064x", 0x2))
	parentHash := primitives.MustData32FromHex(fmt.Sprintf("0x%064x", 0x1))
	filterHash := common.MustHexStringToHash(blockHash.Hex())

	data20 := primitives.MustData20FromHex("0x11")
	data32 := primitives.MustData32FromHex("0x22")
	data256 := primitives.MustData256FromHex("0x33")
	quantity := primitives.QuantityFromHex("0x44")
	nearTxn := indexer.NearTransaction{
		Hash:        nil,
		ReceiptHash: indexer.NearHash(data32),
	}

	var blockData = indexer.Block{
		ChainId:          1313161554,
		Height:           1,
		Hash:             blockHash,
		ParentHash:       parentHash,
		Miner:            data20,
		TransactionsRoot: data32,
		ReceiptsRoot:     data32,
		StateRoot:        data32,
		GasLimit:         quantity,
		GasUsed:          quantity,
		LogsBloom:        data256,
		Transactions: []*indexer.Transaction{
			{ContractAddress: &ca1, BlockHash: blockHash, TransactionIndex: 0, Hash: primitives.MustData32FromHex("0x1"), From: data20, Nonce: quantity, GasPrice: quantity, GasLimit: quantity, MaxFeePerGas: quantity, MaxPriorityFeePerGas: quantity, Value: quantity, S: quantity, R: quantity, NearTransaction: nearTxn, LogsBloom: data256},
			{ContractAddress: &ca2, BlockHash: blockHash, TransactionIndex: 1, Hash: primitives.MustData32FromHex("0x2"), From: data20, Nonce: quantity, GasPrice: quantity, GasLimit: quantity, MaxFeePerGas: quantity, MaxPriorityFeePerGas: quantity, Value: quantity, S: quantity, R: quantity, NearTransaction: nearTxn, LogsBloom: data256},
			{ContractAddress: &ca3, BlockHash: blockHash, TransactionIndex: 2, Hash: primitives.MustData32FromHex("0x3"), From: data20, Nonce: quantity, GasPrice: quantity, GasLimit: quantity, MaxFeePerGas: quantity, MaxPriorityFeePerGas: quantity, Value: quantity, S: quantity, R: quantity, NearTransaction: nearTxn, LogsBloom: data256},
			{ContractAddress: &ca4, BlockHash: blockHash, TransactionIndex: 3, Hash: primitives.MustData32FromHex("0x4"), From: data20, Nonce: quantity, GasPrice: quantity, GasLimit: quantity, MaxFeePerGas: quantity, MaxPriorityFeePerGas: quantity, Value: quantity, S: quantity, R: quantity, NearTransaction: nearTxn, LogsBloom: data256},
			{ContractAddress: &ca5, BlockHash: blockHash, TransactionIndex: 4, Hash: primitives.MustData32FromHex("0x5"), From: data20, Nonce: quantity, GasPrice: quantity, GasLimit: quantity, MaxFeePerGas: quantity, MaxPriorityFeePerGas: quantity, Value: quantity, S: quantity, R: quantity, NearTransaction: nearTxn, LogsBloom: data256},
		},
	}

	reqFrom := common.IntToBN64(10)
	reqTo := common.IntToBN64(20)

	ttable := []struct {
		name             string
		data             request.Filter
		wantFrom, wantTo *uint64
		wantAddress      []primitives.Data20
		wantTopics       [][]primitives.Data32
		wantErr          string
	}{
		{
			name:        "empty options",
			data:        request.Filter{},
			wantFrom:    &blockData.Height,
			wantTo:      nil,
			wantAddress: []primitives.Data20{},
			wantTopics:  nil,
		},
		{
			name: "blockHash is added",
			data: request.Filter{
				BlockHash: &filterHash,
			},
			wantFrom:    &blockData.Height,
			wantTo:      &blockData.Height,
			wantAddress: []primitives.Data20{},
			wantTopics:  nil,
		},
		{
			name: "block range is not overwritten",
			data: request.Filter{
				FromBlock: &reqFrom,
				ToBlock:   &reqTo,
			},
			wantFrom:    reqFrom.Uint64(),
			wantTo:      reqTo.Uint64(),
			wantAddress: []primitives.Data20{},
			wantTopics:  nil,
		},
		{
			name: "addresses get added once",
			data: request.Filter{
				Addresses: []common.Address{
					common.BytesToAddress(primitives.MustData20FromHex(fmt.Sprintf("0x%040x", 0x2)).Bytes()),
					common.BytesToAddress(primitives.MustData20FromHex(fmt.Sprintf("0x%040x", 0x1)).Bytes()),
					common.BytesToAddress(primitives.MustData20FromHex(fmt.Sprintf("0x%040x", 0x2)).Bytes()),
					common.BytesToAddress(primitives.MustData20FromHex(fmt.Sprintf("0x%040x", 0x1)).Bytes()),
					common.BytesToAddress(primitives.MustData20FromHex(fmt.Sprintf("0x%040x", 0x3)).Bytes()),
				},
			},
			wantFrom: &blockData.Height,
			wantTo:   nil,
			wantAddress: []primitives.Data20{
				primitives.MustData20FromHex(fmt.Sprintf("0x%040x", 0x2)),
				primitives.MustData20FromHex(fmt.Sprintf("0x%040x", 0x1)),
				primitives.MustData20FromHex(fmt.Sprintf("0x%040x", 0x3)),
			},
			wantTopics: nil,
		},
		{
			name: "topics are added as is", // TODO: add stronger topics validation/restrict the type from []byte when unmarshalling?
			data: request.Filter{
				Topics: request.Topics{
					{primitives.MustData32FromHex(fmt.Sprintf("0x%064x", 0x1111))},
					{primitives.MustData32FromHex(fmt.Sprintf("0x%064x", 0x2222))},
					{primitives.MustData32FromHex(fmt.Sprintf("0x%064x", 0x3333))},
				},
			},
			wantFrom:    &blockData.Height,
			wantTo:      nil,
			wantAddress: []primitives.Data20{},
			wantTopics: [][]primitives.Data32{
				{primitives.MustData32FromHex(fmt.Sprintf("0x%064x", 0x1111))},
				{primitives.MustData32FromHex(fmt.Sprintf("0x%064x", 0x2222))},
				{primitives.MustData32FromHex(fmt.Sprintf("0x%064x", 0x3333))},
			},
		},
	}
	for _, tc := range ttable {
		t.Run(tc.name, func(t *testing.T) {

			bh, err := badger.NewBlockHandler()
			if err != nil {
				panic(err)
			}
			fh, err := badger.NewFilterHandler()
			if err != nil {
				panic(err)
			}

			handler := db.StoreHandler{
				BlockHandler:  bh,
				FilterHandler: fh,
			}
			defer handler.Close()

			err = bh.InsertBlock(&blockData)
			assert.Nil(t, err)

			baseEndpoint := New(handler)
			eth = NewEth(baseEndpoint)

			want := &types.Filter{
				FromBlock: tc.wantFrom,
				ToBlock:   tc.wantTo,
				Addresses: tc.wantAddress,
				Topics:    tc.wantTopics,
			}

			got, err := eth.parseRequestFilter(ctx, &tc.data)
			assert.Nil(t, err)
			assert.Equal(t, want, got)
		})
	}
}

func TestTopicsUnmarshalJSON(t *testing.T) {
	ttable := []struct {
		data    string
		want    request.Topics
		wantErr string
	}{
		{
			data: `["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",null,["0xabc","0x123"]]`,
			want: request.Topics{
				{primitives.MustData32FromHex(`0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef`)},
				{},
				{primitives.MustData32FromHex("0xabc"), primitives.MustData32FromHex("0x123")},
				{},
			},
		},
		{
			data: `["0x1","0x2","0x3","0x4","0x5"]`,
			want: request.Topics{
				{primitives.MustData32FromHex("0x1")},
				{primitives.MustData32FromHex("0x2")},
				{primitives.MustData32FromHex("0x3")},
				{primitives.MustData32FromHex("0x4")},
			},
		},
		{
			data:    `[`,
			want:    nil,
			wantErr: "do not know how to skip",
		},
	}
	for _, tc := range ttable {
		t.Run(tc.data, func(t *testing.T) {
			var ts request.Topics
			err := jsoniter.Unmarshal([]byte(tc.data), &ts)
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tc.want, ts)
		})
	}
}
