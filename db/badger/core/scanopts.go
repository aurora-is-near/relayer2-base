package core

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
