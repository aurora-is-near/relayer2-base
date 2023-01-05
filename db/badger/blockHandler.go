package badger

import (
	dbh "aurora-relayer-go-common/db"
	"aurora-relayer-go-common/db/badger/core"
	"aurora-relayer-go-common/db/badger/core/dbkey"
	"aurora-relayer-go-common/db/codec"
	"aurora-relayer-go-common/types/common"
	dbt "aurora-relayer-go-common/types/db"
	"aurora-relayer-go-common/types/indexer"
	"aurora-relayer-go-common/types/primitives"
	"aurora-relayer-go-common/types/response"
	"aurora-relayer-go-common/utils"
	"context"
	"github.com/pkg/errors"
)

var (
	keyNotFoundError = errors.New("key not found")
)

type BlockHandler struct {
	db     *core.DB
	config *Config
}

func NewBlockHandler() (dbh.BlockHandler, error) {
	return NewBlockHandlerWithCodec(codec.NewTinypackCodec())
}

func NewBlockHandlerWithCodec(codec codec.Codec) (dbh.BlockHandler, error) {
	config := GetConfig()
	db, err := core.NewDB(config.Core, codec)
	if err != nil {
		return nil, err
	}
	return &BlockHandler{
		db:     db,
		config: config,
	}, nil
}

func (h *BlockHandler) Close() error {
	return h.db.Close()
}

func (h *BlockHandler) BlockNumber(ctx context.Context) (*primitives.HexUint, error) {
	var bn primitives.HexUint
	err := h.db.View(func(txn *core.ViewTxn) error {
		key, err := txn.ReadLatestBlockKey(utils.GetChainId(ctx))
		if err != nil {
			return err
		}
		if key == nil {
			return keyNotFoundError
		}
		bn = primitives.HexUint(key.Height)
		return nil
	})
	return &bn, err
}

func (h *BlockHandler) GetBlockByHash(ctx context.Context, hash common.H256, isFull bool) (*response.Block, error) {
	var resp *response.Block
	var err error
	bh := primitives.DataFromHex[primitives.Len32](hash.String())
	err = h.db.View(func(txn *core.ViewTxn) error {
		chainId := utils.GetChainId(ctx)
		key, err := txn.ReadBlockKey(chainId, bh)
		if err != nil {
			return err
		}
		if key == nil {
			return keyNotFoundError
		}
		resp, err = txn.ReadBlock(chainId, *key, isFull)
		return err
	})
	return resp, err
}

func (h *BlockHandler) GetBlockByNumber(ctx context.Context, number common.BN64, isFull bool) (*response.Block, error) {
	var resp *response.Block
	var err error
	err = h.db.View(func(txn *core.ViewTxn) error {
		var key *dbt.BlockKey
		bn := number.Uint64()
		chainId := utils.GetChainId(ctx)
		if bn == nil {
			key, err = txn.ReadLatestBlockKey(chainId)
			if err != nil {
				return err
			}
		} else if *bn == 0 {
			key, err = txn.ReadEarliestBlockKey(chainId)
			if err != nil {
				return err
			}
		} else {
			key = &dbt.BlockKey{Height: *bn}
		}
		if key != nil {
			resp, err = txn.ReadBlock(chainId, *key, isFull)
		}
		return err
	})
	return resp, err
}

func (h *BlockHandler) GetBlockTransactionCountByHash(ctx context.Context, hash common.H256) (*primitives.HexUint, error) {
	var resp primitives.HexUint
	var err error
	bh := primitives.DataFromHex[primitives.Len32](hash.String())
	err = h.db.View(func(txn *core.ViewTxn) error {
		chainId := utils.GetChainId(ctx)
		key, err := txn.ReadBlockKey(chainId, bh)
		if err != nil {
			return err
		}
		if key == nil {
			return keyNotFoundError
		}
		resp, err = txn.ReadBlockTxCount(chainId, *key)
		return err
	})
	return &resp, err
}

