package codec

import (
	"aurora-relayer-go-common/log"
	"github.com/fxamacker/cbor/v2"
)

type CborCodec struct {
	Encoder
	Decoder
}

func NewCborCodec() Codec {
	return CborCodec{
		Encoder: CborEncoder(),
		Decoder: CborDecoder(),
	}
}

func CborDecoder() cbor.DecMode {
	dec, err := cbor.DecOptions{
		MaxArrayElements: 2147483647,
	}.DecMode()
	if err != nil {
		log.Log().Fatal().Err(err).Msg("failed to initialize CBOR decoder")
	}
	return dec
}

func CborEncoder() cbor.EncMode {
	enc, err := cbor.EncOptions{
		BigIntConvert: cbor.BigIntConvertShortest,
	}.EncMode()
	if err != nil {
		log.Log().Fatal().Err(err).Msg("failed to initialize CBOR encoder")
	}
	return enc
}
