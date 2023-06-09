package primitives

import (
	"encoding/binary"
	"encoding/json"
	"math/big"
	"reflect"

	tp "github.com/aurora-is-near/relayer2-base/tinypack"

	"github.com/fxamacker/cbor/v2"
)

type Quantity struct {
	tp.Data[Len32]
}

func (q Quantity) Bytes() []byte {
	return q.Content
}

func (q Quantity) Hex() string {
	return string(q.WriteHexBytes(make([]byte, 0, 3)))
}

func (d Quantity) WriteHexBytes(dst []byte) []byte {
	return writeQuantityHex(dst, d.Content)
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

func (q Quantity) MarshalJSON() ([]byte, error) {
	buf := make([]byte, 0, 5)
	buf = append(buf, '"')
	buf = q.WriteHexBytes(buf)
	buf = append(buf, '"')
	return buf, nil
}

func (q *Quantity) UnmarshalJSON(b []byte) error {
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return &json.UnmarshalTypeError{Value: "non-string", Type: reflect.ValueOf(Quantity{}).Type()}
	}

	var in string
	err := json.Unmarshal(b, &in)
	if err != nil {
		return err
	}

	*q = QuantityFromHex(in)
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
	b, _ := hexToByte(s)
	return QuantityFromBytes(b)
}

func QuantityFromBigInt(v *big.Int) Quantity {
	return QuantityFromBytes(v.Bytes())
}

func QuantityFromUint64(v uint64) Quantity {
	buf := make([]byte, 32)
	binary.BigEndian.PutUint64(buf[24:], v)
	return QuantityFromBytes(buf)
}
