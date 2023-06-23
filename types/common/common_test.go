package common

import (
	"math"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func TestBlockNumberOrHash_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		mustFail bool
		expected BlockNumberOrHash
	}{

		0:  {`"0x"`, true, BlockNumberOrHash{}},
		1:  {`"0x0"`, false, BlockNumberOrHashWithBN64(0)},
		2:  {`"0X1"`, false, BlockNumberOrHashWithBN64(1)},
		3:  {`"0x00"`, false, BlockNumberOrHashWithBN64(0)},
		4:  {`"0x01"`, false, BlockNumberOrHashWithBN64(1)},
		5:  {`"0x1"`, false, BlockNumberOrHashWithBN64(1)},
		6:  {`"0x12"`, false, BlockNumberOrHashWithBN64(18)},
		7:  {`"0x7fffffffffffffff"`, false, BlockNumberOrHashWithBN64(math.MaxInt64)},
		8:  {`"0x8000000000000000"`, true, BlockNumberOrHash{}},
		9:  {"0", false, BlockNumberOrHashWithBN64(0)},
		10: {`"ff"`, false, BlockNumberOrHashWithBN64(255)},
		11: {`"pending"`, false, BlockNumberOrHashWithBN64(PendingBlockNumber)},
		12: {`"latest"`, false, BlockNumberOrHashWithBN64(LatestBlockNumber)},
		13: {`"earliest"`, false, BlockNumberOrHashWithBN64(EarliestBlockNumber)},
		14: {`someString`, true, BlockNumberOrHash{}},
		15: {`""`, true, BlockNumberOrHash{}},
		16: {``, true, BlockNumberOrHash{}},
		17: {`{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000"}`, false, BlockNumberOrHashWithHash(HexStringToHash("0x0000000000000000000000000000000000000000000000000000000000000000"), false)},
		18: {`{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","requireCanonical":false}`, false, BlockNumberOrHashWithHash(HexStringToHash("0x0000000000000000000000000000000000000000000000000000000000000000"), false)},
		19: {`{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","requireCanonical":true}`, false, BlockNumberOrHashWithHash(HexStringToHash("0x0000000000000000000000000000000000000000000000000000000000000000"), true)},
		20: {`{"blockNumber":"0x1"}`, false, BlockNumberOrHashWithBN64(1)},
		21: {`{"blockNumber":"pending"}`, false, BlockNumberOrHashWithBN64(PendingBlockNumber)},
		22: {`{"blockNumber":"latest"}`, false, BlockNumberOrHashWithBN64(LatestBlockNumber)},
		23: {`{"blockNumber":"earliest"}`, false, BlockNumberOrHashWithBN64(EarliestBlockNumber)},
		24: {`{"blockNumber":"0x1", "blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000"}`, true, BlockNumberOrHash{}},
	}

	for i, test := range tests {
		var bnh BlockNumberOrHash
		err := jsoniter.Unmarshal([]byte(test.input), &bnh)
		if test.mustFail && err == nil {
			t.Errorf("Test %d should fail", i)
			continue
		}
		if !test.mustFail && err != nil {
			t.Errorf("Test %d should pass but got err: %v", i, err)
			continue
		}
		hash, hashOk := bnh.Hash()
		expectedHash, expectedHashOk := test.expected.Hash()
		num, numOk := bnh.Number()
		expectedNum, expectedNumOk := test.expected.Number()
		if bnh.RequireCanonical != test.expected.RequireCanonical ||
			hash.String() != expectedHash.String() || hashOk != expectedHashOk ||
			num != expectedNum || numOk != expectedNumOk {
			t.Errorf("Test %d got unexpected value, want %v, got %v", i, test.expected, bnh)
		}
	}
}
