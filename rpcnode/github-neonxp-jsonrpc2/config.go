package github_neonxp_jsonrpc2

import (
	"aurora-relayer-go-common/log"
	"github.com/spf13/viper"
)

const (
	defaultHttpHost = "localhost" // Default host interface for the HTTP RPC server
	defaultHttpPort = 8545        // Default TCP port for the HTTP RPC server

	configPath = "RpcNode.JsonRpc2"
)

type Config struct {
	HttpPort int16
	HttpHost string
}

func defaultConfig() *Config {
	return &Config{
		HttpPort: defaultHttpPort,
		HttpHost: defaultHttpHost,
	}
}

func GetConfig() *Config {
	config := defaultConfig()
	sub := viper.Sub(configPath)
	if sub != nil {
		if err := sub.Unmarshal(&config); err != nil {
			log.Log().Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
				"falling back to defaults", configPath, viper.ConfigFileUsed())
		}
	}
	return config
}