func (h *BlockHandler) GetBlockTransactionCountByNumber(ctx context.Context, number common.BN64) (*primitives.HexUint, error) {
	var resp primitives.HexUint
	var err error
	err = h.db.View(func(txn *core.ViewTxn) error {
		var key *dbt.BlockKey
		bn := number.Uint64()
		chainId := utils.GetChainId(ctx)
		if bn == nil {
			key, err = txn.ReadLatestBlockKey(chainId)
			if err != nil {
				return err
			}
		} else {
			key = &dbt.BlockKey{Height: *bn}
		}
		resp, err = txn.ReadBlockTxCount(chainId, *key)
		return err
	})
	return &resp, err
}

func (h *BlockHandler) GetTransactionByHash(ctx context.Context, hash common.H256) (*response.Transaction, error) {
	var resp *response.Transaction
	var err error
	th := primitives.DataFromHex[primitives.Len32](hash.String())
	err = h.db.View(func(txn *core.ViewTxn) error {
		var key *dbt.TransactionKey
		chainId := utils.GetChainId(ctx)
		key, err = txn.ReadTxKey(chainId, th)
		if err != nil {
			return err
		}
		if key == nil {
			return keyNotFoundError
		}
		resp, err = txn.ReadTx(chainId, *key)
		return err
	})
	return resp, err
}

func (h *BlockHandler) GetTransactionByBlockHashAndIndex(ctx context.Context, hash common.H256, index common.Uint64) (*response.Transaction, error) {
	var resp *response.Transaction
	var err error
	bh := primitives.DataFromHex[primitives.Len32](hash.String())
	err = h.db.View(func(txn *core.ViewTxn) error {
		chainId := utils.GetChainId(ctx)
		key, err := txn.ReadBlockKey(chainId, bh)
		if err != nil {
			return err
		}
		if key == nil {
			return keyNotFoundError
		}
		resp, err = txn.ReadTx(chainId, dbt.TransactionKey{
			BlockHeight:      key.Height,
			TransactionIndex: index.Uint64(),
		})
		return err
	})
	return resp, err
}

func (h *BlockHandler) GetTransactionByBlockNumberAndIndex(ctx context.Context, number common.BN64, index common.Uint64) (*response.Transaction, error) {
	var resp *response.Transaction
	var err error
	err = h.db.View(func(txn *core.ViewTxn) error {
		var key *dbt.BlockKey
		bn := number.Uint64()
		chainId := utils.GetChainId(ctx)
		if bn == nil {
			key, err = txn.ReadLatestBlockKey(chainId)
			if err != nil {
				return err
			}
			bn = &key.Height
		}
		resp, err = txn.ReadTx(chainId, dbt.TransactionKey{
			BlockHeight:      *bn,
			TransactionIndex: index.Uint64(),
		})
		return err
	})
	return resp, err
}

func (h *BlockHandler) GetTransactionReceipt(ctx context.Context, hash common.H256) (*response.TransactionReceipt, error) {
	var resp *response.TransactionReceipt
	var err error
	th := primitives.DataFromHex[primitives.Len32](hash.String())
	err = h.db.View(func(txn *core.ViewTxn) error {
		var key *dbt.TransactionKey
		chainId := utils.GetChainId(ctx)
		key, err = txn.ReadTxKey(chainId, th)
		if err != nil {
			return err
		}
		if key == nil {
			return keyNotFoundError
		}
		resp, err = txn.ReadTxReceipt(chainId, *key)
		return err
	})
	return resp, err
}

func (h *BlockHandler) GetLogs(ctx context.Context, filter *dbt.LogFilter) ([]*response.Log, error) {
	var resp []*response.Log
	var err error
	err = h.db.View(func(txn *core.ViewTxn) error {
		resp, _, err = h.getLogs(ctx, txn, filter, true)
		return err
	})
	return resp, err
}

func (h *BlockHandler) GetFilterLogs(ctx context.Context, filter *dbt.LogFilter) ([]*response.Log, error) {
	return h.GetLogs(ctx, filter)
}

