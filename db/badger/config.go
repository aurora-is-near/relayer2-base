package badger

import (
	"github.com/dgraph-io/badger/v3"

	"github.com/aurora-is-near/relayer2-base/db/badger/core"
	"github.com/aurora-is-near/relayer2-base/log"
)

const (
	defaultGcIntervalSeconds     = 10
	defaultLogFilterTtlMinutes   = 15
	defaultLogScanRangeThreshold = 3000
	defaultLogMaxScanIterators   = 10000
	defaultDataPath              = "/tmp/badger/data"
)

type Config struct {
	Core core.Config `mapstructure:"core"`
}

func DefaultConfig() *Config {
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
