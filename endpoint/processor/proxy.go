package processor

import (
	"context"
	"github.com/aurora-is-near/relayer2-base/endpoint"
	"github.com/aurora-is-near/relayer2-base/syncutils"
	"reflect"

	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/exp/slices"
)

type Proxy struct {
	clientPtr syncutils.LockablePtr[rpc.Client]
}

func NewProxy() endpoint.Processor {
	return &Proxy{}
}

func (p *Proxy) Pre(ctx context.Context, name string, endpoint *endpoint.Endpoint, response *any, args ...any) (context.Context, bool, error) {
	if endpoint.Config.ProxyEndpoints[name] {
		endpoint.Logger.Info().Msgf("relaying request: [%s] to remote server", name)
		var err error
		client, unlock := p.clientPtr.LockIfNil()
		if unlock != nil {
			client, err = clientConnection(ctx, endpoint.Config.ProxyUrl)
			unlock(client)
		}
		if err != nil {
			return ctx, true, err
		}
		// Delete nil values (empty optional parameters) from the parameter array
		for i, v := range args {
			rv := reflect.ValueOf(v)
			if rv.Kind() == reflect.Ptr && rv.IsNil() {
				args = slices.Delete(args, i, i+1)
			}
		}
		err = client.CallContext(ctx, response, name, args...)
		if err != nil {
			endpoint.Logger.Error().Err(err).Msgf("failed to call remote server for request: [%s]", name)
			return ctx, true, err
		}
		endpoint.Logger.Debug().Msgf("response received from remote server for request: [%s], response: [%v]", name, response)
		return ctx, true, nil
	}
	return ctx, false, nil
}

func (p *Proxy) Post(ctx context.Context, _ string, _ *any, _ *error) context.Context {
	return ctx
}

func clientConnection(ctx context.Context, url string) (*rpc.Client, error) {
	return rpc.DialContext(ctx, url)
}
