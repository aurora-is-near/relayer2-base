package dbprimitives

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
	assert.Equal(t, "0x000000", bytesToHex([]byte{0, 0, 0}, false))
	assert.Equal(t, "0x0", bytesToHex([]byte{0, 0, 0}, true))

	assert.Equal(t, "0x010203", bytesToHex([]byte{1, 2, 3}, false))
	assert.Equal(t, "0x10203", bytesToHex([]byte{1, 2, 3}, true))

	assert.Equal(t, "0x0000000010", bytesToHex([]byte{0, 0, 0, 0, 16}, false))
	assert.Equal(t, "0x10", bytesToHex([]byte{0, 0, 0, 0, 16}, true))
}
