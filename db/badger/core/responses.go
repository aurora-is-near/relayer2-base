package core

import (
	dbt "aurora-relayer-go-common/types/db"
	"aurora-relayer-go-common/types/primitives"
	"aurora-relayer-go-common/types/response"
)

func makeBlockResponse(height uint64, hash primitives.Data32, data dbt.Block, txs []any) *response.Block {
	return &response.Block{
		Number:           primitives.HexUint(height),
		Hash:             hash,
		ParentHash:       data.ParentHash,
		Nonce:            primitives.Data8FromBytes(nil),
		Sha3Uncles:       primitives.Data32FromBytes(nil),
		LogsBloom:        data.LogsBloom,
		TransactionsRoot: data.TransactionsRoot,
		StateRoot:        data.StateRoot,
		ReceiptsRoot:     data.ReceiptsRoot,
		Miner:            data.Miner,
		Difficulty:       primitives.HexUint(0),
		TotalDifficulty:  primitives.HexUint(0),
		ExtraData:        primitives.VarDataFromBytes(nil),
		Size:             primitives.HexUint(data.Size),
		GasLimit:         data.GasLimit,
		GasUsed:          data.GasUsed,
		Timestamp:        primitives.HexUint(data.Timestamp),
		Transactions:     txs,
		Uncles:           []primitives.Data[primitives.Len32]{},
	}
}

func makeTransactionResponse(
	chainId uint64,
	height uint64,
	index uint64,
	blockHash primitives.Data32,
	hash primitives.Data32,
	data *dbt.Transaction,
) *response.Transaction {

	tx := &response.Transaction{
		BlockHash:        blockHash,
		BlockNumber:      primitives.HexUint(height),
		From:             data.From,
		Gas:              primitives.HexUint(data.GasLimit.Uint64()),
		GasPrice:         data.GasPrice,
		Hash:             hash,
		Input:            data.Input,
		Nonce:            data.Nonce,
		V:                primitives.HexUint(data.V),
		R:                data.R,
		S:                data.S,
		TransactionIndex: primitives.HexUint(index),
		Type:             primitives.HexUint(data.Type),
		Value:            data.Value,
	}
	if !data.IsContractDeployment {
		tx.To = &data.ToOrContract
	}
	if data.Type >= 1 {
		accessList := make([]*response.AccessListEntry, len(data.AccessList.Content))
		for i, alEntry := range data.AccessList.Content {
			accessList[i] = &response.AccessListEntry{
				Address:     alEntry.Address,
				StorageKeys: alEntry.StorageKeys.Content,
			}
		}
		tx.AccessList = &accessList
	}
	if data.Type >= 2 {
		chainId := primitives.HexUint(chainId)
		tx.ChainID = &chainId
		tx.MaxPriorityFeePerGas = &data.MaxPriorityFeePerGas
		tx.MaxFeePerGas = &data.MaxFeePerGas
	}

	return tx
}

func makeLogResponse(
	height uint64,
	txIndex uint64,
	logIndex uint64,
	blockHash primitives.Data32,
	txHash primitives.Data32,
	data *dbt.Log,
) *response.Log {

	return &response.Log{
		Removed:          false,
		LogIndex:         primitives.HexUint(logIndex),
		TransactionIndex: primitives.HexUint(txIndex),
		TransactionHash:  txHash,
		BlockHash:        blockHash,
		BlockNumber:      primitives.HexUint(height),
		Address:          data.Address,
		Topics:           data.Topics.Content,
	}
}

func makeTransactionReceiptResponse(
	height uint64,
	txIndex uint64,
	blockHash primitives.Data32,
	txHash primitives.Data32,
	txData *dbt.Transaction,
	Logs []*response.Log,
) *response.TransactionReceipt {

	txReceipt := &response.TransactionReceipt{
		BlockHash:         blockHash,
		BlockNumber:       primitives.HexUint(height),
		CumulativeGasUsed: primitives.QuantityFromUint64(0), // TODO: check
		EffectiveGasPrice: txData.GasPrice,
		From:              txData.From,
		GasUsed:           primitives.HexUint(txData.GasUsed),
		Logs:              Logs,
		LogsBloom:         txData.LogsBloom,
		TransactionHash:   txHash,
		TransactionIndex:  primitives.HexUint(txIndex),
		Type:              primitives.HexUint(txData.Type),
	}
	if txData.IsContractDeployment {
		txReceipt.ContractAddress = &txData.ToOrContract
	} else {
		txReceipt.To = &txData.ToOrContract
	}
	if txData.Status {
		txReceipt.Status = 1
	} else {
		txReceipt.Status = 0
	}

	return txReceipt
}
