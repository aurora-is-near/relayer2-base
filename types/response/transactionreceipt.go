package response

import (
	"aurora-relayer-go-common/types/primitives"
)

// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_gettransactionreceipt
// https://docs.infura.io/infura/networks/ethereum/json-rpc-methods/eth_gettransactionreceipt
type TransactionReceipt struct {
	BlockHash         primitives.Data32   `json:"blockHash"`
	BlockNumber       primitives.HexUint  `json:"blockNumber"`
	ContractAddress   *primitives.Data20  `json:"contractAddress"`
	CumulativeGasUsed primitives.Quantity `json:"cumulativeGasUsed"`
	EffectiveGasPrice primitives.Quantity `json:"effectiveGasPrice"`
	From              primitives.Data20   `json:"from"`
	GasUsed           primitives.HexUint  `json:"gasUsed"`
	Logs              []*Log              `json:"logs"`
	LogsBloom         primitives.Data256  `json:"logsBloom"`
	Status            primitives.HexUint  `json:"status"`
	To                *primitives.Data20  `json:"to"`
	TransactionHash   primitives.Data32   `json:"transactionHash"`
	TransactionIndex  primitives.HexUint  `json:"transactionIndex"`
	Type              primitives.HexUint  `json:"type"`
}
