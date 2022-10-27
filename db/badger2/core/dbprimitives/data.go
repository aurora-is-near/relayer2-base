package dbprimitives

import (
	tp "aurora-relayer-go-common/tinypack"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
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

func DataFromBytes[LD tp.LengthDescriptor](b []byte) Data[LD] {
	var ld LD
	var d Data[LD]
	d.Content = alignBytes(b, ld.GetTinyPackLength(), false)
	return d
}

func DataFromHex[LD tp.LengthDescriptor](s string) Data[LD] {
	return DataFromBytes[LD](common.FromHex(s))
}
