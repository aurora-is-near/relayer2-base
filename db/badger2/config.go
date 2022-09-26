package badger

import (
	"aurora-relayer-go-common/db/badger/core"
	"aurora-relayer-go-common/log"
	"github.com/dgraph-io/badger/v3"
	"github.com/spf13/viper"
)

const (
	defaultGcIntervalSeconds       = 10
	defaultLogFilterTtlMinutes     = 15
	defaultMaxJumps                = 1000
	defaultRangeScanners           = 4
	defaultValueFetchers           = 4
	defaultKeysOnly                = false
	defaultIterationTimeoutSeconds = 5
	defaultIterationMaxItems       = 10000
	defaultDataPath                = "/tmp/badger/data"

	configPath = "db.badger"
)

type Config struct {
	GcIntervalSeconds       int            `mapstructure:"gcIntervalSeconds"`
	IterationTimeoutSeconds uint           `mapstructure:"iterationTimeoutSeconds"`
	IterationMaxItems       uint           `mapstructure:"iterationMaxItems"`
	LogFilterTtlMinutes     int            `mapstructure:"logFilterTtlMinutes"`
	ScanConfig              core.ScanOpts  `mapstructure:"index"`
	BadgerConfig            badger.Options `mapstructure:"options"`
}

func defaultConfig() *Config {
	badgerOptions := badger.DefaultOptions(defaultDataPath)
	badgerOptions.Logger = NewBadgerLogger(log.Log())
	return &Config{
		GcIntervalSeconds:       defaultGcIntervalSeconds,
		LogFilterTtlMinutes:     defaultLogFilterTtlMinutes,
		IterationTimeoutSeconds: defaultIterationTimeoutSeconds,
		IterationMaxItems:       defaultIterationMaxItems,
		ScanConfig: core.ScanOpts{
			MaxJumps:         defaultMaxJumps,
			MaxRangeScanners: defaultRangeScanners,
			MaxValueFetchers: defaultValueFetchers,
			KeysOnly:         defaultKeysOnly,
		},
		BadgerConfig: badgerOptions,
	}
}

func GetConfig() *Config {
	config := defaultConfig()
	sub := viper.Sub(configPath)
	if sub != nil {
		if err := sub.Unmarshal(&config); err != nil {
			log.Log().Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
				"falling back to defaults", configPath, viper.ConfigFileUsed())
		}
	}
	return config
}
