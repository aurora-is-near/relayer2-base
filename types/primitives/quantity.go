package primitives

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/fxamacker/cbor/v2"
	"math/big"
	"reflect"
	tp "relayer2-base/tinypack"
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

func (q *Quantity) UnmarshalJSON(b []byte) error {

	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return &json.UnmarshalTypeError{Value: "non-string", Type: reflect.ValueOf(Quantity{}).Type()}
	}

	var bi hexutil.Big
	err := json.Unmarshal(b, &bi)
	if err != nil {
		return err
	}

	*q = QuantityFromHex(bi.String())
	return nil
}

func (q *Quantity) UnmarshalCBOR(b []byte) error {

	var in string
	err := cbor.Unmarshal(b, &in)
	if err != nil {
		return err
	}

	if len(in) < 2 {
		return &cbor.UnmarshalTypeError{GoType: "Quantity", CBORType: "[" + in + "] is not hex string"}
	}

	if in[0] != '0' || (in[1] != 'x' && in[1] != 'X') {
		return &cbor.UnmarshalTypeError{GoType: "Quantity", CBORType: "[" + in + "] is not hex string"}
	}

	if len(in[2:]) > 64 {
		return &cbor.UnmarshalTypeError{GoType: "Quantity", CBORType: "length of" + "[" + in + "] exceeds 32 bytes"}
	}

	*q = QuantityFromHex(in)
	return nil
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
