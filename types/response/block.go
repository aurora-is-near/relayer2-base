package response

import (
	"aurora-relayer-go-common/types/primitives"
)

// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_getblockbyhash
// https://docs.infura.io/infura/networks/ethereum/json-rpc-methods/eth_getblockbyhash
//easyjson:json
type Block struct {
	Number           primitives.HexUint  `json:"number"`
	Hash             primitives.Data32   `json:"hash"`
	ParentHash       primitives.Data32   `json:"parentHash"`
	Nonce            primitives.Data8    `json:"nonce"`
	Sha3Uncles       primitives.Data32   `json:"sha3Uncles"`
	LogsBloom        primitives.Data256  `json:"logsBloom"`
	TransactionsRoot primitives.Data32   `json:"transactionsRoot"`
	StateRoot        primitives.Data32   `json:"stateRoot"`
	ReceiptsRoot     primitives.Data32   `json:"receiptsRoot"`
	Miner            primitives.Data20   `json:"miner"`
	Difficulty       primitives.HexUint  `json:"difficulty"`
	TotalDifficulty  primitives.HexUint  `json:"totalDifficulty"`
	ExtraData        primitives.VarData  `json:"extraData"`
	Size             primitives.HexUint  `json:"size"`
	GasLimit         primitives.Quantity `json:"gasLimit"`
	GasUsed          primitives.Quantity `json:"gasUsed"`
	Timestamp        primitives.HexUint  `json:"timestamp"`
	Transactions     []any               `json:"transactions"`
	Uncles           []primitives.Data32 `json:"uncles"`
}
