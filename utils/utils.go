package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Integer interface {
	~int | ~int16 | ~int32 | ~int64
}

type Uinteger interface {
	~uint | ~uint16 | ~uint32 | ~uint64
}

func ParseHexStringToAddress(s string) (*Address, error) {
	var address common.Address
	err := address.UnmarshalText([]byte(s))
	if err != nil {
		return nil, err
	}
	return &Address{address}, nil
}

func HexStringToAddress(s string) Address {
	return Address{common.HexToAddress(s)}
}

func HexStringToHash(s string) H256 {
	return H256{common.HexToHash(s)}
}

func IntToHex[T Integer](i T) string {
	return fmt.Sprintf("0x%x", i)
}

func IntToUint256[T Integer](i T) Uint256 {
	return Uint256{big.NewInt(0).SetInt64(int64(i))}
}

func UintToUint256[T Uinteger](i T) Uint256 {
	return Uint256{big.NewInt(0).SetUint64(uint64(i))}
}

func RandomUint256() (*Uint256, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	n := big.NewInt(0).SetBytes(b)
	return &Uint256{n}, nil
}

// Parses the block argument and returns
//   - nil for pending and latest block or zero address for earliest block
//   - validated block number or the error thrown
func ParseBlockArgument(block string) (*Uint256, error) {
	switch block {
	case "", "pending", "latest":
		return nil, nil
	case "earliest":
		zero := IntToUint256(0)
		return &zero, nil // TODO - genesis or zero?
	default:
		val := IntToUint256(0)
		err := val.FromHexString(block)
		if err != nil {
			return nil, err
		}
		return &val, nil
	}
}
