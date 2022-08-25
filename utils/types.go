package utils

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Uint256 struct{ V *big.Int }
type H256 struct{ common.Hash }
type Address struct{ common.Address }
type Bytea []byte
type TxnData string

type Block struct {
	ChainId          uint64        `cbor:"chain_id"`
	Hash             H256          `cbor:"hash"`
	ParentHash       H256          `cbor:"parent_hash"`
	Height           uint64        `cbor:"height"`
	Miner            Address       `cbor:"miner"`
	Timestamp        int64         `cbor:"timestamp"`
	GasLimit         Uint256       `cbor:"gas_limit"`
	GasUsed          Uint256       `cbor:"gas_used"`
	LogsBloom        string        `cbor:"logs_bloom"`
	TransactionsRoot H256          `cbor:"transactions_root"`
	ReceiptsRoot     H256          `cbor:"receipts_root"`
	Transactions     []Transaction `cbor:"transactions"`
	NearBlock        any           `cbor:"near_metadata"`
	StateRoot        string        `cbor:"state_root"`
	Size             Uint256       `cbor:"size"`
	Sequence         Uint256       `cbor:"sequence"`
}

type Transaction struct {
	Hash                 H256            `cbor:"hash"`
	BlockHash            H256            `cbor:"block_hash"`
	BlockHeight          uint64          `cbor:"block_height"`
	ChainId              uint64          `cbor:"chain_id"`
	TransactionIndex     uint32          `cbor:"transaction_index"`
	From                 Address         `cbor:"from"`
	To                   *Address        `cbor:"to"`
	Nonce                Uint256         `cbor:"nonce"`
	GasPrice             Uint256         `cbor:"gas_price"`
	GasLimit             Uint256         `cbor:"gas_limit"`
	GasUsed              Uint256         `cbor:"gas_used"`
	MaxPriorityFeePerGas Uint256         `cbor:"max_priority_fee_per_gas"`
	MaxFeePerGas         Uint256         `cbor:"max_fee_per_gas"`
	Value                Uint256         `cbor:"value"`
	Input                Bytea           `cbor:"input"`
	Output               Bytea           `cbor:"output"`
	AccessList           []AccessList    `cbor:"access_list"`
	TxType               uint8           `cbor:"tx_type"`
	Status               bool            `cbor:"status"`
	Logs                 []Log           `cbor:"logs"`
	ContractAddress      Address         `cbor:"contract_address"`
	V                    uint64          `cbor:"v"`
	R                    Uint256         `cbor:"r"`
	S                    Uint256         `cbor:"s"`
	NearTransaction      NearTransaction `cbor:"near_metadata"`
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

type FilterOptions struct {
	Address   interface{} `json:"address"`
	FromBlock interface{} `json:"fromBlock"`
	ToBlock   interface{} `json:"toBlock"`
	Topics    [][][]byte  `json:"topics"`
	BlockHash *H256       `json:"blockhash"`
}

type LogFilter struct {
	Address   map[Address]bool `json:"address"`
	FromBlock *Uint256         `json:"fromBlock"`
	ToBlock   *Uint256         `json:"toBlock"`
	BlockHash *H256            `json:"blockhash"`
	Topics    [][][]byte       `json:"topics"`
}

type StoredFilter struct {
	Type      string
	CreatedBy string
	PollBlock Uint256
	BlockHash *H256
	FromBlock *Uint256
	ToBlock   *Uint256
	Addresses map[Address]bool
	Topics    [][][]byte
}

type ExistingBlock struct {
	NearHash       string `cbor:"near_hash"`
	NearParentHash string `cbor:"near_parent_hash"`
	Author         string `cbor:"author"`
}

type NearTransaction struct {
	Hash        string `cbor:"hash"`
	ReceiptHash string `cbor:"receipt_hash"`
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
}

type TransactionResponse struct {
	BlockHash        H256     `json:"blockHash"`
	BlockNumber      Uint256  `json:"blockNumber"`
	Hash             H256     `json:"hash"`
	From             Address  `json:"from"`
	To               *Address `json:"to,omitempty"`
	Gas              Uint256  `json:"gas"`
	GasPrice         Uint256  `json:"gasPrice"`
	Value            Uint256  `json:"value"`
	Input            Bytea    `json:"input"`
	Nonce            Uint256  `json:"nonce"`
	TransactionIndex Uint256  `json:"transactionIndex"`
	V                Uint256  `json:"v"`
	R                Bytea    `json:"r"`
	S                Bytea    `json:"s"`
	// AccessList       *[]AccessListResponse `json:"accessList"`        // not in original relayer
	// ChainID          Uint256               `json:"chainID,omitempty"` // not in original relayer
	// Type             Uint256               `json:"type"`              // not in original relayer
}

type AccessListResponse struct{}

type LogResponse struct {
	Removed          bool    `json:"removed"`          // true when the log was removed, due to a chain reorganization. false if it's a valid log.
	LogIndex         Uint256 `json:"logIndex"`         // hexadecimal of the log index position in the block. null when its pending log.
	TransactionIndex Uint256 `json:"transactionIndex"` // hexadecimal of the transactions index position log was created from. null when its pending log.
	TransactionHash  H256    `json:"transactionHash"`  // 32 Bytes - hash of the transactions this log was created from. null when its pending log.
	BlockHash        H256    `json:"blockHash"`        // 32 Bytes - hash of the block where this log was in. null when it's pending. null when its pending log.
	BlockNumber      Uint256 `json:"blockNumber"`      // the block number where this log was in. null when it's pending. null when its pending log.
	Address          Address `json:"address"`          // 20 Bytes - address from which this log originated.
	Data             Bytea   `json:"data"`             // contains one or more 32 Bytes non-indexed arguments of the log.
	Topics           []Bytea `json:"topics"`           // Array of 0 to 4 32 Bytes of indexed log arguments. (In solidity: The first topic is the hash of the signature of the event (e.g. Deposit(address,bytes32,uint256)), except you declared the event with the anonymous specifier.)
}

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
		Number:           bl.Sequence,
		Difficulty:       IntToUint256(0),
		ExtraData:        Bytea("0"),
		GasLimit:         IntToUint256(0),
		GasUsed:          IntToUint256(0),
		Hash:             bl.Hash,
		LogsBloom:        []byte(fmt.Sprintf("%0x", 0)),
		Miner:            bl.Miner,
		Nonce:            []byte(fmt.Sprintf("%016x", 0)),
		ParentHash:       bl.ParentHash,
		ReceiptsRoot:     bl.ReceiptsRoot,
		Sha3Uncles:       HexStringToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
		Size:             bl.Size,
		StateRoot:        HexStringToHash(bl.StateRoot),
		Timestamp:        IntToUint256(bl.Timestamp),
		Transactions:     transactions,
		TransactionsRoot: bl.TransactionsRoot,
		Uncles:           []H256{},
	}
}

