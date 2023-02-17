package core

import (
	"bytes"
	"container/heap"
	"github.com/aurora-is-near/relayer2-base/db/badger/core/dbkey"
	"github.com/aurora-is-near/relayer2-base/db/badger/core/logscan"
	"github.com/aurora-is-near/relayer2-base/types/db"
	"sync"

	"github.com/dgraph-io/badger/v3"
)

const logScannerInitGoroutines = 8
const logScannerBuffer = 500

type logHashIterator struct {
	it          *badger.Iterator
	keyIsCached bool
	cachedKey   *db.LogKey
}

func (it *logHashIterator) getLogKey() *db.LogKey {
	if it.keyIsCached {
		return it.cachedKey
	}
	key := it.it.Item().Key()
	if !dbkey.LogScanEntry.Matches(key) {
		it.cachedKey = nil
	} else {
		it.cachedKey = &db.LogKey{
			BlockHeight:      dbkey.LogScanEntry.ReadUintVar(key, 3),
			TransactionIndex: dbkey.LogScanEntry.ReadUintVar(key, 4),
			LogIndex:         dbkey.LogScanEntry.ReadUintVar(key, 5),
		}
	}
	it.keyIsCached = true
	return it.cachedKey
}

func (it *logHashIterator) next() {
	it.keyIsCached = false
	it.it.Next()
}

type logHashIteratorHeap []*logHashIterator

func (h logHashIteratorHeap) Len() int {
	return len(h)
}

func (h logHashIteratorHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *logHashIteratorHeap) Push(x any) {
	*h = append(*h, x.(*logHashIterator))
}

func (h *logHashIteratorHeap) Pop() any {
	x := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return x
}

func (h logHashIteratorHeap) Less(i, j int) bool {
	if !h[i].it.Valid() {
		return true
	}
	if !h[j].it.Valid() {
		return false
	}
	iKey := h[i].getLogKey()
	if iKey == nil {
		return true
	}
	jKey := h[j].getLogKey()
	if jKey == nil {
		return false
	}
	return iKey.CompareTo(jKey) < 0
}

type logHashScanner struct {
	txn     *ViewTxn
	chainId uint64
	from    *db.LogKey
	to      *db.LogKey
	bitmask uint64

	hashes       chan string
	iterators    logHashIteratorHeap
	iteratorsMtx sync.Mutex

	out      chan *logFetch
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func startLogHashScanner(
	txn *ViewTxn,
	chainId uint64,
	from *db.LogKey,
	to *db.LogKey,
	featureFilters [][][]byte,
	bitmask int,
) *logHashScanner {

	hashes := logscan.GenerateSearchHashes(featureFilters, bitmask)
	ls := &logHashScanner{
		txn:       txn,
		chainId:   chainId,
		from:      from,
		to:        to,
		bitmask:   uint64(bitmask),
		hashes:    make(chan string, len(hashes)),
		iterators: make(logHashIteratorHeap, 0, len(hashes)),
		out:       make(chan *logFetch, logScannerBuffer),
		stopChan:  make(chan struct{}),
	}
	for hash := range hashes {
		ls.hashes <- hash
	}
	close(ls.hashes)
	ls.wg.Add(1)
	go ls.run()
	return ls
}

func (ls *logHashScanner) output() <-chan *logFetch {
	return ls.out
}

func (ls *logHashScanner) stop() {
	close(ls.stopChan)
	ls.wg.Wait()
	for _, it := range ls.iterators {
		it.it.Close()
	}
}

func (ls *logHashScanner) run() {
	defer ls.wg.Done()
	defer close(ls.out)

	var initWg sync.WaitGroup
	initWg.Add(logScannerInitGoroutines)
	for i := 0; i < logScannerInitGoroutines; i++ {
		go ls.runIteratorInit(&initWg)
	}
	initWg.Wait()

	select {
	case <-ls.stopChan:
		return
	default:
	}
	heap.Init(&ls.iterators)

	for {
		select {
		case <-ls.stopChan:
			return
		default:
		}
		if ls.iterators.Len() == 0 {
			return
		}
		it := ls.iterators[0]
		if !it.it.Valid() {
			it.it.Close()
			heap.Pop(&ls.iterators)
			continue
		}
		key := it.getLogKey()
		if key == nil {
			ls.txn.db.logger.Errorf("DB: got corrupted LogScanEntry key while iterating, will ignore")
		} else {
			if key.CompareTo(ls.to) > 0 {
				it.it.Close()
				heap.Pop(&ls.iterators)
				continue
			}
			select {
			case <-ls.stopChan:
				return
			case ls.out <- &logFetch{key: key}:
			}
		}
		it.next()
		heap.Fix(&ls.iterators, 0)
	}
}

func (ls *logHashScanner) runIteratorInit(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ls.stopChan:
			return
		default:
		}
		hash, ok := <-ls.hashes
		if !ok {
			return
		}

		from := dbkey.LogScanEntry.Get(ls.chainId, ls.bitmask, hash, ls.from.BlockHeight, ls.from.TransactionIndex, ls.from.LogIndex)
		to := dbkey.LogScanEntry.Get(ls.chainId, ls.bitmask, hash, ls.to.BlockHeight, ls.to.TransactionIndex, ls.to.LogIndex)
		it := ls.txn.txn.NewIterator(badger.IteratorOptions{
			Prefix: getCommonPrefix(from, to),
		})
		it.Seek(from)
		if !it.Valid() || bytes.Compare(it.Item().Key(), to) > 0 {
			it.Close()
			continue
		}

		ls.iteratorsMtx.Lock()
		ls.iterators = append(ls.iterators, &logHashIterator{it: it})
		ls.iteratorsMtx.Unlock()
	}
}
