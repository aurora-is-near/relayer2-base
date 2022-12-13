package endpoint

import (
	"aurora-relayer-go-common/db"
	"aurora-relayer-go-common/db/badger"
	"aurora-relayer-go-common/types"
	"aurora-relayer-go-common/types/common"
	"aurora-relayer-go-common/types/indexer"
	"aurora-relayer-go-common/types/primitives"
	"aurora-relayer-go-common/types/request"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
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
	wantHash := common.HexStringToHash("0xabcdf")
	wantFromBlock := common.IntToBN64(0xb1)
	wantToBlock := common.IntToBN64(0xb2)
	want := request.Filter{
		Addresses: []common.Address{common.HexStringToAddress("1"), common.HexStringToAddress("2")},
		FromBlock: &wantFromBlock,
		ToBlock:   &wantToBlock,
		BlockHash: &wantHash,
	}
	var result request.Filter
	err := json.Unmarshal([]byte(data), &result)
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

	ca1 := primitives.Data20FromHex("0x2")
	ca2 := primitives.Data20FromHex("0x1")
	ca3 := primitives.Data20FromHex("0x2")
	ca4 := primitives.Data20FromHex("0x1")
	ca5 := primitives.Data20FromHex("0x3")

	blockHash := primitives.Data32FromHex("0xa")
	filterHash := common.HexStringToHash(blockHash.Hex())
	var blockData = indexer.Block{
		ChainId: 1313161554,
		Height:  1,
		Hash:    blockHash,
		Transactions: []*indexer.Transaction{
			{ContractAddress: &ca1},
			{ContractAddress: &ca2},
			{ContractAddress: &ca3},
			{ContractAddress: &ca4},
			{ContractAddress: &ca5},
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
			wantFrom:    nil,
			wantTo:      nil,
			wantAddress: []primitives.Data20{},
			wantTopics:  [][]primitives.Data32{},
		},
		{
			name: "blockHash is added",
			data: request.Filter{
				BlockHash: &filterHash,
			},
			wantFrom:    nil,
			wantTo:      nil,
			wantAddress: []primitives.Data20{},
			wantTopics:  [][]primitives.Data32{},
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
			wantTopics:  [][]primitives.Data32{},
		},
		{
			name: "addresses get added once",
			data: request.Filter{
				Addresses: []common.Address{
					common.BytesToAddress(primitives.Data20FromHex("0x2").Bytes()),
					common.BytesToAddress(primitives.Data20FromHex("0x1").Bytes()),
					common.BytesToAddress(primitives.Data20FromHex("0x2").Bytes()),
					common.BytesToAddress(primitives.Data20FromHex("0x1").Bytes()),
					common.BytesToAddress(primitives.Data20FromHex("0x3").Bytes()),
				},
			},
			wantFrom: nil,
			wantTo:   nil,
			wantAddress: []primitives.Data20{
				primitives.Data20FromHex("0x2"),
				primitives.Data20FromHex("0x1"),
				primitives.Data20FromHex("0x3"),
			},
			wantTopics: [][]primitives.Data32{},
		},
		{
			name: "topics are added as is", // TODO: add stronger topics validation/restrict the type from []byte when unmarshalling?
			data: request.Filter{
				Topics: request.Topics{{primitives.Data32FromHex("0x1111").Bytes()}, {primitives.Data32FromHex("0x2222").Bytes()}, {primitives.Data32FromHex("0x3333").Bytes()}},
			},
			wantFrom:    nil,
			wantTo:      nil,
			wantAddress: []primitives.Data20{},
			wantTopics:  [][]primitives.Data32{{primitives.Data32FromHex("0x1111")}, {primitives.Data32FromHex("0x2222")}, {primitives.Data32FromHex("0x3333")}},
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
				{[]byte(`0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef`)},
				{},
				{[]byte("0xabc"), []byte("0x123")},
				{},
			},
		},
		{
			data: `["0x1","0x2","0x3","0x4","0x5"]`,
			want: request.Topics{
				{[]byte("0x1")},
				{[]byte("0x2")},
				{[]byte("0x3")},
				{[]byte("0x4")},
			},
		},
		{
			data:    `[`,
			want:    nil,
			wantErr: "unexpected end of JSON input",
		},
	}
	for _, tc := range ttable {
		t.Run(tc.data, func(t *testing.T) {
			var ts request.Topics
			err := json.Unmarshal([]byte(tc.data), &ts)
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tc.want, ts)
		})
	}
}
