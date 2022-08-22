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
