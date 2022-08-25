package nats

import (
	"aurora-relayer-go-common/db"
	"aurora-relayer-go-common/db/nats/core"
	"aurora-relayer-go-common/utils"
	"github.com/nats-io/nats.go"
)

type FilterHandler struct {
	kv     *nats.KeyValue
	codec  db.Codec
	config *Config
}

func NewFilterHandler() (db.FilterHandler, error) {
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
		codec:  db.NewCborCodec(),
		config: config,
	}, nil
}

func (h *FilterHandler) StoreFilter(id utils.Uint256, filter *utils.StoredFilter) error {
	buf, err := h.codec.Marshal(filter)
	if err != nil {
		return err
	}
	return h.put(id.String(), buf)
}

func (h *FilterHandler) GetFilter(id utils.Uint256) (*utils.StoredFilter, error) {
	buf, err := h.get(id.String())
	if err != nil {
		return nil, err
	}
	var storedFilter utils.StoredFilter
	err = h.codec.Unmarshal(buf, &storedFilter)
	return &storedFilter, err
}

func (h *FilterHandler) DeleteFilter(id utils.Uint256) error {
	return h.del(id.String())
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
