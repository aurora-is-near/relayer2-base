package primitives

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Special type for convenient uint64-to-hex conversion
// Only used in responses, not for db storage
type HexUint uint64

func (h HexUint) Hex() string {
	return string(h.WriteHexBytes(make([]byte, 0, 3)))
}

func (h HexUint) WriteHexBytes(dst []byte) []byte {
	dst = append(dst, '0', 'x')
	return strconv.AppendUint(dst, uint64(h), 16)
}

func (h HexUint) MarshalJSON() ([]byte, error) {
	buf := make([]byte, 0, 5)
	buf = append(buf, '"')
	buf = h.WriteHexBytes(buf)
	buf = append(buf, '"')
	return buf, nil
}

func (h HexUint) AppendJSON(buf []byte) ([]byte, error) {
	buf = append(buf, '"')
	buf = h.WriteHexBytes(buf)
	buf = append(buf, '"')
	return buf, nil
}

func (h *HexUint) UnmarshalJSON(b []byte) error {
	ui64, err := hexutil.DecodeUint64(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*h = HexUint(ui64)
	return nil
}
