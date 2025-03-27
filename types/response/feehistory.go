package response

import (
	"github.com/aurora-is-near/relayer2-base/types/primitives"
)

// https://docs.metamask.io/services/reference/ethereum/json-rpc-methods/eth_feehistory/
//
//easyjson:json
type FeeHistory struct {
	BaseFeePerBlobGas []primitives.Quantity   `json:"baseFeePerBlobGas"`
	BaseFeePerGas     []primitives.Quantity   `json:"baseFeePerGas"`
	Reward            [][]primitives.Quantity `json:"reward"`
	BlobGasUsedRatio  []float32               `json:"blobGasUsedRatio"`
	GasUsedRatio      []float32               `json:"gasUsedRatio"`
	OldestBlock       primitives.HexUint      `json:"oldestBlock"`
}
