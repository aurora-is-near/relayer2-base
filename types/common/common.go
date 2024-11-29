package common

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/aurora-is-near/relayer2-base/types/utils"
	"github.com/holiman/uint256"
	jsoniter "github.com/json-iterator/go"
)

const (
	FilterIdByteSize = 16

	SafeBlockNumber      = BN64(-4)
	FinalizedBlockNumber = BN64(-3)
	PendingBlockNumber   = BN64(-2)
	LatestBlockNumber    = BN64(-1)
	EarliestBlockNumber  = BN64(0)
)

type Uinteger interface {
	~uint | ~uint16 | ~uint32 | ~uint64
}

type Integer interface {
	~int | ~int16 | ~int32 | ~int64
}

type Uint256 big.Int

type BN64 int64

type Uint64 struct{ uint64 }

type DataVec struct{ primitives.VarData }

type H256 struct{ primitives.Data32 }

type Address struct{ primitives.Data20 }

type BlockNumberOrHash struct {
	BlockNumber      *BN64 `json:"blockNumber,omitempty"`
	BlockHash        *H256 `json:"blockHash,omitempty"`
	RequireCanonical bool  `json:"requireCanonical,omitempty"`
}

func (h256 H256) String() string {
	return h256.Hex()
}

func (ui64 *Uint64) UnmarshalJSON(data []byte) error {
	input := strings.TrimSpace(string(data))
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		if len(input) >= 3 && input[1:3] == "0x" {
			value, err := strconv.ParseUint(input[3:len(input)-1], 10, 64)
			if err != nil {
				return err
			}
			if value > math.MaxInt64 {
				return fmt.Errorf("block number larger than int64")
			}
			ui64.uint64 = value
		} else {
			value, err := strconv.ParseUint(input[1:len(input)-1], 10, 64)
			if err != nil {
				return err
			}
			ui64.uint64 = value
		}
	} else {
		value, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			return err
		}
		ui64.uint64 = value
	}
	return nil
}

func (ui64 *Uint64) Uint64() uint64 {
	return ui64.uint64
}

func (bnh *BlockNumberOrHash) UnmarshalJSON(data []byte) error {
	type erased BlockNumberOrHash
	e := erased{}
	err := jsoniter.Unmarshal(data, &e)
	if err == nil {
		if e.BlockNumber != nil && e.BlockHash != nil {
			return fmt.Errorf("cannot specify both BlockHash and BlockNumber, choose one or the other")
		}
		bnh.BlockNumber = e.BlockNumber
		bnh.BlockHash = e.BlockHash
		bnh.RequireCanonical = e.RequireCanonical
		return nil
	}
	var input BN64
	err = jsoniter.Unmarshal(data, &input)
	if err != nil {
		return err
	}
	bnh.BlockNumber = &input
	return nil
}

func (bnh *BlockNumberOrHash) Number() (BN64, bool) {
	if bnh.BlockNumber != nil {
		return *bnh.BlockNumber, true
	}
	return BN64(0), false
}

func (bnh *BlockNumberOrHash) String() string {
	if bnh.BlockNumber != nil {
		return strconv.Itoa(int(*bnh.BlockNumber))
	}
	if bnh.BlockHash != nil {
		return bnh.BlockHash.String()
	}
	return "nil"
}

func (bnh *BlockNumberOrHash) Hash() (H256, bool) {
	if bnh.BlockHash != nil {
		return *bnh.BlockHash, true
	}
	return H256{}, false
}

func BlockNumberOrHashWithBN64(bn BN64) BlockNumberOrHash {
	return BlockNumberOrHash{
		BlockNumber:      &bn,
		BlockHash:        nil,
		RequireCanonical: false,
	}
}

func BlockNumberOrHashWithHash(hash H256, canonical bool) BlockNumberOrHash {
	return BlockNumberOrHash{
		BlockNumber:      nil,
		BlockHash:        &hash,
		RequireCanonical: canonical,
	}
}

func (n *BN64) UnmarshalJSON(data []byte) error {
	input := strings.TrimSpace(string(data))
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		input = strings.ToLower(input[1 : len(input)-1])
		value, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			switch input {
			case "earliest":
				*n = EarliestBlockNumber
			case "latest":
				*n = LatestBlockNumber
			case "pending":
				*n = PendingBlockNumber
			case "finalized":
				*n = FinalizedBlockNumber
			case "safe":
				*n = SafeBlockNumber
			default:
				ui64, err := utils.HexStringToUint64(input)
				if err != nil {
					return err
				}
				if ui64 > math.MaxInt64 {
					return fmt.Errorf("hex number [%s] is greater than max int64 [%d]", input, math.MaxInt64)
				}
				*n = BN64(ui64)
			}
		} else {
			*n = BN64(value)
		}
	} else {
		value, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			return err
		}
		*n = BN64(value)
	}
	return nil
}