func (h *BlockHandler) GetFilterChanges(ctx context.Context, filter any) (*[]interface{}, error) {

	var err error
	filterChanges := make([]interface{}, 0)
	if bf, ok := filter.(*dbt.BlockFilter); ok {
		var blockHashes []primitives.Data32
		var lastKey *dbt.BlockKey
		err = h.db.View(func(txn *core.ViewTxn) error {
			blockHashes, lastKey, err = h.getBlockHashes(ctx, txn, bf)
			if lastKey != nil && lastKey.CompareTo(&bf.From) > -1 {
				bf.Next = *lastKey.Next()
			}
			return err
		})
		for _, log := range blockHashes {
			filterChanges = append(filterChanges, log)
		}
	} else if tf, ok := filter.(*dbt.TransactionFilter); ok {
		var txnHashes []any
		var lastKey *dbt.TransactionKey
		err = h.db.View(func(txn *core.ViewTxn) error {
			txnHashes, lastKey, err = h.getTransactionHashes(ctx, txn, tf)
			if lastKey != nil && lastKey.CompareTo(&tf.From) > -1 {
				tf.Next = *lastKey.Next()
			}
			if txnHashes != nil && len(txnHashes) > 0 {
				filterChanges = txnHashes
			}
			return err
		})
	} else if lf, ok := filter.(*dbt.LogFilter); ok {
		var logs []*response.Log
		var lastKey *dbt.LogKey
		err = h.db.View(func(txn *core.ViewTxn) error {
			logs, lastKey, err = h.getLogs(ctx, txn, lf, false)
			if lastKey != nil && lastKey.CompareTo(&lf.From) > -1 {
				lf.Next = *lastKey.Next()
			}
			return err
		})
		for _, log := range logs {
			filterChanges = append(filterChanges, log)
		}
	}

	return &filterChanges, nil
}

func (h *BlockHandler) BlockHashToNumber(ctx context.Context, hash common.H256) (*uint64, error) {
	var resp uint64
	var err error
	bh := primitives.DataFromHex[primitives.Len32](hash.String())
	err = h.db.View(func(txn *core.ViewTxn) error {
		chainId := utils.GetChainId(ctx)
		key, err := txn.ReadBlockKey(chainId, bh)
		if err != nil {
			return err
		}
		if key == nil {
			return keyNotFoundError
		}
		resp = key.Height
		return err
	})
	return &resp, err
}

func (h *BlockHandler) BlockNumberToHash(ctx context.Context, number common.BN64) (*string, error) {
	var resp string
	var err error
	err = h.db.View(func(txn *core.ViewTxn) error {
		var b *response.Block
		chainId := utils.GetChainId(ctx)
		b, err = txn.ReadBlock(chainId, dbt.BlockKey{
			Height: *number.Uint64(),
		}, false)
		if err != nil {
			return err
		}
		if b == nil {
			return keyNotFoundError
		}
		resp = b.Hash.Hex()
		return err
	})
	return &resp, err
}

