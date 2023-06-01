package response

import (
	"github.com/aurora-is-near/relayer2-base/types/primitives"
)

// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_gettransactionreceipt
// https://docs.infura.io/infura/networks/ethereum/json-rpc-methods/eth_gettransactionreceipt
type TransactionReceipt struct {
	BlockHash           primitives.Data32   `json:"blockHash"`
	BlockNumber         primitives.HexUint  `json:"blockNumber"`
	ContractAddress     *primitives.Data20  `json:"contractAddress"`
	CumulativeGasUsed   primitives.Quantity `json:"cumulativeGasUsed"`
	From                primitives.Data20   `json:"from"`
	GasUsed             primitives.HexUint  `json:"gasUsed"`
	EffectiveGasPrice   primitives.Quantity `json:"effectiveGasPrice"`
	Logs                []*Log              `json:"logs"`
	LogsBloom           primitives.Data256  `json:"logsBloom"`
	NearReceiptHash     primitives.Data32   `json:"nearReceiptHash"`
	NearTransactionHash primitives.Data32   `json:"nearTransactionHash"`
	Status              primitives.HexUint  `json:"status"`
	To                  *primitives.Data20  `json:"to"`
	Type                primitives.HexUint  `json:"type"`
	TransactionHash     primitives.Data32   `json:"transactionHash"`
	TransactionIndex    primitives.HexUint  `json:"transactionIndex"`
}
