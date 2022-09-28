package badger

import (
	"aurora-relayer-go-common/db/badger/core"
	"aurora-relayer-go-common/db/codec"
	"github.com/dgraph-io/badger/v3"
	"golang.org/x/net/context"
)

type txnKey struct{}

func GetTxn(ctx context.Context) *badger.Txn {
	if txn, ok := ctx.Value(txnKey{}).(*badger.Txn); ok {
		return txn
	}
	return nil
}

func PutTxn(ctx context.Context, txn *badger.Txn) context.Context {
	return context.WithValue(ctx, txnKey{}, txn)
}

func fetch[T any](ctx context.Context, codec codec.Codec, key []byte) (*T, error) {
	res := new(T)
	buf, err := core.Fetch(key, GetTxn(ctx))
	if err != nil {
		return nil, err
	}
	err = codec.Unmarshal(*buf, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func fetchPrefixedWithLimitAndTimeout[T any](ctx context.Context, codec codec.Codec, limit uint, timeout uint, prefix []byte) ([]*T, error) {
	res := make([]*T, 0)
	buf, err := core.FetchPrefixedWithLimitAndTimeout(limit, timeout, prefix, GetTxn(ctx))
	if err != nil {
		return nil, err
	}
	for _, t := range buf {
		r := new(T)
		err = codec.Unmarshal(t, r)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func insert[T any](ctx context.Context, codec codec.Codec, key []byte, value T) error {
	buf, err := codec.Marshal(value)
	if err != nil {
		return err
	}
	return core.Insert(key, buf, GetTxn(ctx))
}

func insertBatch[T any](writer *badger.WriteBatch, codec codec.Codec, key []byte, value T) error {
	buf, err := codec.Marshal(value)
	if err != nil {
		return err
	}
	return core.InsertBatch(writer, key, buf)
}

func dlt(ctx context.Context, key []byte) error {
	return core.Delete(key, GetTxn(ctx))
}
