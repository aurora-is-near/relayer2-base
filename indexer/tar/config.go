package tar

import (
	"github.com/spf13/viper"
	"relayer2-base/cmd"
	"relayer2-base/log"
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
		cmd.BindSubViper(sub, configPath)
		if err := sub.Unmarshal(&config); err != nil {
			log.Log().Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
				"falling back to defaults", configPath, viper.ConfigFileUsed())
		}
	}
	return config
}
