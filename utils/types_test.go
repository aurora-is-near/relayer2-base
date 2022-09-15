package utils_test

import (
	"aurora-relayer-go-common/utils"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"math"
	"math/big"
	"testing"
	"time"
)

func TestUint256Cbor(t *testing.T) {
	before := utils.IntToUint256(1337)
	bs, err := cbor.Marshal(&before)
	assert.Nil(t, err)

	var after utils.Uint256
	err = cbor.Unmarshal(bs, &after)
	assert.Nil(t, err)
	assert.Equal(t, before, after)
}

func TestBlockCbor(t *testing.T) {
	before := utils.Block{
		ChainId:          math.MaxUint64,
		Hash:             randomHash(),
		ParentHash:       randomHash(),
		Height:           math.MaxUint64,
		Miner:            randomAddress(),
		Timestamp:        time.Now().UnixNano(),
		GasLimit:         randomUint256(),
		GasUsed:          randomUint256(),
		LogsBloom:        "",
		TransactionsRoot: randomHash(),
		ReceiptsRoot:     randomHash(),
		Transactions: []*utils.Transaction{
			{
				Hash:                 randomHash(),
				BlockHash:            randomHash(),
				BlockHeight:          math.MaxUint64,
				ChainId:              math.MaxUint64,
				TransactionIndex:     math.MaxUint32,
				From:                 randomAddress(),
				To:                   nil,
				Nonce:                randomUint256(),
				GasPrice:             randomUint256(),
				GasLimit:             randomUint256(),
				GasUsed:              randomUint256(),
				MaxPriorityFeePerGas: randomUint256(),
				MaxFeePerGas:         randomUint256(),
				Value:                randomUint256(),
				Input:                randomBytes(10),
				Output:               randomBytes(10),
				AccessList:           []utils.AccessList{},
				TxType:               math.MaxUint8,
				Status:               true,
				Logs: []*utils.Log{
					{
						Address: randomAddress(),
						Topics: []utils.Bytea{
							randomBytea(),
						},
						Data: randomBytes(10),
					},
				},
				ContractAddress: randomAddress(),
				V:               math.MaxUint64,
				R:               randomUint256(),
				S:               randomUint256(),
				NearTransaction: utils.NearTransaction{
					Hash:        randomHash(),
					ReceiptHash: randomHash(),
				},
			},
		},
		NearBlock: map[interface{}]interface{}{"foo": "bar"},
		StateRoot: randomHash().String(),
		Size:      randomUint256(),
		Sequence:  randomUint256().Uint64(),
	}
	bs, err := cbor.Marshal(&before)
	assert.Nil(t, err)

	blank := utils.Block{}

	var after utils.Block
	err = cbor.Unmarshal(bs, &after)
	after.NearBlock = map[interface{}]interface{}{"foo": "bar"}
	assert.Nil(t, err)
	assert.Equal(t, before, after, "data survives CBOR marshaling")
	assert.NotEqual(t, blank, after, "unmarshalled data doesn't equal default")
}

func TestTransactionCbor(t *testing.T) {
	before := utils.Transaction{
		Hash:                 randomHash(),
		BlockHash:            randomHash(),
		BlockHeight:          math.MaxUint64,
		ChainId:              math.MaxUint64,
		TransactionIndex:     math.MaxUint32,
		From:                 randomAddress(),
		To:                   nil,
		Nonce:                randomUint256(),
		GasPrice:             randomUint256(),
		GasLimit:             randomUint256(),
		GasUsed:              randomUint256(),
		MaxPriorityFeePerGas: randomUint256(),
		MaxFeePerGas:         randomUint256(),
		Value:                randomUint256(),
		Input:                randomBytes(10),
		Output:               randomBytes(10),
		AccessList:           []utils.AccessList{},
		TxType:               math.MaxUint8,
		Status:               true,
		Logs:                 nil,
		ContractAddress:      randomAddress(),
		V:                    math.MaxUint64,
		R:                    randomUint256(),
		S:                    randomUint256(),
		NearTransaction: utils.NearTransaction{
			Hash:        randomHash(),
			ReceiptHash: randomHash(),
		},
	}
	bs, err := cbor.Marshal(&before)
	assert.Nil(t, err)

	blank := utils.Transaction{}

	var after utils.Transaction
	err = cbor.Unmarshal(bs, &after)
	assert.Nil(t, err)
	assert.Equal(t, before, after, "data survives CBOR marshaling")
	assert.NotEqual(t, blank, after, "unmarshalled data doesn't equal default")
}

