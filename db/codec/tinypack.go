package codec

import "github.com/aurora-is-near/relayer2-base/tinypack"

type Tinypack struct {
	Encoder
	Decoder
}

func NewTinypackCodec() Codec {
	return Tinypack{
		Encoder: tinypack.DefaultEncoder(),
		Decoder: tinypack.DefaultDecoder(),
	}
}
