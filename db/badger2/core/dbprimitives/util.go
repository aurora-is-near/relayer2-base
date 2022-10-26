package dbprimitives

import "encoding/hex"

var zeroChar = "0"[0]

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

func bytesToHex(b []byte, trimLeadingZeroes bool) string {
	hex := hex.EncodeToString(b)
	i := 0
	if trimLeadingZeroes {
		for i < len(hex)-1 && hex[i] == zeroChar {
			i++
		}
	}
	return "0x" + hex[i:]
}
