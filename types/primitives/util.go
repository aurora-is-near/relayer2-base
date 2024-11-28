package primitives

import (
	"encoding/hex"
)

var hextable = "0123456789abcdef"

const nonHexMarker = 100

var reverseHexTable [256]byte

func init() {
	// Pre-fill the reverseHexTable
	for i := 0; i < len(reverseHexTable); i++ {
		switch {
		case i >= '0' && i <= '9':
			reverseHexTable[i] = byte(i) - '0'
		case i >= 'a' && i <= 'f':
			reverseHexTable[i] = byte(i) - 'a' + 10
		case i >= 'A' && i <= 'F':
			reverseHexTable[i] = byte(i) - 'A' + 10
		default:
			reverseHexTable[i] = nonHexMarker
		}
	}
}

func alignBytes(b []byte, length int, bigEndian bool) []byte {
	if length < 0 {
		if b == nil {
			return []byte{}
		} else {
			return b
		}
	}
	if len(b) < length {
		if bigEndian {
			b = append(make([]byte, length-len(b)), b...)
		} else {
			b = append(b, make([]byte, length-len(b))...)
		}
	}
	if len(b) > length {
		if bigEndian {
			b = b[len(b)-length:]
		} else {
			b = b[:length]
		}
	}
	return b
}

func byteToHex(v byte) (byte, byte) {
	return hextable[v>>4], hextable[v&0x0f]
}

func writeDataHex(dst []byte, b []byte) []byte {
	dst = append(dst, '0', 'x')
	for _, v := range b {
		l, r := byteToHex(v)
		dst = append(dst, l, r)
	}
	return dst
}

func writeQuantityHex(dst []byte, b []byte) []byte {
	dst = append(dst, '0', 'x')
	i := 0

	for ; i < len(b) && b[i] == 0; i++ {
	}
	if i == len(b) {
		return append(dst, '0')
	}

	l, r := byteToHex(b[i])
	if l != '0' {
		dst = append(dst, l)
	}
	dst = append(dst, r)

	for i++; i < len(b); i++ {
		l, r := byteToHex(b[i])
		dst = append(dst, l, r)
	}

	return dst
}

func hexToBytes(s string) ([]byte, error) {
	if len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		s = s[2:]
	}

	if len(s)%2 == 0 {
		return hex.DecodeString(s)
	}

	dst := make([]byte, hex.DecodedLen(1+len(s)))

	if reverseHexTable[s[0]] == nonHexMarker {
		return nil, hex.InvalidByteError(s[0])
	}
	dst[0] = reverseHexTable[s[0]]

	// Allocation-free string -> []byte conversion there (compiler will optimize it out)
	if _, err := hex.Decode(dst[1:], []byte(s[1:])); err != nil {
		return nil, err
	}

	return dst, nil
}
