package core

import (
	"aurora-relayer-go-common/db/badger/core/dbkey"
	"aurora-relayer-go-common/db/badger/core/logscan"
	tp "aurora-relayer-go-common/tinypack"
	"aurora-relayer-go-common/types/db"
	"aurora-relayer-go-common/types/primitives"
	"aurora-relayer-go-common/types/response"
	"context"
	"fmt"

	badger "github.com/dgraph-io/badger/v3"
)

func (txn *ViewTxn) ReadEarliestLogKey(chainId uint64) (*db.LogKey, error) {
	it := txn.txn.NewIterator(badger.IteratorOptions{
		Prefix: dbkey.Logs.Get(chainId),
	})
	defer it.Close()
	it.Rewind()
	key, err := readLogKeyFromLogIterator(it)
	if err != nil {
		txn.db.logger.Errorf("DB: can't read earliest log key: %v", err)
	}
	return key, err
}

func (txn *ViewTxn) ReadLatestLogKey(chainId uint64) (*db.LogKey, error) {
	it := txn.txn.NewIterator(badger.IteratorOptions{
		Reverse: true,
		Prefix:  dbkey.Logs.Get(chainId),
	})
	defer it.Close()
	it.Seek(dbkey.Log.Get(chainId, dbkey.MaxBlockHeight, dbkey.MaxTxIndex, dbkey.MaxLogIndex))
	key, err := readLogKeyFromLogIterator(it)
	if err != nil {
		txn.db.logger.Errorf("DB: can't read latest log key: %v", err)
	}
	return key, err
}

func readLogKeyFromLogIterator(it *badger.Iterator) (*db.LogKey, error) {
	if !it.Valid() {
		return nil, nil
	}
	if !dbkey.Log.Matches(it.Item().Key()) {
		err := fmt.Errorf("found unexpected key format (expected to match dbkey.Log)")
		return nil, err
	}
	return &db.LogKey{
		BlockHeight:      dbkey.Log.ReadUintVar(it.Item().Key(), 1),
		TransactionIndex: dbkey.Log.ReadUintVar(it.Item().Key(), 2),
		LogIndex:         dbkey.Log.ReadUintVar(it.Item().Key(), 3),
	}, nil
}

func (txn *ViewTxn) ReadLogs(
	ctx context.Context,
	chainId uint64,
	from *db.LogKey,
	to *db.LogKey,
	addresses []primitives.Data20,
	topics [][]primitives.Data32,
	limit int,
) ([]*response.Log, *db.LogKey, error) {

	if limit < 0 {
		limit = 100_000
	}

	if from.CompareTo(to) > 0 {
		return nil, nil, fmt.Errorf("from > to")
	}

	var addressFilter map[string]struct{}
	topicFilters := make([]map[string]struct{}, len(topics))
	featureFilters := make([][][]byte, len(topics)+1)
	addressFilter, featureFilters[0] = processFilter(addresses)
	for i, topicFilter := range topics {
		topicFilters[i], featureFilters[i+1] = processFilter(topicFilter)
	}
	scanBitmask := logscan.SelectSearchBitmask(featureFilters, txn.db.maxLogScanIterators)

	var iterator *logIterator
	var hashScanner *logHashScanner
	var fetcher *logFetcher

	if to.BlockHeight-from.BlockHeight <= uint64(txn.db.logScanRangeThreshold) || scanBitmask == 0 {
		iterator = startLogIterator(txn, chainId, from, to)
		defer iterator.stop()
		fetcher = startLogFetcher(txn, chainId, addressFilter, topicFilters, iterator.output())
	} else {
		hashScanner = startLogHashScanner(txn, chainId, from, to, featureFilters, scanBitmask)
		defer hashScanner.stop()
		fetcher = startLogFetcher(txn, chainId, addressFilter, topicFilters, hashScanner.output())
	}
	defer fetcher.stop()

	responses := []*response.Log{}

	getLastKey := func() *db.LogKey {
		if len(responses) == 0 {
			return from.Prev()
		}
		last := responses[len(responses)-1]
		return &db.LogKey{
			BlockHeight:      uint64(last.BlockNumber),
			TransactionIndex: uint64(last.TransactionIndex),
			LogIndex:         uint64(last.LogIndex),
		}
	}

	for {
		select {
		case <-ctx.Done():
			return responses, getLastKey(), ctx.Err()
		default:
		}

		select {
		case <-ctx.Done():
			return responses, getLastKey(), ctx.Err()
		case out, ok := <-fetcher.output():
			if !ok {
				return responses, to, nil
			}
			if len(responses) == limit {
				return responses, getLastKey(), ErrLimited
			}
			responses = append(responses, out)
		}
	}
}

func processFilter[LD tp.LengthDescriptor](filter []primitives.Data[LD]) (map[string]struct{}, [][]byte) {
	mapValues := make(map[string]struct{}, len(filter))
	values := make([][]byte, 0, len(filter))
	for _, value := range filter {
		if _, ok := mapValues[string(value.Content)]; !ok {
			mapValues[string(value.Content)] = struct{}{}
			values = append(values, value.Content)
		}
	}
	return mapValues, values
}
