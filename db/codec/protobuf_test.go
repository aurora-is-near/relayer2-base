package codec

import (
	"testing"

	"github.com/aurora-is-near/relayer2-base/types/db"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/stretchr/testify/require"
)

func roundTrip[T any](t *testing.T, c Codec, original *T) *T {
	t.Helper()
	data, err := c.Marshal(original)
	require.NoError(t, err, "marshal must succeed")
	result := new(T)
	err = c.Unmarshal(data, result)
	require.NoError(t, err, "unmarshal must succeed")
	return result
}

// TestBlockRoundTrip tests protobuf round-trip for Block with all fields populated.
func TestBlockRoundTrip(t *testing.T) {
	c := NewProtobufCodec()
	block := &db.Block{
		ParentHash:       primitives.Data32FromBytes([]byte{1, 2, 3}),
		Miner:            primitives.Data20FromBytes([]byte{4, 5, 6}),
		Timestamp:        1234567890,
		GasLimit:         primitives.QuantityFromUint64(8000000),
		GasUsed:          primitives.QuantityFromUint64(5000000),
		LogsBloom:        primitives.Data256FromBytes(make([]byte, 256)),
		TransactionsRoot: primitives.Data32FromBytes([]byte{7, 8, 9}),
		StateRoot:        primitives.Data32FromBytes([]byte{10, 11, 12}),
		ReceiptsRoot:     primitives.Data32FromBytes([]byte{13, 14, 15}),
		Size:             42,
	}
	result := roundTrip(t, c, block)
	require.Equal(t, block, result)
}

// TestBlockZeroValueRoundTrip tests round-trip with zero-value block fields.
func TestBlockZeroValueRoundTrip(t *testing.T) {
	c := NewProtobufCodec()
	// Zero-value block: all primitive Data fields have nil Content, uint64 fields are 0.
	// After round-trip, nil Content fields should remain as zero-value structs.
	block := &db.Block{}
	result := roundTrip(t, c, block)
	require.Equal(t, block, result)
}

// TestTransactionType0RoundTrip tests Type 0 transaction (legacy, no access list).
func TestTransactionType0RoundTrip(t *testing.T) {
	c := NewProtobufCodec()
	addr := primitives.Data20FromBytes([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})
	hash := primitives.Data32FromBytes([]byte{0xaa, 0xbb, 0xcc})
	tx := &db.Transaction{
		Type:                 0,
		From:                 primitives.Data20FromBytes([]byte{10, 20, 30}),
		IsContractDeployment: false,
		ToOrContract:         &addr,
		Nonce:                primitives.QuantityFromUint64(42),
		GasPrice:             primitives.QuantityFromUint64(20000000000),
		GasLimit:             primitives.QuantityFromUint64(21000),
		GasUsed:              21000,
		CumulativeGasUsed:    primitives.QuantityFromUint64(21000),
		Value:                primitives.QuantityFromUint64(1000000000000000000),
		Input:                primitives.VarDataFromBytes([]byte{0xde, 0xad}),
		NearHash:             &hash,
		NearReceiptHash:      primitives.Data32FromBytes([]byte{0xdd, 0xee}),
		Status:               true,
		V:                    27,
		R:                    primitives.QuantityFromBytes([]byte{0x11, 0x22}),
		S:                    primitives.QuantityFromBytes([]byte{0x33, 0x44}),
		LogsBloom:            primitives.Data256FromBytes(make([]byte, 256)),
	}
	result := roundTrip(t, c, tx)
	require.Equal(t, tx, result)
}

// TestTransactionType1RoundTrip tests Type 1 (EIP-2930) with access list.
func TestTransactionType1RoundTrip(t *testing.T) {
	c := NewProtobufCodec()
	addr := primitives.Data20FromBytes([]byte{1, 2, 3, 4, 5})
	tx := &db.Transaction{
		Type:         1,
		From:         primitives.Data20FromBytes([]byte{10, 20}),
		ToOrContract: &addr,
		Nonce:        primitives.QuantityFromUint64(1),
		GasPrice:     primitives.QuantityFromUint64(100),
		GasLimit:     primitives.QuantityFromUint64(21000),
		GasUsed:      21000,
		Value:        primitives.QuantityFromUint64(0),
		Input:        primitives.VarDataFromBytes(nil),
		Status:       true,
		V:            1,
		R:            primitives.QuantityFromBytes([]byte{0x55}),
		S:            primitives.QuantityFromBytes([]byte{0x66}),
		LogsBloom:    primitives.Data256FromBytes(make([]byte, 256)),
		AccessList: []db.AccessListEntry{
			{
				Address: primitives.Data20FromBytes([]byte{0xaa}),
				StorageKeys: []primitives.Data32{
					primitives.Data32FromBytes([]byte{0xbb}),
					primitives.Data32FromBytes([]byte{0xcc}),
				},
			},
			{
				Address:     primitives.Data20FromBytes([]byte{0xdd}),
				StorageKeys: []primitives.Data32{},
			},
		},
	}
	result := roundTrip(t, c, tx)
	require.Equal(t, tx, result)
}

