package primitives

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/aurora-is-near/relayer2-base/db/codec"
	tp "github.com/aurora-is-near/relayer2-base/tinypack"
)

type Data[LD tp.LengthDescriptor] struct {
	tp.Data[LD]
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
	*d, err = hexStringToData[LD](string(b[1 : len(b)-1]))
	return err
}

func (d *Data[LD]) UnmarshalCBOR(b []byte) error {
	var in string
	var err error
	err = codec.CborDecoder().Unmarshal(b, &in)
	if err != nil {
		return err
	}
	*d, err = hexStringToData[LD](in)
	return err
}

func DataFromBytes[LD tp.LengthDescriptor](b []byte) Data[LD] {
	var ld LD
	var d Data[LD]
	d.Content = alignBytes(b, ld.GetTinyPackLength(), false)
	return d
}

func DataFromHex[LD tp.LengthDescriptor](s string) Data[LD] {
	data, _ := hexStringToData[LD](s)
	return data
}

func hexStringToData[LD tp.LengthDescriptor](in string) (Data[LD], error) {
	if len(in) == 2 && in[0] == '0' && (in[1] == 'x' || in[2] == 'X') {
		return Data[LD]{}, nil
	}
	bytes, err := hexToByte(in)
	if err != nil {
		return Data[LD]{}, err
	}
	var ld LD
	l := ld.GetTinyPackLength()
	if l > 0 && l != len(bytes) {
		return Data[LD]{}, fmt.Errorf("hex string of length [%d] is expected, found [%d]", l, len(bytes))
	}
	return DataFromBytes[LD](bytes), nil
}
