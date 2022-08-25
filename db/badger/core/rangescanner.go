package core

import (
	"bytes"
	"github.com/dgraph-io/badger/v3"
	"sync"
)

type rangeScanner struct {
	scan *Scan

	input  <-chan []byte
	output chan<- []byte
	done   chan<- error

	stopChan chan bool
	wg       sync.WaitGroup
}

func (scan *Scan) startRangeScanner(input <-chan []byte, output chan<- []byte, done chan<- error) *rangeScanner {
	rs := &rangeScanner{
		scan:     scan,
		input:    input,
		output:   output,
		done:     done,
		stopChan: make(chan bool, 1),
	}

	rs.wg.Add(1)
	go rs.run()

	return rs
}

func (rs *rangeScanner) stop() {
	rs.stopChan <- true
	rs.wg.Wait()
}

func (rs *rangeScanner) run() {
	defer rs.wg.Done()

	it := rs.scan.txn.NewIterator(badger.IteratorOptions{
		PrefetchSize:   100,
		PrefetchValues: len(rs.scan.downgradedFilters) > 0,
	})
	defer it.Close()

	for {
		// Prioritized stop check
		select {
		case <-rs.stopChan:
			return
		default:
		}

		select {
		case scanPrefix, haveInput := <-rs.input:
			if !haveInput {
				rs.done <- nil
				return
			}
			if !rs.iterate(it, scanPrefix) {
				return
			}
		case <-rs.stopChan:
			return
		}
	}
}

func (rs *rangeScanner) iterate(it *badger.Iterator, scanPrefix []byte) bool {
	startKey := concatBytes(scanPrefix, rs.scan.minPrimaryKey)

	for it.Seek(startKey); it.ValidForPrefix(scanPrefix); it.Next() {
		primaryKey := it.Item().Key()[len(scanPrefix):]
		if rs.scan.maxPrimaryKey != nil && bytes.Compare(primaryKey, rs.scan.maxPrimaryKey) > 0 {
			break
		}

		// Prioritized stop check
		select {
		case <-rs.stopChan:
			return false
		default:
		}

		if len(rs.scan.downgradedFilters) > 0 {
			var fits bool
			err := it.Item().Value(func(val []byte) error {
				fits = rs.checkDowngradedFilters(val)
				return nil
			})
			if err != nil {
				rs.done <- err
				return false
			}
			if !fits {
				continue
			}
		}

		primaryKeyCopy := make([]byte, len(primaryKey))
		copy(primaryKeyCopy, primaryKey)

		select {
		case rs.output <- primaryKeyCopy:
		case <-rs.stopChan:
			return false
		}
	}

	return true
}

func (rs *rangeScanner) checkDowngradedFilters(value []byte) bool {
	parseCursor := 0
	for pos := 0; pos < int(rs.scan.index.numFields); pos++ {
		if rs.scan.bitmask&(1<<pos) > 0 {
			continue
		}

		if parseCursor >= len(value) {
			return false
		}

		fieldSize := int(value[parseCursor])
		fieldStart := parseCursor + 1
		fieldEnd := fieldStart + fieldSize

		if fieldEnd > len(value) {
			return false
		}

		if filter, ok := rs.scan.downgradedFilters[pos]; ok {
			fieldValue := []byte{}
			if fieldSize > 0 {
				fieldValue = value[fieldStart:fieldEnd]
			}
			if _, ok := filter[string(fieldValue)]; !ok {
				return false
			}
		}

		parseCursor = fieldEnd
	}

	return true
}
