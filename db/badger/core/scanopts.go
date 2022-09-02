package core

/*
	Can be reused
	Also can be placed right in JSON-config to have tunable behavior
*/
type ScanOpts struct {
	MaxJumps         uint `mapstructure:"maxJumps"`         // Maximum number of iterator-seeks (default 10000)
	MaxRangeScanners uint `mapstructure:"maxRangeScanners"` // Maximum number of key-iteration goroutines (default 4)
	MaxValueFetchers uint `mapstructure:"maxValueFetchers"` // Maximum number of value-fetching goroutines (default 4)
	KeysOnly         bool `mapstructure:"keysOnly"`         // default false
}
