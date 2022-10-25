package utils

import "github.com/ethereum/go-ethereum/rpc"

func BlockNumToEngine(b *BlockNum) *int64 {
	if b == nil {
		return nil
	}
	bInt64 := b.Int64()
	return &bInt64
}

func IntToBlockNum[T Integer](i T) *BlockNum {
	return &BlockNum{rpc.BlockNumber(i)}
}