func TestLogResponseCbor(t *testing.T) {
	before := utils.LogResponse{ // we save LogResponse{} instead of Log{} to DB
		Removed:          false,
		LogIndex:         randomUint256(),
		TransactionIndex: randomUint256(),
		TransactionHash:  randomHash(),
		BlockHash:        randomHash(),
		BlockNumber:      randomUint256(),
		Address:          randomAddress(),
		Data:             randomBytes(10),
		Topics:           []utils.Bytea{randomBytea()},
	}
	bs, err := cbor.Marshal(&before)
	assert.Nil(t, err)

	blank := utils.LogResponse{}

	var after utils.LogResponse
	err = cbor.Unmarshal(bs, &after)
	assert.Nil(t, err)
	assert.Equal(t, before, after, "data survives CBOR marshaling")
	assert.NotEqual(t, blank, after, "unmarshalled data doesn't equal default")
}

func TestUint256UnmarshalJSON(t *testing.T) {
	ttable := []struct {
		data    string
		want    utils.Uint256
		wantErr string
	}{
		{"1", utils.IntToUint256(1), ""},
		{"0x1", utils.IntToUint256(1), ""},
		{"0x01", utils.IntToUint256(1), ""},
		{"10", utils.IntToUint256(16), ""},
		{"0x10", utils.IntToUint256(16), ""},
		{"0x100", utils.IntToUint256(256), ""},
		{"0xffff", utils.IntToUint256(65535), ""},
		{"0xFFFF", utils.IntToUint256(65535), ""},
		{"0x", utils.Uint256{}, "failed to parse"},
		{"0xG", utils.Uint256{}, "failed to parse"},
		{"0xZ", utils.Uint256{}, "failed to parse"},
	}
	for _, tc := range ttable {
		t.Run(tc.data, func(t *testing.T) {
			b := []byte(`"` + tc.data + `"`)
			var parsed utils.Uint256
			err := json.Unmarshal(b, &parsed)
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.want, parsed)
			}
		})
	}
}

func TestUint256MarshalJSON(t *testing.T) {
	ttable := []struct {
		data *big.Int
		want string
	}{
		{nil, `"0x0"`},
		{big.NewInt(0), `"0x0"`},
		{big.NewInt(1), `"0x1"`},
		{big.NewInt(10), `"0xa"`},
		{big.NewInt(16), `"0x10"`},
		{big.NewInt(65535), `"0xffff"`},
		{big.NewInt(65536), `"0x10000"`},
	}
	for _, tc := range ttable {
		t.Run(tc.want, func(t *testing.T) {
			value := utils.Uint256{Int: tc.data}
			res, err := json.Marshal(&value)
			assert.Nil(t, err)
			assert.Equal(t, tc.want, string(res))
		})
	}
}

var addressParseTests = []struct {
	data    string
	want    utils.Address
	wantErr string
}{
	{"0x0", utils.HexStringToAddress("0x0"), "hex string of odd length"},
	{"0x00", utils.HexStringToAddress("0x0"), "hex string has length 2, want 40"},
	{fmt.Sprintf("0x%040x", 0), utils.HexStringToAddress("0x0"), ""},
	{fmt.Sprintf("0x%040x", 0x54d), utils.HexStringToAddress("0x54d"), ""},
}

func TestAddressUnmarshalJSON(t *testing.T) {
	for _, tc := range addressParseTests {
		t.Run(tc.data, func(t *testing.T) {
			var parsed utils.Address
			err := json.Unmarshal([]byte(`"`+tc.data+`"`), &parsed)
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.Nil(t, err)
				want := tc.want.Address.Hash().Big()
				got := parsed.Hash().Big()
				assert.Equal(t, got, want)
			}
		})
	}
}

func TestUint256ToUint32Key(t *testing.T) {
	ttable := []struct {
		name    string
		data    uint64
		want    []byte
		wantErr string
	}{
		{
			name: "zero",
			data: 0,
			want: []byte{0, 0, 0, 0},
		},
		{
			name: "one",
			data: 1,
			want: []byte{0, 0, 0, 0x1},
		},
		{
			name: "large number",
			data: 1234567890,
			want: []byte{0x49, 0x96, 0x2, 0xd2},
		},
		{
			name: "max number",
			data: math.MaxUint32,
			want: []byte{0xff, 0xff, 0xff, 0xff},
		},
		{
			name:    "max number plus one",
			data:    math.MaxUint32 + 1,
			wantErr: "u256 doesn't fit in a u32",
		},
	}
	for _, tc := range ttable {
		t.Run(tc.name, func(t *testing.T) {
			bi := utils.UintToUint256(tc.data)
			key, err := bi.ToUint32Key()
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
				assert.Nil(t, key)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.want, key.KeyBytes())
			}
		})
	}
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func randomHash() utils.H256 {
	return utils.H256{Hash: common.BigToHash(big.NewInt(0).SetBytes(randomBytes(32)))}
}

func randomUint256() utils.Uint256 {
	return utils.Uint256{Int: big.NewInt(0).SetBytes(randomBytes(32))}
}

func randomAddress() utils.Address {
	return utils.Address{Address: common.BigToAddress(big.NewInt(0).SetBytes(randomBytes(20)))}
}

func randomBytea() utils.Bytea {
	return utils.Bytea(randomBytes(10))
}
