package badger

const (
	DefaultMaxJumps         = 10000
	DefaultMaxRangeScanners = 4
	DefaultMaxValueFetchers = 4
)

/*
	Can be reused
	Also can be placed right in JSON-config to have tunable behavior
*/
type ScanOpts struct {
	MaxJumps         uint // Maximum number of iterator-seeks (default 10000)
	MaxRangeScanners uint // Maximum number of key-iteration goroutines (default 4)
	MaxValueFetchers uint // Maximum number of value-fetching goroutines (default 4)
	KeysOnly         bool // default false
}

func (opts ScanOpts) FillMissingFields() ScanOpts {
	if opts.MaxJumps == 0 {
		opts.MaxJumps = DefaultMaxJumps
	}
	if opts.MaxRangeScanners == 0 {
		opts.MaxRangeScanners = DefaultMaxRangeScanners
	}
	if opts.MaxValueFetchers == 0 {
		opts.MaxValueFetchers = DefaultMaxValueFetchers
	}
	return opts
}
