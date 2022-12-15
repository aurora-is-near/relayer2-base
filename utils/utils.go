package utils

import (
	"aurora-relayer-go-common/log"
	"aurora-relayer-go-common/tinypack"
	dbt "aurora-relayer-go-common/types/db"
	"aurora-relayer-go-common/types/indexer"
	"aurora-relayer-go-common/types/primitives"
	"crypto/sha256"
	"encoding/binary"
)

const (
	AccountId = "aurora"
)

func IndexerBlockToDbBlock(block *indexer.Block) *dbt.Block {
	b := dbt.Block{
		ParentHash:       block.ParentHash,
		Miner:            block.Miner,
		Timestamp:        uint64(block.Timestamp),
		GasLimit:         block.GasLimit,
		GasUsed:          block.GasUsed,
		LogsBloom:        block.LogsBloom,
		TransactionsRoot: block.TransactionsRoot,
		StateRoot:        block.StateRoot,
		ReceiptsRoot:     block.ReceiptsRoot,
		Size:             uint64(block.Size),
	}
	return &b
}

func IndexerTxnToDbTxn(txn *indexer.Transaction) *dbt.Transaction {

	toOrContract := tinypack.CreateNullable[primitives.Data20](nil)
	isContractDeployment := false
	if txn.ContractAddress != nil {
		if txn.To != nil {
			log.Log().Warn().Msgf("both contract address and to address is set for txn: [%v], to: [%s], contract: [%s]",
				txn.Hash, txn.To.Hex(), txn.ContractAddress.Hex())
		}
		isContractDeployment = true
		toOrContract = tinypack.CreateNullable[primitives.Data20](txn.ContractAddress)
	} else if txn.To != nil {
		toOrContract = tinypack.CreateNullable[primitives.Data20](txn.To)
	} else {
		log.Log().Warn().Msgf("both contract address and to address is null for txn: [%v]", txn.Hash.Hex())
	}

	var accessListEntries []dbt.AccessListEntry
	for _, al := range txn.AccessList {
		var storageKeys []primitives.Data32
		for _, sk := range al.StorageKeys {
			storageKeys = append(storageKeys, sk)
		}
		sk := tinypack.CreateList[primitives.VarLen, primitives.Data32](storageKeys...)
		accessListEntry := dbt.AccessListEntry{
			Address:     al.Address,
			StorageKeys: tinypack.VarList[primitives.Data32]{sk},
		}
		accessListEntries = append(accessListEntries, accessListEntry)
	}
	ake := tinypack.CreateList[primitives.VarLen, dbt.AccessListEntry](accessListEntries...)

	nearHash := tinypack.CreateNullable[primitives.Data32](nil)
	if txn.NearTransaction.Hash != nil {
		nh := primitives.Data32(*txn.NearTransaction.Hash)
		nearHash = tinypack.CreateNullable[primitives.Data32](&nh)
	}

	t := dbt.Transaction{
		Type:                 txn.TxType,
		From:                 txn.From,
		IsContractDeployment: isContractDeployment,
		ToOrContract:         toOrContract,
		Nonce:                txn.Nonce,
		GasPrice:             txn.GasPrice,
		GasLimit:             txn.GasLimit,
		GasUsed:              txn.GasUsed,
		Value:                txn.Value,
		Input:                primitives.VarData(txn.Input),
		NearHash:             nearHash,
		NearReceiptHash:      primitives.Data32(txn.NearTransaction.ReceiptHash),
		Status:               txn.Status,
		V:                    txn.V,
		R:                    txn.R,
		S:                    txn.S,
		LogsBloom:            txn.LogsBloom,
		AccessList:           tinypack.VarList[dbt.AccessListEntry]{ake},
		MaxPriorityFeePerGas: txn.MaxPriorityFeePerGas,
		MaxFeePerGas:         txn.MaxFeePerGas,
	}
	return &t
}

func IndexerLogToDbLog(log *indexer.Log) *dbt.Log {
	var topics []primitives.Data32
	for _, t := range log.Topics {
		topics = append(topics, primitives.Data32(t))
	}
	t := tinypack.CreateList[primitives.VarLen, primitives.Data32](topics...)

	l := dbt.Log{
		Address: log.Address,
		Data:    primitives.VarData(log.Data),
		Topics:  tinypack.VarList[primitives.Data32]{t},
	}
	return &l
}

func ComputeBlockHash(bHeight, chainId uint64) []byte {
	bufEmpty25 := make([]byte, 25)

	bufCId := make([]byte, 8)
	binary.BigEndian.PutUint64(bufCId, chainId)

	bufAId := []byte(AccountId)

	bufBH := make([]byte, 8)
	binary.BigEndian.PutUint64(bufBH, bHeight)

	bufHash := append(bufEmpty25, bufCId...)
	bufHash = append(bufHash, bufAId...)
	bufHash = append(bufHash, bufBH...)
	hash := sha256.Sum256(bufHash)

	return hash[:]
}
