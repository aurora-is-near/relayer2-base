package primitives

type LengthDescriptor interface {
	GetLength() int
}

type VarLen struct{}

func (VarLen) GetLength() int {
	return -1
}

type Len8 struct{}

func (Len8) GetLength() int {
	return 8
}

type Len20 struct{}

func (Len20) GetLength() int {
	return 20
}

type Len32 struct{}

func (Len32) GetLength() int {
	return 32
}

type Len256 struct{}

func (Len256) GetLength() int {
	return 256
}
