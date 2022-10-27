package dbprimitives

import (
	"encoding/binary"
	"fmt"
	"math/big"

	tp "aurora-relayer-go-common/tinypack"
	"github.com/ethereum/go-ethereum/common"
)

type Quantity struct {
	tp.Data[Len32]
}

func (q Quantity) Bytes() []byte {
	return q.Content
}

func (q Quantity) Hex() string {
	return bytesToHex(q.Content, true)
}

func (q Quantity) BigInt() *big.Int {
	x := big.NewInt(0)
	x.SetBytes(q.Content)
	return x
}

func (q Quantity) IsUint64() bool {
	for i := 0; i < 24; i++ {
		if q.Content[i] > 0 {
			return false
		}
	}
	return true
}

func (q Quantity) Uint64() uint64 {
	return binary.BigEndian.Uint64(q.Content[24:])
}

// Can (and must) be dramatically optimized
func (q Quantity) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, q.Hex())), nil
}

func QuantityFromBytes(b []byte) Quantity {
	var q Quantity
	q.Content = alignBytes(b, 32, true)
	return q
}

func QuantityFromHex(s string) Quantity {
	return QuantityFromBytes(common.FromHex(s))
}

func QuantityFromBigInt(v *big.Int) Quantity {
	return QuantityFromBytes(v.Bytes())
}

func QuantityFromUint64(v uint64) Quantity {
	buf := make([]byte, 32)
	binary.BigEndian.PutUint64(buf[24:], v)
	return QuantityFromBytes(buf)
}
