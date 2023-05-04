package endpoint

import (
	"fmt"
	"math/big"

	"github.com/aurora-is-near/relayer2-base/types/common"
)

const (
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
	ProxyConfig       `mapstructure:"proxyEndpoints"`
	DisabledEndpoints map[string]bool `mapstructure:"disabledEndpoints"`
	EthConfig         EthConfig       `mapstructure:"eth"`
	EngineConfig      EngineConfig    `mapstructure:"engine"`
}

type ProxyConfig struct {
	ProxyUrl       string          `mapstructure:"url"`
	ProxyEndpoints map[string]bool `mapstructure:"endpoints"`
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

func DefaultConfig() *Config {
	return &Config{
		ProxyConfig: ProxyConfig{
			ProxyUrl:       "https://testnet.aurora.dev:443",
			ProxyEndpoints: map[string]bool{},
		},
		DisabledEndpoints: map[string]bool{},
		EthConfig: EthConfig{
			ProtocolVersion: common.IntToUint256(0x41),
			Hashrate:        common.IntToUint256(0),
			GasEstimate:     common.IntToUint256(0x6691b7),
			GasPrice:        common.IntToUint256(0x42c1d80),
			ZeroAddress:     fmt.Sprintf("0x%040x", 0),
		},
		EngineConfig: EngineConfig{
			NearNetworkID:                 "",
			NearNodeURL:                   "",
			Signer:                        "",
			SignerKey:                     "",
			AsyncSendRawTxs:               asyncSendRawTxsDefault,
			MinGasPrice:                   big.NewInt(minGasPriceDefault),
			MinGasLimit:                   minGasLimitDefault,
			GasForNearTxsCall:             GasForNearTxsCallDefault,
			DepositForNearTxsCall:         big.NewInt(DepositForNearTxsCallDefault),
			RetryWaitTimeMsForNearTxsCall: retryWaitTimeMsForNearTxsCallDefault,
			RetryNumberForNearTxsCall:     retryNumberForNearTxsCallDefault,
		},
	}
}
