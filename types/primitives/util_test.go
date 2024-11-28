package primitives

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlignBytes(t *testing.T) {
	assert.Equal(t, []byte{0, 0, 0}, alignBytes(nil, 3, false))
	assert.Equal(t, []byte{0, 0, 0}, alignBytes(nil, 3, true))

	assert.Equal(t, []byte{1, 2, 0, 0}, alignBytes([]byte{1, 2}, 4, false))
	assert.Equal(t, []byte{0, 0, 1, 2}, alignBytes([]byte{1, 2}, 4, true))

	assert.Equal(t, []byte{1, 2}, alignBytes([]byte{1, 2, 3, 4}, 2, false))
	assert.Equal(t, []byte{3, 4}, alignBytes([]byte{1, 2, 3, 4}, 2, true))
}

func TestBytesToHex(t *testing.T) {
	assert.Equal(t, "0x000000", string(writeDataHex(nil, []byte{0, 0, 0})))
	assert.Equal(t, "0x0", string(writeQuantityHex(nil, []byte{0, 0, 0})))

	assert.Equal(t, "0x010203", string(writeDataHex(nil, []byte{1, 2, 3})))
	assert.Equal(t, "0x10203", string(writeQuantityHex(nil, []byte{1, 2, 3})))

	assert.Equal(t, "0x0000000010", string(writeDataHex(nil, []byte{0, 0, 0, 0, 16})))
	assert.Equal(t, "0x10", string(writeQuantityHex(nil, []byte{0, 0, 0, 0, 16})))
}

func TestHexToBytes(t *testing.T) {
	testCases := []struct {
		name string
		in   string
		out  []byte
		err  bool
	}{{
		name: "empty",
		in:   "",
		out:  []byte{},
	}, {
		name: "empty with prefix",
		in:   "0x",
		out:  []byte{},
	}, {
		name: "zero without prefix",
		in:   "00",
		out:  []byte{0x00},
	}, {
		name: "even length without prefix",
		in:   "1a",
		out:  []byte{0x1a},
	}, {
		name: "uppercase without prefix",
		in:   "FC",
		out:  []byte{0xfc},
	}, {
		name: "longer, without prefix",
		in:   "0a1b2c",
		out:  []byte{0x0a, 0x1b, 0x2c},
	}, {
		name: "zero with prefix",
		in:   "0x00",
		out:  []byte{0x00},
	}, {
		name: "even length with prefix",
		in:   "0x1a",
		out:  []byte{0x1a},
	}, {
		name: "uppercase with prefix",
		in:   "0XFC",
		out:  []byte{0xfc},
	}, {
		name: "zero with odd length",
		in:   "0x0",
		out:  []byte{0x00},
	}, {
		name: "odd length without prefix",
		in:   "a",
		out:  []byte{0x0a},
	}, {
		name: "uppercase with odd length",
		in:   "0XFC1",
		out:  []byte{0x0f, 0xc1},
	}, {
		name: "longer, with odd length",
		in:   "0xa1b2c",
		out:  []byte{0x0a, 0x1b, 0x2c},
	}, {
		name: "only invalid characters",
		in:   "0xgg",
		err:  true,
	}, {
		name: "only invalid characters without prefix",
		in:   "gg",
		err:  true,
	}, {
		name: "mixed invalid characters, with odd length",
		in:   "0x1b0gg",
		err:  true,
	}, {
		name: "with whitespace",
		in:   "0x1 a",
		err:  true,
	}, {
		name: "non-ascii characters",
		in:   "0x（╯°□°）╯ ┻━┻",
		err:  true,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := hexToBytes(tc.in)
			if tc.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.out, res)
			}
		})
	}
}
