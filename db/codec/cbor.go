package codec

import (
	"github.com/aurora-is-near/relayer2-base/log"
	"github.com/fxamacker/cbor/v2"
)

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
