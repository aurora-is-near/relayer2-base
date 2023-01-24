package codec

import "relayer2-base/tinypack"

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
