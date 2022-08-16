package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Uint256 string
type H256 string
type Address string
type Bytea []byte

type Block struct {
	ChainId          uint64        `cbor:"chain_id"`
	Hash             H256          `cbor:"hash"`        // H256
	ParentHash       H256          `cbor:"parent_hash"` // H256
	Height           uint64        `cbor:"height"`
	Miner            Address       `cbor:"miner"` // Address
	Timestamp        int64         `cbor:"timestamp"`
	GasLimit         Uint256       `cbor:"gas_limit"`         // U256
	GasUsed          Uint256       `cbor:"gas_used"`          // U256
	LogsBloom        string        `cbor:"logs_bloom"`        // U256
	TransactionsRoot H256          `cbor:"transactions_root"` // H256
	ReceiptsRoot     H256          `cbor:"receipts_root"`     // H256
	Transactions     []Transaction `cbor:"transactions"`      // Vec<AuroraTransaction>,
	NearBlock        any           `cbor:"near_metadata"`     // pub near_metadata: NearBlock,
	StateRoot        string        `cbor:"state_root"`
	Size             Uint256       `cbor:"size"`
	Sequence         uint64
}

type Transaction struct {
	Hash                 H256            `cbor:"hash"`       // H256
	BlockHash            H256            `cbor:"block_hash"` // H256
	BlockHeight          uint64          `cbor:"block_height"`
	ChainId              uint64          `cbor:"chain_id"`
	TransactionIndex     uint32          `cbor:"transaction_index"`
	From                 Address         `cbor:"from"`                     // Address
	To                   Address         `cbor:"to"`                       // Address
	Nonce                Uint256         `cbor:"nonce"`                    // U256
	GasPrice             Uint256         `cbor:"gas_price"`                // U256
	GasLimit             Uint256         `cbor:"gas_limit"`                // U256
	GasUsed              uint64          `cbor:"gas_used"`                 // U256
	MaxPriorityFeePerGas Uint256         `cbor:"max_priority_fee_per_gas"` // U256
	MaxFeePerGas         Uint256         `cbor:"max_fee_per_gas"`          // U256
	Value                Uint256         `cbor:"value"`                    // Wei
	Input                Bytea           `cbor:"input"`                    // Vec<u8>
	Output               Bytea           `cbor:"output"`                   // Vec<u8>
	AccessList           []AccessList    `cbor:"access_list"`              // Vec<AccessTuple>
	TxType               uint8           `cbor:"tx_type"`                  // u8
	Status               bool            `cbor:"status"`                   // bool
	Logs                 []Log           `cbor:"logs"`                     // Vec<ResultLog>,
	ContractAddress      Address         `cbor:"contract_address"`         // Address
	V                    uint64          `cbor:"v"`                        // U64
	R                    Uint256         `cbor:"r"`                        // U256
	S                    Uint256         `cbor:"s"`                        // U256
	NearTransaction      NearTransaction `cbor:"near_metadata"`            // pub near_metadata: NearTransaction,
}

type AccessList struct {
	Address     Address `json:"address"`
	StorageKeys []H256  `json:"storageKeys"`
}

type Log struct {
	Address Address `cbor:"Address"`
	Topics  []Bytea `cbor:"Topics"`
	Data    Bytea   `cbor:"data"`
}

type ExistingBlock struct {
	NearHash       string `cbor:"near_hash"`        // CryptoHash
	NearParentHash string `cbor:"near_parent_hash"` // CryptoHash
	Author         string `cbor:"author"`           // AuroraId
}

type NearTransaction struct {
	Hash        string `cbor:"hash"`         // Vec<AccessTuple>
	ReceiptHash string `cbor:"receipt_hash"` // Vec<AccessTuple>
}