func (h *BlockHandler) InsertBlock(block *indexer.Block) error {

	writer := h.db.NewWriter()
	defer writer.Cancel()

	chainId := block.ChainId
	height := block.Height
	hash := block.Hash
	err := writer.InsertBlock(chainId, height, hash, utils.IndexerBlockToDbBlock(block))
	if err != nil {
		return err
	}

	for i, t := range block.Transactions {
		txnIndex := uint64(i)
		err = writer.InsertTransaction(chainId, height, txnIndex, t.Hash, utils.IndexerTxnToDbTxn(t))
		if err != nil {
			return err
		}
		for j, l := range t.Logs {
			err = writer.InsertLog(chainId, height, txnIndex, uint64(j), utils.IndexerLogToDbLog(l))
			if err != nil {
				return err
			}
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (h *BlockHandler) getLogs(ctx context.Context, txn *core.ViewTxn, filter *dbt.LogFilter, ignoreNext bool) ([]*response.Log, *dbt.LogKey, error) {
	var from, to, lastKey *dbt.LogKey
	var err error
	var resp []*response.Log
	var bk *dbt.BlockKey

	chainId := utils.GetChainId(ctx)

	if ignoreNext || (filter.Next.BlockHeight == 0 && filter.Next.TransactionIndex == 0 && filter.Next.LogIndex == 0) {
		// for GetLogs and GetFilterLogs, use the initial filter definition, i.e.: ignoreNext = true => use initial 'from'
		// or in case 'next' is all zero, then also use initial 'from'
		if filter.From.BlockHeight == 0 && filter.From.TransactionIndex == 0 && filter.From.LogIndex == 0 {
			// use the latest block key if initial 'from' is all zero
			bk, err = txn.ReadLatestBlockKey(chainId)
			if err != nil {
				return nil, nil, err
			}
			from = &dbt.LogKey{BlockHeight: bk.Height, TransactionIndex: 0, LogIndex: 0}
		} else {
			from = &filter.From
		}
	} else {
		// for GetFilterChanges (i.e.: ignoreNext = false) and non-zero 'next' case use next as 'from'
		from = &filter.Next
	}

	if filter.To.BlockHeight == 0 && filter.To.TransactionIndex == dbkey.MaxTxIndex && filter.To.BlockHeight == dbkey.MaxLogIndex {
		// use the latest block key if initial 'to' is all set to defaults
		if bk == nil {
			bk, err = txn.ReadLatestBlockKey(chainId)
			if err != nil {
				return nil, nil, err
			}
		}
		to = &dbt.LogKey{BlockHeight: bk.Height, TransactionIndex: dbkey.MaxTxIndex, LogIndex: dbkey.MaxLogIndex}
	} else {
		to = &filter.To
	}

	addresses := filter.Addresses.Content
	var topics [][]primitives.Data32
	for i, t := range filter.Topics.Content {
		topics[i] = t.Content
	}

	if from.CompareTo(to) > 0 {
		// unlikely but can happen if, for the same filter, getFilterChanges called twice within one block indexing time.
		// i.e.: latest block not changed on DB since the last call, in such case returns next.prev() which is the last
		// log's key returned. => the caller increments it to set next again if necessary, see BlockHandler.GetLogs and
		// BlockHandler.GetFilterChanges.
		return resp, filter.Next.Prev(), nil
	} else {
		resp, lastKey, err = txn.ReadLogs(ctx, chainId, from, to, addresses, topics, 1000)
		if err != nil && len(resp) > 0 {
			return resp, lastKey, nil
		}
	}
	return resp, lastKey, err
}

func (h *BlockHandler) getBlockHashes(ctx context.Context, txn *core.ViewTxn, filter *dbt.BlockFilter) ([]primitives.Data32, *dbt.BlockKey, error) {
	var resp []primitives.Data32
	var err error
	var from, to, lastKey *dbt.BlockKey

	chainId := utils.GetChainId(ctx)

	if filter.Next.Height == 0 {
		from = &filter.From
	} else {
		from = &filter.Next
	}

	if filter.To.Height == 0 {
		var bk *dbt.BlockKey
		bk, err = txn.ReadLatestBlockKey(chainId)
		if err != nil {
			return nil, nil, err
		}
		to = &dbt.BlockKey{Height: bk.Height}
	} else {
		to = &filter.To
	}

	resp, lastKey, err = txn.ReadBlockHashes(ctx, chainId, from, to, 1000)
	if err != nil && len(resp) > 0 {
		return resp, lastKey, nil
	}
	return resp, lastKey, err
}

func (h *BlockHandler) getTransactionHashes(ctx context.Context, txn *core.ViewTxn, filter *dbt.TransactionFilter) ([]any, *dbt.TransactionKey, error) {
	var resp []any
	var err error
	var from, to, lastKey *dbt.TransactionKey

	chainId := utils.GetChainId(ctx)

	if filter.Next.BlockHeight == 0 && filter.Next.TransactionIndex == 0 {
		from = &filter.From
	} else {
		from = &filter.Next
	}

	if filter.To.BlockHeight == 0 && filter.To.TransactionIndex == 0 {
		var bk *dbt.BlockKey
		bk, err = txn.ReadLatestBlockKey(chainId)
		if err != nil {
			return nil, nil, err
		}
		to = &dbt.TransactionKey{BlockHeight: bk.Height, TransactionIndex: 0}
	} else {
		to = &filter.To
	}

	resp, lastKey, err = txn.ReadTransactions(ctx, chainId, from, to, false, 1000)
	if err != nil && len(resp) > 0 {
		return resp, lastKey, nil
	}
	return resp, lastKey, err
}