// TestTransactionType2RoundTrip tests Type 2 (EIP-1559) with all fields.
func TestTransactionType2RoundTrip(t *testing.T) {
	c := NewProtobufCodec()
	addr := primitives.Data20FromBytes([]byte{0xff})
	tx := &db.Transaction{
		Type:                 2,
		From:                 primitives.Data20FromBytes([]byte{0xee}),
		ToOrContract:         &addr,
		Nonce:                primitives.QuantityFromUint64(5),
		GasPrice:             primitives.QuantityFromUint64(100),
		GasLimit:             primitives.QuantityFromUint64(21000),
		GasUsed:              21000,
		Value:                primitives.QuantityFromUint64(0),
		Input:                primitives.VarDataFromBytes([]byte{0x01, 0x02, 0x03}),
		Status:               true,
		V:                    0,
		R:                    primitives.QuantityFromBytes([]byte{0x77}),
		S:                    primitives.QuantityFromBytes([]byte{0x88}),
		LogsBloom:            primitives.Data256FromBytes(make([]byte, 256)),
		MaxPriorityFeePerGas: primitives.QuantityFromUint64(2000000000),
		MaxFeePerGas:         primitives.QuantityFromUint64(100000000000),
		// AccessList left nil: proto3 cannot distinguish empty from nil for repeated fields.
		// This is acceptable because consumers check data.Type >= 1 to decide whether to include AccessList.
	}
	result := roundTrip(t, c, tx)
	require.Equal(t, tx, result)
}

// TestTransactionNilPointerFields tests that nil ToOrContract and NearHash are preserved.
func TestTransactionNilPointerFields(t *testing.T) {
	c := NewProtobufCodec()
	tx := &db.Transaction{
		Type:     0,
		From:     primitives.Data20FromBytes([]byte{0xab}),
		Nonce:    primitives.QuantityFromUint64(0),
		GasPrice: primitives.QuantityFromUint64(0),
		GasLimit: primitives.QuantityFromUint64(0),
		Value:    primitives.QuantityFromUint64(0),
		Input:    primitives.VarDataFromBytes(nil),
		R:        primitives.QuantityFromBytes(nil),
		S:        primitives.QuantityFromBytes(nil),
		// ToOrContract and NearHash intentionally nil
	}
	result := roundTrip(t, c, tx)
	require.Nil(t, result.ToOrContract, "nil ToOrContract must be preserved")
	require.Nil(t, result.NearHash, "nil NearHash must be preserved")
	require.Equal(t, tx, result)
}

// TestTransactionZeroValuePointerFields tests that non-nil but zero-content pointers are preserved.
func TestTransactionZeroValuePointerFields(t *testing.T) {
	c := NewProtobufCodec()
	zeroAddr := primitives.Data20FromBytes(nil) // 20 zero bytes
	zeroHash := primitives.Data32FromBytes(nil) // 32 zero bytes
	tx := &db.Transaction{
		Type:         0,
		From:         primitives.Data20FromBytes([]byte{0xab}),
		ToOrContract: &zeroAddr,
		NearHash:     &zeroHash,
		Nonce:        primitives.QuantityFromUint64(0),
		GasPrice:     primitives.QuantityFromUint64(0),
		GasLimit:     primitives.QuantityFromUint64(0),
		Value:        primitives.QuantityFromUint64(0),
		Input:        primitives.VarDataFromBytes(nil),
		R:            primitives.QuantityFromBytes(nil),
		S:            primitives.QuantityFromBytes(nil),
	}
	result := roundTrip(t, c, tx)
	require.NotNil(t, result.ToOrContract, "non-nil ToOrContract must be preserved")
	require.NotNil(t, result.NearHash, "non-nil NearHash must be preserved")
	require.Equal(t, tx, result)
}

