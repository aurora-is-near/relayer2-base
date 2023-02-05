package jsonbench

import (
	"relayer2-base/types/primitives"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

const jsoniterEncoders = true

func registerJsoniterEncoders() {
	if !jsoniterEncoders {
		return
	}

	jsoniter.RegisterTypeEncoderFunc(
		reflect2.TypeOf(primitives.HexUint(0)).String(),
		func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			x := *((*primitives.HexUint)(ptr))
			buf, _ := x.AppendJSON(stream.Buffer())
			stream.SetBuffer(buf)
		},
		func(p unsafe.Pointer) bool {
			return false
		},
	)

	jsoniter.RegisterTypeEncoderFunc(
		reflect2.TypeOf(primitives.Quantity{}).String(),
		func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			x := *((*primitives.Quantity)(ptr))
			buf, _ := x.AppendJSON(stream.Buffer())
			stream.SetBuffer(buf)
		},
		func(p unsafe.Pointer) bool {
			return false
		},
	)

	jsoniter.RegisterTypeEncoderFunc(
		reflect2.TypeOf(primitives.VarData{}).String(),
		func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			x := *((*primitives.VarData)(ptr))
			buf, _ := x.AppendJSON(stream.Buffer())
			stream.SetBuffer(buf)
		},
		func(p unsafe.Pointer) bool {
			return false
		},
	)

	jsoniter.RegisterTypeEncoderFunc(
		reflect2.TypeOf(primitives.Data8{}).String(),
		func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			x := *((*primitives.Data8)(ptr))
			buf, _ := x.AppendJSON(stream.Buffer())
			stream.SetBuffer(buf)
		},
		func(p unsafe.Pointer) bool {
			return false
		},
	)

	jsoniter.RegisterTypeEncoderFunc(
		reflect2.TypeOf(primitives.Data20{}).String(),
		func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			x := *((*primitives.Data20)(ptr))
			buf, _ := x.AppendJSON(stream.Buffer())
			stream.SetBuffer(buf)
		},
		func(p unsafe.Pointer) bool {
			return false
		},
	)

	jsoniter.RegisterTypeEncoderFunc(
		reflect2.TypeOf(primitives.Data32{}).String(),
		func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			x := *((*primitives.Data32)(ptr))
			buf, _ := x.AppendJSON(stream.Buffer())
			stream.SetBuffer(buf)
		},
		func(p unsafe.Pointer) bool {
			return false
		},
	)

	jsoniter.RegisterTypeEncoderFunc(
		reflect2.TypeOf(primitives.Data256{}).String(),
		func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			x := *((*primitives.Data256)(ptr))
			buf, _ := x.AppendJSON(stream.Buffer())
			stream.SetBuffer(buf)
		},
		func(p unsafe.Pointer) bool {
			return false
		},
	)
}
