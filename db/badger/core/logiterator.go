package core

import (
	"aurora-relayer-go-common/db/badger/core/dbkey"
	dbt "aurora-relayer-go-common/types/db"
	"bytes"
	"sync"

	"github.com/dgraph-io/badger/v3"
)

const logIteratorBufferSize = 500

type logIterator struct {
	txn     *ViewTxn
	chainId uint64
	from    *dbt.LogKey
	to      *dbt.LogKey

	out      chan *logFetch
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func startLogIterator(txn *ViewTxn, chainId uint64, from *dbt.LogKey, to *dbt.LogKey) *logIterator {
	lit := &logIterator{
		txn:      txn,
		chainId:  chainId,
		from:     from,
		to:       to,
		out:      make(chan *logFetch, logIteratorBufferSize),
		stopChan: make(chan struct{}),
	}
	lit.wg.Add(1)
	go lit.run()
	return lit
}

func (lit *logIterator) output() <-chan *logFetch {
	return lit.out
}

func (lit *logIterator) stop() {
	close(lit.stopChan)
	lit.wg.Wait()
}

func (lit *logIterator) run() {
	defer lit.wg.Done()
	defer close(lit.out)

	fromKey := dbkey.Log.Get(lit.chainId, lit.from.BlockHeight, lit.from.TransactionIndex, lit.from.LogIndex)
	toKey := dbkey.Log.Get(lit.chainId, lit.to.BlockHeight, lit.to.TransactionIndex, lit.to.LogIndex)
	it := lit.txn.txn.NewIterator(badger.IteratorOptions{
		PrefetchValues: true,
		PrefetchSize:   1000,
		Prefix:         getCommonPrefix(fromKey, toKey),
	})
	defer it.Close()

	for it.Seek(fromKey); it.Valid(); it.Next() {
		select {
		case <-lit.stopChan:
			return
		default:
		}
		if bytes.Compare(it.Item().Key(), toKey) > 0 {
			return
		}

		logData, err := readItem[dbt.Log](lit.txn.db, it.Item())
		if err != nil || logData == nil {
			lit.txn.db.logger.Errorf("DB: unable to read log, will skip: %v", err)
			continue
		}

		curKey := &dbt.LogKey{
			BlockHeight:      dbkey.Log.ReadUintVar(it.Item().Key(), 1),
			TransactionIndex: dbkey.Log.ReadUintVar(it.Item().Key(), 2),
			LogIndex:         dbkey.Log.ReadUintVar(it.Item().Key(), 3),
		}

		select {
		case <-lit.stopChan:
			return
		case lit.out <- &logFetch{data: logData, key: curKey}:
		}
	}
}
