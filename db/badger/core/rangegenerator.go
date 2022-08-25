package core

import (
	"bytes"
	"sync"
)

type rangeGenerator struct {
	scan *Scan

	output   chan []byte
	stopChan chan bool
	wg       sync.WaitGroup
}

func (scan *Scan) startRangeGenerator() *rangeGenerator {
	rg := &rangeGenerator{
		scan:     scan,
		output:   make(chan []byte, 500),
		stopChan: make(chan bool, 1),
	}

	rg.wg.Add(1)
	go rg.run()

	return rg
}

func (rg *rangeGenerator) getOutput() <-chan []byte {
	return rg.output
}

func (rg *rangeGenerator) stop() {
	rg.stopChan <- true
	rg.wg.Wait()
}

func (rg *rangeGenerator) run() {
	defer rg.wg.Done()
	defer close(rg.output)
	pieces := [][]byte{rg.scan.indexTablePrefix, {rg.scan.bitmask}}
	rg.iterate(0, pieces)
}

func (rg *rangeGenerator) iterate(pos uint, pieces [][]byte) bool {
	if pos == rg.scan.index.numFields {
		// Prioritized stop check
		select {
		case <-rg.stopChan:
			return false
		default:
		}

		select {
		case rg.output <- concatBytes(pieces...):
			return true
		case <-rg.stopChan:
			return false
		}
	}

	if rg.scan.bitmask&(1<<pos) == 0 {
		return rg.iterate(pos+1, pieces)
	}

	for _, fieldValue := range rg.scan.fieldFilters[pos] {
		value := fieldValue
		if len(value) > 255 {
			value = value[:255]
		}
		nextPieces := append(pieces, []byte{uint8(len(value))})
		nextPieces = append(nextPieces, value)
		if !rg.iterate(pos+1, nextPieces) {
			return false
		}
	}

	return true
}

func concatBytes(pieces ...[]byte) []byte {
	var buf bytes.Buffer
	for _, piece := range pieces {
		buf.Write(piece)
	}
	return buf.Bytes()
}
