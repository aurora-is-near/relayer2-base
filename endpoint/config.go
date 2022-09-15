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

type ethConfig struct {
	protocolVersion int `mapstructure:"protocolVersion"`
	hashrate        int `mapstructure:"hashrate"`
	chainId         int `mapstructure:"chainId"`
	zeroAddress     int `mapstructure:"zeroAddress"`
}

type Config struct {
	ProxyEndpoints    map[string]bool `mapstructure:"ProxyEndpoints"`
	DisabledEndpoints map[string]bool `mapstructure:"DisabledEndpoints"`
	EthConfig         EthConfig       `mapstructure:"eth"`
}

type config struct {
	ProxyEndpoints    []string  `mapstructure:"ProxyEndpoints"`
	DisabledEndpoints []string  `mapstructure:"DisabledEndpoints"`
	EthConfig         ethConfig `mapstructure:"eth"`
}

func defaultConfig() *config {
	return &config{
		ProxyEndpoints:    []string{},
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
		ProxyEndpoints:    make(map[string]bool, len(c.ProxyEndpoints)),
	}

	for _, de := range c.DisabledEndpoints {
		config.DisabledEndpoints[de] = true
	}

	for _, pe := range c.ProxyEndpoints {
		config.ProxyEndpoints[pe] = true
	}

	return config
}
