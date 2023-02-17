package primitives

import "github.com/aurora-is-near/relayer2-base/tinypack"

type VarLen = tinypack.VariadicLength

type Len8 struct{}

func (Len8) GetTinyPackLength() int {
	return 8
}

type Len20 struct{}

func (Len20) GetTinyPackLength() int {
	return 20
}

type Len32 struct{}

func (Len32) GetTinyPackLength() int {
	return 32
}

type Len256 struct{}

func (Len256) GetTinyPackLength() int {
	return 256
}