// TestTransactionEmptyAccessListBecomesNil tests that empty AccessList becomes nil via proto3.
// Proto3 cannot distinguish empty repeated fields from absent ones.
// This is acceptable because consumers use len()/range which work the same for nil and empty.
func TestTransactionEmptyAccessListBecomesNil(t *testing.T) {
	c := NewProtobufCodec()
	tx := &db.Transaction{
		Type:       1,
		From:       primitives.Data20FromBytes([]byte{0xab}),
		Nonce:      primitives.QuantityFromUint64(0),
		GasPrice:   primitives.QuantityFromUint64(0),
		GasLimit:   primitives.QuantityFromUint64(0),
		Value:      primitives.QuantityFromUint64(0),
		Input:      primitives.VarDataFromBytes(nil),
		R:          primitives.QuantityFromBytes(nil),
		S:          primitives.QuantityFromBytes(nil),
		AccessList: []db.AccessListEntry{},
	}
	result := roundTrip(t, c, tx)
	require.Nil(t, result.AccessList, "empty AccessList becomes nil via proto3")
	require.Empty(t, result.AccessList, "empty/nil AccessList has len 0")
}

// TestLogRoundTrip tests Log with topics.
func TestLogRoundTrip(t *testing.T) {
	c := NewProtobufCodec()
	log := &db.Log{
		Address: primitives.Data20FromBytes([]byte{0x01, 0x02}),
		Data:    primitives.VarDataFromBytes([]byte{0xaa, 0xbb, 0xcc, 0xdd}),
		Topics: []primitives.Data32{
			primitives.Data32FromBytes([]byte{0x11}),
			primitives.Data32FromBytes([]byte{0x22}),
			primitives.Data32FromBytes([]byte{0x33}),
		},
	}
	result := roundTrip(t, c, log)
	require.Equal(t, log, result)
}

// TestLogEmptyTopicsRoundTrip tests that empty Topics round-trips as empty (not nil).
// Log.Topics always gets allocated with make() to ensure consistent behavior in JSON serialization
// (nil serializes as null, empty as []).
func TestLogEmptyTopicsRoundTrip(t *testing.T) {
	c := NewProtobufCodec()
	log := &db.Log{
		Address: primitives.Data20FromBytes([]byte{0x01}),
		Data:    primitives.VarDataFromBytes([]byte{0xaa}),
		Topics:  []primitives.Data32{},
	}
	result := roundTrip(t, c, log)
	require.NotNil(t, result.Topics, "empty Topics must remain non-nil after round-trip")
	require.Empty(t, result.Topics, "empty Topics must remain empty after round-trip")
}

// TestLogEmptyDataRoundTrip tests Log with empty VarData.
func TestLogEmptyDataRoundTrip(t *testing.T) {
	c := NewProtobufCodec()
	log := &db.Log{
		Address: primitives.Data20FromBytes([]byte{0x01}),
		Data:    primitives.VarDataFromBytes(nil), // empty non-nil Content
		Topics:  []primitives.Data32{primitives.Data32FromBytes([]byte{0x42})},
	}
	result := roundTrip(t, c, log)
	require.Equal(t, log, result)
	require.NotNil(t, result.Data.Content, "VarData Content must be non-nil after round-trip")
}

// TestKeyRoundTrips tests all key types.
func TestKeyRoundTrips(t *testing.T) {
	c := NewProtobufCodec()

	t.Run("BlockKey", func(t *testing.T) {
		for _, height := range []uint64{0, 1, 100, 1<<32 - 1, 1<<64 - 1} {
			key := &db.BlockKey{Height: height}
			result := roundTrip(t, c, key)
			require.Equal(t, key, result)
		}
	})

	t.Run("TransactionKey", func(t *testing.T) {
		key := &db.TransactionKey{BlockHeight: 42, TransactionIndex: 7}
		result := roundTrip(t, c, key)
		require.Equal(t, key, result)
	})

	t.Run("TransactionKeyZero", func(t *testing.T) {
		key := &db.TransactionKey{}
		result := roundTrip(t, c, key)
		require.Equal(t, key, result)
	})

	t.Run("LogKey", func(t *testing.T) {
		key := &db.LogKey{BlockHeight: 100, TransactionIndex: 5, LogIndex: 3}
		result := roundTrip(t, c, key)
		require.Equal(t, key, result)
	})

	t.Run("LogKeyZero", func(t *testing.T) {
		key := &db.LogKey{}
		result := roundTrip(t, c, key)
		require.Equal(t, key, result)
	})
}

