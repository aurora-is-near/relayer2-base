package dbresponses

import (
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
)

// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_getfilterchanges
// https://docs.infura.io/infura/networks/ethereum/json-rpc-methods/eth_getlogs
type Log struct {
	Removed          bool         `json:"removed"`
	LogIndex         dbp.HexUint  `json:"logIndex"`
	TransactionIndex dbp.HexUint  `json:"transactionIndex"`
	TransactionHash  dbp.Data32   `json:"transactionHash"`
	BlockHash        dbp.Data32   `json:"blockHash"`
	BlockNumber      dbp.HexUint  `json:"blockNumber"`
	Address          dbp.Data20   `json:"address"`
	Data             dbp.VarData  `json:"data"`
	Topics           []dbp.Data32 `json:"topics"`
}
