package db

import "aurora-relayer-go-common/db/badger/core/dbkey"

type BlockKey struct {
	Height uint64
}

func (bk *BlockKey) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&bk.Height,
	}, nil
}

func (bk BlockKey) Prev() *BlockKey {
	if bk.Height == 0 {
		return nil
	}
	bk.Height--
	return &bk
}

func (bk BlockKey) Next() *BlockKey {
	if bk.Height == dbkey.MaxBlockHeight {
		return nil
	}
	bk.Height++
	return &bk
}

// CompareTo compares receiver block key with the argument block key and returns
//	 0 if they are equal
//	 1 if receiver is greater than argument
//	-1 if receiver is smaller than argument
func (bk *BlockKey) CompareTo(other *BlockKey) int {
	if bk.Height < other.Height {
		return -1
	}
	if bk.Height > other.Height {
		return 1
	}
	return 0
}
