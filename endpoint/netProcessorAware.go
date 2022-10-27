package endpoint

import (
	"golang.org/x/net/context"
)

type NetProcessorAware struct {
	*Net
}

func NewNetProcessorAware(net *Net) *NetProcessorAware {
	return &NetProcessorAware{net}
}

func (e *NetProcessorAware) Listening(ctx context.Context) (*bool, error) {
	return Process(ctx, "net_listening", e.Endpoint, func(ctx context.Context) (*bool, error) {
		return e.Net.Listening(ctx)
	})
}

func (e *NetProcessorAware) PeerCount(ctx context.Context) (*string, error) {
	return Process(ctx, "net_listening", e.Endpoint, func(ctx context.Context) (*string, error) {
		return e.Net.PeerCount(ctx)
	})
}
