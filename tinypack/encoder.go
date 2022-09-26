package tinypack

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"capnproto.org/go/capnp/v3/packed"
)

type Encoder struct{}

func DefaultEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := e.Write(&buf, v); err != nil {
		return nil, err
	}

	padding := (8 - buf.Len()%8) % 8
	buf.Grow(padding)
	for i := 0; i < padding; i++ {
		if err := buf.WriteByte(0); err != nil {
			return nil, fmt.Errorf("tinypack: can't write padding bytes: %w", err)
		}
	}

	return packed.Pack(nil, buf.Bytes()), nil
}

func (e *Encoder) Write(w Writer, v ...any) error {
	for _, cur := range v {
		if err := e.write(w, cur); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) write(w Writer, v any) error {
	if v == nil {
		return fmt.Errorf("tinypack: can't write type %T: provided pointer is nil", v)
	}

	switch vt := v.(type) {
	case *bool:
		return e.WriteBool(w, *vt)
	case *int64:
		return e.WriteVarint(w, *vt)
	case *uint64:
		return e.WriteUvarint(w, *vt)
	case *float64:
		return e.WriteFloat(w, *vt)
	default:
	}

	if packable, ok := v.(TinyPackable); ok {
		if err := packable.WriteTinyPack(w, e); err != nil {
			return fmt.Errorf("tinypack: can't write TinyPackable type %T: %w", v, err)
		}
		return nil
	}

	if composite, ok := v.(Composite); ok {
		children, err := composite.GetTinyPackChildrenPointers()
		if err != nil {
			return fmt.Errorf("tinypack: can't get children for composite type %T: %w", v, err)
		}
		return e.Write(w, children...)
	}

	return fmt.Errorf("tinypack: no idea how to write type %T", v)
}

func (e *Encoder) WriteBool(w Writer, v bool) error {
	boolByte := byte(0)
	if v {
		boolByte = 1
	}
	if err := w.WriteByte(boolByte); err != nil {
		return fmt.Errorf("tinypack: can't write bool: %w", err)
	}
	return nil
}

func (e *Encoder) WriteVarint(w Writer, v int64) error {
	var buf [binary.MaxVarintLen64]byte
	len := binary.PutVarint(buf[:], v)
	if err := ensureWrite(w, buf[:len]); err != nil {
		return fmt.Errorf("tinypack: can't write varint: %w", err)
	}
	return nil
}

func (e *Encoder) WriteUvarint(w Writer, v uint64) error {
	var buf [binary.MaxVarintLen64]byte
	len := binary.PutUvarint(buf[:], v)
	if err := ensureWrite(w, buf[:len]); err != nil {
		return fmt.Errorf("tinypack: can't write uvarint: %w", err)
	}
	return nil
}

func (e *Encoder) WriteFloat(w Writer, v float64) error {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(v))
	if err := ensureWrite(w, buf[:]); err != nil {
		return fmt.Errorf("tinypack: can't write float64: %w", err)
	}
	return nil
}
