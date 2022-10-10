package endpoint

import (
	"aurora-relayer-go-common/log"
	"aurora-relayer-go-common/utils"
	"fmt"
	"math/big"

	"github.com/spf13/viper"
)

const (
	configPath                           = "endpoint"
	minGasPriceDefault                   = 0
	minGasLimitDefault                   = 21000
	GasForNearTxsCallDefault             = 300000000000000
	DepositForNearTxsCallDefault         = 0
	retryWaitTimeMsForNearTxsCallDefault = 3000
	retryNumberForNearTxsCallDefault     = 3
)

type EthConfig struct {
	protocolVersion utils.Uint256 `mapstructure:"protocolVersion"`
	hashrate        utils.Uint256 `mapstructure:"hashrate"`
	zeroAddress     string        `mapstructure:"zeroAddress"`
}

type Config struct {
	ProxyUrl          string
	ProxyEndpoints    map[string]bool `mapstructure:"proxyEndpoints"`
	DisabledEndpoints map[string]bool `mapstructure:"disabledEndpoints"`
	EthConfig         EthConfig       `mapstructure:"eth"`
	EngineConfig      EngineConfig    `mapstructure:"engine"`
}

type ethConfig struct {
	protocolVersion int `mapstructure:"protocolVersion"`
	hashrate        int `mapstructure:"hashrate"`
	zeroAddress     int `mapstructure:"zeroAddress"`
}

type EngineConfig struct {
	NearNetworkID                 string   `mapstructure:"nearNetworkID"`
	NearNodeURL                   string   `mapstructure:"nearNodeURL"`
	NearReceiverID                string   `mapstructure:"nearReceiverID"`
	MinGasPrice                   *big.Int `mapstructure:"minGasPrice"`
	MinGasLimit                   uint64   `mapstructure:"minGasLimit"`
	GasForNearTxsCall             uint64   `mapstructure:"gasForNearTxsCall"`
	DepositForNearTxsCall         *big.Int `mapstructure:"depositForNearTxsCall"`
	RetryWaitTimeMsForNearTxsCall int      `mapstructure:"retryWaitTimeMsForNearTxsCall"`
	RetryNumberForNearTxsCall     int      `mapstructure:"retryNumberForNearTxsCall"`
}

type engineConfig struct {
	NearNetworkID                 string `mapstructure:"nearNetworkID"`
	NearNodeURL                   string `mapstructure:"nearNodeURL"`
	NearReceiverID                string `mapstructure:"nearReceiverID"`
	MinGasPrice                   uint64 `mapstructure:"minGasPrice"`
	MinGasLimit                   uint64 `mapstructure:"minGasLimit"`
	GasForNearTxsCall             uint64 `mapstructure:"gasForNearTxsCall"`
	DepositForNearTxsCall         uint64 `mapstructure:"depositForNearTxsCall"`
	RetryWaitTimeMsForNearTxsCall int    `mapstructure:"retryWaitTimeMsForNearTxsCall"`
	RetryNumberForNearTxsCall     int    `mapstructure:"retryNumberForNearTxsCall"`
}

type proxyConfig struct {
	Url       string   `mapstructure:"url"`
	Endpoints []string `mapstructure:"endpoints"`
}

type config struct {
	ProxyConfig       proxyConfig  `mapstructure:"proxyEndpoints"`
	DisabledEndpoints []string     `mapstructure:"disabledEndpoints"`
	EthConfig         ethConfig    `mapstructure:"eth"`
	EngineConfig      engineConfig `mapstructure:"engine"`
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
			zeroAddress:     0,
		},
		EngineConfig: engineConfig{
			NearNetworkID:                 "",
			NearNodeURL:                   "",
			NearReceiverID:                "",
			MinGasPrice:                   minGasPriceDefault,
			MinGasLimit:                   minGasLimitDefault,
			GasForNearTxsCall:             GasForNearTxsCallDefault,
			DepositForNearTxsCall:         DepositForNearTxsCallDefault,
			RetryWaitTimeMsForNearTxsCall: retryWaitTimeMsForNearTxsCallDefault,
			RetryNumberForNearTxsCall:     retryNumberForNearTxsCallDefault,
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
			zeroAddress:     fmt.Sprintf("0x%040x", c.EthConfig.zeroAddress),
		},
		EngineConfig: EngineConfig{
			NearNetworkID:                 c.EngineConfig.NearNetworkID,
			NearNodeURL:                   c.EngineConfig.NearNodeURL,
			NearReceiverID:                c.EngineConfig.NearReceiverID,
			MinGasPrice:                   big.NewInt(0).SetUint64(c.EngineConfig.MinGasPrice),
			MinGasLimit:                   c.EngineConfig.MinGasLimit,
			GasForNearTxsCall:             c.EngineConfig.GasForNearTxsCall,
			DepositForNearTxsCall:         big.NewInt(0).SetUint64(c.EngineConfig.DepositForNearTxsCall),
			RetryWaitTimeMsForNearTxsCall: c.EngineConfig.RetryWaitTimeMsForNearTxsCall,
			RetryNumberForNearTxsCall:     c.EngineConfig.RetryNumberForNearTxsCall,
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
