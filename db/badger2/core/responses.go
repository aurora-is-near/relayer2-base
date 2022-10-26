package db

import (
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
	"aurora-relayer-go-common/db/badger2/core/dbresponses"
	"aurora-relayer-go-common/db/badger2/core/dbtypes"
)

func makeBlockResponse(height uint64, hash dbp.Data32, data dbtypes.Block, txs []any) *dbresponses.Block {
	return &dbresponses.Block{
		Number:           dbp.HexUint(height),
		Hash:             hash,
		ParentHash:       data.ParentHash,
		Nonce:            dbp.Data8FromBytes(nil),
		Sha3Uncles:       dbp.Data32FromBytes(nil),
		LogsBloom:        data.LogsBloom,
		TransactionsRoot: data.TransactionsRoot,
		StateRoot:        data.StateRoot,
		ReceiptsRoot:     data.ReceiptsRoot,
		Miner:            data.Miner,
		Difficulty:       dbp.HexUint(0),
		TotalDifficulty:  dbp.HexUint(0),
		ExtraData:        dbp.VarDataFromBytes(nil),
		Size:             dbp.HexUint(data.Size),
		GasLimit:         data.GasLimit,
		GasUsed:          data.GasUsed,
		Timestamp:        dbp.HexUint(data.Timestamp),
		Transactions:     txs,
		Uncles:           []dbp.Data[dbp.Len32]{},
	}
}

func makeTransactionResponse(
	chainId uint64,
	height uint64,
	index uint64,
	blockHash dbp.Data32,
	hash dbp.Data32,
	data *dbtypes.Transaction,
) *dbresponses.Transaction {

	tx := &dbresponses.Transaction{
		BlockHash:        blockHash,
		BlockNumber:      dbp.HexUint(height),
		From:             data.From,
		Gas:              dbp.HexUint(data.GasUsed), // Or data.GasLimit?
		GasPrice:         data.GasPrice,
		Hash:             hash,
		Input:            data.Input,
		Nonce:            data.Nonce,
		V:                dbp.HexUint(data.V),
		R:                data.R,
		S:                data.S,
		TransactionIndex: dbp.HexUint(index),
		Type:             dbp.HexUint(data.Type),
		Value:            data.Value,
	}
	if !data.IsContractDeployment {
		tx.To = &data.ToOrContract
	}
	if data.Type >= 1 {
		accessList := make([]*dbresponses.AccessListEntry, len(data.AccessList.Content))
		for i, alEntry := range data.AccessList.Content {
			accessList[i] = &dbresponses.AccessListEntry{
				Address:     alEntry.Address,
				StorageKeys: alEntry.StorageKeys.Content,
			}
		}
		tx.AccessList = &accessList
	}
	if data.Type >= 2 {
		chainId := dbp.HexUint(chainId)
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
	blockHash dbp.Data32,
	txHash dbp.Data32,
	data *dbtypes.Log,
) *dbresponses.Log {

	return &dbresponses.Log{
		Removed:          false,
		LogIndex:         dbp.HexUint(logIndex),
		TransactionIndex: dbp.HexUint(txIndex),
		TransactionHash:  txHash,
		BlockHash:        blockHash,
		BlockNumber:      dbp.HexUint(height),
		Address:          data.Address,
		Topics:           data.Topics.Content,
	}
}

func makeTransactionReceiptResponse(
	height uint64,
	txIndex uint64,
	blockHash dbp.Data32,
	txHash dbp.Data32,
	txData *dbtypes.Transaction,
	Logs []*dbresponses.Log,
) *dbresponses.TransactionReceipt {

	txReceipt := &dbresponses.TransactionReceipt{
		BlockHash:         blockHash,
		BlockNumber:       dbp.HexUint(height),
		CumulativeGasUsed: dbp.QuantityFromUint64(0), // TODO: check
		EffectiveGasPrice: txData.GasPrice,
		From:              txData.From,
		GasUsed:           dbp.HexUint(txData.GasUsed),
		Logs:              Logs,
		LogsBloom:         txData.LogsBloom,
		TransactionHash:   txHash,
		TransactionIndex:  dbp.HexUint(txIndex),
		Type:              dbp.HexUint(txData.Type),
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
