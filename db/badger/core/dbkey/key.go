package dbkey

import (
	dbs "aurora-relayer-go-common/db/badger/core/dbkey/dbschema"
	"aurora-relayer-go-common/db/badger/core/logscan"
)

const (
	bytesPerBlockHeight = 4
	bytesPerTxIndex     = 2
	bytesPerLogIndex    = 2
)

const (
	MaxBlockHeight uint64 = (1 << (bytesPerBlockHeight * 8)) - 1
	MaxTxIndex     uint64 = (1 << (bytesPerTxIndex * 8)) - 1
	MaxLogIndex    uint64 = (1 << (bytesPerLogIndex * 8)) - 1
)

var (
	chainId     = dbs.Var(8)
	blockHeight = dbs.Var(bytesPerBlockHeight)
	txIndex     = dbs.Var(bytesPerTxIndex)
	logIndex    = dbs.Var(bytesPerLogIndex)
	hash        = dbs.Var(32)
	logScanMask = dbs.Var(1)
	logScanHash = dbs.Var(logscan.HashSize)
	filterId    = dbs.Var(32)
)

var (
	Chains           = dbs.Path(dbs.Const(0))
	Chain            = dbs.Path(dbs.Const(0), chainId)
	BlockHashes      = dbs.Path(dbs.Const(0), chainId, dbs.Const(0))
	BlockHash        = dbs.Path(dbs.Const(0), chainId, dbs.Const(0), blockHeight)
	BlocksData       = dbs.Path(dbs.Const(0), chainId, dbs.Const(1))
	BlockData        = dbs.Path(dbs.Const(0), chainId, dbs.Const(1), blockHeight)
	BlockKeysByHash  = dbs.Path(dbs.Const(0), chainId, dbs.Const(2))
	BlockKeyByHash   = dbs.Path(dbs.Const(0), chainId, dbs.Const(2), hash)
	TxHashes         = dbs.Path(dbs.Const(0), chainId, dbs.Const(3))
	TxHashesForBlock = dbs.Path(dbs.Const(0), chainId, dbs.Const(3), blockHeight)
	TxHash           = dbs.Path(dbs.Const(0), chainId, dbs.Const(3), blockHeight, txIndex)
	TxsData          = dbs.Path(dbs.Const(0), chainId, dbs.Const(4))
	TxsDataForBlock  = dbs.Path(dbs.Const(0), chainId, dbs.Const(4), blockHeight)
	TxData           = dbs.Path(dbs.Const(0), chainId, dbs.Const(4), blockHeight, txIndex)
	TxKeysByHash     = dbs.Path(dbs.Const(0), chainId, dbs.Const(5))
	TxKeyByHash      = dbs.Path(dbs.Const(0), chainId, dbs.Const(5), hash)
	Logs             = dbs.Path(dbs.Const(0), chainId, dbs.Const(6))
	LogsForBlock     = dbs.Path(dbs.Const(0), chainId, dbs.Const(6), blockHeight)
	LogsForTx        = dbs.Path(dbs.Const(0), chainId, dbs.Const(6), blockHeight, txIndex)
	Log              = dbs.Path(dbs.Const(0), chainId, dbs.Const(6), blockHeight, txIndex, logIndex)
	LogScan          = dbs.Path(dbs.Const(0), chainId, dbs.Const(7))
	LogScanForMask   = dbs.Path(dbs.Const(0), chainId, dbs.Const(7), logScanMask)
	LogScanForHash   = dbs.Path(dbs.Const(0), chainId, dbs.Const(7), logScanMask, logScanHash)
	LogScanForBlock  = dbs.Path(dbs.Const(0), chainId, dbs.Const(7), logScanMask, logScanHash, blockHeight)
	LogScanForTx     = dbs.Path(dbs.Const(0), chainId, dbs.Const(7), logScanMask, logScanHash, blockHeight, txIndex)
	LogScanEntry     = dbs.Path(dbs.Const(0), chainId, dbs.Const(7), logScanMask, logScanHash, blockHeight, txIndex, logIndex)
	Filters          = dbs.Path(dbs.Const(0), chainId, dbs.Const(8))
	BlockFilters     = dbs.Path(dbs.Const(0), chainId, dbs.Const(8), dbs.Const(0))
	BlockFilter      = dbs.Path(dbs.Const(0), chainId, dbs.Const(8), dbs.Const(0), filterId)
	TxFilters        = dbs.Path(dbs.Const(0), chainId, dbs.Const(8), dbs.Const(1))
	TxFilter         = dbs.Path(dbs.Const(0), chainId, dbs.Const(8), dbs.Const(1), filterId)
	LogFilters       = dbs.Path(dbs.Const(0), chainId, dbs.Const(8), dbs.Const(2))
	LogFilter        = dbs.Path(dbs.Const(0), chainId, dbs.Const(8), dbs.Const(2), filterId)
	IndexerState     = dbs.Path(dbs.Const(1))
)
