package dbtypes

import (
	"fmt"

	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
	tp "aurora-relayer-go-common/tinypack"
)

//
// https://docs.infura.io/infura/networks/ethereum/concepts/transaction-types
//
type Transaction struct {
	Type                 uint64
	From                 dbp.Data20
	IsContractDeployment bool
	ToOrContract         dbp.Data20
	Nonce                dbp.Quantity
	GasPrice             dbp.Quantity
	GasLimit             dbp.Quantity
	GasUsed              uint64
	Value                dbp.Quantity
	Input                dbp.VarData
	NearHash             dbp.Data32 // Is it 32 bytes long?
	NearReceiptHash      dbp.Data32 // Is it 32 bytes long?
	Status               bool
	V                    uint64
	R                    dbp.Quantity
	S                    dbp.Quantity
	LogsBloom            dbp.Data256
	AccessList           tp.VarList[AccessListEntry] // Only for post-EIP-2930 transactions (TxType >= 0x1)
	MaxPriorityFeePerGas dbp.Quantity                // Only for post-EIP-1559 transactions (TxType >= 0x2)
	MaxFeePerGas         dbp.Quantity                // Only for post-EIP-1559 transactions (TxType >= 0x2)
}

func (tx *Transaction) getFields() []any {
	fields := []any{
		&tx.From,
		&tx.IsContractDeployment,
		&tx.ToOrContract,
		&tx.Nonce,
		&tx.GasPrice,
		&tx.GasLimit,
		&tx.GasUsed,
		&tx.Value,
		&tx.Input,
		&tx.NearHash,
		&tx.NearReceiptHash,
		&tx.Status,
		&tx.V,
		&tx.R,
		&tx.S,
		&tx.LogsBloom,
	}

	if tx.Type >= 1 {
		fields = append(fields, &tx.AccessList)
	}

	if tx.Type >= 2 {
		fields = append(fields, &tx.MaxPriorityFeePerGas, &tx.MaxFeePerGas)
	}

	return fields
}

func (tx *Transaction) WriteTinyPack(w tp.Writer, e *tp.Encoder) error {
	if err := e.WriteUvarint(w, tx.Type); err != nil {
		return fmt.Errorf("can't write tx-type: %v", err)
	}
	return e.Write(w, tx.getFields()...)
}

func (tx *Transaction) ReadTinyPack(r tp.Reader, d *tp.Decoder) error {
	var err error
	tx.Type, err = d.ReadUvarint(r)
	if err != nil {
		return fmt.Errorf("can't read tx-type: %v", err)
	}
	return d.Read(r, tx.getFields()...)
}