// TestBlockFilterRoundTrip tests BlockFilter protobuf round-trip.
func TestBlockFilterRoundTrip(t *testing.T) {
	c := NewProtobufCodec()
	filter := &db.BlockFilter{
		CreatedAt: 1234567890,
		Metadata:  primitives.VarDataFromBytes([]byte{0x01, 0x02}),
		From:      db.BlockKey{Height: 100},
		Next:      db.BlockKey{Height: 150},
		To:        db.BlockKey{Height: 200},
	}
	result := roundTrip(t, c, filter)
	require.Equal(t, filter, result)
}

// TestBlockFilterZeroKeysRoundTrip tests BlockFilter with zero-value keys.
func TestBlockFilterZeroKeysRoundTrip(t *testing.T) {
	c := NewProtobufCodec()
	filter := &db.BlockFilter{
		CreatedAt: 999,
		Metadata:  primitives.VarDataFromBytes(nil),
		From:      db.BlockKey{Height: 0},
		Next:      db.BlockKey{Height: 0},
		To:        db.BlockKey{Height: 0},
	}
	result := roundTrip(t, c, filter)
	require.Equal(t, filter, result)
}

// TestTransactionFilterRoundTrip tests TransactionFilter protobuf round-trip.
func TestTransactionFilterRoundTrip(t *testing.T) {
	c := NewProtobufCodec()
	filter := &db.TransactionFilter{
		CreatedAt: 9876543210,
		Metadata:  primitives.VarDataFromBytes([]byte{0xff}),
		From:      db.TransactionKey{BlockHeight: 50, TransactionIndex: 0},
		Next:      db.TransactionKey{BlockHeight: 75, TransactionIndex: 3},
		To:        db.TransactionKey{BlockHeight: 100, TransactionIndex: 10},
	}
	result := roundTrip(t, c, filter)
	require.Equal(t, filter, result)
}

// TestLogFilterRoundTrip tests LogFilter with addresses and topics.
func TestLogFilterRoundTrip(t *testing.T) {
	c := NewProtobufCodec()
	filter := &db.LogFilter{
		CreatedAt: 111222333,
		Metadata:  primitives.VarDataFromBytes(nil),
		From:      db.LogKey{BlockHeight: 10, TransactionIndex: 0, LogIndex: 0},
		Next:      db.LogKey{BlockHeight: 15, TransactionIndex: 2, LogIndex: 1},
		To:        db.LogKey{BlockHeight: 20, TransactionIndex: 5, LogIndex: 3},
		Addresses: []primitives.Data20{
			primitives.Data20FromBytes([]byte{0x01}),
			primitives.Data20FromBytes([]byte{0x02}),
		},
		Topics: [][]primitives.Data32{
			{
				primitives.Data32FromBytes([]byte{0xaa}),
				primitives.Data32FromBytes([]byte{0xbb}),
			},
			{
				primitives.Data32FromBytes([]byte{0xcc}),
			},
		},
	}
	result := roundTrip(t, c, filter)
	require.Equal(t, filter, result)
}

// TestLogFilterNilAddressesAndTopics tests LogFilter with nil Addresses and Topics.
func TestLogFilterNilAddressesAndTopics(t *testing.T) {
	c := NewProtobufCodec()
	filter := &db.LogFilter{
		CreatedAt: 555,
		Metadata:  primitives.VarDataFromBytes(nil),
		From:      db.LogKey{BlockHeight: 0},
		Next:      db.LogKey{},
		To:        db.LogKey{BlockHeight: 100},
		// Addresses and Topics intentionally nil
	}
	result := roundTrip(t, c, filter)
	require.Nil(t, result.Addresses, "nil Addresses must be preserved")
	require.Nil(t, result.Topics, "nil Topics must be preserved")
	require.Equal(t, filter, result)
}

