package dbtypes

import "aurora-relayer-go-common/db/badger2/core/dbkey"

type TransactionKey struct {
	BlockHeight      uint64
	TransactionIndex uint64
}

func (tk *TransactionKey) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&tk.BlockHeight,
		&tk.TransactionIndex,
	}, nil
}

func (tk TransactionKey) Prev() *TransactionKey {
	if tk.TransactionIndex == 0 {
		if tk.BlockHeight == 0 {
			return nil
		}
		tk.BlockHeight--
		tk.TransactionIndex = dbkey.MaxTxIndex
	} else {
		tk.TransactionIndex--
	}
	return &tk
}

func (tk TransactionKey) Next() *TransactionKey {
	if tk.TransactionIndex == dbkey.MaxTxIndex {
		if tk.BlockHeight == dbkey.MaxBlockHeight {
			return nil
		}
		tk.BlockHeight++
		tk.TransactionIndex = 0
	} else {
		tk.TransactionIndex++
	}
	return &tk
}

func (tk *TransactionKey) CompareTo(other *TransactionKey) int {
	if tk.BlockHeight < other.BlockHeight {
		return -1
	}
	if tk.BlockHeight > other.BlockHeight {
		return 1
	}
	if tk.TransactionIndex < other.TransactionIndex {
		return -1
	}
	if tk.TransactionIndex > other.TransactionIndex {
		return 1
	}
	return 0
}
