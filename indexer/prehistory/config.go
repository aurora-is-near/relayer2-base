package prehistory

import (
	"github.com/aurora-is-near/relayer2-base/cmd"
	"github.com/aurora-is-near/relayer2-base/log"

	"github.com/spf13/viper"
)

const (
	defaultIndexFromPrehistory = false
	defaultFromBlock           = 0
	defaultBatchSize           = 10000
	defaultChainId             = 1313161554

	configPath = "prehistoryIndexer"
)

type Config struct {
	IndexFromPrehistory bool   `mapstructure:"indexFromPrehistory"`
	PrehistoryHeight    uint64 `mapstructure:"prehistoryHeight"`
	From                uint64 `mapstructure:"from"`
	To                  uint64 `mapstructure:"to"`
	ArchiveURL          string `mapstructure:"archiveURL"`
	BatchSize           uint64 `mapstructure:"batchSize"`
	ChainId             uint64 `mapstructure:"chainId"`
}

func defaultConfig() *Config {
	return &Config{
		IndexFromPrehistory: defaultIndexFromPrehistory,
		From:                defaultFromBlock,
		BatchSize:           defaultBatchSize,
		ChainId:             defaultChainId,
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