// TestLogFilterEmptyAddressesAndTopics tests LogFilter with empty (non-nil) Addresses and Topics.
// Note: empty slices marshal to empty repeated fields, which proto3 deserializes as nil.
// The codec guards against this for Addresses and Topics using nil checks.
func TestLogFilterEmptyAddressesAndTopics(t *testing.T) {
	c := NewProtobufCodec()
	filter := &db.LogFilter{
		CreatedAt: 555,
		Metadata:  primitives.VarDataFromBytes(nil),
		From:      db.LogKey{},
		Next:      db.LogKey{},
		To:        db.LogKey{BlockHeight: 100},
		Addresses: []primitives.Data20{},
		Topics:    [][]primitives.Data32{},
	}
	data, err := c.Marshal(filter)
	require.NoError(t, err)

	result := new(db.LogFilter)
	err = c.Unmarshal(data, result)
	require.NoError(t, err)

	// Proto3 turns empty repeated fields into nil. This is expected.
	// The key behavior is that logFilterFromProto guards Addresses/Topics with nil checks,
	// so empty slices become nil on round-trip. This is acceptable because:
	// - nil and empty slices are functionally equivalent for filter matching (no filtering applied)
	// - The ToLogFilter() constructor in types/filter.go passes through the slices from Filter
	require.Nil(t, result.Addresses, "empty Addresses becomes nil via proto3")
	require.Nil(t, result.Topics, "empty Topics becomes nil via proto3")
}

// TestStandaloneDataTypes tests round-trip for standalone Data20, Data32, VarData.
func TestStandaloneDataTypes(t *testing.T) {
	c := NewProtobufCodec()

	t.Run("Data32", func(t *testing.T) {
		d := primitives.Data32FromBytes([]byte{0x01, 0x02, 0x03})
		result := roundTrip(t, c, &d)
		require.Equal(t, d, *result)
	})

	t.Run("Data20", func(t *testing.T) {
		d := primitives.Data20FromBytes([]byte{0x04, 0x05, 0x06})
		result := roundTrip(t, c, &d)
		require.Equal(t, d, *result)
	})

	t.Run("VarData", func(t *testing.T) {
		d := primitives.VarDataFromBytes([]byte{0x07, 0x08, 0x09})
		result := roundTrip(t, c, &d)
		require.Equal(t, d, *result)
	})

	t.Run("VarDataEmpty", func(t *testing.T) {
		d := primitives.VarDataFromBytes(nil) // Content = []byte{}
		result := roundTrip(t, c, &d)
		require.Equal(t, d, *result)
		require.NotNil(t, result.Content, "VarData empty Content must round-trip as non-nil")
	})

	t.Run("Data32Zero", func(t *testing.T) {
		d := primitives.Data32FromBytes(nil) // 32 zero bytes
		result := roundTrip(t, c, &d)
		require.Equal(t, d, *result)
	})
}

// TestFilterFromTypesFilter tests that filters created via types.Filter constructors
// round-trip correctly through protobuf.
func TestFilterFromTypesFilter(t *testing.T) {
	c := NewProtobufCodec()

	t.Run("LogFilterWithAddresses", func(t *testing.T) {
		fromBlock := uint64(100)
		toBlock := uint64(200)
		addr1 := primitives.Data20FromBytes([]byte{0x01})
		addr2 := primitives.Data20FromBytes([]byte{0x02})
		topic1 := primitives.Data32FromBytes([]byte{0xaa})
		topic2 := primitives.Data32FromBytes([]byte{0xbb})

		filter := &db.LogFilter{
			CreatedAt: 12345,
			Metadata:  primitives.VarDataFromBytes(nil),
			From:      db.LogKey{BlockHeight: fromBlock},
			To:        db.LogKey{BlockHeight: toBlock, TransactionIndex: 0xFFFF, LogIndex: 0xFFFF},
			Addresses: []primitives.Data20{addr1, addr2},
			Topics:    [][]primitives.Data32{{topic1, topic2}},
		}
		result := roundTrip(t, c, filter)
		require.Equal(t, filter, result)
	})

	t.Run("BlockFilter", func(t *testing.T) {
		filter := &db.BlockFilter{
			CreatedAt: 12345,
			Metadata:  primitives.VarDataFromBytes(nil),
			From:      db.BlockKey{Height: 100},
			To:        db.BlockKey{Height: 200},
		}
		result := roundTrip(t, c, filter)
		require.Equal(t, filter, result)
	})

	t.Run("TransactionFilter", func(t *testing.T) {
		filter := &db.TransactionFilter{
			CreatedAt: 12345,
			Metadata:  primitives.VarDataFromBytes(nil),
			From:      db.TransactionKey{BlockHeight: 100},
			To:        db.TransactionKey{BlockHeight: 200, TransactionIndex: 0xFFFF},
		}
		result := roundTrip(t, c, filter)
		require.Equal(t, filter, result)
	})
}

