package core

import "github.com/dgraph-io/badger/v3"

type Config struct {
	MaxScanIterators   uint           `mapstructure:"maxScanIterators"`
	ScanRangeThreshold uint           `mapstructure:"scanRangeThreshold"`
	FilterTtlMinutes   int            `mapstructure:"filterTtlMinutes"`
	GcIntervalSeconds  int            `mapstructure:"gcIntervalSeconds"`
	BadgerConfig       badger.Options `mapstructure:"options"`
}
