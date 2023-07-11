package tinypack

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"capnproto.org/go/capnp/v3/packed"
)

type Decoder struct {
	MaxVariadicLength int
}

func DefaultDecoder() *Decoder {
	return &Decoder{
		MaxVariadicLength: 10_000_000,
	}
}

func (d *Decoder) Unmarshal(data []byte, v interface{}) error {
	reader := bufio.NewReaderSize(packed.NewReader(bufio.NewReaderSize(bytes.NewReader(data), 16)), 16)
	return d.Read(reader, v)
}

func (d *Decoder) Read(r Reader, v ...any) error {
	for _, cur := range v {
		if err := d.read(r, cur); err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) read(r Reader, v any) error {
	if v == nil {
		return fmt.Errorf("tinypack: can't read type %T: provided pointer is nil", v)
	}

	var err error
	switch vt := v.(type) {
	case *bool:
		*vt, err = d.ReadBool(r)
		return err
	case *int64:
		*vt, err = d.ReadVarint(r)
		return err
	case *uint64:
		*vt, err = d.ReadUvarint(r)
		return err
	case *float64:
		*vt, err = d.ReadFloat(r)
		return err
	default:
	}

	if packable, ok := v.(TinyPackable); ok {
		if err := packable.ReadTinyPack(r, d); err != nil {
			return fmt.Errorf("tinypack: can't read TinyPackable type %T: %w", v, err)
		}
		return nil
	}

	if composite, ok := v.(Composite); ok {
		children, err := composite.GetTinyPackChildrenPointers()
		if err != nil {
			return fmt.Errorf("tinypack: can't get children for composite type %T: %w", v, err)
		}
		return d.Read(r, children...)
	}

	return fmt.Errorf("tinypack: no idea how to read type %T", v)
}

func (d *Decoder) ReadBool(r Reader) (bool, error) {
	boolByte, err := r.ReadByte()
	if err != nil {
		return false, fmt.Errorf("tinypack: can't read bool: %w", err)
	}
	switch boolByte {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, fmt.Errorf("tinypack: can't read bool: got wrong byte %v", boolByte)
	}
}

func (d *Decoder) ReadFloat(r Reader) (float64, error) {
	var buf [8]byte
	if err := ensureRead(r, buf[:]); err != nil {
		return 0, fmt.Errorf("tinypack: can't read Float: %w", err)
	}
	return math.Float64frombits(binary.BigEndian.Uint64(buf[:])), nil
}

func (d *Decoder) ReadVarint(r Reader) (int64, error) {
	v, err := binary.ReadVarint(r)
	if err != nil {
		return 0, fmt.Errorf("tinypack: can't read varint: %w", err)
	}
	return v, nil
}

func (d *Decoder) ReadUvarint(r Reader) (uint64, error) {
	v, err := binary.ReadUvarint(r)
	if err != nil {
		return 0, fmt.Errorf("tinypack: can't read uvarint: %w", err)
	}
	return v, nil
}
