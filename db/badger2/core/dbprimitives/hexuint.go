package dbprimitives

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mailru/easyjson/jwriter"
)

// Special type for convenient uint64-to-hex conversion
// Only used in responses, not for db storage
type HexUint uint64

func (h HexUint) Hex() string {
	return hexutil.EncodeUint64(uint64(h))
}

// Can (and must) be dramatically optimized
func (h HexUint) MarshalEasyJSON(w *jwriter.Writer) {
	w.String(h.Hex())
}
