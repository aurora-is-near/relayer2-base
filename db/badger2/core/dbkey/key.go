package dbkey

import dbs "aurora-relayer-go-common/db/badger2/core/dbkey/dbschema"

var (
	chainId     = dbs.Var(8)
	blockHeight = dbs.Var(4)
	txIndex     = dbs.Var(2)
	logIndex    = dbs.Var(2)
	hash        = dbs.Var(32)
)

var (
	Chains          = dbs.Path(dbs.Const(0))
	Chain           = dbs.Path(dbs.Const(0), chainId)
	Blocks          = dbs.Path(dbs.Const(0), chainId, dbs.Const(0))
	Block           = dbs.Path(dbs.Const(0), chainId, dbs.Const(0), blockHeight)
	BlockHash       = dbs.Path(dbs.Const(0), chainId, dbs.Const(0), blockHeight, dbs.Const(0))
	BlockData       = dbs.Path(dbs.Const(0), chainId, dbs.Const(0), blockHeight, dbs.Const(1))
	BlockKeysByHash = dbs.Path(dbs.Const(0), chainId, dbs.Const(1))
	BlockKeyByHash  = dbs.Path(dbs.Const(0), chainId, dbs.Const(1), hash)
	Txs             = dbs.Path(dbs.Const(0), chainId, dbs.Const(2))
	TxsForBlock     = dbs.Path(dbs.Const(0), chainId, dbs.Const(2), blockHeight)
	Tx              = dbs.Path(dbs.Const(0), chainId, dbs.Const(2), blockHeight, txIndex)
	TxHash          = dbs.Path(dbs.Const(0), chainId, dbs.Const(2), blockHeight, txIndex, dbs.Const(0))
	TxData          = dbs.Path(dbs.Const(0), chainId, dbs.Const(2), blockHeight, txIndex, dbs.Const(1))
	TxKeysByHash    = dbs.Path(dbs.Const(0), chainId, dbs.Const(3))
	TxKeyByHash     = dbs.Path(dbs.Const(0), chainId, dbs.Const(3), hash)
	Logs            = dbs.Path(dbs.Const(0), chainId, dbs.Const(4))
	LogsForBlock    = dbs.Path(dbs.Const(0), chainId, dbs.Const(4), blockHeight)
	LogsForTx       = dbs.Path(dbs.Const(0), chainId, dbs.Const(4), blockHeight, txIndex)
	Log             = dbs.Path(dbs.Const(0), chainId, dbs.Const(4), blockHeight, txIndex, logIndex)
	LogIndex        = dbs.Path(dbs.Const(0), chainId, dbs.Const(5))
	IndexerState    = dbs.Path(dbs.Const(1))
)
