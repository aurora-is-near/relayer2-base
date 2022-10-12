package utils

func BlockNumToEngine(b *BlockNum) *int64 {
	if b == nil {
		return nil
	}
	bInt64 := b.Int64()
	return &bInt64
}