// TestUnsupportedType tests that marshaling/unmarshaling an unsupported type returns an error.
func TestUnsupportedType(t *testing.T) {
	c := NewProtobufCodec()

	type Unknown struct{}
	_, err := c.Marshal(&Unknown{})
	require.Error(t, err, "marshal of unknown type must fail")

	err = c.Unmarshal([]byte{}, &Unknown{})
	require.Error(t, err, "unmarshal of unknown type must fail")
}

// TestTransactionLargeAccessList tests round-trip with a large access list.
func TestTransactionLargeAccessList(t *testing.T) {
	c := NewProtobufCodec()

	accessList := make([]db.AccessListEntry, 50)
	for i := range accessList {
		keys := make([]primitives.Data32, 10)
		for j := range keys {
			keys[j] = primitives.Data32FromBytes([]byte{byte(i), byte(j)})
		}
		accessList[i] = db.AccessListEntry{
			Address:     primitives.Data20FromBytes([]byte{byte(i)}),
			StorageKeys: keys,
		}
	}

	tx := &db.Transaction{
		Type:       1,
		From:       primitives.Data20FromBytes([]byte{0xab}),
		Nonce:      primitives.QuantityFromUint64(0),
		GasPrice:   primitives.QuantityFromUint64(0),
		GasLimit:   primitives.QuantityFromUint64(0),
		Value:      primitives.QuantityFromUint64(0),
		Input:      primitives.VarDataFromBytes(nil),
		R:          primitives.QuantityFromBytes(nil),
		S:          primitives.QuantityFromBytes(nil),
		AccessList: accessList,
	}
	result := roundTrip(t, c, tx)
	require.Equal(t, tx, result)
}

// TestLogFilterManyTopicGroups tests LogFilter with multiple topic groups including empty ones.
func TestLogFilterManyTopicGroups(t *testing.T) {
	c := NewProtobufCodec()
	filter := &db.LogFilter{
		CreatedAt: 777,
		Metadata:  primitives.VarDataFromBytes(nil),
		From:      db.LogKey{},
		Next:      db.LogKey{},
		To:        db.LogKey{BlockHeight: 1000},
		Addresses: []primitives.Data20{
			primitives.Data20FromBytes([]byte{0x01}),
		},
		Topics: [][]primitives.Data32{
			{primitives.Data32FromBytes([]byte{0xaa})},
			{}, // empty topic group
			{primitives.Data32FromBytes([]byte{0xbb}), primitives.Data32FromBytes([]byte{0xcc})},
			{}, // another empty topic group
		},
	}
	result := roundTrip(t, c, filter)
	require.Equal(t, filter, result)
}

// TestVarDataNilVsEmpty verifies that VarData always has non-nil Content after round-trip,
// even when the original Content is nil (VarDataFromBytes always ensures non-nil).
func TestVarDataNilVsEmpty(t *testing.T) {
	c := NewProtobufCodec()

	// VarDataFromBytes(nil) produces Content = []byte{}
	d := primitives.VarDataFromBytes(nil)
	require.NotNil(t, d.Content, "VarDataFromBytes(nil) must produce non-nil Content")
	require.Empty(t, d.Content, "VarDataFromBytes(nil) must produce empty Content")

	result := roundTrip(t, c, &d)
	require.NotNil(t, result.Content, "round-tripped VarData must have non-nil Content")
	require.Empty(t, result.Content, "round-tripped VarData must have empty Content")
}

