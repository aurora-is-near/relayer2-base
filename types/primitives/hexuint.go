package primitives

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Special type for convenient uint64-to-hex conversion
// Only used in responses, not for db storage
type HexUint uint64

func (h HexUint) Hex() string {
	return hexutil.EncodeUint64(uint64(h))
}

// Can (and must) be dramatically optimized
func (h HexUint) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, h.Hex())), nil
}

func (h *HexUint) UnmarshalJSON(b []byte) error {
	ui64, err := hexutil.DecodeUint64(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*h = HexUint(ui64)
	return nil
}