type BlockResponse struct {
	Difficulty       Uint256       `json:"difficulty"`
	ExtraData        Bytea         `json:"extraData"`
	GasLimit         Uint256       `json:"gasLimit"`
	GasUsed          Uint256       `json:"gasUsed"`
	Hash             H256          `json:"hash"`
	LogsBloom        Bytea         `json:"logsBloom"`
	Miner            Address       `json:"miner"`
	Nonce            Bytea         `json:"nonce"`
	Number           Uint256       `json:"number"`
	ParentHash       H256          `json:"parentHash"`
	ReceiptsRoot     H256          `json:"receiptsRoot"`
	Sha3Uncles       H256          `json:"sha3Uncles"`
	Size             Uint256       `json:"size"`
	StateRoot        H256          `json:"stateRoot"`
	Timestamp        Uint256       `json:"timestamp"`
	TotalDifficulty  Uint256       `json:"totalDifficulty"`
	Transactions     []interface{} `json:"transactions"`
	TransactionsRoot H256          `json:"transactionsRoot"`
	Uncles           []H256        `json:"uncles"`

	// difficulty 0x0,
	// extraData 0x,
	// gasLimit 0x0,
	// gasUsed 0x0,
	// hash 0xd93b6e4c7fe00b396adcd9ee4546d3d9e270788069ec1158cbb5b64bbbd31644,
	// logsBloom 0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000,
	// miner 0x0000000000000000000000000000000000000000,
	// nonce 0x0000000000000000,
	// number 0x59f0bf3,
	// parentHash 0x83d3382c11f3b65b574f674c87ded1027dc3777ef831295ce0f2fbcfacd0377e,
	// receiptsRoot 0x45f5311f2f0bc736546b9b21e04ece0e33775d2d757fa27ee0a1dc560f527315,
	// sha3Uncles 0x0000000000000000000000000000000000000000000000000000000000000000,
	// size 0xd3d,
	// stateRoot 0xa190208e2d45a746603753d8125c9cac8dfa8ee8483b466b3b9690df55b4ca89,
	// timestamp 0x62c7099e,
	// totalDifficulty 0x0,
	// transactions [],
	// transactionsRoot 0x1223349a40d2ee10bd1bebb5889ef8018c8bc13359ed94b387810af96c6e4268,
	// uncles []
}

type TransactionResponse struct {
	AccessList       *[]AccessListResponse `json:"accessList"`
	BlockHash        H256                  `json:"blockHash"`
	BlockNumber      Uint256               `json:"blockNumber"`
	ChainID          Uint256               `json:"chainID"`
	Hash             H256                  `json:"hash"`
	From             Address               `json:"from"`
	To               *Address              `json:"to"`
	Gas              Uint256               `json:"gas"`
	GasPrice         Uint256               `json:"gasPrice"`
	Value            Uint256               `json:"value"`
	Input            Bytea                 `json:"input"`
	Nonce            Uint256               `json:"nonce"`
	TransactionIndex Uint256               `json:"transactionIndex"`
	V                Uint256               `json:"v"`
	R                Bytea                 `json:"r"`
	S                Bytea                 `json:"s"`
	// Type             Uint256  `json:"type"` // not in original relayer?
}

type AccessListResponse struct{}

func (bl *Block) TxCount() int64 {
	return int64(len(bl.Transactions))
}

func (bl *Block) ToResponse(full bool) *BlockResponse {
	transactions := make([]interface{}, 0, len(bl.Transactions))
	for _, tx := range bl.Transactions {
		if !full {
			transactions = append(transactions, tx.Hash)
		} else {
			transactions = append(transactions, tx.ToResponse())
		}
	}
	return &BlockResponse{
		Difficulty:       "0",
		ExtraData:        Bytea("0x"),
		GasLimit:         "0",
		GasUsed:          "0",
		Hash:             bl.Hash,
		LogsBloom:        []byte(fmt.Sprintf("%0x", 0)),
		Miner:            bl.Miner,
		Nonce:            []byte(fmt.Sprintf("%016x", 0)),
		ParentHash:       bl.ParentHash,
		ReceiptsRoot:     bl.ReceiptsRoot,
		Sha3Uncles:       H256("x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
		Size:             bl.Size,
		StateRoot:        H256(bl.StateRoot),
		Timestamp:        IntToUint256(bl.Timestamp),
		Transactions:     transactions,
		TransactionsRoot: bl.TransactionsRoot,
		Uncles:           []H256{},
	}
}

func (tx *Transaction) ToResponse() *TransactionResponse {
	return &TransactionResponse{}
}

func (a Address) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"0x%s"`, a)), nil
}

func (b Bytea) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"0x%s"`, b)), nil
}

func (h H256) MarshalJSON() ([]byte, error) {
	if h == "" {
		h = H256(fmt.Sprintf("%064x", 0))
	}
	return []byte(fmt.Sprintf(`"0x%s"`, h)), nil
}

func (i Uint256) MarshalJSON() ([]byte, error) {
	if i == "" {
		i = "0"
	}
	return []byte(fmt.Sprintf(`"0x%s"`, i)), nil
}

func (h *H256) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*h = ParseHexString[H256](s)
	return nil
}

func (i *Uint256) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*i = ParseHexString[Uint256](s)
	return nil
}

func (i Uint256) ToInt64() (int64, error) {
	return strconv.ParseInt(string(i), 16, 64)
}
