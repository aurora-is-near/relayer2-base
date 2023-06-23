package utils

import (
	"errors"
	"math/big"
	"strings"

	"github.com/holiman/uint256"
)

func HexStringToBigInt(input string) (*big.Int, error) {
	number := SanitizeStringForNumber(input)
	i, err := uint256.FromHex(number)
	if err != nil {
		return nil, err
	}
	return i.ToBig(), nil
}

func HexStringToUint64(input string) (uint64, error) {
	var out uint64
	number := SanitizeStringForNumber(input)
	i, err := uint256.FromHex(number)
	if err != nil {
		return out, err
	}
	if !i.IsUint64() {
		return out, errors.New("hex number is larger than 64 bits")
	}
	return i.Uint64(), nil
}

// SanitizeStringForNumber removes '"', 'any leading redundant 0 before number part' and 'space' from the input string,
// returns a string with hex prefix. See outputs for sample inputs:
//
// 0x0000a -> 0xa
//
// 0x0000a0 -> 0xa0
//
// " 0x0000a " -> 0xa
//
// 0x000 -> 0x0
//
// 0x0 -> 0x0
func SanitizeStringForNumber(in string) string {
	rawString := strings.Trim(in, "\"")
	rawString = strings.TrimSpace(rawString)
	numPart := strings.TrimPrefix(rawString, "0x")
	if len(numPart) == 0 {
		return "0x"
	}
	trimmed := strings.TrimLeft(numPart, "0")
	if len(trimmed) == 0 {
		trimmed = "0"
	}
	return "0x" + trimmed
}
