package db

import (
	"github.com/aurora-is-near/relayer2-base/types/primitives"
)

// https://docs.infura.io/infura/networks/ethereum/concepts/transaction-types
type Transaction struct {
	Type                 uint64
	From                 primitives.Data20
	IsContractDeployment bool
	ToOrContract         *primitives.Data20
	Nonce                primitives.Quantity
	GasPrice             primitives.Quantity
	GasLimit             primitives.Quantity
	GasUsed              uint64
	CumulativeGasUsed    primitives.Quantity
	Value                primitives.Quantity
	Input                primitives.VarData
	NearHash             *primitives.Data32
	NearReceiptHash      primitives.Data32
	Status               bool
	V                    uint64
	R                    primitives.Quantity
	S                    primitives.Quantity
	YParity              *uint8
	LogsBloom            primitives.Data256
	AccessList           []AccessListEntry   // Only for post-EIP-2930 transactions (TxType >= 0x1)
	MaxPriorityFeePerGas primitives.Quantity // Only for post-EIP-1559 transactions (TxType >= 0x2)
	MaxFeePerGas         primitives.Quantity // Only for post-EIP-1559 transactions (TxType >= 0x2)
	AuthorizationList    []Authorization     // EIP-7702: authorization list (TxType == 0x04)
}
