package utils

import (
	"aurora-relayer-go-common/log"
	"aurora-relayer-go-common/tinypack"
	dbt "aurora-relayer-go-common/types/db"
	"aurora-relayer-go-common/types/indexer"
	"aurora-relayer-go-common/types/primitives"
)

func IndexerBlockToDbBlock(block *indexer.Block) *dbt.Block {
	b := dbt.Block{
		ParentHash:       primitives.Data32FromBytes(block.ParentHash.Bytes()),
		Miner:            primitives.Data20FromBytes(block.Miner.Bytes()),
		Timestamp:        uint64(block.Timestamp.Uint),
		GasLimit:         primitives.QuantityFromHex(block.GasLimit.String()),
		GasUsed:          primitives.QuantityFromHex(block.GasUsed.String()),
		LogsBloom:        primitives.Data256FromBytes(block.LogsBloom.Bytes()),
		TransactionsRoot: primitives.Data32FromBytes(block.TransactionsRoot.Bytes()),
		StateRoot:        primitives.Data32FromBytes(block.StateRoot.Bytes()),
		ReceiptsRoot:     primitives.Data32FromBytes(block.ReceiptsRoot.Bytes()),
		Size:             block.Size.ToInt().Uint64(),
	}
	return &b
}

func IndexerTxnToDbTxn(txn *indexer.Transaction) *dbt.Transaction {

	toOrContract := tinypack.CreateNullable[primitives.Data20](nil)
	isContractDeployment := false
	if txn.ContractAddress != nil {
		if txn.To != nil {
			log.Log().Warn().Msgf("both contract address and to address is set for txn: [%v], to: [%s], contract: [%s]",
				txn.Hash, txn.To.String(), txn.ContractAddress.String())
		}
		isContractDeployment = true
		contract := primitives.Data20FromBytes(txn.ContractAddress.Bytes())
		toOrContract = tinypack.CreateNullable[primitives.Data20](&contract)
	} else if txn.To != nil {
		to := primitives.Data20FromBytes(txn.To.Bytes())
		toOrContract = tinypack.CreateNullable[primitives.Data20](&to)
	} else {
		log.Log().Warn().Msgf("both contract address and to address is null for txn: [%v]", txn.Hash)
	}

	var accessListEntries []dbt.AccessListEntry
	for _, al := range txn.AccessList {
		var storageKeys []primitives.Data32
		for _, sk := range al.StorageKeys {
			storageKey := primitives.Data32FromBytes(sk.Bytes())
			storageKeys = append(storageKeys, storageKey)
		}
		sk := tinypack.CreateList[primitives.VarLen, primitives.Data32](storageKeys...)
		accessListEntry := dbt.AccessListEntry{
			Address:     primitives.Data20FromBytes(al.Address.Bytes()),
			StorageKeys: tinypack.VarList[primitives.Data32]{sk},
		}
		accessListEntries = append(accessListEntries, accessListEntry)
	}
	ake := tinypack.CreateList[primitives.VarLen, dbt.AccessListEntry](accessListEntries...)

	t := dbt.Transaction{
		Type:                 txn.TxType.Uint64(),
		From:                 primitives.Data20FromBytes(txn.From.Bytes()),
		IsContractDeployment: isContractDeployment,
		ToOrContract:         toOrContract,
		Nonce:                primitives.QuantityFromHex(txn.Nonce.String()),
		GasPrice:             primitives.QuantityFromHex(txn.GasPrice.String()),
		GasLimit:             primitives.QuantityFromHex(txn.GasLimit.String()),
		GasUsed:              txn.GasUsed.ToInt().Uint64(),
		Value:                primitives.QuantityFromHex(txn.Value.String()),
		Input:                primitives.VarDataFromBytes(txn.Input.Bytes),
		NearHash:             primitives.Data32FromBytes([]byte(txn.NearTransaction.Hash)),
		NearReceiptHash:      primitives.Data32FromBytes([]byte(txn.NearTransaction.ReceiptHash)),
		Status:               txn.Status,
		V:                    txn.V.ToInt().Uint64(),
		R:                    primitives.QuantityFromHex(txn.R.String()),
		S:                    primitives.QuantityFromHex(txn.S.String()),
		LogsBloom:            primitives.Data256FromBytes(txn.LogsBloom.Bytes()),
		AccessList:           tinypack.VarList[dbt.AccessListEntry]{ake},
		MaxPriorityFeePerGas: primitives.QuantityFromHex(txn.MaxPriorityFeePerGas.String()),
		MaxFeePerGas:         primitives.QuantityFromHex(txn.MaxFeePerGas.String()),
	}
	return &t
}

func IndexerLogToDbLog(log *indexer.Log) *dbt.Log {
	var topics []primitives.Data32
	for _, t := range log.Topics {
		topic := primitives.Data32FromHex(t.String())
		topics = append(topics, topic)
	}
	t := tinypack.CreateList[primitives.VarLen, primitives.Data32](topics...)

	l := dbt.Log{
		Address: primitives.Data20FromBytes(log.Address.Bytes()),
		Data:    primitives.VarDataFromHex(log.Data.String()),
		Topics:  tinypack.VarList[primitives.Data32]{t},
	}
	return &l
}
