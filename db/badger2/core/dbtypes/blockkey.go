package dbtypes

type BlockKey struct {
	Height uint64
}

func (bk *BlockKey) GetTinyPackChildrenPointers() ([]any, error) {
	return []any{
		&bk.Height,
	}, nil
}
