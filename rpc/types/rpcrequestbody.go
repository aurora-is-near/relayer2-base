package types

import (
	"strings"

	"github.com/buger/jsonparser"
)

const (
	subscribeMethod   = "eth_subscribe"
	unsubscribeMethod = "eth_unsubscribe"
)

type RPCRequestBody struct {
	ID     *JsonValue
	Method *JsonValue
	Params *JsonValue
}

func (rb *RPCRequestBody) IsNotification() bool {
	return rb == nil || (rb.ID.Value == nil && len(rb.Method.Value) > 0)
}

func (rb *RPCRequestBody) HasValidID() bool {
	return rb.ID != nil
}

func (rb *RPCRequestBody) IsMethodCall() bool {
	return rb != nil && rb.HasValidID() && len(rb.Method.Value) > 0
}

func (rb *RPCRequestBody) IsSubscribe() bool {
	return strings.EqualFold(rb.Method.Str(), subscribeMethod)
}

func (rb *RPCRequestBody) IsUnsubscribe() bool {
	return strings.EqualFold(rb.Method.Str(), unsubscribeMethod)
}

var (
	rpcRequestBodySpec = NewJsonSpec(
		&JsonFieldSpec{
			Path: []string{"id"},
		},
		&JsonFieldSpec{
			Path: []string{"method"},
			AllowedTypes: map[jsonparser.ValueType]bool{
				jsonparser.String: true,
			},
		},
		&JsonFieldSpec{
			Path: []string{"params"},
		},
	)
)

func ParseRPCRequestBody(data []byte) (*RPCRequestBody, error) {
	fields, err := rpcRequestBodySpec.Parse(data)
	if err != nil {
		return nil, err
	}
	return &RPCRequestBody{
		ID:     fields[0],
		Method: fields[1],
		Params: fields[2],
	}, nil
}
