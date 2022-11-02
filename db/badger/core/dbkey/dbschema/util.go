package dbschema

func putBigEndian(dst []byte, value uint64) {
	for i := len(dst) - 1; i >= 0; i-- {
		dst[i] = byte(value)
		value >>= 8
	}
}

func readBigEndian(src []byte) uint64 {
	value := uint64(0)
	for _, c := range src {
		value = (value << 8) | uint64(c)
	}
	return value
}