// TestTransactionType4RoundTrip tests Type 4 (EIP-7702 set-code) transaction with authorization list.
func TestTransactionType4RoundTrip(t *testing.T) {
	c := NewProtobufCodec()
	addr := primitives.Data20FromBytes([]byte{0xaa})
	tx := &db.Transaction{
		Type:                 4,
		From:                 primitives.Data20FromBytes([]byte{0xbb}),
		ToOrContract:         &addr,
		Nonce:                primitives.QuantityFromUint64(7),
		GasPrice:             primitives.QuantityFromUint64(100),
		GasLimit:             primitives.QuantityFromUint64(50000),
		GasUsed:              30000,
		Value:                primitives.QuantityFromUint64(0),
		Input:                primitives.VarDataFromBytes([]byte{0xde, 0xad}),
		Status:               true,
		V:                    1,
		R:                    primitives.QuantityFromBytes([]byte{0x99}),
		S:                    primitives.QuantityFromBytes([]byte{0xaa}),
		LogsBloom:            primitives.Data256FromBytes(make([]byte, 256)),
		MaxPriorityFeePerGas: primitives.QuantityFromUint64(1500000000),
		MaxFeePerGas:         primitives.QuantityFromUint64(50000000000),
		AuthorizationList: []db.Authorization{
			{
				ChainId: 1313161554,
				Address: primitives.Data20FromBytes([]byte{0xcc, 0xdd}),
				Nonce:   42,
				YParity: 1,
				R:       primitives.QuantityFromBytes([]byte{0x11, 0x22, 0x33}),
				S:       primitives.QuantityFromBytes([]byte{0x44, 0x55, 0x66}),
			},
			{
				ChainId: 1,
				Address: primitives.Data20FromBytes([]byte{0xee, 0xff}),
				Nonce:   0,
				YParity: 0,
				R:       primitives.QuantityFromBytes([]byte{0x77}),
				S:       primitives.QuantityFromBytes([]byte{0x88}),
			},
		},
	}
	result := roundTrip(t, c, tx)
	require.Equal(t, tx, result)
}

// TestTransactionType0BackwardCompat tests that Type 0 transactions with zero-value
// new fields still round-trip correctly (backward compatibility).
func TestTransactionType0BackwardCompat(t *testing.T) {
	c := NewProtobufCodec()
	addr := primitives.Data20FromBytes([]byte{1, 2, 3})
	tx := &db.Transaction{
		Type:         0,
		From:         primitives.Data20FromBytes([]byte{10, 20, 30}),
		ToOrContract: &addr,
		Nonce:        primitives.QuantityFromUint64(42),
		GasPrice:     primitives.QuantityFromUint64(20000000000),
		GasLimit:     primitives.QuantityFromUint64(21000),
		GasUsed:      21000,
		Value:        primitives.QuantityFromUint64(1000000000000000000),
		Input:        primitives.VarDataFromBytes(nil),
		Status:       true,
		V:            27,
		R:            primitives.QuantityFromBytes([]byte{0x11}),
		S:            primitives.QuantityFromBytes([]byte{0x22}),
		LogsBloom:    primitives.Data256FromBytes(make([]byte, 256)),
		// All new EIP-7702 fields are zero-value (default)
	}
	result := roundTrip(t, c, tx)
	require.Equal(t, tx, result)
	// Verify new fields remain zero-value
	require.Nil(t, result.AuthorizationList)
}