func (tx *Transaction) ToResponse() *TransactionResponse {
	return &TransactionResponse{
		Hash:             tx.Hash,
		BlockHash:        tx.BlockHash,
		From:             tx.From,
		To:               tx.To,
		Gas:              tx.GasUsed,
		GasPrice:         tx.GasPrice,
		Value:            tx.Value,
		Nonce:            tx.Nonce,
		TransactionIndex: UintToUint256(tx.TransactionIndex),
		V:                UintToUint256(tx.V),
		R:                common.Hex2Bytes(tx.R.String()),
		S:                common.Hex2Bytes(tx.S.String()),
	}
}

func (b Bytea) MarshalJSON() ([]byte, error) {
	if b == nil || len(b) == 0 {
		return []byte(`"0x0"`), nil
	}
	return []byte(fmt.Sprintf(`"0x%s"`, b)), nil
}

func (i *Uint256) Int64() int64 {
	return i.V.Int64()
}

func (i *Uint256) Uint64() uint64 {
	return i.V.Uint64()
}

func (i Uint256) String() string {
	return i.V.String()
}

func (i *Uint256) Bytes() []byte {
	return i.V.Bytes()
}

func (i Uint256) Add(x int64) *Uint256 {
	n := big.NewInt(0).Set(i.V)
	n.Add(n, big.NewInt(x))
	return &Uint256{n}
}

func (i *Uint256) FromBytes(b []byte) {
	i.V = big.NewInt(0).SetBytes(b)
}

func (i *Uint256) FromHexString(s string) error {
	if val, success := big.NewInt(0).SetString(s, 16); !success {
		return errors.New("failed to parse hexadecimal")
	} else {
		i.V = val
	}
	return nil
}

func (i Uint256) MarshalJSON() ([]byte, error) {
	if i.V == nil {
		return []byte(`"0x0"`), nil
	}
	return []byte(fmt.Sprintf(`"0x%s"`, i.V.Text(16))), nil
}

func (i *Uint256) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = strings.Trim(s, `"`)
	s = strings.TrimPrefix(s, "0x")
	return i.FromHexString(s)
}
