package tinypack

import "fmt"

type Pointer[T any] struct {
	Ptr *T
}

func (p *Pointer[T]) WriteTinyPack(w Writer, e *Encoder) error {
	if p.Ptr == nil {
		return fmt.Errorf("tinypack: Ptr field of Pointer can't be nil")
	}
	if err := e.Write(w, p.Ptr); err != nil {
		return fmt.Errorf("tinypack: can't write content of Pointer: %w", err)
	}
	return nil
}

func (p *Pointer[T]) ReadTinyPack(r Reader, d *Decoder) error {
	p.Ptr = new(T)
	if err := d.Read(r, p.Ptr); err != nil {
		return fmt.Errorf("tinypack: can't read content of Pointer: %w", err)
	}
	return nil
}

func CreatePointer[T any](ptr *T) (result Pointer[T]) {
	result.Ptr = ptr
	return
}
