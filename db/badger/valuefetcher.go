package badger

import "sync"

type valueFetcher struct {
	scan *Scan

	input  <-chan []byte
	output chan<- Entry
	done   chan<- error

	stopChan chan bool
	wg       sync.WaitGroup
}

func (scan *Scan) startValueFetcher(input <-chan []byte, output chan<- Entry, done chan<- error) *valueFetcher {
	vf := &valueFetcher{
		scan:     scan,
		input:    input,
		output:   output,
		done:     done,
		stopChan: make(chan bool, 1),
	}

	vf.wg.Add(1)
	go vf.run()

	return vf
}

func (vf *valueFetcher) stop() {
	vf.stopChan <- true
	vf.wg.Wait()
}

func (vf *valueFetcher) run() {
	defer vf.wg.Done()

	for {
		// Prioritized stop check
		select {
		case <-vf.stopChan:
			return
		default:
		}

		select {
		case primaryKey, hasInput := <-vf.input:
			if !hasInput {
				vf.done <- nil
				return
			}
			entry := Entry{
				PrimaryKey: primaryKey,
			}
			if !vf.scan.opts.KeysOnly {
				item, err := vf.scan.txn.Get(concatBytes(vf.scan.originTablePrefix, primaryKey))
				if err != nil {
					vf.done <- err
					return
				}
				entry.Value, err = item.ValueCopy(nil)
				if err != nil {
					vf.done <- err
					return
				}
			}
			select {
			case vf.output <- entry:
			case <-vf.stopChan:
				return
			}
		case <-vf.stopChan:
			return
		}
	}
}
