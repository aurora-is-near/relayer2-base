package core

import (
	"github.com/aurora-is-near/relayer2-base/tweaks"
	dbt "github.com/aurora-is-near/relayer2-base/types/db"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/aurora-is-near/relayer2-base/types/response"
	"github.com/aurora-is-near/relayer2-base/utils"
)

var (
	sha3Uncles = primitives.MustData32FromHex("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347")
	nonce      = primitives.Data8FromBytes(nil)
	difficulty = primitives.HexUint(0)
	extraData  = primitives.VarDataFromBytes(nil)
	uncles     = []primitives.Data[primitives.Len32]{}
	gasLimit   = primitives.QuantityFromUint64(*utils.Constants.GasLimit())
	miner      = primitives.MustData20FromHex(*utils.Constants.ZeroStrUint160())
	mixHash    = primitives.MustData32FromHex(*utils.Constants.ZeroStrUint256())
)

func makeBlockResponse(height uint64, hash primitives.Data32, data dbt.Block, txs []any) *response.Block {
	// It is not allowed to burn more gas than the gas limit. However, because of an error there were
	// invalid gas consumption records. Below control is added to support the following fix.
	// https://github.com/aurora-is-near/aurora-relayer/pull/348/files
	gasUsed := data.GasUsed
	if data.GasUsed.Uint64() > *utils.Constants.GasLimit() {
		gasUsed = primitives.QuantityFromUint64(*utils.Constants.GasLimit())
	}

	block := &response.Block{
		BaseFeePerGas:    tweaks.BaseFeePerGas(),
		Number:           primitives.HexUint(height),
		Hash:             hash,
		ParentHash:       data.ParentHash,
		Nonce:            nonce,
		Sha3Uncles:       sha3Uncles,
		LogsBloom:        data.LogsBloom,
		TransactionsRoot: data.TransactionsRoot,
		StateRoot:        data.StateRoot,
		ReceiptsRoot:     data.ReceiptsRoot,
		Miner:            miner,
		MixHash:          mixHash,
		Difficulty:       difficulty,
		TotalDifficulty:  difficulty,
		ExtraData:        extraData,
		Size:             primitives.HexUint(data.Size),
		GasLimit:         gasLimit,
		GasUsed:          gasUsed,
		Timestamp:        primitives.HexUint(data.Timestamp),
		Transactions:     txs,
		Uncles:           uncles,
	}
	return block
}

func makeTransactionResponse(
	chainId uint64,
	height uint64,
	index uint64,
	blockHash primitives.Data32,
	hash primitives.Data32,
	data *dbt.Transaction,
) *response.Transaction {
	// It is not allowed to burn more gas than the gas limit. However, because of an error there were
	// invalid gas consumption records. Below control is added to support the following fix.
	// https://github.com/aurora-is-near/aurora-relayer/pull/348/files
	gas := min(data.GasLimit.Uint64(), *utils.Constants.GasLimit())

	tx := &response.Transaction{
		BlockHash:        blockHash,
		BlockNumber:      primitives.HexUint(height),
		From:             data.From,
		Gas:              primitives.HexUint(gas),
		GasPrice:         data.GasPrice,
		Hash:             hash,
		Input:            data.Input,
		Nonce:            data.Nonce,
		V:                primitives.HexUint(data.V),
		R:                data.R,
		S:                data.S,
		TransactionIndex: primitives.HexUint(index),
		Value:            data.Value,
		Type:             primitives.HexUint(data.Type),
	}
	if !data.IsContractDeployment {
		tx.To = data.ToOrContract
	}
	switch {
	case data.Type > 4:
		tx.Type = primitives.HexUint(0)
	case data.Type >= 1:
		cid := primitives.HexUint(chainId)
		tx.ChainID = &cid
		if data.YParity != nil {
			yParity := primitives.HexUint(*data.YParity)
			tx.YParity = &yParity
		} else if data.V == 0 || data.V == 1 {
			yParity := primitives.HexUint(data.V)
			tx.YParity = &yParity
		}
		accessList := make([]*response.AccessListEntry, len(data.AccessList))
		for i, alEntry := range data.AccessList {
			accessList[i] = &response.AccessListEntry{
				Address:     alEntry.Address,
				StorageKeys: alEntry.StorageKeys,
			}
		}
		tx.AccessList = &accessList
		if data.Type >= 2 {
			tx.MaxPriorityFeePerGas = &data.MaxPriorityFeePerGas
			tx.MaxFeePerGas = &data.MaxFeePerGas
		}
		if data.Type == 4 {
			authList := make([]*response.AuthorizationEntry, len(data.AuthorizationList))
			for i, auth := range data.AuthorizationList {
				authList[i] = &response.AuthorizationEntry{
					ChainID: primitives.HexUint(auth.ChainId),
					Address: auth.Address,
					Nonce:   primitives.HexUint(auth.Nonce),
					YParity: primitives.HexUint(auth.YParity),
					R:       auth.R,
					S:       auth.S,
				}
			}
			tx.AuthorizationList = authList
		}
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
		Topics:           data.Topics,
		Data:             data.Data,
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
	// It is not allowed to burn more gas than the gas limit. However, because of an error there were
	// invalid gas consumption records. Below control is added to support the following fix.
	// https://github.com/aurora-is-near/aurora-relayer/pull/348/files
	gasUsed := min(txData.GasUsed, *utils.Constants.GasLimit())

	txReceipt := &response.TransactionReceipt{
		BlockHash:         blockHash,
		BlockNumber:       primitives.HexUint(height),
		CumulativeGasUsed: txData.CumulativeGasUsed,
		From:              txData.From,
		GasUsed:           primitives.HexUint(gasUsed),
		EffectiveGasPrice: txData.GasPrice,
		Logs:              Logs,
		LogsBloom:         txData.LogsBloom,
		NearReceiptHash:   txData.NearReceiptHash,
		TransactionHash:   txHash,
		TransactionIndex:  primitives.HexUint(txIndex),
		Type:              primitives.HexUint(txData.Type),
	}
	if txData.IsContractDeployment {
		txReceipt.ContractAddress = txData.ToOrContract
	} else {
		txReceipt.To = txData.ToOrContract
	}
	if txData.Status {
		txReceipt.Status = 1
	} else {
		txReceipt.Status = 0
	}
	if txData.NearHash != nil {
		txReceipt.NearTransactionHash = *txData.NearHash
	}
	if txData.Type > 4 {
		txReceipt.Type = primitives.HexUint(0)
	}

	return txReceipt
}
