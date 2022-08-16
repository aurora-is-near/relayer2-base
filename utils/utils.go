package utils

import (
	"fmt"
	"strings"
)

type Integer interface {
	~int | ~int16 | ~int32 | ~int64
}

func IntToHex[T Integer](i T) string {
	return fmt.Sprintf("0x%x", i)
}

func IntToUint256[T Integer](i T) Uint256 {
	return Uint256(fmt.Sprintf("%x", i))
}

func ParseHexString[T ~string](s string) T {
	return T(strings.TrimPrefix(s, "0x"))
}