// TestTransactionYParityRoundTrip tests that the optional YParity field on Transaction
// round-trips correctly through protobuf for both nil and non-nil cases.
func TestTransactionYParityRoundTrip(t *testing.T) {
	c := NewProtobufCodec()
	addr := primitives.Data20FromBytes([]byte{0x01})

	t.Run("YParityNil", func(t *testing.T) {
		tx := &db.Transaction{
			Type:         1,
			From:         primitives.Data20FromBytes([]byte{0xab}),
			ToOrContract: &addr,
			Nonce:        primitives.QuantityFromUint64(1),
			GasPrice:     primitives.QuantityFromUint64(100),
			GasLimit:     primitives.QuantityFromUint64(21000),
			GasUsed:      21000,
			Value:        primitives.QuantityFromUint64(0),
			Input:        primitives.VarDataFromBytes(nil),
			Status:       true,
			V:            1,
			R:            primitives.QuantityFromBytes([]byte{0x55}),
			S:            primitives.QuantityFromBytes([]byte{0x66}),
			LogsBloom:    primitives.Data256FromBytes(make([]byte, 256)),
			// YParity intentionally nil
		}
		result := roundTrip(t, c, tx)
		require.Nil(t, result.YParity, "nil YParity must be preserved after round-trip")
		require.Equal(t, tx, result)
	})

	t.Run("YParityZero", func(t *testing.T) {
		yp := uint8(0)
		tx := &db.Transaction{
			Type:         2,
			From:         primitives.Data20FromBytes([]byte{0xcd}),
			ToOrContract: &addr,
			Nonce:        primitives.QuantityFromUint64(2),
			GasPrice:     primitives.QuantityFromUint64(200),
			GasLimit:     primitives.QuantityFromUint64(21000),
			GasUsed:      21000,
			Value:        primitives.QuantityFromUint64(0),
			Input:        primitives.VarDataFromBytes(nil),
			Status:       true,
			V:            0,
			R:            primitives.QuantityFromBytes([]byte{0x77}),
			S:            primitives.QuantityFromBytes([]byte{0x88}),
			LogsBloom:    primitives.Data256FromBytes(make([]byte, 256)),
			YParity:      &yp,
		}
		result := roundTrip(t, c, tx)
		require.NotNil(t, result.YParity, "non-nil YParity=0 must be preserved after round-trip")
		require.Equal(t, uint8(0), *result.YParity)
		require.Equal(t, tx, result)
	})

	t.Run("YParityOne", func(t *testing.T) {
		yp := uint8(1)
		tx := &db.Transaction{
			Type:         2,
			From:         primitives.Data20FromBytes([]byte{0xef}),
			ToOrContract: &addr,
			Nonce:        primitives.QuantityFromUint64(3),
			GasPrice:     primitives.QuantityFromUint64(300),
			GasLimit:     primitives.QuantityFromUint64(21000),
			GasUsed:      21000,
			Value:        primitives.QuantityFromUint64(0),
			Input:        primitives.VarDataFromBytes(nil),
			Status:       true,
			V:            1,
			R:            primitives.QuantityFromBytes([]byte{0x99}),
			S:            primitives.QuantityFromBytes([]byte{0xaa}),
			LogsBloom:    primitives.Data256FromBytes(make([]byte, 256)),
			YParity:      &yp,
		}
		result := roundTrip(t, c, tx)
		require.NotNil(t, result.YParity, "non-nil YParity=1 must be preserved after round-trip")
		require.Equal(t, uint8(1), *result.YParity)
		require.Equal(t, tx, result)
	})
}

// TestUnmarshalMalformedData tests that Unmarshal returns errors for corrupted data.
func TestUnmarshalMalformedData(t *testing.T) {
	c := NewProtobufCodec()
	garbage := []byte{0xff, 0xff, 0xff, 0xff}

	t.Run("block", func(t *testing.T) {
		var b db.Block
		require.Error(t, c.Unmarshal(garbage, &b))
	})
	t.Run("transaction", func(t *testing.T) {
		var tx db.Transaction
		require.Error(t, c.Unmarshal(garbage, &tx))
	})
	t.Run("log", func(t *testing.T) {
		var l db.Log
		require.Error(t, c.Unmarshal(garbage, &l))
	})
	t.Run("blockKey", func(t *testing.T) {
		var bk db.BlockKey
		require.Error(t, c.Unmarshal(garbage, &bk))
	})
	t.Run("transactionKey", func(t *testing.T) {
		var tk db.TransactionKey
		require.Error(t, c.Unmarshal(garbage, &tk))
	})
	t.Run("logKey", func(t *testing.T) {
		var lk db.LogKey
		require.Error(t, c.Unmarshal(garbage, &lk))
	})
	t.Run("blockFilter", func(t *testing.T) {
		var bf db.BlockFilter
		require.Error(t, c.Unmarshal(garbage, &bf))
	})
	t.Run("transactionFilter", func(t *testing.T) {
		var tf db.TransactionFilter
		require.Error(t, c.Unmarshal(garbage, &tf))
	})
	t.Run("logFilter", func(t *testing.T) {
		var lf db.LogFilter
		require.Error(t, c.Unmarshal(garbage, &lf))
	})
	t.Run("data32", func(t *testing.T) {
		var d primitives.Data32
		require.Error(t, c.Unmarshal(garbage, &d))
	})
	t.Run("data20", func(t *testing.T) {
		var d primitives.Data20
		require.Error(t, c.Unmarshal(garbage, &d))
	})
	t.Run("varData", func(t *testing.T) {
		var d primitives.VarData
		require.Error(t, c.Unmarshal(garbage, &d))
	})
}
