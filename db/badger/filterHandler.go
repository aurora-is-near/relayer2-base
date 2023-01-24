package badger

import (
	"context"
	"errors"
	dbh "relayer2-base/db"
	"relayer2-base/db/badger/core"
	"relayer2-base/db/codec"
	dbt "relayer2-base/types/db"
	"relayer2-base/types/primitives"
	"relayer2-base/utils"
)

type FilterHandler struct {
	db     *core.DB
	config *Config
}

func NewFilterHandler() (dbh.FilterHandler, error) {
	return NewFilterHandlerWithCodec(codec.NewTinypackCodec())
}

func NewFilterHandlerWithCodec(codec codec.Codec) (dbh.FilterHandler, error) {
	config := GetConfig()
	db, err := core.NewDB(config.Core, codec)
	if err != nil {
		return nil, err
	}
	return &FilterHandler{
		db:     db,
		config: config,
	}, nil
}

func (h *FilterHandler) GetFilter(ctx context.Context, filterId primitives.Data32) (any, error) {
	var resp any
	var err error
	err = h.db.View(func(txn *core.ViewTxn) error {
		resp, err = txn.ReadFilter(utils.GetChainId(ctx), filterId)
		if resp == nil && err == nil {
			err = errors.New("filter not found")
		}
		return err
	})
	return resp, err
}

func (h *FilterHandler) GetBlockFilter(ctx context.Context, filterId primitives.Data32) (*dbt.BlockFilter, error) {
	var resp *dbt.BlockFilter
	var err error
	err = h.db.View(func(txn *core.ViewTxn) error {
		resp, err = txn.ReadBlockFilter(utils.GetChainId(ctx), filterId)
		if resp == nil && err == nil {
			err = errors.New("filter not found")
		}
		return err
	})
	return resp, err
}

func (h *FilterHandler) GetTransactionFilter(ctx context.Context, filterId primitives.Data32) (*dbt.TransactionFilter, error) {
	var resp *dbt.TransactionFilter
	var err error
	err = h.db.View(func(txn *core.ViewTxn) error {
		resp, err = txn.ReadTransactionFilter(utils.GetChainId(ctx), filterId)
		if resp == nil && err == nil {
			err = errors.New("filter not found")
		}
		return err
	})
	return resp, err
}

func (h *FilterHandler) GetLogFilter(ctx context.Context, filterId primitives.Data32) (*dbt.LogFilter, error) {
	var resp *dbt.LogFilter
	var err error
	err = h.db.View(func(txn *core.ViewTxn) error {
		resp, err = txn.ReadLogFilter(utils.GetChainId(ctx), filterId)
		if resp == nil && err == nil {
			err = errors.New("filter not found")
		}
		return err
	})
	return resp, err
}

func (h *FilterHandler) StoreFilter(ctx context.Context, filterId primitives.Data32, filter any) error {
	if bf, ok := filter.(*dbt.BlockFilter); ok {
		return h.StoreBlockFilter(ctx, filterId, bf)
	} else if tf, ok := filter.(*dbt.TransactionFilter); ok {
		return h.StoreTransactionFilter(ctx, filterId, tf)
	} else if lf, ok := filter.(*dbt.LogFilter); ok {
		return h.StoreLogFilter(ctx, filterId, lf)
	}
	return errors.New("unknown filter type")
}

func (h *FilterHandler) StoreBlockFilter(ctx context.Context, filterId primitives.Data32, filter *dbt.BlockFilter) error {
	return h.db.InsertBlockFilter(utils.GetChainId(ctx), filterId, filter)
}

func (h *FilterHandler) StoreTransactionFilter(ctx context.Context, filterId primitives.Data32, filter *dbt.TransactionFilter) error {
	return h.db.InsertTransactionFilter(utils.GetChainId(ctx), filterId, filter)
}

func (h *FilterHandler) StoreLogFilter(ctx context.Context, filterId primitives.Data32, filter *dbt.LogFilter) error {
	return h.db.InsertLogFilter(utils.GetChainId(ctx), filterId, filter)
}

func (h *FilterHandler) DeleteFilter(ctx context.Context, filterId primitives.Data32) error {
	return h.db.Update(func(txn *core.ViewTxn) error {
		chainId := utils.GetChainId(ctx)
		filter, err := txn.ReadFilter(chainId, filterId)
		if err != nil {
			return err
		}
		if filter != nil {
			if _, ok := filter.(*dbt.BlockFilter); ok {
				return txn.DeleteBlockFilter(chainId, filterId)
			} else if _, ok := filter.(*dbt.TransactionFilter); ok {
				return txn.DeleteTransactionFilter(chainId, filterId)
			} else if _, ok := filter.(*dbt.LogFilter); ok {
				return txn.DeleteLogFilter(chainId, filterId)
			} else {
				return errors.New("unknown filter type")
			}
		}
		return errors.New("filter not found")

	})
}

func (h *FilterHandler) Close() error {
	return h.db.Close()
}
