package core

import (
	"testing"

	dbt "github.com/aurora-is-near/relayer2-base/types/db"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/aurora-is-near/relayer2-base/types/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func baseTx() *dbt.Transaction {
	addr := genAddress(100)
	nearHash := genHash(101)
	return &dbt.Transaction{
		Type:              0,
		From:              genAddress(102),
		ToOrContract:      &addr,
		Nonce:             genQuantity(103),
		GasPrice:          genQuantity(104),
		GasLimit:          primitives.QuantityFromUint64(21000),
		GasUsed:           21000,
		CumulativeGasUsed: primitives.QuantityFromUint64(21000),
		Value:             genQuantity(105),
		Input:             genVarData(0, 10, 106),
		NearHash:          &nearHash,
		NearReceiptHash:   genHash(107),
		Status:            true,
		V:                 27,
		R:                 genQuantity(108),
		S:                 genQuantity(109),
		LogsBloom:         genLogsBloom(110),
	}
}

func TestMakeTransactionResponse_Type2_ChainIDAndFees(t *testing.T) {
	tx := baseTx()
	tx.Type = 2
	tx.V = 1
	tx.AccessList = []dbt.AccessListEntry{}
	tx.MaxPriorityFeePerGas = genQuantity(200)
	tx.MaxFeePerGas = genQuantity(201)

	resp := makeTransactionResponse(testChainId, 100, 0, genHash(10), genHash(11), tx)
	require.NotNil(t, resp.ChainID)
	assert.Equal(t, primitives.HexUint(testChainId), *resp.ChainID)
	require.NotNil(t, resp.MaxPriorityFeePerGas)
	require.NotNil(t, resp.MaxFeePerGas)
	require.NotNil(t, resp.YParity, "Type 2 should have YParity")
	assert.Equal(t, primitives.HexUint(1), *resp.YParity)
	assert.Nil(t, resp.AuthorizationList, "Type 2 should not have AuthorizationList")
	assert.Equal(t, primitives.HexUint(2), resp.Type)
}

func TestMakeTransactionResponse_Type4_SetCode(t *testing.T) {
	tx := baseTx()
	tx.Type = 4
	tx.V = 1
	tx.AccessList = []dbt.AccessListEntry{}
	tx.MaxPriorityFeePerGas = genQuantity(400)
	tx.MaxFeePerGas = genQuantity(401)
	tx.AuthorizationList = []dbt.Authorization{
		{
			ChainId: 1,
			Address: genAddress(402),
			Nonce:   42,
			YParity: 1,
			R:       genQuantity(403),
			S:       genQuantity(404),
		},
		{
			ChainId: 1313161554,
			Address: genAddress(405),
			Nonce:   99,
			YParity: 0,
			R:       genQuantity(406),
			S:       genQuantity(407),
		},
	}

	resp := makeTransactionResponse(testChainId, 100, 0, genHash(10), genHash(11), tx)

	// Type 4 inherits type 2 fields
	require.NotNil(t, resp.ChainID, "Type 4 should have ChainID")
	assert.Equal(t, primitives.HexUint(testChainId), *resp.ChainID)
	require.NotNil(t, resp.MaxPriorityFeePerGas, "Type 4 should have MaxPriorityFeePerGas")
	require.NotNil(t, resp.MaxFeePerGas, "Type 4 should have MaxFeePerGas")
	require.NotNil(t, resp.AccessList, "Type 4 should have AccessList")
	require.NotNil(t, resp.YParity, "Type 4 should have YParity")
	assert.Equal(t, primitives.HexUint(1), *resp.YParity)

	// Type 4 specific: AuthorizationList
	require.Len(t, resp.AuthorizationList, 2)
	assert.Equal(t, primitives.HexUint(1), resp.AuthorizationList[0].ChainID)
	assert.Equal(t, genAddress(402), resp.AuthorizationList[0].Address)
	assert.Equal(t, primitives.HexUint(42), resp.AuthorizationList[0].Nonce)
	assert.Equal(t, primitives.HexUint(1), resp.AuthorizationList[0].YParity)
	assert.Equal(t, genQuantity(403), resp.AuthorizationList[0].R)
	assert.Equal(t, genQuantity(404), resp.AuthorizationList[0].S)

	assert.Equal(t, primitives.HexUint(1313161554), resp.AuthorizationList[1].ChainID)
	assert.Equal(t, genAddress(405), resp.AuthorizationList[1].Address)
	assert.Equal(t, primitives.HexUint(99), resp.AuthorizationList[1].Nonce)
	assert.Equal(t, primitives.HexUint(0), resp.AuthorizationList[1].YParity)
	assert.Equal(t, genQuantity(406), resp.AuthorizationList[1].R)
	assert.Equal(t, genQuantity(407), resp.AuthorizationList[1].S)

	assert.Equal(t, primitives.HexUint(4), resp.Type, "Type should be preserved as 4")
}

func TestMakeTransactionResponse_UnknownType_FallsBackToZero(t *testing.T) {
	tx := baseTx()
	tx.Type = 5
	resp := makeTransactionResponse(testChainId, 100, 0, genHash(10), genHash(11), tx)
	assert.Equal(t, primitives.HexUint(0), resp.Type, "Unknown type > 4 should fall back to 0")
	// Unknown types are treated as legacy (type 0): no type-conditional fields are set
	assert.Nil(t, resp.AccessList, "Unknown type should not have AccessList")
	assert.Nil(t, resp.ChainID, "Unknown type should not have ChainID")
	assert.Nil(t, resp.YParity, "Unknown type should not have YParity")
	assert.Nil(t, resp.MaxPriorityFeePerGas, "Unknown type should not have MaxPriorityFeePerGas")
	assert.Nil(t, resp.MaxFeePerGas, "Unknown type should not have MaxFeePerGas")
	assert.Nil(t, resp.AuthorizationList, "Unknown type should not have AuthorizationList")
}

func TestMakeTransactionResponse_Type0and1_Unchanged(t *testing.T) {
	t.Run("type 0", func(t *testing.T) {
		tx := baseTx()
		tx.Type = 0
		resp := makeTransactionResponse(testChainId, 100, 0, genHash(10), genHash(11), tx)
		assert.Equal(t, primitives.HexUint(0), resp.Type)
		assert.Nil(t, resp.AccessList, "Type 0 should not have AccessList")
		assert.Nil(t, resp.ChainID, "Type 0 should not have ChainID")
		assert.Nil(t, resp.YParity, "Type 0 should not have YParity")
		assert.Nil(t, resp.AuthorizationList)
	})

	t.Run("type 1", func(t *testing.T) {
		tx := baseTx()
		tx.Type = 1
		tx.V = 0
		tx.AccessList = []dbt.AccessListEntry{
			{Address: genAddress(500), StorageKeys: []primitives.Data32{genHash(501)}},
		}
		resp := makeTransactionResponse(testChainId, 100, 0, genHash(10), genHash(11), tx)
		assert.Equal(t, primitives.HexUint(1), resp.Type)
		require.NotNil(t, resp.AccessList, "Type 1 should have AccessList")
		require.NotNil(t, resp.ChainID, "Type 1 should have ChainID per EIP-2930")
		assert.Equal(t, primitives.HexUint(testChainId), *resp.ChainID)
		require.NotNil(t, resp.YParity, "Type 1 should have YParity")
		assert.Equal(t, primitives.HexUint(0), *resp.YParity)
		assert.Nil(t, resp.MaxPriorityFeePerGas, "Type 1 should not have MaxPriorityFeePerGas")
		assert.Nil(t, resp.MaxFeePerGas, "Type 1 should not have MaxFeePerGas")
		assert.Nil(t, resp.AuthorizationList)
	})
}

func TestMakeTransactionResponse_YParityDerivation(t *testing.T) {
	t.Run("explicit YParity present is used as-is", func(t *testing.T) {
		tx := baseTx()
		tx.Type = 2
		tx.AccessList = []dbt.AccessListEntry{}
		tx.MaxPriorityFeePerGas = genQuantity(200)
		tx.MaxFeePerGas = genQuantity(201)
		tx.V = 27 // legacy-range V, should be ignored
		yp := uint8(1)
		tx.YParity = &yp

		resp := makeTransactionResponse(testChainId, 100, 0, genHash(10), genHash(11), tx)
		require.NotNil(t, resp.YParity)
		assert.Equal(t, primitives.HexUint(1), *resp.YParity)
	})

	t.Run("YParity nil V=0 yields yParity=0", func(t *testing.T) {
		tx := baseTx()
		tx.Type = 2
		tx.AccessList = []dbt.AccessListEntry{}
		tx.MaxPriorityFeePerGas = genQuantity(200)
		tx.MaxFeePerGas = genQuantity(201)
		tx.V = 0
		tx.YParity = nil

		resp := makeTransactionResponse(testChainId, 100, 0, genHash(10), genHash(11), tx)
		require.NotNil(t, resp.YParity)
		assert.Equal(t, primitives.HexUint(0), *resp.YParity)
	})

	t.Run("YParity nil V=1 yields yParity=1", func(t *testing.T) {
		tx := baseTx()
		tx.Type = 2
		tx.AccessList = []dbt.AccessListEntry{}
		tx.MaxPriorityFeePerGas = genQuantity(200)
		tx.MaxFeePerGas = genQuantity(201)
		tx.V = 1
		tx.YParity = nil

		resp := makeTransactionResponse(testChainId, 100, 0, genHash(10), genHash(11), tx)
		require.NotNil(t, resp.YParity)
		assert.Equal(t, primitives.HexUint(1), *resp.YParity)
	})

	t.Run("YParity nil V=27 yields yParity=null", func(t *testing.T) {
		tx := baseTx()
		tx.Type = 2
		tx.AccessList = []dbt.AccessListEntry{}
		tx.MaxPriorityFeePerGas = genQuantity(200)
		tx.MaxFeePerGas = genQuantity(201)
		tx.V = 27
		tx.YParity = nil

		resp := makeTransactionResponse(testChainId, 100, 0, genHash(10), genHash(11), tx)
		assert.Nil(t, resp.YParity)
	})
}

func TestMakeTransactionReceiptResponse_Type4(t *testing.T) {
	tx := baseTx()
	tx.Type = 4

	resp := makeTransactionReceiptResponse(100, 0, genHash(10), genHash(11), tx, nil)

	assert.Equal(t, primitives.HexUint(4), resp.Type, "Type should be preserved as 4")
}

func TestMakeTransactionReceiptResponse_Type2(t *testing.T) {
	tx := baseTx()
	tx.Type = 2

	resp := makeTransactionReceiptResponse(100, 0, genHash(10), genHash(11), tx, nil)

	assert.Equal(t, primitives.HexUint(2), resp.Type)
}

func TestMakeTransactionReceiptResponse_UnknownType_FallsBackToZero(t *testing.T) {
	tx := baseTx()
	tx.Type = 5

	resp := makeTransactionReceiptResponse(100, 0, genHash(10), genHash(11), tx, nil)

	assert.Equal(t, primitives.HexUint(0), resp.Type, "Unknown type > 4 should fall back to 0")
}

func TestMakeTransactionReceiptResponse_Type0_Unchanged(t *testing.T) {
	tx := baseTx()
	tx.Type = 0

	resp := makeTransactionReceiptResponse(100, 0, genHash(10), genHash(11), tx, nil)

	assert.Equal(t, primitives.HexUint(0), resp.Type)
}

func TestMakeTransactionResponse_Type4_EmptyAuthList(t *testing.T) {
	tx := baseTx()
	tx.Type = 4
	tx.AccessList = []dbt.AccessListEntry{}
	tx.MaxPriorityFeePerGas = genQuantity(800)
	tx.MaxFeePerGas = genQuantity(801)
	tx.AuthorizationList = []dbt.Authorization{}

	resp := makeTransactionResponse(testChainId, 100, 0, genHash(10), genHash(11), tx)

	assert.Empty(t, resp.AuthorizationList, "Empty auth list should be preserved as empty slice")
}

func TestMakeTransactionResponse_Type0_NoYParity(t *testing.T) {
	tx := baseTx() // Type=0, V=27 by default
	resp := makeTransactionResponse(testChainId, 100, 0, genHash(10), genHash(11), tx)
	assert.Nil(t, resp.YParity, "Type 0 should never have yParity")
}

func TestMakeTransactionResponse_Type1_V1_YParity1(t *testing.T) {
	tx := baseTx()
	tx.Type = 1
	tx.V = 1
	tx.YParity = nil
	tx.AccessList = []dbt.AccessListEntry{}
	resp := makeTransactionResponse(testChainId, 100, 0, genHash(10), genHash(11), tx)
	require.NotNil(t, resp.YParity, "Type 1 with V=1 should have yParity")
	assert.Equal(t, primitives.HexUint(1), *resp.YParity)
}

// TestMakeTransactionResponse_Type4_FullFields tests that type 4 transactions produce
// correct transaction and receipt responses with authorization list fields.
func TestMakeTransactionResponse_Type4_FullFields(t *testing.T) {
	tx := baseTx()
	tx.Type = 4
	tx.AccessList = []dbt.AccessListEntry{}
	tx.MaxPriorityFeePerGas = genQuantity(1000)
	tx.MaxFeePerGas = genQuantity(1001)
	tx.AuthorizationList = []dbt.Authorization{
		{
			ChainId: 1313161554,
			Address: genAddress(1002),
			Nonce:   7,
			YParity: 1,
			R:       genQuantity(1003),
			S:       genQuantity(1004),
		},
	}

	blockHash := genHash(1050)
	txHash := genHash(1051)
	height := uint64(1000)
	index := uint64(0)

	// Test transaction response
	txResp := makeTransactionResponse(testChainId, height, index, blockHash, txHash, tx)
	assert.Equal(t, primitives.HexUint(4), txResp.Type)
	require.NotNil(t, txResp.ChainID)
	require.Len(t, txResp.AuthorizationList, 1)
	assert.Equal(t, primitives.HexUint(1313161554), txResp.AuthorizationList[0].ChainID)
	assert.Equal(t, primitives.HexUint(7), txResp.AuthorizationList[0].Nonce)
	assert.Equal(t, primitives.HexUint(1), txResp.AuthorizationList[0].YParity)

	// Test receipt response
	receiptResp := makeTransactionReceiptResponse(height, index, blockHash, txHash, tx, []*response.Log{})
	assert.Equal(t, primitives.HexUint(4), receiptResp.Type)
}
