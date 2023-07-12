package tar

import (
	"github.com/aurora-is-near/relayer2-base/cmdutils"
	"github.com/aurora-is-near/relayer2-base/log"
	"github.com/spf13/viper"
)

const (
	DefaultIndexFromBackup = false

	configPath = "backupIndexer"
)

type Config struct {
	IndexFromBackup bool   `mapstructure:"indexFromBackup"`
	Dir             string `mapstructure:"backupDir"`
	NamePrefix      string `mapstructure:"backupNamePrefix"`
	From            uint64 `mapstructure:"from"`
	To              uint64 `mapstructure:"to"`
}

func defaultConfig() *Config {
	return &Config{
		IndexFromBackup: DefaultIndexFromBackup,
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
