package badger

import (
	"aurora-relayer-go-common/db"
	"aurora-relayer-go-common/db/badger/core"
	"aurora-relayer-go-common/db/codec"
	"aurora-relayer-go-common/utils"
	"context"
	"github.com/dgraph-io/badger/v3"
)

type FilterHandler struct {
	db     *badger.DB
	codec  codec.Codec
	config *Config
}

func NewFilterHandler() (db.FilterHandler, error) {
	return NewFilterHandlerWithCodec(codec.NewCborCodec())
}

func NewFilterHandlerWithCodec(codec codec.Codec) (db.FilterHandler, error) {
	config := GetConfig()
	bdb, err := core.Open(config.BadgerConfig, config.GcIntervalSeconds)
	if err != nil {
		return nil, err
	}
	return &FilterHandler{db: bdb,
		codec:  codec,
		config: config,
	}, nil
}

func (h *FilterHandler) StoreFilter(ctx context.Context, id utils.Uint256, filter *utils.StoredFilter) error {
	// TODO store with TTL
	return insert(h.codec, filterByIdKey(id), filter)
}

func (h *FilterHandler) GetFilter(ctx context.Context, id utils.Uint256) (*utils.StoredFilter, error) {
	return fetch[utils.StoredFilter](ctx, h.codec, filterByIdKey(id))
}

func (h *FilterHandler) DeleteFilter(ctx context.Context, id utils.Uint256) error {
	return dlt(filterByIdKey(id))
}

func (h *FilterHandler) Close() error {
	return core.Close()
}
