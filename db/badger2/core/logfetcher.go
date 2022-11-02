package db

import (
	"aurora-relayer-go-common/db/badger/dbkey"
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
	"aurora-relayer-go-common/db/badger2/core/dbresponses"
	"aurora-relayer-go-common/db/badger2/core/dbtypes"
	"fmt"
	"sync"
)

const logFetchGoroutines = 6
const logFetchBufferSize = 200

type logFetch struct {
	data     *dbtypes.Log
	key      *dbtypes.LogKey
	response chan *dbresponses.Log
}

type logFetcher struct {
	txn           *ViewTxn
	chainId       uint64
	addressFilter map[string]struct{}
	topicFilters  []map[string]struct{}

	input        <-chan *logFetch
	processQueue chan *logFetch
	outQueue     chan *logFetch
	out          chan *dbresponses.Log

	wg       sync.WaitGroup
	stopChan chan struct{}
}

func startLogFetcher(
	txn *ViewTxn,
	chainId uint64,
	addressFilter map[string]struct{},
	topicFilters []map[string]struct{},
	input <-chan *logFetch,
) *logFetcher {

	lf := &logFetcher{
		txn:           txn,
		chainId:       chainId,
		addressFilter: addressFilter,
		topicFilters:  topicFilters,
		input:         input,
		processQueue:  make(chan *logFetch, logFetchBufferSize),
		outQueue:      make(chan *logFetch, logFetchBufferSize),
		out:           make(chan *dbresponses.Log, logFetchBufferSize),
		stopChan:      make(chan struct{}),
	}
	lf.wg.Add(logFetchGoroutines + 2)
	go lf.runInputter()
	go lf.runOutputter()
	for i := 0; i < logFetchGoroutines; i++ {
		go lf.runProcessor()
	}
	return lf
}

func (lf *logFetcher) output() <-chan *dbresponses.Log {
	return lf.out
}

func (lf *logFetcher) stop() {
	close(lf.stopChan)
	lf.wg.Wait()
}

func (lf *logFetcher) runInputter() {
	defer lf.wg.Done()
	defer close(lf.processQueue)
	defer close(lf.outQueue)

	for {
		if lf.isStopped() {
			return
		}

		select {
		case <-lf.stopChan:
			return
		case in, ok := <-lf.input:
			if !ok {
				return
			}
			in.response = make(chan *dbresponses.Log, 1)
			select {
			case <-lf.stopChan:
				return
			case lf.outQueue <- in:
			}
			select {
			case <-lf.stopChan:
				return
			case lf.processQueue <- in:
			}
		}
	}
}

func (lf *logFetcher) runOutputter() {
	defer lf.wg.Done()
	defer close(lf.out)

	for {
		if lf.isStopped() {
			return
		}

		select {
		case <-lf.stopChan:
			return
		case out, ok := <-lf.outQueue:
			if !ok {
				return
			}
			select {
			case <-lf.stopChan:
				return
			case response := <-out.response:
				if response != nil {
					select {
					case <-lf.stopChan:
						return
					case lf.out <- response:
					}
				}
			}
		}
	}
}

func (lf *logFetcher) runProcessor() {
	defer lf.wg.Done()

	for {
		if lf.isStopped() {
			return
		}

		select {
		case <-lf.stopChan:
			return
		case item, ok := <-lf.processQueue:
			if !ok {
				return
			}
			response, err := lf.processItem(item)
			if err != nil {
				lf.txn.db.logger.Errorf("DB: unable to fetch log, will ignore it: %v", err)
			} else {
				item.response <- response
			}
		}
	}
}

func (lf *logFetcher) processItem(item *logFetch) (*dbresponses.Log, error) {
	if item.data == nil {
		var err error
		key := dbkey.Log.Get(lf.chainId, item.key.BlockHeight, item.key.TransactionIndex, item.key.LogIndex)
		item.data, err = read[dbtypes.Log](lf.txn, key)
		if err != nil || item.data == nil {
			return nil, fmt.Errorf("unable to fetch log: %v", err)
		}
	}

	if len(lf.addressFilter) > 0 {
		if _, ok := lf.addressFilter[string(item.data.Address.Bytes())]; !ok {
			return nil, nil
		}
	}
	for i, topicFilter := range lf.topicFilters {
		if len(topicFilter) == 0 {
			continue
		}
		if i >= len(item.data.Topics.Content) {
			return nil, nil
		}
		if _, ok := topicFilter[string(item.data.Topics.Content[i].Bytes())]; !ok {
			return nil, nil
		}
	}

	if lf.isStopped() {
		return nil, nil
	}

	blockHash, err := readCached[dbp.Data32](lf.txn, dbkey.BlockHash.Get(lf.chainId, item.key.BlockHeight))
	if err != nil || blockHash == nil {
		return nil, fmt.Errorf("unable to fetch blockHash: %v", err)
	}

	if lf.isStopped() {
		return nil, nil
	}

	txHash, err := readCached[dbp.Data32](lf.txn, dbkey.TxHash.Get(lf.chainId, item.key.BlockHeight, item.key.TransactionIndex))
	if err != nil || txHash == nil {
		return nil, fmt.Errorf("unable to fetch txHash: %v", err)
	}

	return makeLogResponse(
		item.key.BlockHeight,
		item.key.TransactionIndex,
		item.key.LogIndex,
		*blockHash,
		*txHash,
		item.data,
	), nil
}

func (lf *logFetcher) isStopped() bool {
	select {
	case <-lf.stopChan:
		return true
	default:
		return false
	}
}
