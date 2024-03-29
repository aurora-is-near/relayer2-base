package endpoint

import (
	"golang.org/x/net/context"
)

var (
	peerCount = "0x0"
	listening = true
)

type Net struct {
	*Endpoint
}

func NewNet(endpoint *Endpoint) *Net {
	return &Net{endpoint}
}

// Listening always returns true on success
//
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
func (e *Net) Listening(_ context.Context) (*bool, error) {
	return &listening, nil
}

// PeerCount always returns hex zero
//
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
func (e *Net) PeerCount(_ context.Context) (*string, error) {
	return &peerCount, nil
}
