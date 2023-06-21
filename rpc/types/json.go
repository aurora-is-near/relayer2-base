package types

import (
	"fmt"

	"github.com/buger/jsonparser"
)

type JsonFieldSpec struct {
	Path         []string
	AllowedTypes map[jsonparser.ValueType]bool // empty for anything
}

type JsonSpec struct {
	fields []*JsonFieldSpec
	paths  [][]string
}

type JsonValue struct {
	Value []byte
	Type  jsonparser.ValueType

	strCached     string
	reprCached    []byte
	strReprCached string
}

func NewJsonSpec(fields ...*JsonFieldSpec) *JsonSpec {
	spec := &JsonSpec{fields: fields, paths: make([][]string, len(fields))}
	for i, field := range fields {
		spec.paths[i] = field.Path
	}
	return spec
}

// Parse parses the incoming byte slice to a slice of JsonValues
func (spec *JsonSpec) Parse(data []byte) ([]*JsonValue, error) {
	values := make([]*JsonValue, len(spec.fields))

	var anyErr error
	jsonparser.EachKey(data, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
		if anyErr != nil {
			return
		}
		if err != nil {
			anyErr = err
			return
		}
		values[idx] = &JsonValue{
			Value: value,
			Type:  vt,
		}
	}, spec.paths...)

	if anyErr != nil {
		return nil, anyErr
	}

	for i := range values {
		if values[i] == nil {
			values[i] = &JsonValue{Type: jsonparser.NotExist}
		}
		if spec.fields[i].AllowedTypes != nil && len(spec.fields[i].AllowedTypes) > 0 {
			if allowed, ok := spec.fields[i].AllowedTypes[values[i].Type]; !allowed || !ok {
				return nil, fmt.Errorf("invalid type of field %v: %v", spec.paths[i], values[i].Type)
			}
		}
	}

	return values, nil
}

func (val *JsonValue) Str() string {
	if len(val.strCached) == 0 {
		val.strCached = string(val.Value)
	}
	return val.strCached
}

// Repr returns string representation of the JsonValue as a byte array
// If the representation was alredy cached, method directly returns the cached value
// If not, then the value first cached and then returned
func (val *JsonValue) Repr() []byte {
	if val.reprCached == nil {
		if val.Type != jsonparser.String {
			val.reprCached = val.Value
		} else {
			val.reprCached = make([]byte, len(val.Value)+2)
			val.reprCached[0] = `"`[0]
			copy(val.reprCached[1:len(val.reprCached)-1], val.Value)
			val.reprCached[len(val.reprCached)-1] = val.reprCached[0]
		}
	}
	return val.reprCached
}

// StrRepr returns string representation of the JsonValue as string
// If the representation was alredy cached,  method directly returns the cached value
// If not, then the value first cached and then returned
func (val *JsonValue) StrRepr() string {
	if len(val.strReprCached) == 0 {
		if val.Type != jsonparser.String {
			val.strReprCached = val.Str()
		} else {
			val.strReprCached = fmt.Sprintf(`"%s"`, val.Value)
		}
	}
	return val.strReprCached
}
