package db

import (
	"aurora-relayer-common/utils"
)

type Handler interface {
	BlockNumber() (*uint64, error)
	GetBlockByHash(hash utils.H256) (*utils.Block, error)
	GetBlockByNumber(number utils.Uint256) (*utils.Block, error)
	GetBlockTransactionCountByHash(hash utils.H256) (*uint64, error)
	GetBlockTransactionCountByNumber(number utils.Uint256) (*uint64, error)
	GetTransactionByHash(hash utils.H256) (*utils.Transaction, error)
	GetTransactionByBlockHashAndIndex(bh utils.H256, idx int64) (*utils.Transaction, error)
	GetTransactionByBlockNumberAndIndex(bn utils.Uint256, idx int64) (*utils.Transaction, error)
	GetLogs(addr utils.Address, bn utils.Uint256, topic ...[]string) (*utils.Log, error)

	InsertBlock(block utils.Block) error
	InsertTransaction(block *utils.Block, idx int, tx *utils.Transaction) error
	InsertLog(tx *utils.Transaction) error

	Close() error
}
