package dbkey

import (
	dbs "aurora-relayer-go-common/db/badger2/core/dbkey/dbschema"
	"aurora-relayer-go-common/db/badger2/core/logscan"
)

const (
	bytesPerBlockHeight = 4
	bytesPerTxIndex     = 2
	bytesPerLogIndex    = 2
)

const (
	MaxBlockHeight uint64 = (1 << bytesPerBlockHeight) - 1
	MaxTxIndex     uint64 = (1 << bytesPerTxIndex) - 1
	MaxLogIndex    uint64 = (1 << bytesPerLogIndex) - 1
)

var (
	chainId     = dbs.Var(8)
	blockHeight = dbs.Var(bytesPerBlockHeight)
	txIndex     = dbs.Var(bytesPerTxIndex)
	logIndex    = dbs.Var(bytesPerLogIndex)
	hash        = dbs.Var(32)
	logScanMask = dbs.Var(1)
	logScanHash = dbs.Var(logscan.HashSize)
)

var (
	Chains          = dbs.Path(dbs.Const(0))
	Chain           = dbs.Path(dbs.Const(0), chainId)
	BlockHashes     = dbs.Path(dbs.Const(0), chainId, dbs.Const(0))
	BlockHash       = dbs.Path(dbs.Const(0), chainId, dbs.Const(0), blockHeight)
	BlocksData      = dbs.Path(dbs.Const(0), chainId, dbs.Const(1))
	BlockData       = dbs.Path(dbs.Const(0), chainId, dbs.Const(1), blockHeight)
	BlockKeysByHash = dbs.Path(dbs.Const(0), chainId, dbs.Const(2))
	BlockKeyByHash  = dbs.Path(dbs.Const(0), chainId, dbs.Const(2), hash)
	Txs             = dbs.Path(dbs.Const(0), chainId, dbs.Const(3))
	TxsForBlock     = dbs.Path(dbs.Const(0), chainId, dbs.Const(3), blockHeight)
	Tx              = dbs.Path(dbs.Const(0), chainId, dbs.Const(3), blockHeight, txIndex)
	TxHash          = dbs.Path(dbs.Const(0), chainId, dbs.Const(3), blockHeight, txIndex, dbs.Const(0))
	TxData          = dbs.Path(dbs.Const(0), chainId, dbs.Const(3), blockHeight, txIndex, dbs.Const(1))
	TxKeysByHash    = dbs.Path(dbs.Const(0), chainId, dbs.Const(4))
	TxKeyByHash     = dbs.Path(dbs.Const(0), chainId, dbs.Const(4), hash)
	Logs            = dbs.Path(dbs.Const(0), chainId, dbs.Const(5))
	LogsForBlock    = dbs.Path(dbs.Const(0), chainId, dbs.Const(5), blockHeight)
	LogsForTx       = dbs.Path(dbs.Const(0), chainId, dbs.Const(5), blockHeight, txIndex)
	Log             = dbs.Path(dbs.Const(0), chainId, dbs.Const(5), blockHeight, txIndex, logIndex)
	LogScan         = dbs.Path(dbs.Const(0), chainId, dbs.Const(6))
	LogScanForMask  = dbs.Path(dbs.Const(0), chainId, dbs.Const(6), logScanMask)
	LogScanForHash  = dbs.Path(dbs.Const(0), chainId, dbs.Const(6), logScanMask, logScanHash)
	LogScanForBlock = dbs.Path(dbs.Const(0), chainId, dbs.Const(6), logScanMask, logScanHash, blockHeight)
	LogScanForTx    = dbs.Path(dbs.Const(0), chainId, dbs.Const(6), logScanMask, logScanHash, blockHeight, txIndex)
	LogScanEntry    = dbs.Path(dbs.Const(0), chainId, dbs.Const(6), logScanMask, logScanHash, blockHeight, txIndex, logIndex)
	IndexerState    = dbs.Path(dbs.Const(1))
)
