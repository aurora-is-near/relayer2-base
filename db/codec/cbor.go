package codec

import "github.com/fxamacker/cbor/v2"

type CborCodec struct {
	Encoder
	Decoder
}

func NewCborCodec() Codec {
	enc, dec := cborEncDec()
	return CborCodec{enc, dec}
}

func cborEncDec() (Encoder, Decoder) {
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
