package dbresponses

import (
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
)

// https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_gettransactionreceipt
// https://docs.infura.io/infura/networks/ethereum/json-rpc-methods/eth_gettransactionreceipt
type TransactionReceipt struct {
	BlockHash         dbp.Data32   `json:"blockHash"`
	BlockNumber       dbp.HexUint  `json:"blockNumber"`
	ContractAddress   *dbp.Data20  `json:"contractAddress"`
	CumulativeGasUsed dbp.Quantity `json:"cumulativeGasUsed"`
	EffectiveGasPrice dbp.Quantity `json:"effectiveGasPrice"`
	From              dbp.Data20   `json:"from"`
	GasUsed           dbp.HexUint  `json:"gasUsed"`
	Logs              []*Log       `json:"logs"`
	LogsBloom         dbp.Data256  `json:"logsBloom"`
	Status            dbp.HexUint  `json:"status"`
	To                *dbp.Data20  `json:"to"`
	TransactionHash   dbp.Data32   `json:"transactionHash"`
	TransactionIndex  dbp.HexUint  `json:"transactionIndex"`
	Type              dbp.HexUint  `json:"type"`
}
