package db

import (
	"aurora-relayer-go-common/utils"
	"context"
	"fmt"
	"strings"
)

type Handler interface {
	Close() error
	BlockHandler
	FilterHandler
}

type BlockHandler interface {
	BlockNumber(ctx context.Context) (*utils.Uint256, error)
	GetBlockByHash(ctx context.Context, hash utils.H256) (*utils.Block, error)
	GetBlockByNumber(ctx context.Context, number utils.Uint256) (*utils.Block, error)
	GetBlockTransactionCountByHash(ctx context.Context, hash utils.H256) (int64, error)
	GetBlockTransactionCountByNumber(ctx context.Context, number utils.Uint256) (int64, error)
	GetTransactionByHash(ctx context.Context, hash utils.H256) (*utils.Transaction, error)
	GetTransactionByBlockHashAndIndex(ctx context.Context, hash utils.H256, index utils.Uint256) (*utils.Transaction, error)
	GetTransactionByBlockNumberAndIndex(ctx context.Context, number utils.Uint256, index utils.Uint256) (*utils.Transaction, error)
	GetLogs(ctx context.Context, filter utils.LogFilter) (*[]utils.LogResponse, error)
	GetBlockHashesSinceNumber(ctx context.Context, number utils.Uint256) ([]utils.H256, error)
	GetLogsForTransaction(ctx context.Context, tx *utils.Transaction) ([]*utils.LogResponse, error)
	GetTransactionsForBlock(ctx context.Context, block *utils.Block) ([]*utils.Transaction, error)

	BlockHashToNumber(ctx context.Context, hash utils.H256) (*utils.Uint256, error)
	BlockNumberToHash(ctx context.Context, number utils.Uint256) (*utils.H256, error)
	CurrentBlockSequence(ctx context.Context) uint64

	InsertBlock(block *utils.Block) error

	Close() error
}

type FilterHandler interface {
	GetFilter(ctx context.Context, id utils.Uint256) (*utils.StoredFilter, error)
	StoreFilter(ctx context.Context, id utils.Uint256, filter *utils.StoredFilter) error
	DeleteFilter(ctx context.Context, id utils.Uint256) error

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
