package tinypack

import "fmt"

type Nullable[T any] struct {
	Ptr *T
}

func (opt *Nullable[T]) WriteTinyPack(w Writer, e *Encoder) error {
	header := byte(0)
	if opt.Ptr != nil {
		header = 1
	}
	if err := w.WriteByte(header); err != nil {
		return fmt.Errorf("tinypack: can't write header of Nullable: %w", err)
	}
	if opt.Ptr != nil {
		if err := e.Write(w, opt.Ptr); err != nil {
			return fmt.Errorf("tinypack: can't write content of Nullable: %w", err)
		}
	}
	return nil
}

func (opt *Nullable[T]) ReadTinyPack(r Reader, d *Decoder) error {
	header, err := r.ReadByte()
	if err != nil {
		return fmt.Errorf("tinypack: can't read header of Nullable: %w", err)
	}
	if header == 0 {
		opt.Ptr = nil
		return nil
	}
	opt.Ptr = new(T)
	if err := d.Read(r, opt.Ptr); err != nil {
		return fmt.Errorf("tinypack: can't read content of Nullable: %w", err)
	}
	return nil
}

func CreateNullable[T any](ptr *T) (result Nullable[T]) {
	result.Ptr = ptr
	return
}
