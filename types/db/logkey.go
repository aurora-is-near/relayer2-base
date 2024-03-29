package db

import "github.com/aurora-is-near/relayer2-base/db/badger/core/dbkey"

type LogKey struct {
	BlockHeight      uint64
	TransactionIndex uint64
	LogIndex         uint64
}

func (lk *LogKey) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&lk.BlockHeight,
		&lk.TransactionIndex,
		&lk.LogIndex,
	}, nil
}

func (lk LogKey) Prev() *LogKey {
	if lk.LogIndex == 0 {
		if lk.TransactionIndex == 0 {
			if lk.BlockHeight == 0 {
				return nil
			}
			lk.BlockHeight--
			lk.TransactionIndex = dbkey.MaxTxIndex
		} else {
			lk.TransactionIndex--
		}
		lk.LogIndex = dbkey.MaxLogIndex
	} else {
		lk.LogIndex--
	}
	return &lk
}

func (lk LogKey) Next() *LogKey {
	if lk.LogIndex == dbkey.MaxLogIndex {
		if lk.TransactionIndex == dbkey.MaxTxIndex {
			if lk.BlockHeight == dbkey.MaxBlockHeight {
				return nil
			}
			lk.BlockHeight++
			lk.TransactionIndex = 0
		} else {
			lk.TransactionIndex++
		}
		lk.LogIndex = 0
	} else {
		lk.LogIndex++
	}
	return &lk
}

// CompareTo compares receiver log key with the argument log key and returns
//	 0 if they are equal
//	 1 if receiver is greater than argument
//	-1 if receiver is smaller than argument
func (lk *LogKey) CompareTo(other *LogKey) int {
	if lk.BlockHeight < other.BlockHeight {
		return -1
	}
	if lk.BlockHeight > other.BlockHeight {
		return 1
	}
	if lk.TransactionIndex < other.TransactionIndex {
		return -1
	}
	if lk.TransactionIndex > other.TransactionIndex {
		return 1
	}
	if lk.LogIndex < other.LogIndex {
		return -1
	}
	if lk.LogIndex > other.LogIndex {
		return 1
	}
	return 0
}
