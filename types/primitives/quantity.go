package primitives

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	tp "github.com/aurora-is-near/relayer2-base/tinypack"
	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"

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

func (q Quantity) IsZero() bool {
	tmpBig := q.BigInt()
	// length of bits in BigInt for value `0` is 0
	return len(tmpBig.Bits()) == 0
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
	err := jsoniter.Unmarshal(b, &in)
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
	b, _ := hexToBytes(s)
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

func QuantityDecodeHook() mapstructure.DecodeHookFuncType {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(Quantity{}) {
			return data, nil
		}

		result := big.NewInt(0)

		switch f.Kind() {
		case reflect.String:
			var ok bool
			result, ok = big.NewInt(0).SetString(data.(string), 0)
			if !ok {
				return nil, fmt.Errorf("unable to parse unsigned quantity from string '%s'", data)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result.SetInt64(reflect.ValueOf(data).Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result.SetUint64(reflect.ValueOf(data).Uint())
		default:
			return data, nil
		}

		if result.Sign() < 0 {
			return nil, fmt.Errorf("can't parse negative number %s into unsigned quantity", result.Text(10))
		}
		return QuantityFromBigInt(result), nil
	}
}
