package github_neonxp_jsonrpc2

import (
	"relayer2-base/cmd"
	"relayer2-base/log"

	"github.com/spf13/viper"
)

const (
	defaultHttpHost = "localhost" // Default host interface for the HTTP RPC server
	defaultHttpPort = 8545        // Default TCP port for the HTTP RPC server

	configPath = "rpcNode.jsonRpc2"
)

type Config struct {
	HttpPort int16  `mapstructure:"httpPort"`
	HttpHost string `mapstructure:"httpHost"`
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
		cmd.BindSubViper(sub, configPath)
		if err := sub.Unmarshal(&config); err != nil {
			log.Log().Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
				"falling back to defaults", configPath, viper.ConfigFileUsed())
		}
	}
	return config
}
