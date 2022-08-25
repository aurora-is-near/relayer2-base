package db

import (
	"aurora-relayer-go-common/utils"
	"fmt"
	"strings"
)

type Handler interface {
	Close() error
	BlockHandler
	FilterHandler
}

type BlockHandler interface {
	BlockNumber() (*utils.Uint256, error)
	GetBlockByHash(hash utils.H256) (*utils.Block, error)
	GetBlockByNumber(number utils.Uint256) (*utils.Block, error)
	GetBlockTransactionCountByHash(hash utils.H256) (int64, error)
	GetBlockTransactionCountByNumber(number utils.Uint256) (int64, error)
	GetTransactionByHash(hash utils.H256) (*utils.Transaction, error)
	GetTransactionByBlockHashAndIndex(hash utils.H256, index utils.Uint256) (*utils.Transaction, error)
	GetTransactionByBlockNumberAndIndex(number utils.Uint256, index utils.Uint256) (*utils.Transaction, error)
	GetLogs(filter utils.LogFilter) (*[]utils.LogResponse, error)
	GetBlockHashesSinceNumber(number utils.Uint256) ([]utils.H256, error)
	GetLogsForTransaction(tx *utils.Transaction) ([]*utils.LogResponse, error)
	GetTransactionsForBlock(block *utils.Block) ([]*utils.Transaction, error)

	BlockHashToNumber(hash utils.H256) (*utils.Uint256, error)
	BlockNumberToHash(number utils.Uint256) (*utils.H256, error)

	InsertBlock(block *utils.Block) error
	InsertTransaction(tx *utils.Transaction, index int, block *utils.Block) error
	InsertLog(log *utils.Log, idx int, tx *utils.Transaction, txIdx int, block *utils.Block) error

	Close() error
}

type FilterHandler interface {
	GetFilter(id utils.Uint256) (*utils.StoredFilter, error)
	StoreFilter(id utils.Uint256, filter *utils.StoredFilter) error
	DeleteFilter(id utils.Uint256) error

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
