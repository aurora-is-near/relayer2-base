package endpoint

import (
	"aurora-relayer-go-common/log"
	"aurora-relayer-go-common/utils"
	"fmt"
	"github.com/spf13/viper"
)

const (
	configPath = "endpoint"
)

type EthConfig struct {
	protocolVersion utils.Uint256 `mapstructure:"protocolVersion"`
	hashrate        utils.Uint256 `mapstructure:"hashrate"`
	chainId         utils.Uint256 `mapstructure:"chainId"`
	zeroAddress     string        `mapstructure:"zeroAddress"`
}

type Config struct {
	ProxyUrl          string
	ProxyEndpoints    map[string]bool `mapstructure:"proxyEndpoints"`
	DisabledEndpoints map[string]bool `mapstructure:"disabledEndpoints"`
	EthConfig         EthConfig       `mapstructure:"eth"`
}

type ethConfig struct {
	protocolVersion int `mapstructure:"protocolVersion"`
	hashrate        int `mapstructure:"hashrate"`
	chainId         int `mapstructure:"chainId"`
	zeroAddress     int `mapstructure:"zeroAddress"`
}

type proxyConfig struct {
	Url       string   `mapstructure:"url"`
	Endpoints []string `mapstructure:"endpoints"`
}

type config struct {
	ProxyConfig       proxyConfig `mapstructure:"proxyEndpoints"`
	DisabledEndpoints []string    `mapstructure:"disabledEndpoints"`
	EthConfig         ethConfig   `mapstructure:"eth"`
}

func defaultConfig() *config {
	return &config{
		ProxyConfig: proxyConfig{
			Url:       "https://testnet.aurora.dev:443",
			Endpoints: []string{},
		},
		DisabledEndpoints: []string{},
		EthConfig: ethConfig{
			protocolVersion: 0x41,
			hashrate:        0,
			chainId:         1,
			zeroAddress:     0,
		},
	}
}

func GetConfig() *Config {
	c := defaultConfig()
	sub := viper.Sub(configPath)
	if sub != nil {
		if err := sub.Unmarshal(&c); err != nil {
			log.Log().Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
				"falling back to defaults", configPath, viper.ConfigFileUsed())
		}
	}

	config := &Config{
		EthConfig: EthConfig{
			protocolVersion: utils.IntToUint256(c.EthConfig.protocolVersion),
			hashrate:        utils.IntToUint256(c.EthConfig.hashrate),
			chainId:         utils.IntToUint256(c.EthConfig.chainId),
			zeroAddress:     fmt.Sprintf("0x%040x", c.EthConfig.zeroAddress),
		},
		DisabledEndpoints: make(map[string]bool, len(c.DisabledEndpoints)),
		ProxyEndpoints:    make(map[string]bool, len(c.ProxyConfig.Endpoints)),
		ProxyUrl:          c.ProxyConfig.Url,
	}

	for _, de := range c.DisabledEndpoints {
		config.DisabledEndpoints[de] = true
	}

	for _, pe := range c.ProxyConfig.Endpoints {
		config.ProxyEndpoints[pe] = true
	}

	return config
}
