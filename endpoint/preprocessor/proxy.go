package preprocessor

import (
	"aurora-relayer-go-common/endpoint"
	"aurora-relayer-go-common/log"
	"github.com/spf13/viper"
)

const (
	configPathProxy = "Endpoint.ProxyEndpoints"
)

type Proxy struct {
	proxyEndpoints map[string]bool
}

func NewProxy() endpoint.Preprocessor {
	conf := make([]string, 0)
	sub := viper.Sub(configPathProxy)
	if sub != nil {
		if err := sub.Unmarshal(&conf); err != nil {
			log.New().Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
				"falling back to defaults", configPathProxy, viper.ConfigFileUsed())
		}
	}

	ed := Proxy{proxyEndpoints: setProxyEndpoints(conf)}

	// TODO config change and json-rpc client implementation

	return func(name string, _ *endpoint.Endpoint, _ ...any) (bool, *any, error) {
		if ed.proxyEndpoints[name] {
			// TODO send request to proxy server
			return true, nil, nil
		}
		return false, nil, nil
	}
}

func setProxyEndpoints(conf []string) map[string]bool {
	de := make(map[string]bool, len(conf))
	for _, e := range conf {
		de[e] = true
	}
	return de
}
