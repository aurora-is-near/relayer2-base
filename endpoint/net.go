package endpoint

import (
	"fmt"
	"golang.org/x/net/context"
)

type Net struct {
	*Endpoint
}

func NewNet(endpoint *Endpoint) *Eth {
	return &Eth{endpoint}
}

// Listening always returns true on success
//
//	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
func (e *Net) Listening(_ context.Context) (bool, error) {
	return true, nil
}

// PeerCount always returns hex zero
//
//	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
func (e *Net) PeerCount(_ context.Context) (string, error) {
	return fmt.Sprintf("0x%x", 0), nil
}
