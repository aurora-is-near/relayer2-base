package tinypack

import (
	"fmt"
)

type LengthDescriptor interface {
	GetTinyPackLength() int
}

type VariadicLength struct{}

func (VariadicLength) GetTinyPackLength() int {
	return -1
}

type List[LD LengthDescriptor, T any] struct {
	Content []T
}

type VarList[T any] struct {
	List[VariadicLength, T]
}

type Data[LD LengthDescriptor] struct {
	List[LD, byte]
}

type VarData struct {
	List[VariadicLength, byte]
}

func (l *List[LD, T]) getDescriptorLength() int {
	var ld LD
	return ld.GetTinyPackLength()
}

func (l *List[LD, T]) WriteTinyPack(w Writer, e *Encoder) error {
	length := l.getDescriptorLength()

	if length < 0 {
		length = len(l.Content)
		if err := e.WriteUvarint(w, uint64(length)); err != nil {
			return fmt.Errorf("tinypack: can't write variadic list length: %w", err)
		}
	} else {
		if len(l.Content) != length {
			return fmt.Errorf("tinypack: can't write list: len(Content) (%v) != descriptorLength (%v)", len(l.Content), length)
		}
	}

	if length == 0 {
		return nil
	}

	if byteData, isByteData := any(l.Content).([]byte); isByteData {
		if err := ensureWrite(w, byteData); err != nil {
			return fmt.Errorf("tinypack: can't write list byte-data: %w", err)
		}
		return nil
	}

	for i := range l.Content {
		if err := e.Write(w, &l.Content[i]); err != nil {
			return fmt.Errorf("tinypack: can't write list child [%v]: %w", i, err)
		}
	}
	return nil
}

func (l *List[LD, T]) ReadTinyPack(r Reader, d *Decoder) error {
	length := l.getDescriptorLength()

	if length < 0 {
		uint64Length, err := d.ReadUvarint(r)
		if err != nil {
			return fmt.Errorf("tinypack: can't read variadic list length: %w", err)
		}
		if uint64Length > uint64(d.MaxVariadicLength) {
			return fmt.Errorf("tinypack: read list length (%v) > MaxVariadicLength (%v)", uint64Length, d.MaxVariadicLength)
		}
		length = int(uint64Length)
	}
	l.Content = make([]T, length)

	if length == 0 {
		return nil
	}

	if byteData, isByteData := any(l.Content).([]byte); isByteData {
		if err := ensureRead(r, byteData); err != nil {
			return fmt.Errorf("tinypack: can't read list byte-data: %w", err)
		}
		return nil
	}

	for i := range l.Content {
		if err := d.Read(r, &l.Content[i]); err != nil {
			return fmt.Errorf("tinypack: can't read list child [%v]: %w", i, err)
		}
	}
	return nil
}

func CreateList[LD LengthDescriptor, T any](content ...T) (result List[LD, T]) {
	if len(content) == 0 {
		result.Content = make([]T, 0)
	} else {
		result.Content = content
	}
	return
}

func CreateVarList[T any](content ...T) (result VarList[T]) {
	if len(content) == 0 {
		result.Content = make([]T, 0)
	} else {
		result.Content = content
	}
	return
}

func CreateData[LD LengthDescriptor](content ...byte) (result Data[LD]) {
	if len(content) == 0 {
		result.Content = make([]byte, 0)
	} else {
		result.Content = content
	}
	return
}

func CreateVarData(content ...byte) (result VarData) {
	if len(content) == 0 {
		result.Content = make([]byte, 0)
	} else {
		result.Content = content
	}
	return
}
