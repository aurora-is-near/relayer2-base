package endpoint

import (
	"aurora-relayer-go-common/db"
	"aurora-relayer-go-common/db/badger"
	"aurora-relayer-go-common/utils"
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

func TestLogFilterUnmarshalJSON(t *testing.T) {
	data := `{"address":["` +
		fmt.Sprintf("0x%040x", 0x1) +
		`","` +
		fmt.Sprintf("0x%040x", 0x2) +
		`"],"fromBlock":"0xb1","toBlock":"0xb2","blockHash":"` +
		fmt.Sprintf("0x%064x", 0xabcdf) +
		`"}`
	wantHash := utils.HexStringToHash("0xabcdf")
	want := utils.FilterOptions{
		Address:   utils.Addresses{utils.HexStringToAddress("1"), utils.HexStringToAddress("2")},
		FromBlock: "0x" + utils.IntToUint256(0xb1).Text(16),
		ToBlock:   "0x" + utils.IntToUint256(0xb2).Text(16),
		BlockHash: &wantHash,
	}

	var result utils.FilterOptions
	err := json.Unmarshal([]byte(data), &result)
	assert.Nil(t, err)
	assert.Equal(t, want, result)
}

func TestFormatFilterOptions(t *testing.T) {

	var eth *Eth

	viper.SetConfigType("yml")
	err := viper.ReadConfig(strings.NewReader(ethTestYaml))
	if err != nil {
		panic(err)
	}

	var blockData = utils.Block{
		Sequence: utils.IntToUint256(1),
		Hash:     utils.HexStringToHash("a"),
		Transactions: []*utils.Transaction{
			{ContractAddress: utils.HexStringToAddress("0x2")},
			{ContractAddress: utils.HexStringToAddress("0x1")},
			{ContractAddress: utils.HexStringToAddress("0x2")},
			{ContractAddress: utils.HexStringToAddress("0x1")},
			{ContractAddress: utils.HexStringToAddress("0x3")},
		}}

	ttable := []struct {
		name             string
		data             utils.FilterOptions
		wantFrom, wantTo utils.Uint256
		wantAddress      [][]byte
		wantTopics       [][][]byte
		wantErr          string
	}{
		{
			name:     "empty options",
			data:     utils.FilterOptions{},
			wantFrom: blockData.Sequence,
			wantTo:   blockData.Sequence,
		},
		{
			name: "blockHash is added",
			data: utils.FilterOptions{
				BlockHash: &blockData.Hash,
			},
			wantFrom: blockData.Sequence,
			wantTo:   blockData.Sequence,
		},
		{
			name: "block range is not overwritten",
			data: utils.FilterOptions{
				FromBlock: utils.IntToHex(10),
				ToBlock:   utils.IntToHex(20),
			},
			wantFrom: utils.IntToUint256(10),
			wantTo:   utils.IntToUint256(20),
		},
		{
			name: "addresses get added once",
			data: utils.FilterOptions{
				Address: utils.Addresses{
					utils.HexStringToAddress("0x2"),
					utils.HexStringToAddress("0x1"),
					utils.HexStringToAddress("0x2"),
					utils.HexStringToAddress("0x1"),
					utils.HexStringToAddress("0x3"),
				},
			},
			wantFrom: blockData.Sequence,
			wantTo:   blockData.Sequence,
			wantAddress: [][]byte{
				utils.HexStringToAddress("0x2").Bytes(),
				utils.HexStringToAddress("0x1").Bytes(),
				utils.HexStringToAddress("0x3").Bytes(),
			},
		},
		{
			name: "topics are added as is", // TODO: add stronger topics validation/restrict the type from []byte when unmarshalling?
			data: utils.FilterOptions{
				Topics: utils.Topics{{[]byte("foo")}, {[]byte("bar")}, {[]byte("bazz")}},
			},
			wantFrom:   blockData.Sequence,
			wantTo:     blockData.Sequence,
			wantTopics: [][][]byte{{[]byte("foo")}, {[]byte("bar")}, {[]byte("bazz")}},
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

			want := &utils.LogFilter{
				Address:   tc.wantAddress,
				FromBlock: &tc.wantFrom,
				ToBlock:   &tc.wantTo,
				Topics:    tc.wantTopics,
			}
			got, err := eth.formatFilterOptions(&tc.data)
			assert.Nil(t, err)
			assert.Equal(t, want, got)
		})
	}
}

func TestTopicsUnmarshalJSON(t *testing.T) {
	ttable := []struct {
		data    string
		want    utils.Topics
		wantErr string
	}{
		{
			data: `["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",null,["0xabc","0x123"]]`,
			want: utils.Topics{
				{[]byte(`0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef`)},
				{},
				{[]byte("0xabc"), []byte("0x123")},
				{},
			},
		},
		{
			data: `["0x1","0x2","0x3","0x4","0x5"]`,
			want: utils.Topics{
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
			var ts utils.Topics
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
