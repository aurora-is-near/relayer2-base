package nats

import (
	"aurora-relayer-go-common/db"
	"aurora-relayer-go-common/db/codec"
	"aurora-relayer-go-common/db/nats/core"
	dbt "aurora-relayer-go-common/types/db"
	"aurora-relayer-go-common/types/primitives"
	"context"
	"github.com/nats-io/nats.go"
)

type FilterHandler struct {
	kv     *nats.KeyValue
	codec  codec.Codec
	config *Config
}

func NewFilterHandler() (db.FilterHandler, error) {
	return NewFilterHandlerWithCodec(codec.NewCborCodec())
}

func NewFilterHandlerWithCodec(codec codec.Codec) (db.FilterHandler, error) {
	config := GetConfig()
	conn, err := core.Open(config.NatsConfig)
	if err != nil {
		return nil, err
	}

	jetStream, err := conn.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	kv, err := jetStream.CreateKeyValue(&nats.KeyValueConfig{Bucket: config.Bucket})
	if err != nil {
		return nil, err
	}

	return &FilterHandler{
		kv:     &kv,
		codec:  codec,
		config: config,
	}, nil
}

// TODO: implement
func (h *FilterHandler) GetFilter(ctx context.Context, filterId primitives.Data32) (any, error) {
	return nil, nil
}

// TODO: implement
func (h *FilterHandler) GetBlockFilter(ctx context.Context, filterId primitives.Data32) (*dbt.BlockFilter, error) {
	return nil, nil
}

// TODO: implement
func (h *FilterHandler) GetTransactionFilter(ctx context.Context, filterId primitives.Data32) (*dbt.TransactionFilter, error) {
	return nil, nil
}

// TODO: implement
func (h *FilterHandler) GetLogFilter(ctx context.Context, filterId primitives.Data32) (*dbt.LogFilter, error) {
	return nil, nil
}

// TODO: implement
func (h *FilterHandler) StoreFilter(ctx context.Context, filterId primitives.Data32, filter any) error {
	return nil
}

// TODO: implement
func (h *FilterHandler) StoreBlockFilter(ctx context.Context, filterId primitives.Data32, filter *dbt.BlockFilter) error {
	return nil
}

// TODO: implement
func (h *FilterHandler) StoreTransactionFilter(ctx context.Context, filterId primitives.Data32, filter *dbt.TransactionFilter) error {
	return nil
}

// TODO: implement
func (h *FilterHandler) StoreLogFilter(ctx context.Context, filterId primitives.Data32, filter *dbt.LogFilter) error {
	return nil
}

// TODO: implement
func (h *FilterHandler) DeleteFilter(ctx context.Context, filterId primitives.Data32) error {
	return nil
}

func (h *FilterHandler) Close() error {
	return core.Close()
}

func (h *FilterHandler) put(key string, value []byte) error {
	_, err := (*h.kv).Put(key, value)
	return err
}

func (h *FilterHandler) get(key string) ([]byte, error) {
	v, err := (*h.kv).Get(key)
	if err != nil {
		return nil, err
	}
	return v.Value(), nil
}

func (h *FilterHandler) del(key string) error {
	return (*h.kv).Purge(key)
}
