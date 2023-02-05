package primitives

var hextable = "0123456789abcdef"

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
