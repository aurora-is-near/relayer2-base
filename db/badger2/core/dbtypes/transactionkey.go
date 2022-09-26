package dbtypes

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
