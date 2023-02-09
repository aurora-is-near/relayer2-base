package response

import (
	"relayer2-base/types/primitives"
)

// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_gettransactionbyhash
// https://docs.infura.io/infura/networks/ethereum/json-rpc-methods/eth_gettransactionbyhash
type Transaction struct {
	AccessList           *[]*AccessListEntry  `json:"accessList,omitempty"`
	BlockHash            primitives.Data32    `json:"blockHash"`
	BlockNumber          primitives.HexUint   `json:"blockNumber"`
	ChainID              *primitives.HexUint  `json:"chainID,omitempty"`
	From                 primitives.Data20    `json:"from"`
	Gas                  primitives.HexUint   `json:"gas"`
	GasPrice             primitives.Quantity  `json:"gasPrice"`
	Hash                 primitives.Data32    `json:"hash"`
	Input                primitives.VarData   `json:"input"`
	MaxPriorityFeePerGas *primitives.Quantity `json:"maxPriorityFeePerGas,omitempty"`
	MaxFeePerGas         *primitives.Quantity `json:"maxFeePerGas,omitempty"`
	Nonce                primitives.Quantity  `json:"nonce"`
	V                    primitives.HexUint   `json:"v"`
	R                    primitives.Quantity  `json:"r"`
	S                    primitives.Quantity  `json:"s"`
	To                   *primitives.Data20   `json:"to"`
	Type                 primitives.HexUint   `json:"type"`
	TransactionIndex     primitives.HexUint   `json:"transactionIndex"`
	Value                primitives.Quantity  `json:"value"`
}
