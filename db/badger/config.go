package badger

import (
	"github.com/aurora-is-near/relayer2-base/cmdutils"
	"github.com/aurora-is-near/relayer2-base/db/badger/core"
	"github.com/aurora-is-near/relayer2-base/log"

	"github.com/dgraph-io/badger/v3"
	"github.com/spf13/viper"
)

const (
	defaultGcIntervalSeconds     = 10
	defaultLogFilterTtlMinutes   = 15
	defaultLogScanRangeThreshold = 3000
	defaultLogMaxScanIterators   = 10000
	defaultDataPath              = "/tmp/badger/data"

	configPath = "db.badger"
)

type Config struct {
	Core core.Config `mapstructure:"core"`
}

func defaultConfig() *Config {
	badgerOptions := badger.DefaultOptions(defaultDataPath)
	badgerOptions.Logger = NewBadgerLogger(log.Log())
	return &Config{
		Core: core.Config{
			MaxScanIterators:   defaultLogMaxScanIterators,
			ScanRangeThreshold: defaultLogScanRangeThreshold,
			FilterTtlMinutes:   defaultLogFilterTtlMinutes,
			GcIntervalSeconds:  defaultGcIntervalSeconds,
			BadgerConfig:       badgerOptions,
		},
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
