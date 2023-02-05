package jsonbench

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"testing"

	"github.com/bytedance/sonic"
	goccy "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
	segj "github.com/segmentio/encoding/json"
	"github.com/wI2L/jettison"
)

var sideEffectHash uint32

func run(b *testing.B, fn func(payload any, hasher io.Writer)) {
	b.StopTimer()
	b.ReportAllocs()
	hasher := crc32.NewIEEE()
	r := NewRandom()

	for i := 0; i < b.N; i++ {
		payloads := generatePayloads(r)
		b.StartTimer()
		for _, payload := range payloads {
			fn(payload, hasher)
		}
		b.StopTimer()
	}

	sideEffectHash = hasher.Sum32()
}

func runMarshal(b *testing.B, marshalFn func(payload any) ([]byte, error)) {
	run(b, func(payload any, hasher io.Writer) {
		data, err := marshalFn(payload)
		if err != nil {
			b.Fatal(err)
		}
		hasher.Write(data)
	})
}

func runEncode(b *testing.B, encodeFn func(writer io.Writer, payload any) error) {
	run(b, func(payload any, hasher io.Writer) {
		if err := encodeFn(hasher, payload); err != nil {
			b.Fatal(err)
		}
	})
}

func BenchmarkStdMarshal(b *testing.B) {
	runMarshal(b, json.Marshal)
}

func BenchmarkStdEncode(b *testing.B) {
	runEncode(b, func(writer io.Writer, payload any) error {
		return json.NewEncoder(writer).Encode(payload)
	})
}

func BenchmarkStdEncodeNoEscape(b *testing.B) {
	runEncode(b, func(writer io.Writer, payload any) error {
		encoder := json.NewEncoder(writer)
		encoder.SetEscapeHTML(false)
		return encoder.Encode(payload)
	})
}

func BenchmarkJettison(b *testing.B) {
	runMarshal(b, jettison.Marshal)
}

func BenchmarkJettisonFast(b *testing.B) {
	opts := []jettison.Option{
		jettison.NoCompact(),
		jettison.NoNumberValidation(),
		jettison.UnsortedMap(),
	}

	runMarshal(b, func(payload any) ([]byte, error) {
		return jettison.MarshalOpts(payload, opts...)
	})
}

func BenchmarkJettisonSuperFast(b *testing.B) {
	opts := []jettison.Option{
		jettison.NoCompact(),
		jettison.NoHTMLEscaping(),
		jettison.NoNumberValidation(),
		jettison.NoStringEscaping(),
		jettison.NoUTF8Coercion(),
		jettison.UnsortedMap(),
	}

	runMarshal(b, func(payload any) ([]byte, error) {
		return jettison.MarshalOpts(payload, opts...)
	})
}

func BenchmarkSegmentioMarshal(b *testing.B) {
	runMarshal(b, segj.Marshal)
}

func BenchmarkSegmentioEncode(b *testing.B) {
	runEncode(b, func(writer io.Writer, payload any) error {
		return segj.NewEncoder(writer).Encode(payload)
	})
}

func BenchmarkSegmentioEncodeFast(b *testing.B) {
	runEncode(b, func(writer io.Writer, payload any) error {
		encoder := segj.NewEncoder(writer)
		encoder.SetEscapeHTML(false)
		encoder.SetSortMapKeys(false)
		encoder.SetTrustRawMessage(true)
		return encoder.Encode(payload)
	})
}

func BenchmarkJsoniterDefaultMarshal(b *testing.B) {
	registerJsoniterEncoders()
	runMarshal(b, jsoniter.ConfigDefault.Marshal)
}

func BenchmarkJsoniterDefaultEncode(b *testing.B) {
	registerJsoniterEncoders()
	runEncode(b, func(writer io.Writer, payload any) error {
		return jsoniter.ConfigDefault.NewEncoder(writer).Encode(payload)
	})
}

func BenchmarkJsoniterCompatMarshal(b *testing.B) {
	registerJsoniterEncoders()
	runMarshal(b, jsoniter.ConfigCompatibleWithStandardLibrary.Marshal)
}

func BenchmarkJsoniterCompatEncode(b *testing.B) {
	registerJsoniterEncoders()
	runEncode(b, func(writer io.Writer, payload any) error {
		return jsoniter.ConfigCompatibleWithStandardLibrary.NewEncoder(writer).Encode(payload)
	})
}

func BenchmarkJsoniterFastestMarshal(b *testing.B) {
	registerJsoniterEncoders()
	runMarshal(b, jsoniter.ConfigFastest.Marshal)
}

func BenchmarkJsoniterFastestEncode(b *testing.B) {
	registerJsoniterEncoders()
	runEncode(b, func(writer io.Writer, payload any) error {
		return jsoniter.ConfigFastest.NewEncoder(writer).Encode(payload)
	})
}

func BenchmarkGoccyMarshal(b *testing.B) {
	runMarshal(b, goccy.Marshal)
}

func BenchmarkGoccyMarshalNoEscape(b *testing.B) {
	runMarshal(b, goccy.MarshalNoEscape)
}

func BenchmarkGoccyMarshalFast(b *testing.B) {
	opts := []goccy.EncodeOptionFunc{
		goccy.DisableHTMLEscape(),
		goccy.DisableNormalizeUTF8(),
		goccy.UnorderedMap(),
	}

	runMarshal(b, func(payload any) ([]byte, error) {
		return goccy.MarshalWithOption(payload, opts...)
	})
}

func BenchmarkGoccyEncode(b *testing.B) {
	runEncode(b, func(writer io.Writer, payload any) error {
		return goccy.NewEncoder(writer).Encode(payload)
	})
}

func BenchmarkGoccyEncodeFast(b *testing.B) {
	opts := []goccy.EncodeOptionFunc{
		goccy.DisableHTMLEscape(),
		goccy.DisableNormalizeUTF8(),
		goccy.UnorderedMap(),
	}

	runEncode(b, func(writer io.Writer, payload any) error {
		return goccy.NewEncoder(writer).EncodeWithOption(payload, opts...)
	})
}

func BenchmarkSonicMarhsalDefault(b *testing.B) {
	runMarshal(b, sonic.ConfigDefault.Marshal)
}

func BenchmarkSonicMarhsalStd(b *testing.B) {
	runMarshal(b, sonic.ConfigStd.Marshal)
}

func BenchmarkSonicMarhsalFastest(b *testing.B) {
	runMarshal(b, sonic.ConfigFastest.Marshal)
}

func BenchmarkSonicEncodeDefault(b *testing.B) {
	runEncode(b, func(writer io.Writer, payload any) error {
		return sonic.ConfigDefault.NewEncoder(writer).Encode(payload)
	})
}

func BenchmarkSonicEncodeStd(b *testing.B) {
	runEncode(b, func(writer io.Writer, payload any) error {
		return sonic.ConfigStd.NewEncoder(writer).Encode(payload)
	})
}

func BenchmarkSonicEncodeFastest(b *testing.B) {
	runEncode(b, func(writer io.Writer, payload any) error {
		return sonic.ConfigFastest.NewEncoder(writer).Encode(payload)
	})
}

func TestValidateJettison(t *testing.T) {
	r := NewRandom()
	for i := 0; i < 1000; i++ {
		payloads := generatePayloads(r)
		for _, payload := range payloads {
			stdData, err := json.Marshal(payload)
			if err != nil {
				t.Fatal(err)
			}
			jtsData, err := jettison.Marshal(payload)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(stdData, jtsData) {
				fmt.Printf("STD: %s\n\n", stdData)
				fmt.Printf("JTS: %s\n", jtsData)
				t.FailNow()
			}
		}
	}
}
