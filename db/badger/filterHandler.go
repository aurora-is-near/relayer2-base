package badger

import (
	"aurora-relayer-go-common/db"
	"aurora-relayer-go-common/db/badger/core"
	"aurora-relayer-go-common/utils"
	"github.com/dgraph-io/badger/v3"
)

type FilterHandler struct {
	db     *badger.DB
	codec  db.Codec
	config *Config
}

func NewFilterHandler() (db.FilterHandler, error) {
	config := GetConfig()
	codec := db.NewCborCodec()
	bdb, err := core.Open(config.BadgerConfig, config.GcIntervalSeconds)
	if err != nil {
		return nil, err
	}
	return &FilterHandler{db: bdb,
		codec:  codec,
		config: config,
	}, nil
}

func (h *FilterHandler) StoreFilter(id utils.Uint256, filter *utils.StoredFilter) error {
	// TODO store with TTL
	return insert(h.codec, filterByIdKey(id), filter)
}

func (h *FilterHandler) GetFilter(id utils.Uint256) (*utils.StoredFilter, error) {
	return fetch[utils.StoredFilter](h.codec, filterByIdKey(id))
}

func (h *FilterHandler) DeleteFilter(id utils.Uint256) error {
	return dlt(filterByIdKey(id))
}

func (h *FilterHandler) Close() error {
	return core.Close()
}
