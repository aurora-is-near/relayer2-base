package indexer

import (
	"aurora-relayer-go-common/types/common"
	"encoding/hex"
	"encoding/json"
	"fmt"
	gethcom "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"strconv"
	"strings"
)

type Uint32 uint32
type NearHash string

type Bloom struct{ types.Bloom }
type Timestamp struct{ hexutil.Uint }
type GasUsed struct{ hexutil.Big }
type GasLimit struct{ hexutil.Big }
type InputOutput struct{ hexutil.Bytes }
type Topic struct{ hexutil.Bytes }
type Data struct{ hexutil.Bytes }
type V struct{ hexutil.Big }
type Miner struct{ gethcom.Address }

type Block struct {
	ChainId          common.Uint64  `cbor:"chain_id"          json:"chain_id"`
	Height           common.Uint64  `cbor:"height"            json:"height"`
	Sequence         common.Uint64  `cbor:"sequence"          json:"sequence"`
	GasLimit         GasLimit       `cbor:"gas_limit"         json:"gas_limit"`
	GasUsed          common.Uint256 `cbor:"gas_used"          json:"gas_used"`
	Timestamp        Timestamp      `cbor:"timestamp"         json:"timestamp"`
	Hash             common.H256    `cbor:"hash"              json:"hash"`
	ParentHash       common.H256    `cbor:"parent_hash"       json:"parent_hash"`
	TransactionsRoot common.H256    `cbor:"transactions_root" json:"transactions_root"`
	ReceiptsRoot     common.H256    `cbor:"receipts_root"     json:"receipts_root"`
	StateRoot        common.H256    `cbor:"state_root"        json:"state_root"`
	Size             common.Uint256 `cbor:"size"              json:"size"`
	Miner            Miner          `cbor:"miner"             json:"miner"`
	LogsBloom        Bloom          `cbor:"logs_bloom"        json:"logs_bloom"`
	Transactions     []*Transaction `cbor:"transactions"      json:"transactions"`
	NearBlock        any            `cbor:"near_metadata"     json:"near_metadata"`
}

type Transaction struct {
	Hash                 common.H256     `cbor:"hash"                     json:"hash"`
	BlockHash            common.H256     `cbor:"block_hash"               json:"block_hash"`
	BlockHeight          common.Uint64   `cbor:"block_height"             json:"block_height"`
	ChainId              common.Uint64   `cbor:"chain_id"                 json:"chain_id"`
	TransactionIndex     common.Uint64   `cbor:"transaction_index"        json:"transaction_index"`
	From                 common.Address  `cbor:"from"                     json:"from"`
	To                   *common.Address `cbor:"to"                       json:"to"`
	Nonce                common.Uint256  `cbor:"nonce"                    json:"nonce"`
	GasPrice             common.Uint256  `cbor:"gas_price"                json:"gas_price"`
	GasLimit             common.Uint256  `cbor:"gas_limit"                json:"gas_limit"`
	GasUsed              GasUsed         `cbor:"gas_used"                 json:"gas_used"`
	MaxPriorityFeePerGas common.Uint256  `cbor:"max_priority_fee_per_gas" json:"max_priority_fee_per_gas"`
	MaxFeePerGas         common.Uint256  `cbor:"max_fee_per_gas"          json:"max_fee_per_gas"`
	Value                common.Uint256  `cbor:"value"                    json:"value"`
	Input                InputOutput     `cbor:"input"                    json:"input"`
	Output               InputOutput     `cbor:"output"                   json:"output"`
	AccessList           []AccessList    `cbor:"access_list"              json:"access_list"`
	TxType               common.Uint64   `cbor:"tx_type"                  json:"tx_type"`
	Status               bool            `cbor:"status"                   json:"status"`
	Logs                 []*Log          `cbor:"logs"                     json:"logs"`
	LogsBloom            Bloom           `cbor:"logs_bloom"               json:"logs_bloom"`
	ContractAddress      *common.Address `cbor:"contract_address"         json:"contract_address"`
	V                    V               `cbor:"v"                        json:"v"`
	R                    common.Uint256  `cbor:"r"                        json:"r"`
	S                    common.Uint256  `cbor:"s"                        json:"s"`
	NearTransaction      NearTransaction `cbor:"near_metadata"            json:"near_metadata"`
}

