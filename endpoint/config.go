package endpoint

import (
	"fmt"
	"math/big"

	"github.com/aurora-is-near/relayer2-base/cmdutils"
	"github.com/aurora-is-near/relayer2-base/log"
	"github.com/aurora-is-near/relayer2-base/types/common"

	"github.com/spf13/viper"
)

const (
	configPath                           = "endpoint"
	asyncSendRawTxsDefault               = true
	minGasPriceDefault                   = 0
	minGasLimitDefault                   = 21000
	GasForNearTxsCallDefault             = 300000000000000
	DepositForNearTxsCallDefault         = 0
	retryWaitTimeMsForNearTxsCallDefault = 3000
	retryNumberForNearTxsCallDefault     = 3
)

type EthConfig struct {
	ProtocolVersion common.Uint256 `mapstructure:"protocolVersion"`
	Hashrate        common.Uint256 `mapstructure:"hashrate"`
	GasEstimate     common.Uint256 `mapstructure:"gasEstimate"`
	GasPrice        common.Uint256 `mapstructure:"gasPrice"`
	ZeroAddress     string         `mapstructure:"zeroAddress"`
}

type Config struct {
	ProxyUrl          string
	ProxyEndpoints    map[string]bool `mapstructure:"proxyEndpoints"`
	DisabledEndpoints map[string]bool `mapstructure:"disabledEndpoints"`
	EthConfig         EthConfig       `mapstructure:"eth"`
	EngineConfig      EngineConfig    `mapstructure:"engine"`
}

type ethConfig struct {
	ProtocolVersion int `mapstructure:"protocolVersion"`
	Hashrate        int `mapstructure:"hashrate"`
	GasEstimate     int `mapstructure:"gasEstimate"`
	GasPrice        int `mapstructure:"gasPrice"`
	ZeroAddress     int `mapstructure:"zeroAddress"`
}

type EngineConfig struct {
	NearNetworkID                 string   `mapstructure:"nearNetworkID"`
	NearNodeURL                   string   `mapstructure:"nearNodeURL"`
	Signer                        string   `mapstructure:"signer"`
	SignerKey                     string   `mapstructure:"signerKey"`
	AsyncSendRawTxs               bool     `mapstructure:"asyncSendRawTxs"`
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
	Signer                        string `mapstructure:"signer"`
	SignerKey                     string `mapstructure:"signerKey"`
	AsyncSendRawTxs               bool   `mapstructure:"asyncSendRawTxs"`
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
			ProtocolVersion: 0x41,
			Hashrate:        0,
			GasEstimate:     0x6691b7,
			GasPrice:        0x42c1d80,
			ZeroAddress:     0,
		},
		EngineConfig: engineConfig{
			NearNetworkID:                 "",
			NearNodeURL:                   "",
			Signer:                        "",
			SignerKey:                     "",
			AsyncSendRawTxs:               asyncSendRawTxsDefault,
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
		cmdutils.BindSubViper(sub, configPath)
		if err := sub.Unmarshal(&c); err != nil {
			log.Log().Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
				"falling back to defaults", configPath, viper.ConfigFileUsed())
		}
	}

	config := &Config{
		EthConfig: EthConfig{
			ProtocolVersion: common.IntToUint256(c.EthConfig.ProtocolVersion),
			Hashrate:        common.IntToUint256(c.EthConfig.Hashrate),
			GasEstimate:     common.IntToUint256(c.EthConfig.GasEstimate),
			GasPrice:        common.IntToUint256(c.EthConfig.GasPrice),
			ZeroAddress:     fmt.Sprintf("0x%040x", c.EthConfig.ZeroAddress),
		},
		EngineConfig: EngineConfig{
			NearNetworkID:                 c.EngineConfig.NearNetworkID,
			NearNodeURL:                   c.EngineConfig.NearNodeURL,
			Signer:                        c.EngineConfig.Signer,
			SignerKey:                     c.EngineConfig.SignerKey,
			AsyncSendRawTxs:               c.EngineConfig.AsyncSendRawTxs,
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
