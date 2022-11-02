package codec

import "aurora-relayer-go-common/tinypack"

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
