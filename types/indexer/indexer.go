package indexer

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aurora-is-near/relayer2-base/db/codec"
	"github.com/aurora-is-near/relayer2-base/tinypack"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/aurora-is-near/relayer2-base/types/utils"
	"github.com/btcsuite/btcutil/base58"
	jsoniter "github.com/json-iterator/go"
)

type NearHash primitives.Data32
type InputOutputData primitives.VarData
type Topic primitives.Data32
type Size uint64
type Timestamp uint64

type Block struct {
	ChainId          uint64              `cbor:"chain_id"          json:"chain_id"`
	Height           uint64              `cbor:"height"            json:"height"`
	Sequence         uint64              `cbor:"sequence"          json:"sequence"`
	GasLimit         primitives.Quantity `cbor:"gas_limit"         json:"gas_limit"`
	GasUsed          primitives.Quantity `cbor:"gas_used"          json:"gas_used"`
	Timestamp        Timestamp           `cbor:"timestamp"         json:"timestamp"`
	Hash             primitives.Data32   `cbor:"hash"              json:"hash"`
	ParentHash       primitives.Data32   `cbor:"parent_hash"       json:"parent_hash"`
	TransactionsRoot primitives.Data32   `cbor:"transactions_root" json:"transactions_root"`
	ReceiptsRoot     primitives.Data32   `cbor:"receipts_root"     json:"receipts_root"`
	StateRoot        primitives.Data32   `cbor:"state_root"        json:"state_root"`
	Size             Size                `cbor:"size"              json:"size"`
	Miner            primitives.Data20   `cbor:"miner"             json:"miner"`
	LogsBloom        primitives.Data256  `cbor:"logs_bloom"        json:"logs_bloom"`
	Transactions     []*Transaction      `cbor:"transactions"      json:"transactions"`
	NearBlock        any                 `cbor:"near_metadata"     json:"near_metadata"`
}

type Transaction struct {
	Hash                 primitives.Data32   `cbor:"hash"                     json:"hash"`
	BlockHash            primitives.Data32   `cbor:"block_hash"               json:"block_hash"`
	BlockHeight          uint64              `cbor:"block_height"             json:"block_height"`
	ChainId              uint64              `cbor:"chain_id"                 json:"chain_id"`
	TransactionIndex     uint64              `cbor:"transaction_index"        json:"transaction_index"`
	From                 primitives.Data20   `cbor:"from"                     json:"from"`
	To                   *primitives.Data20  `cbor:"to"                       json:"to"`
	Nonce                primitives.Quantity `cbor:"nonce"                    json:"nonce"`
	GasPrice             primitives.Quantity `cbor:"gas_price"                json:"gas_price"`
	GasLimit             primitives.Quantity `cbor:"gas_limit"                json:"gas_limit"`
	GasUsed              uint64              `cbor:"gas_used"                 json:"gas_used"`
	MaxPriorityFeePerGas primitives.Quantity `cbor:"max_priority_fee_per_gas" json:"max_priority_fee_per_gas"`
	MaxFeePerGas         primitives.Quantity `cbor:"max_fee_per_gas"          json:"max_fee_per_gas"`
	Value                primitives.Quantity `cbor:"value"                    json:"value"`
	Input                InputOutputData     `cbor:"input"                    json:"input"`
	Output               InputOutputData     `cbor:"output"                   json:"output"`
	AccessList           []AccessList        `cbor:"access_list"              json:"access_list"`
	TxType               uint64              `cbor:"tx_type"                  json:"tx_type"`
	Status               bool                `cbor:"status"                   json:"status"`
	Logs                 []*Log              `cbor:"logs"                     json:"logs"`
	LogsBloom            primitives.Data256  `cbor:"logs_bloom"               json:"logs_bloom"`
	ContractAddress      *primitives.Data20  `cbor:"contract_address"         json:"contract_address"`
	V                    uint64              `cbor:"v"                        json:"v"`
	R                    primitives.Quantity `cbor:"r"                        json:"r"`
	S                    primitives.Quantity `cbor:"s"                        json:"s"`
	NearTransaction      NearTransaction     `cbor:"near_metadata"            json:"near_metadata"`
}

type AccessList struct {
	Address     primitives.Data20   `json:"address"     json:"address"`
	StorageKeys []primitives.Data32 `json:"storageKeys" json:"storageKeys"`
}