// Uint64 returns block number as *uint64.
// 	If BN64 contains following tags; Safe, Finalized, Pending, Latest it returns nil
//  If BN64 contains the Earliest tag, it returns 0
func (n *BN64) Uint64() *uint64 {
	i := n.Int64()
	if *i >= 0 {
		ui := uint64(*i)
		return &ui
	}
	return nil
}

// Int64 returns block number as *int64. If BN64 contains following tags; Earliest, Latest, Pending, Finalized, Safe
// it returns; 0, -1, -2, -3, -4 respectively
func (n *BN64) Int64() *int64 {
	return (*int64)(n)
}

// MarshalText implements encoding.TextMarshaler
func (ui256 Uint256) MarshalText() ([]byte, error) {
	bigint := (big.Int)(ui256)
	if sign := bigint.Sign(); sign == 0 {
		return []byte("0x0"), nil
	} else if sign > 0 {
		return []byte("0x" + bigint.Text(16)), nil
	} else {
		return []byte("-0x" + bigint.Text(16)[1:]), nil
	}
}

func (ui256 *Uint256) UnmarshalJSON(data []byte) error {
	number := utils.SanitizeStringForNumber(string(data))
	i, err := uint256.FromHex(number)
	if err != nil {
		return err
	}
	*ui256 = Uint256(*i.ToBig())
	return nil
}

func (ui256 *Uint256) Bytes() []byte {
	return ((*big.Int)(ui256)).Bytes()
}

func (ui256 *Uint256) Uint64() uint64 {
	return ((*big.Int)(ui256)).Uint64()
}

func (ui256 *Uint256) Text(base int) string {
	return ((*big.Int)(ui256)).Text(base)
}

func (ui256 *Uint256) Cmp(test Uint256) int {
	return ((*big.Int)(ui256)).Cmp((*big.Int)(&test))
}

func (ui256 *Uint256) Data32() primitives.Data32 {
	return primitives.Data32FromBytes(ui256.Bytes())
}

func UintToUint64[T Uinteger](i T) Uint64 {
	return Uint64{uint64(i)}
}

func IntToUint64[T Integer](i T) Uint64 {
	return Uint64{uint64(i)}
}

// RandomUint256 returns a random Uint256 with a size of FilterIdByteSize. The 8 MSBs are BigEndian ordered Unix
// nanoseconds, the rest of the bytes are randomly generated by crypto/rand
func RandomUint256() Uint256 {
	t := make([]byte, FilterIdByteSize)
	binary.BigEndian.PutUint64(t[:8], uint64(time.Now().UnixNano()))
	_, _ = rand.Read(t[8:])
	big := *big.NewInt(0).SetBytes(t)
	return (Uint256)(big)
}

func IntToUint256[T Integer](i T) Uint256 {
	big := *big.NewInt(0).SetInt64(int64(i))
	return (Uint256)(big)
}

func UintToUint256[T Uinteger](i T) Uint256 {
	big := *big.NewInt(0).SetUint64(uint64(i))
	return (Uint256)(big)
}

func IntToBN64[T Integer](i T) BN64 {
	return BN64(i)
}

func UintToBN64[T Uinteger](i T) BN64 {
	return BN64(i)
}

func BytesToAddress(b []byte) Address {
	return Address{primitives.Data20FromBytes(b)}
}

func HexStringToAddress(s string) (Address, error) {
	data, err := primitives.Data20FromHex(s)
	if err != nil {
		return Address{}, err
	}

	return Address{data}, nil
}

func MustHexStringToAddress(s string) Address {
	return Address{primitives.MustData20FromHex(s)}
}

func HexStringToDataVec(s string) (DataVec, error) {
	data, err := hex.DecodeString(s[2:])
	if err != nil {
		return DataVec{}, err
	}

	return DataVec{primitives.VarDataFromBytes(data)}, nil
}

func MustHexStringToDataVec(s string) DataVec {
	dataVec, err := HexStringToDataVec(s)
	if err != nil {
		panic(err)
	}

	return dataVec
}

func HexStringToHash(s string) (H256, error) {
	data, err := primitives.Data32FromHex(s)
	if err != nil {
		return H256{}, err
	}

	return H256{data}, nil
}

func MustHexStringToHash(s string) H256 {
	return H256{primitives.MustData32FromHex(s)}
}

func Uint256FromBytes(b []byte) Uint256 {
	big := *big.NewInt(0).SetBytes(b)
	return (Uint256)(big)
}

func Uint256FromHex(s string) (*Uint256, error) {
	big, err := utils.HexStringToBigInt(s)
	if err != nil {
		return nil, err
	}
	return (*Uint256)(big), nil
}
