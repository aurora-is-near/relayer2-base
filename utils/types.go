package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/fxamacker/cbor/v2"

	"github.com/ethereum/go-ethereum/common"
)

type Uint256 struct{ *big.Int }
type H256 struct{ common.Hash }
type Address struct{ common.Address }
type Addresses []Address
type Bytea []byte
type Topics [][][]byte
type Int32Key struct{ key [4]byte }
type TxnData string

type Block struct {
	ChainId          uint64         `cbor:"chain_id"`
	Hash             H256           `cbor:"hash"`
	ParentHash       H256           `cbor:"parent_hash"`
	Height           uint64         `cbor:"height"`
	Miner            Address        `cbor:"miner"`
	Timestamp        int64          `cbor:"timestamp"`
	GasLimit         Uint256        `cbor:"gas_limit"`
	GasUsed          Uint256        `cbor:"gas_used"`
	LogsBloom        string         `cbor:"logs_bloom"`
	TransactionsRoot H256           `cbor:"transactions_root"`
	ReceiptsRoot     H256           `cbor:"receipts_root"`
	Transactions     []*Transaction `cbor:"transactions"`
	NearBlock        any            `cbor:"near_metadata"`
	StateRoot        string         `cbor:"state_root"`
	Size             Uint256        `cbor:"size"`
	Sequence         uint64         `cbor:"sequence"`
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
	Logs                 []*Log          `cbor:"logs"`
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

type LogSubscriptionOptions struct {
	Address Addresses `json:"address"`
	Topics  Topics    `json:"topics"`
}

type FilterOptions struct {
	Address   Addresses `json:"address"`
	FromBlock string    `json:"fromBlock"`
	ToBlock   string    `json:"toBlock"`
	Topics    Topics    `json:"topics"`
	BlockHash *H256     `json:"blockhash"`
}

type LogFilter struct {
	Address   [][]byte `json:"address"`
	FromBlock *Uint256 `json:"fromBlock"`
	ToBlock   *Uint256 `json:"toBlock"`
	// Blockhash *H256      `json:"blockhash"` // Blockhash is converted to FromBlock = ToBlock on eth_getLogs call
	Topics [][][]byte `json:"topics"`
}

type StoredFilter struct {
	Type      string
	CreatedBy string
	PollBlock Uint256
	// BlockHash *H256
	FromBlock *Uint256
	ToBlock   *Uint256
	Addresses [][]byte
	Topics    [][][]byte
}

type NearTransaction struct {
	Hash        H256 `cbor:"hash"`
	ReceiptHash H256 `cbor:"receipt_hash"`
}

//easyjson:json
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
	MixHash          H256          `json:"mixHash"`
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

//easyjson:json
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

//easyjson:json
type TransactionReceiptResponse struct {
	BlockHash         H256           `json:"blockHash"`         // 32 Bytes - hash of the block including this transaction.
	BlockNumber       Uint256        `json:"blockNumber"`       // block number including this transaction.
	ContractAddress   *Address       `json:"contractAddress"`   // 20 Bytes - the contract address created, if the transaction was a contract creation, otherwise null.
	CumulativeGasUsed Uint256        `json:"cumulativeGasUsed"` // the total amount of gas used when this transaction was executed in the block.
	EffectiveGasPrice Uint256        `json:"effectiveGasPrice"` // the actual value per gas deducted from the sender's account. Before EIP-1559, equal to the gas price.
	From              Address        `json:"from"`              // 20 Bytes - address of the sender.
	To                *Address       `json:"to"`                // 20 Bytes - address of the receiver. null when the transaction is a contract creation transaction.
	GasUsed           Uint256        `json:"gasUsed"`           // the amount of gas used by this specific transaction alone.
	Logs              []*LogResponse `json:"logs"`              // Array - Array of log objects, which this transaction generated.
	LogsBloom         Bytea          `json:"logsBloom"`         // 256 Bytes - Bloom filter for light clients to quickly retrieve related logs.
	Status            Uint256        `json:"status"`            // either 1 (success) or 0 (failure)
	TransactionHash   H256           `json:"transactionHash"`   // 32 Bytes - hash of the transaction.
	TransactionIndex  Uint256        `json:"transactionIndex"`  // hexadecimal of the transaction's index position in the block.
	Type              Uint256        `json:"type"`              // the transaction type.
	NearHash          H256           `json:"nearTransactionHash"`
	NearReceiptHash   H256           `json:"nearReceiptHash"`
}

type AccessListResponse struct{}

//easyjson:json
type LogResponse struct {
	Removed          bool    `json:"removed"`          // true when the log was removed, due to a chain reorganization. false if it's a valid log.
	LogIndex         Uint256 `json:"logIndex"`         // hexadecimal of the log index position in the block. null when its pending log.
	TransactionIndex Uint256 `json:"transactionIndex"` // hexadecimal of the transactions index position log was created from. null when its pending log.
	TransactionHash  H256    `json:"transactionHash"`  // 32 Bytes - hash of the transactions this log was created from. null when its pending log.
	BlockHash        H256    `json:"blockHash"`        // 32 Bytes - hash of the block where this log was in. null when its pending. null when its pending log.
	BlockNumber      Uint256 `json:"blockNumber"`      // the block number where this log was in. null when its pending. null when its pending log.
	Address          Address `json:"address"`          // 20 Bytes - address from which this log originated.
	Data             Bytea   `json:"data"`             // contains one or more 32 Bytes non-indexed arguments of the log.
	Topics           []Bytea `json:"topics"`           // Array of 0 to 4 32 Bytes of indexed log arguments. (In solidity: The first topic is the hash of the signature of the event (e.g. Deposit(address,bytes32,uint256)), except you declared the event with the anonymous specifier.)
}

func (bl *Block) TxCount() int64 {
	return int64(len(bl.Transactions))
}

func (bl *Block) ToResponse(fullTx bool) *BlockResponse {
	transactions := make([]interface{}, 0, len(bl.Transactions))
	for _, tx := range bl.Transactions {
		if !fullTx {
			transactions = append(transactions, tx.Hash)
		} else {
			transactions = append(transactions, tx.ToResponse())
		}
	}
	return &BlockResponse{
		Number:           UintToUint256(bl.Height),
		Difficulty:       IntToUint256(0),
		ExtraData:        Bytea("0"),
		GasLimit:         bl.GasLimit,
		GasUsed:          bl.GasUsed,
		Hash:             bl.Hash,
		LogsBloom:        Bytea(bl.LogsBloom),
		Miner:            bl.Miner,
		Nonce:            []byte(fmt.Sprintf("%016x", 0)),
		ParentHash:       bl.ParentHash,
		MixHash:          HexStringToHash("0"),
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

func (tx *Transaction) ToReceiptResponse() *TransactionReceiptResponse {
	var status Uint256
	if tx.Status {
		status = IntToUint256(1)
	} else {
		status = IntToUint256(0)
	}

	var contractAddress *Address
	if !tx.ContractAddress.IsZero() {
		contractAddress = &tx.ContractAddress
	}

	var logsBloom = Bytea(fmt.Sprintf("%0256x", 0))

	return &TransactionReceiptResponse{
		BlockHash:         tx.BlockHash,
		BlockNumber:       UintToUint256(tx.BlockHeight),
		ContractAddress:   contractAddress,
		CumulativeGasUsed: IntToUint256(0),
		EffectiveGasPrice: tx.GasPrice,
		From:              tx.From,
		To:                tx.To,
		GasUsed:           tx.GasUsed,
		LogsBloom:         logsBloom,
		Status:            status,
		TransactionHash:   tx.Hash,
		TransactionIndex:  UintToUint256(tx.TransactionIndex),
		Type:              IntToUint256(0),
		NearHash:          tx.NearTransaction.Hash,
		NearReceiptHash:   tx.NearTransaction.ReceiptHash,
	}
}

type EstimateGasRequest struct {
	From                 string `json:"from"`
	To                   string `json:"to"`
	Gas                  string `json:"gas"`
	GasPrice             string `json:"gasPrice"`
	MaxPriorityFeePerGas int64  `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         int64  `json:"maxFeePerGas"`
	Value                string `json:"value"`
	Data                 string `json:"data"`
}

func (a *Address) IsZero() bool {
	return a.Hash().Big().Sign() == 0
}

func (a *Addresses) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return fmt.Errorf("no bytes")
	}

	switch b[0] {
	case '"':
		var addr Address
		err := json.Unmarshal(b, &addr)
		if err != nil {
			return nil
		}
		*a = Addresses{addr}
	case '[':
		addrs := make([]Address, 0)
		err := json.Unmarshal(b, &addrs)
		if err != nil {
			return nil
		}
		*a = Addresses(addrs)
	}

	return nil
}

func (t *Topics) UnmarshalJSON(b []byte) error {
	tps := [4]interface{}{}
	err := json.Unmarshal(b, &tps)
	if err != nil {
		return err
	}
	results := Topics{{}, {}, {}, {}}
	for i, t := range tps {
		switch v := t.(type) {
		case string:
			results[i] = append(results[i], []byte(v))
		case []interface{}:
			for _, topic := range v {
				if topic, ok := topic.(string); ok {
					results[i] = append(results[i], []byte(topic))
				}
			}
		case nil:
		default:
		}
	}
	*t = results
	return nil
}

func (b Bytea) MarshalJSON() ([]byte, error) {
	if b == nil || len(b) == 0 {
		return []byte(`"0x0"`), nil
	}
	return []byte(fmt.Sprintf(`"0x%s"`, b)), nil
}

func (b *Bytea) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*b = Bytea(s[2:len(s)])
	return nil
}

func (h H256) KeyBytes() []byte {
	return h.Bytes()
}

func (i Uint256) Add(x int64) *Uint256 {
	n := big.NewInt(0).Set(i.Int)
	n.Add(n, big.NewInt(x))
	return &Uint256{n}
}

func (i Uint256) Bytes() []byte {
	if i.Int == nil {
		return []byte{}
	}
	return i.Int.Bytes()
}

func (i *Uint256) FromHexString(s string) error {
	s = strings.TrimPrefix(s, "0x")
	if val, success := big.NewInt(0).SetString(s, 16); !success {
		return errors.New("failed to parse hexadecimal")
	} else {
		*i = Uint256{val}
	}
	return nil
}

func (i Uint256) KeyBytes() []byte {
	u := i.Uint64()
	b := new(bytes.Buffer)
	_ = binary.Write(b, binary.BigEndian, u)
	return b.Bytes()
}

func (i Uint256) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(i.Int)
}

func (i *Uint256) UnmarshalCBOR(b []byte) error {
	val := big.NewInt(0)
	if err := cbor.Unmarshal(b, &val); err != nil {
		return err
	}
	*i = Uint256{val}
	return nil
}

func (i Uint256) MarshalJSON() ([]byte, error) {
	if i.Int == nil {
		return []byte(`"0x0"`), nil
	}
	return []byte(fmt.Sprintf(`"0x%s"`, i.Text(16))), nil
}

func (i *Uint256) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = strings.Trim(s, `"`)
	return i.FromHexString(s)
}

func (i *Uint256) SetBytes(b []byte) {
	if i.Int == nil {
		i.Int = big.NewInt(0)
	}
	i.Int.SetBytes(b)
}

func (i Uint256) ToUint32Key() (*Int32Key, error) {
	b := i.Bytes()
	if len(b) > 4 {
		return nil, errors.New(fmt.Sprintf("u256 doesn't fit in a u32: 0x%s", i.Text(16)))
	} else if len(b) < 4 {
		b = append([]byte{0, 0, 0, 0}, b...)
	}
	return &Int32Key{*(*[4]byte)(b[len(b)-4:])}, nil
}

func (i Int32Key) KeyBytes() []byte {
	return i.key[:]
}
