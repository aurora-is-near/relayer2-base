package node

import (
	"fmt"
	"time"

	"github.com/aurora-is-near/relayer2-base/cmdutils"
	"github.com/aurora-is-near/relayer2-base/log"

	"github.com/spf13/viper"
)

const (
	configPath = "rpcNode"
)

var (
	defaultHttpPort           int16         = 8545
	defaultHttpHost           string        = "localhost"
	defaultHttpCompress       bool          = true
	defaultHttpTimeout        time.Duration = 300
	defaultWsHandshakeTimeout time.Duration = 10
	defaultMaxBatchRequests   uint          = 1000
)

type Config struct {
	HttpPort           int16         `mapstructure:"httpPort"`
	HttpHost           string        `mapstructure:"httpHost"`
	HttpCors           []string      `mapstructure:"httpCors"`
	HttpCompress       bool          `mapstructure:"httpCompress"`
	HttpTimeout        time.Duration `mapstructure:"httpTimeout"`
	WsPort             int16         `mapstructure:"wsPort"`
	WsHost             string        `mapstructure:"wsHost"`
	WsHandshakeTimeout time.Duration `mapstructure:"wsHandshakeTimeout"`
	MaxBatchRequests   uint          `mapstructure:"maxBatchRequests"`
}

// httpEndpoint resolves an HTTP endpoint based on the configured host interface
// and port parameters.
func (c *Config) httpEndpoint() string {
	if c.HttpHost == "" || c.HttpPort <= 0 {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.HttpHost, c.HttpPort)
}

// wsEndpoint resolves a websocket endpoint based on the configured host interface
// and port parameters.
func (c *Config) wsEndpoint() string {
	if c.WsHost == "" || c.WsPort <= 0 {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.WsHost, c.WsPort)
}

func defaultConfig() *Config {
	return &Config{
		HttpPort:           defaultHttpPort,
		HttpHost:           defaultHttpHost,
		HttpCors:           []string{},
		HttpCompress:       defaultHttpCompress,
		HttpTimeout:        defaultHttpTimeout,
		WsHandshakeTimeout: defaultWsHandshakeTimeout,
		MaxBatchRequests:   defaultMaxBatchRequests,
	}
}

func GetConfig() *Config {
	config := defaultConfig()
	sub := viper.Sub(configPath)
	if sub != nil {
		cmdutils.BindSubViper(sub, configPath)
		if err := sub.Unmarshal(&config); err != nil {
			log.Log().Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
				"falling back to defaults", configPath, viper.ConfigFileUsed())
		}
	}
	return config
}