type NearTransaction struct {
	Hash        *NearHash `cbor:"transaction_hash" json:"transaction_hash"`
	ReceiptHash NearHash  `cbor:"receipt_hash"     json:"receipt_hash"`
}

type Log struct {
	Address primitives.Data20 `cbor:"Address" json:"address"`
	Topics  []Topic           `cbor:"Topics"  json:"topics"`
	Data    InputOutputData   `cbor:"data"    json:"data"`
}

func (iod *InputOutputData) UnmarshalJSON(b []byte) error {
	var i []byte
	err := jsoniter.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	p, err := byteArrayToData[primitives.VarLen](i)
	if err != nil {
		return err
	}
	*iod = InputOutputData(*p)
	return nil
}

func (iod *InputOutputData) UnmarshalCBOR(b []byte) error {
	var i []byte
	err := codec.CborDecoder().Unmarshal(b, &i)
	if err != nil {
		return err
	}
	p, err := byteArrayToData[primitives.VarLen](i)
	if err != nil {
		return err
	}
	*iod = InputOutputData(*p)
	return nil
}

func (t Topic) MarshalJSON() ([]byte, error) {
	return primitives.Data32.MarshalJSON(primitives.Data[primitives.Len32](t))
}

func (t *Topic) UnmarshalJSON(b []byte) error {
	var i []byte
	err := jsoniter.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	p, err := byteArrayToData[primitives.Len32](i)
	if err != nil {
		return err
	}
	*t = Topic(*p)
	return nil
}

func (t *Topic) UnmarshalCBOR(b []byte) error {
	var i []byte
	err := codec.CborDecoder().Unmarshal(b, &i)
	if err != nil {
		return err
	}
	p, err := byteArrayToData[primitives.Len32](i)
	if err != nil {
		return err
	}
	*t = Topic(*p)
	return nil
}

func (s *Size) UnmarshalJSON(b []byte) error {
	var in string
	err := jsoniter.Unmarshal(b, &in)

	ui64, err := utils.HexStringToUint64(in)
	if err != nil {
		return err
	}
	*s = Size(ui64)
	return nil
}

func (s *Size) UnmarshalCBOR(b []byte) error {
	var in string
	err := codec.CborDecoder().Unmarshal(b, &in)
	if err != nil {
		return err
	}

	ui64, err := utils.HexStringToUint64(in)
	if err != nil {
		return err
	}
	*s = Size(ui64)
	return nil
}

func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	var err error
	in := fmt.Sprintf("%s", b)
	*ts, err = timeStringToTimestamp(in)
	return err
}

func (ts *Timestamp) UnmarshalCBOR(b []byte) error {
	var ui64 uint64
	err := codec.CborDecoder().Unmarshal(b, &ui64)
	if err != nil {
		return err
	}
	in := fmt.Sprintf("%d", ui64)
	*ts, err = timeStringToTimestamp(in)
	return err
}

func (nh *NearHash) UnmarshalJSON(b []byte) error {
	var in string
	err := jsoniter.Unmarshal(b, &in)
	if err != nil {
		return err
	}
	*nh, err = base58StringToNearHash(in)
	return err
}

func (nh *NearHash) UnmarshalCBOR(b []byte) error {
	var in string
	err := codec.CborDecoder().Unmarshal(b, &in)
	if err != nil {
		return err
	}
	*nh, err = base58StringToNearHash(in)
	return err
}

func timeStringToTimestamp(time string) (Timestamp, error) {
	if len(time) > 10 {
		time = time[:10]
	}
	t, err := strconv.ParseUint(time, 10, 0)
	if err != nil {
		return Timestamp(0), err
	}
	return Timestamp(t), nil
}

func byteArrayToData[LD tinypack.LengthDescriptor](in []byte) (*primitives.Data[LD], error) {
	var sb strings.Builder
	_, err := sb.WriteString("0x")
	if err != nil {
		return nil, err
	}
	if len(in) > 0 {
		s := hex.EncodeToString(in)
		_, err = sb.WriteString(s)
		if err != nil {
			return nil, err
		}
	}
	p := primitives.DataFromHex[LD](sb.String())
	return &p, nil
}

func base58StringToNearHash(in string) (NearHash, error) {
	nearHash := base58.Decode(in)
	if len(nearHash) > 32 {
		return NearHash{}, errors.New("length of [" + in + "] exceeds 32 bytes")
	}
	return NearHash(primitives.Data32FromBytes(nearHash)), nil
}
