package primitives

import (
	"encoding/json"
	"reflect"

	"github.com/fxamacker/cbor/v2"
)

type Data[LD LengthDescriptor] struct {
	Content []byte
}

func (d Data[LD]) Bytes() []byte {
	return d.Content
}

func (d Data[LD]) Hex() string {
	return string(d.WriteHexBytes(make([]byte, 0, 2+len(d.Content)*2)))
}

func (d Data[LD]) WriteHexBytes(dst []byte) []byte {
	return writeDataHex(dst, d.Content)
}

func (d Data[LD]) MarshalJSON() ([]byte, error) {
	buf := make([]byte, 0, 4+len(d.Content)*2)
	buf = append(buf, '"')
	buf = d.WriteHexBytes(buf)
	buf = append(buf, '"')
	return buf, nil
}

func (d *Data[LD]) UnmarshalJSON(b []byte) error {
	var err error
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return &json.UnmarshalTypeError{Value: "non-string", Type: reflect.ValueOf(Data[LD]{}).Type()}
	}
	*d, err = DataFromHex[LD](string(b[1 : len(b)-1]))
	return err
}

func (d *Data[LD]) UnmarshalCBOR(b []byte) error {
	var in string
	var err error
	err = cbor.Unmarshal(b, &in)
	if err != nil {
		return err
	}
	*d, err = DataFromHex[LD](in)
	return err
}

func DataFromBytes[LD LengthDescriptor](b []byte) Data[LD] {
	var ld LD
	var d Data[LD]
	d.Content = alignBytes(b, ld.GetLength(), true)
	return d
}

func MustDataFromHex[LD LengthDescriptor](s string) Data[LD] {
	data, err := DataFromHex[LD](s)
	if err != nil {
		panic(err)
	}

	return data
}

func DataFromHex[LD LengthDescriptor](in string) (Data[LD], error) {
	if len(in) == 2 && in[0] == '0' && (in[1] == 'x' || in[1] == 'X') {
		return Data[LD]{}, nil
	}
	bytes, err := hexToBytes(in)
	if err != nil {
		return Data[LD]{}, err
	}
	return DataFromBytes[LD](bytes), nil
}
