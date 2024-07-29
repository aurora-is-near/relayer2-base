package core

import "github.com/dgraph-io/badger/v3"

type Config struct {
	MaxScanIterators     uint           `mapstructure:"maxScanIterators"`
	ScanRangeThreshold   uint           `mapstructure:"scanRangeThreshold"`
	FilterTtlMinutes     int            `mapstructure:"filterTtlMinutes"`
	GcIntervalSeconds    int            `mapstructure:"gcIntervalSeconds"`
	RecreateOnCorruption bool           `mapstructure:"recreateOnCorruption"`
	BadgerConfig         badger.Options `mapstructure:"options"`
}