type AccessList struct {
	Address     common.Address `json:"address"     json:"address"`
	StorageKeys []common.H256  `json:"storageKeys" json:"storageKeys"`
}

type NearTransaction struct {
	Hash        NearHash `cbor:"hash"         json:"hash"`
	ReceiptHash NearHash `cbor:"receipt_hash" json:"receipt_hash"`
}

type Log struct {
	Address common.Address `cbor:"Address" json:"address"`
	Topics  []Topic        `cbor:"Topics"  json:"topics"`
	Data    Data           `cbor:"data"    json:"data"`
}

func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	s := fmt.Sprintf("%s", b)
	if len(s) > 10 {
		s = s[:10]
	}
	i, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return err
	}
	err = ts.Uint.UnmarshalText([]byte(fmt.Sprintf("0x%x", i)))
	if err != nil {
		return err
	}
	return nil
}

func (gu *GasUsed) UnmarshalJSON(b []byte) error {
	var i uint64
	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	bi := big.NewInt(0).SetUint64(i)
	err = gu.Big.UnmarshalText([]byte(hexutil.EncodeBig(bi)))
	if err != nil {
		return err
	}
	return nil
}

func (gl *GasLimit) UnmarshalJSON(_ []byte) error {
	err := gl.Big.UnmarshalText([]byte("0xfffffffffffff"))
	if err != nil {
		return err
	}
	return nil
}

func (m *Miner) UnmarshalJSON(_ []byte) error {
	err := m.Address.UnmarshalText([]byte("0x0000000000000000000000000000000000000000"))
	if err != nil {
		return err
	}
	return nil
}

func (io *InputOutput) UnmarshalJSON(b []byte) error {
	var i []int
	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}

	var sb strings.Builder
	_, err = sb.WriteString("0x")
	if err != nil {
		return err
	}

	if len(i) > 0 {
		tmp := make([]byte, len(i))
		for j, v := range i {
			tmp[j] = byte(v)
		}
		s := hex.EncodeToString(tmp)
		_, err = sb.WriteString(s)
		if err != nil {
			return err
		}
	}

	err = io.Bytes.UnmarshalText([]byte(sb.String()))
	if err != nil {
		return err
	}
	return nil
}

func (t *Topic) UnmarshalJSON(b []byte) error {
	var i []int
	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}

	var sb strings.Builder
	_, err = sb.WriteString("0x")
	if err != nil {
		return err
	}

	if len(i) > 0 {
		tmp := make([]byte, len(i))
		for j, v := range i {
			tmp[j] = byte(v)
		}
		s := hex.EncodeToString(tmp)
		_, err = sb.WriteString(s)
		if err != nil {
			return err
		}
	}

	err = t.Bytes.UnmarshalText([]byte(sb.String()))
	if err != nil {
		return err
	}
	return nil
}

func (d *Data) UnmarshalJSON(b []byte) error {
	var i []int
	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}

	var sb strings.Builder
	_, err = sb.WriteString("0x")
	if err != nil {
		return err
	}

	if len(i) > 0 {
		tmp := make([]byte, len(i))
		for j, v := range i {
			tmp[j] = byte(v)
		}
		s := hex.EncodeToString(tmp)
		_, err = sb.WriteString(s)
		if err != nil {
			return err
		}
	}

	err = d.Bytes.UnmarshalText([]byte(sb.String()))
	if err != nil {
		return err
	}
	return nil
}

func (v *V) UnmarshalJSON(b []byte) error {
	var i uint64
	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	bi := big.NewInt(0).SetUint64(i)
	err = v.Big.UnmarshalText([]byte(hexutil.EncodeBig(bi)))
	if err != nil {
		return err
	}
	return nil
}

func (b Bloom) String() string {
	s, _ := b.MarshalText()
	return fmt.Sprintf("%v", s)
}
