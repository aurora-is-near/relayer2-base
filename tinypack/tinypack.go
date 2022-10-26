package tinypack

import (
	"io"
)

type Writer interface {
	io.Writer
	io.ByteWriter
}

type Reader interface {
	io.Reader
	io.ByteReader
}

type TinyPackable interface {
	WriteTinyPack(w Writer, e *Encoder) error
	ReadTinyPack(r Reader, d *Decoder) error
}

type Composite interface {
	GetTinyPackChildrenPointers() ([]any, error)
}
