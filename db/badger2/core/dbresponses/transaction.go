package dbresponses

import (
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
)

// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_gettransactionbyhash
// https://docs.infura.io/infura/networks/ethereum/json-rpc-methods/eth_gettransactionbyhash
type Transaction struct {
	AccessList           *[]*AccessListEntry `json:"accessList,omitempty"`
	BlockHash            dbp.Data32          `json:"blockHash"`
	BlockNumber          dbp.HexUint         `json:"blockNumber"`
	ChainID              *dbp.HexUint        `json:"chainID,omitempty"`
	From                 dbp.Data20          `json:"from"`
	Gas                  dbp.HexUint         `json:"gas"`
	GasPrice             dbp.Quantity        `json:"gasPrice"`
	Hash                 dbp.Data32          `json:"hash"`
	Input                dbp.VarData         `json:"input"`
	MaxPriorityFeePerGas *dbp.Quantity       `json:"maxPriorityFeePerGas,omitempty"`
	MaxFeePerGas         *dbp.Quantity       `json:"maxFeePerGas,omitempty"`
	Nonce                dbp.Quantity        `json:"nonce"`
	V                    dbp.HexUint         `json:"v"`
	R                    dbp.Quantity        `json:"r"`
	S                    dbp.Quantity        `json:"s"`
	To                   *dbp.Data20         `json:"to"`
	TransactionIndex     dbp.HexUint         `json:"transactionIndex"`
	Type                 dbp.HexUint         `json:"type"`
	Value                dbp.Quantity        `json:"value"`
}
