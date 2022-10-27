package dbresponses

import (
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
)

// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_getblockbyhash
// https://docs.infura.io/infura/networks/ethereum/json-rpc-methods/eth_getblockbyhash
//easyjson:json
type Block struct {
	Number           dbp.HexUint  `json:"number"`
	Hash             dbp.Data32   `json:"hash"`
	ParentHash       dbp.Data32   `json:"parentHash"`
	Nonce            dbp.Data8    `json:"nonce"`
	Sha3Uncles       dbp.Data32   `json:"sha3Uncles"`
	LogsBloom        dbp.Data256  `json:"logsBloom"`
	TransactionsRoot dbp.Data32   `json:"transactionsRoot"`
	StateRoot        dbp.Data32   `json:"stateRoot"`
	ReceiptsRoot     dbp.Data32   `json:"receiptsRoot"`
	Miner            dbp.Data20   `json:"miner"`
	Difficulty       dbp.HexUint  `json:"difficulty"`
	TotalDifficulty  dbp.HexUint  `json:"totalDifficulty"`
	ExtraData        dbp.VarData  `json:"extraData"`
	Size             dbp.HexUint  `json:"size"`
	GasLimit         dbp.Quantity `json:"gasLimit"`
	GasUsed          dbp.Quantity `json:"gasUsed"`
	Timestamp        dbp.HexUint  `json:"timestamp"`
	Transactions     []any        `json:"transactions"`
	Uncles           []dbp.Data32 `json:"uncles"`
}
