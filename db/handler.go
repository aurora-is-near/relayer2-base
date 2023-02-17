package db

import (
	"context"
	"fmt"
	"github.com/aurora-is-near/relayer2-base/types/common"
	"github.com/aurora-is-near/relayer2-base/types/db"
	"github.com/aurora-is-near/relayer2-base/types/indexer"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/aurora-is-near/relayer2-base/types/response"
	"strings"
)

type Handler interface {
	Close() error
	BlockHandler
	FilterHandler
}

type BlockHandler interface {
	BlockNumber(ctx context.Context) (*primitives.HexUint, error)
	GetBlockByHash(ctx context.Context, hash common.H256, isFull bool) (*response.Block, error)
	GetBlockByNumber(ctx context.Context, number common.BN64, isFull bool) (*response.Block, error)
	GetBlockTransactionCountByHash(ctx context.Context, hash common.H256) (*primitives.HexUint, error)
	GetBlockTransactionCountByNumber(ctx context.Context, number common.BN64) (*primitives.HexUint, error)
	GetTransactionByHash(ctx context.Context, hash common.H256) (*response.Transaction, error)
	GetTransactionByBlockHashAndIndex(ctx context.Context, hash common.H256, index common.Uint64) (*response.Transaction, error)
	GetTransactionByBlockNumberAndIndex(ctx context.Context, number common.BN64, index common.Uint64) (*response.Transaction, error)
	GetTransactionReceipt(ctx context.Context, hash common.H256) (*response.TransactionReceipt, error)

	GetLogs(ctx context.Context, filter *db.LogFilter) ([]*response.Log, error)
	GetFilterLogs(ctx context.Context, filter *db.LogFilter) ([]*response.Log, error)
	GetFilterChanges(ctx context.Context, filter any) (*[]interface{}, error)

	BlockHashToNumber(ctx context.Context, hash common.H256) (*uint64, error)
	BlockNumberToHash(ctx context.Context, number common.BN64) (*string, error)

	InsertBlock(block *indexer.Block) error

	SetIndexerState(chainId uint64, data []byte) error
	GetIndexerState(chainId uint64) ([]byte, error)

	Close() error
}

type FilterHandler interface {
	GetFilter(ctx context.Context, filterId primitives.Data32) (any, error)
	GetBlockFilter(ctx context.Context, filterId primitives.Data32) (*db.BlockFilter, error)
	GetTransactionFilter(ctx context.Context, filterId primitives.Data32) (*db.TransactionFilter, error)
	GetLogFilter(ctx context.Context, filterId primitives.Data32) (*db.LogFilter, error)

	StoreFilter(ctx context.Context, filterId primitives.Data32, filter any) error
	StoreBlockFilter(ctx context.Context, filterId primitives.Data32, filter *db.BlockFilter) error
	StoreTransactionFilter(ctx context.Context, filterId primitives.Data32, filter *db.TransactionFilter) error
	StoreLogFilter(ctx context.Context, filterId primitives.Data32, filter *db.LogFilter) error

	DeleteFilter(ctx context.Context, filterId primitives.Data32) error
	Close() error
}

type StoreHandler struct {
	BlockHandler
	FilterHandler
}

func (h StoreHandler) Close() error {
	var errs []string
	err := h.BlockHandler.Close()
	if err != nil {
		errs = append(errs, err.Error())
	}
	err = h.FilterHandler.Close()
	if err != nil {
		errs = append(errs, err.Error())
	}
	return fmt.Errorf(strings.Join(errs, "\n"))
}
