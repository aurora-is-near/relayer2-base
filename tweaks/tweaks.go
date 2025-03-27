package tweaks

import (
	"log"
	"sync"

	"github.com/aurora-is-near/relayer2-base/cmdutils"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/spf13/viper"
)

const (
	configPath = "tweaks"
)

var (
	config         *Config = nil
	loadConfigOnce sync.Once
)

type Config struct {
	BaseFeePerGas        *primitives.Quantity `mapstructure:"baseFeePerGas"`
	BaseFeePerBlobGas    *primitives.Quantity `mapstructure:"baseFeePerBlobGas"`
	GasUsedRatio         *float32             `mapstructure:"gasUsedRatio"`
	MaxPriorityFeePerGas *primitives.Quantity `mapstructure:"maxPriorityFeePerGas"`
}

func defaultConfig() *Config {
	return &Config{
		BaseFeePerGas:        nil,
		BaseFeePerBlobGas:    nil,
		GasUsedRatio:         nil,
		MaxPriorityFeePerGas: nil,
	}
}

func loadConfig() {
	loadConfigOnce.Do(func() {
		config = defaultConfig()
		sub := viper.Sub(configPath)
		if sub != nil {
			cmdutils.BindSubViper(sub, configPath)
			err := sub.Unmarshal(&config, viper.DecodeHook(primitives.QuantityDecodeHook()))
			if err != nil {
				log.Printf("tweaks: unable to parse configuration: %v", err)
				sub.Debug()
			}
		}
	})
}

func BaseFeePerGas() *primitives.Quantity {
	loadConfig()
	return config.BaseFeePerGas
}

func BaseFeePerBlobGas() *primitives.Quantity {
	loadConfig()
	return config.BaseFeePerBlobGas
}

func GasUsedRatio() *float32 {
	loadConfig()
	return config.GasUsedRatio
}

func MaxPriorityFeePerGas() *primitives.Quantity {
	loadConfig()
	return config.MaxPriorityFeePerGas
}
