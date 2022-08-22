package db

import (
	"aurora-relayer-go-common/utils"
)

type Handler interface {
	BlockNumber() (*utils.Uint256, error)
	GetBlockByHash(hash utils.H256) (*utils.Block, error)
	GetBlockByNumber(number utils.Uint256) (*utils.Block, error)
	GetBlockTransactionCountByHash(hash utils.H256) (*int64, error)
	GetBlockTransactionCountByNumber(number utils.Uint256) (*int64, error)
	GetTransactionByHash(hash utils.H256) (*utils.Transaction, error)
	GetTransactionByBlockHashAndIndex(bh utils.H256, idx int64) (*utils.Transaction, error)
	GetTransactionByBlockNumberAndIndex(bn utils.Uint256, idx int64) (*utils.Transaction, error)
	GetLogs(filter *utils.LogFilter) (*[]utils.LogResponse, error)
	GetBlockHashesSinceNumber(number utils.Uint256) ([]utils.H256, error)

	InsertBlock(block utils.Block) error
	InsertTransaction(tx utils.Transaction, idx int, block *utils.Block) error
	InsertLog(log utils.Log, idx int, tx *utils.Transaction, txIdx int, block *utils.Block) error

	Close() error
}
