package github_neonxp_jsonrpc2

import (
	"aurora-relayer-common/log"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.neonxp.dev/jsonrpc2/rpc"
	"go.neonxp.dev/jsonrpc2/transport"
)

const (
	configPath = "RpcNode.JsonRpc2"
)

type JsonRpc2 struct {
	rpc.RpcServer
}

func New() (*JsonRpc2, error) {
	logger := log.New()
	conf := DefaultConfig()
	sub := viper.Sub(configPath)
	if sub != nil {
		if err := sub.Unmarshal(&conf); err != nil {
			logger.Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
				"falling back to defaults", configPath, viper.ConfigFileUsed())
		}
	}
	conf.Logger = NewNeonxpJsonRpc2Logger(logger)
	return NewWithConf(conf)
}

func NewWithConf(config *Config) (*JsonRpc2, error) {
	n := rpc.New(
		rpc.WithLogger(config.Logger),
		rpc.WithTransport(&transport.HTTP{Bind: fmt.Sprintf("%d", config.HTTPConfig.Port)}),
	)
	return &JsonRpc2{*n}, nil
}

func (n *JsonRpc2) Start() error {
	return n.Run(context.Background())
}
