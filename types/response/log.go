package response

import (
	"github.com/aurora-is-near/relayer2-base/types/primitives"
)

// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_getfilterchanges
// https://docs.infura.io/infura/networks/ethereum/json-rpc-methods/eth_getlogs
type Log struct {
	Removed          bool                `json:"removed"`
	LogIndex         primitives.HexUint  `json:"logIndex"`
	TransactionIndex primitives.HexUint  `json:"transactionIndex"`
	TransactionHash  primitives.Data32   `json:"transactionHash"`
	BlockHash        primitives.Data32   `json:"blockHash"`
	BlockNumber      primitives.HexUint  `json:"blockNumber"`
	Address          primitives.Data20   `json:"address"`
	Data             primitives.VarData  `json:"data"`
	Topics           []primitives.Data32 `json:"topics"`
}
