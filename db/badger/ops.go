package badger

import (
	"aurora-relayer-go-common/db"
	"aurora-relayer-go-common/db/badger/core"
)

func fetch[T any](codec db.Codec, key []byte) (*T, error) {
	res := new(T)
	buf, err := core.Fetch(key)
	if err != nil {
		return nil, err
	}
	err = codec.Unmarshal(*buf, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func fetchPrefixed[T any](codec db.Codec, prefix []byte) ([]*T, error) {
	res := make([]*T, 0)
	buf, err := core.FetchPrefixed(prefix)
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

func fetchPrefixedWithLimitAndTimeout[T any](codec db.Codec, limit uint, timeout uint, prefix []byte) ([]*T, error) {
	res := make([]*T, 0)
	buf, err := core.FetchPrefixedWithLimitAndTimeout(limit, timeout, prefix)
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

func insert[T any](codec db.Codec, key []byte, value T) error {
	buf, err := codec.Marshal(value)
	if err != nil {
		return err
	}
	return core.Insert(key, buf)
}

func dlt(key []byte) error {
	return core.Delete(key)
}
