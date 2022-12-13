package primitives

import (
	"aurora-relayer-go-common/db/codec"
	tp "aurora-relayer-go-common/tinypack"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"reflect"
	"strconv"
)

type Data[LD tp.LengthDescriptor] struct {
	tp.Data[LD]
}

func (d Data[LD]) Bytes() []byte {
	return d.Content
}

func (d Data[LD]) Hex() string {
	return bytesToHex(d.Content, false)
}

// Can (and must) be dramatically optimized
func (d Data[LD]) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, d.Hex())), nil
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

func hexStringToData[LD tp.LengthDescriptor](in string) (Data[LD], error) {
	if len(in) == 2 && in[0] == '0' && (in[1] == 'x' || in[2] == 'X') {
		return Data[LD]{}, nil
	}
	var ld LD
	l := ld.GetTinyPackLength()
	if l < 0 {
		l = len(in[2:]) / 2 // remove heading heading '0x'
	}
	out := make([]byte, l, l)
	err := hexutil.UnmarshalFixedText("Data"+strconv.Itoa(l), []byte(in), out)
	if err != nil {
		return Data[LD]{}, err
	}
	return DataFromBytes[LD](out), nil
}

func DataFromBytes[LD tp.LengthDescriptor](b []byte) Data[LD] {
	var ld LD
	var d Data[LD]
	d.Content = alignBytes(b, ld.GetTinyPackLength(), false)
	return d
}

func DataFromHex[LD tp.LengthDescriptor](s string) Data[LD] {
	return DataFromBytes[LD](common.FromHex(s))
}
