package db

import "github.com/fxamacker/cbor/v2"

type Encoder interface {
	Marshal(v interface{}) ([]byte, error)
}

type Decoder interface {
	Unmarshal(data []byte, v interface{}) error
}

type Codec interface {
	Encoder
	Decoder
}

type CborCodec struct {
	Encoder
	Decoder
}

func NewCborCodec() Codec {
	enc, dec := NewCborEncDec()
	return CborCodec{enc, dec}
}

func NewCborEncDec() (Encoder, Decoder) {
	enc, err := cbor.EncOptions{
		BigIntConvert: cbor.BigIntConvertShortest,
	}.EncMode()
	if err != nil {
		panic(err)
	}
	dec, err := cbor.DecOptions{}.DecMode()
	if err != nil {
		panic(err)
	}
	return enc, dec
}
