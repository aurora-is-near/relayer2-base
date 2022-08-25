package core

import (
	"sync"
	"sync/atomic"

	"github.com/dgraph-io/badger/v3"
)

type Entry struct {
	PrimaryKey []byte
	Value      []byte
}

type Scan struct {
	// inputs
	index             *Index
	opts              ScanOpts
	txn               *badger.Txn
	originTablePrefix []byte
	indexTablePrefix  []byte
	fieldFilters      [][][]byte
	minPrimaryKey     []byte
	maxPrimaryKey     []byte

	// control-level
	output   chan Entry
	wg       sync.WaitGroup
	stopSent int32 // atomic bool
	stop     chan bool
	err      error

	// pre-count
	bitmask           uint8
	downgradedFilters map[int]map[string]bool
}

/*
	Starts scan asynchronously.

	fieldFilters - i-th filter is an array, which means that i-th field should have one of listed values.
	If i-th filter is empty or nil, scan will accept any value for i-th position.
	Missing filters are considered empty.
	Extra filters are not taken into account.
	Fields longer than 255-characters are chopped.

	minPrimaryKey and maxPrimaryKey are inclusive, both can be nil.
*/
func (index *Index) StartScan(
	opts *ScanOpts,
	txn *badger.Txn,
	originTablePrefix []byte,
	indexTablePrefix []byte,
	fieldFilters [][][]byte,
	minPrimaryKey []byte,
	maxPrimaryKey []byte,
) *Scan {

	s := &Scan{
		index:             index,
		opts:              *opts,
		txn:               txn,
		originTablePrefix: originTablePrefix,
		indexTablePrefix:  indexTablePrefix,
		fieldFilters:      fieldFilters,
		minPrimaryKey:     minPrimaryKey,
		maxPrimaryKey:     maxPrimaryKey,
		output:            make(chan Entry, 500),
		stopSent:          0,
		stop:              make(chan bool, 1),
	}

	s.wg.Add(1)
	go s.run()

	return s
}

/*
	Returns output channel, which will be closed automatically when scan finishes.
*/
func (s *Scan) Output() <-chan Entry {
	return s.output
}

/*
	Stops the scan if it's not finished already.
	Returns scan error if it's present.
	Must not be called concurrently.
*/
func (s *Scan) Stop() error {
	if atomic.CompareAndSwapInt32(&s.stopSent, 0, 1) {
		s.stop <- true
	}
	s.wg.Wait()
	return s.err
}

func (s *Scan) run() {
	defer s.wg.Done()
	defer close(s.output)

	s.generateBitmask()
	s.populateDowngradedFilters()

	generator := s.startRangeGenerator()
	defer generator.stop()

	rangeScannerOutput := make(chan []byte, 500)
	rangeScannerDone := make(chan error, s.opts.MaxRangeScanners)
	rangeScannerDoneCnt := 0
	for i := 0; i < int(s.opts.MaxRangeScanners); i++ {
		rs := s.startRangeScanner(generator.getOutput(), rangeScannerOutput, rangeScannerDone)
		defer rs.stop()
	}

	valueFetcherDone := make(chan error, s.opts.MaxValueFetchers)
	valueFetcherDoneCnt := 0
	for i := 0; i < int(s.opts.MaxValueFetchers); i++ {
		vf := s.startValueFetcher(rangeScannerOutput, s.output, valueFetcherDone)
		defer vf.stop()
	}

	for {
		// Prioritized stop check
		select {
		case <-s.stop:
			return
		default:
		}

		select {
		case err := <-rangeScannerDone:
			if err != nil {
				s.err = err
				return
			}
			rangeScannerDoneCnt++
			if rangeScannerDoneCnt == int(s.opts.MaxRangeScanners) {
				close(rangeScannerOutput)
			}
		case err := <-valueFetcherDone:
			if err != nil {
				s.err = err
				return
			}
			valueFetcherDoneCnt++
			if valueFetcherDoneCnt == int(s.opts.MaxValueFetchers) {
				return
			}
		case <-s.stop:
			return
		}
	}
}

func (s *Scan) generateBitmask() {
	s.bitmask = 0
	jumps, filtersCnt := uint(1), uint(0)
	for curBitmask := uint8(0); curBitmask <= s.index.maxBitmask; curBitmask++ {
		valid, curJumps, curFiltersCnt := s.validateBitmask(curBitmask)
		if valid && (curFiltersCnt > filtersCnt || (curFiltersCnt == filtersCnt && curJumps < jumps)) {
			s.bitmask, jumps, filtersCnt = curBitmask, curJumps, curFiltersCnt
		}
	}
}

func (s *Scan) validateBitmask(bitmask uint8) (bool, uint, uint) {
	jumps, filtersCnt := uint(1), uint(0)
	for pos := 0; pos < int(s.index.numFields); pos++ {
		if bitmask&(1<<pos) == 0 {
			continue
		}
		if len(s.fieldFilters) <= pos || len(s.fieldFilters[pos]) == 0 {
			return false, 0, 0
		}
		if uint64(jumps)*uint64(len(s.fieldFilters[pos])) > uint64(s.opts.MaxJumps) {
			return false, 0, 0
		}
		jumps *= uint(len(s.fieldFilters[pos]))
		filtersCnt++
	}
	return true, jumps, filtersCnt
}

func (s *Scan) populateDowngradedFilters() {
	s.downgradedFilters = make(map[int]map[string]bool)
	for pos := 0; pos < int(s.index.numFields); pos++ {
		if s.bitmask&(1<<pos) > 0 {
			continue
		}
		if len(s.fieldFilters) <= pos || len(s.fieldFilters[pos]) == 0 {
			continue
		}
		s.downgradedFilters[pos] = make(map[string]bool, len(s.fieldFilters[pos]))
		for _, filterValue := range s.fieldFilters[pos] {
			s.downgradedFilters[pos][string(filterValue)] = true
		}
	}
}
