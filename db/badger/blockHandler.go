package badger

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"github.com/aurora-is-near/relayer2-base/db/badger/core"
	"github.com/aurora-is-near/relayer2-base/db/badger/core/dbkey"
	"github.com/aurora-is-near/relayer2-base/db/codec"
	"github.com/aurora-is-near/relayer2-base/types/common"
	dbt "github.com/aurora-is-near/relayer2-base/types/db"
	errs "github.com/aurora-is-near/relayer2-base/types/errors"
	"github.com/aurora-is-near/relayer2-base/types/indexer"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/aurora-is-near/relayer2-base/types/response"
	"github.com/aurora-is-near/relayer2-base/utils"
)

type BlockHandler struct {
	Config *Config
	db     *core.DB
}

func NewBlockHandler() (*BlockHandler, error) {
	return NewBlockHandlerWithCodec(codec.NewTinypackCodec())
}

func NewBlockHandlerWithCodec(codec codec.Codec) (*BlockHandler, error) {
	config := GetConfig()
	db, err := core.NewDB(config.Core, codec)
	if err != nil {
		return nil, err
	}
	return &BlockHandler{
		db:     db,
		Config: config,
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
			return &errs.KeyNotFoundError{}
		}
		bn = primitives.HexUint(key.Height)
		return nil
	})
	return &bn, err
}

func (h *BlockHandler) GetBlockByHash(ctx context.Context, hash common.H256, isFull bool) (*response.Block, error) {
	var resp *response.Block
	var err error
	bh := hash.Data32
	err = h.db.View(func(txn *core.ViewTxn) error {
		chainId := utils.GetChainId(ctx)
		key, err := txn.ReadBlockKey(chainId, bh)
		if err != nil {
			return err
		}
		if key == nil {
			// Provided hash not found in DB, check if this is a prehistory hash
			prehistoryChainId := utils.GetPrehistoryChainId()
			if chainId == prehistoryChainId {
				return &errs.KeyNotFoundError{}
			} else {
				blockHeight, err := decodeBlockHeight(bh.Content)
				if blockHeight == nil || *blockHeight > utils.GetPrehistoryHeight() {
					return &errs.KeyNotFoundError{}
				}
				if err != nil {
					return &errs.InvalidParamsError{Message: err.Error()}
				}

				key = &dbt.BlockKey{Height: *blockHeight}
				resp, err = txn.ReadBlock(prehistoryChainId, *key, isFull)
				if err != nil {
					return &errs.KeyNotFoundError{}
				}
				if resp != nil {
					resp, err = postProcessPrehistoryBlock(resp, key.Height, chainId)
					if err != nil {
						return err
					}
				}
				return nil
			}
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
		var skipPrehistoryChecks bool
		bn := number.Uint64()
		chainId := utils.GetChainId(ctx)
		prehistoryChainId := utils.GetPrehistoryChainId()
		if chainId == prehistoryChainId {
			skipPrehistoryChecks = true
		}
		if bn == nil {
			key, err = txn.ReadLatestBlockKey(chainId)
			if err != nil {
				return err
			}
			// Check prehistory blocks if latest block key is nil and prehistory has a different chainId
			if key == nil && !skipPrehistoryChecks {
				key, err = txn.ReadLatestBlockKey(prehistoryChainId)
				if err != nil {
					return err
				}
			}
		} else if *bn == 0 {
			key = nil
			// Check prehistory blocks first if prehistory has a different chainId
			if !skipPrehistoryChecks {
				key, err = txn.ReadEarliestBlockKey(prehistoryChainId)
				if err != nil {
					return err
				}
			}
			// If prehistory has the same chain with relayer or key not found in prehistory, then check the relayer chain
			if key == nil {
				key, err = txn.ReadEarliestBlockKey(chainId)
				if err != nil {
					return err
				}
			}
		} else {
			key = &dbt.BlockKey{Height: *bn}
		}
		if key != nil {
			if key.Height >= utils.GetPrehistoryHeight() || skipPrehistoryChecks {
				resp, err = txn.ReadBlock(chainId, *key, isFull)
			} else {
				resp, err = txn.ReadBlock(prehistoryChainId, *key, isFull)
				if err != nil {
					return err
				}
				if resp != nil {
					resp, err = postProcessPrehistoryBlock(resp, key.Height, chainId)
				}
			}
		}
		return err
	})
	return resp, err
}

func (h *BlockHandler) GetBlockTransactionCountByHash(ctx context.Context, hash common.H256) (*primitives.HexUint, error) {
	var resp primitives.HexUint
	var err error
	bh := hash.Data32
	err = h.db.View(func(txn *core.ViewTxn) error {
		chainId := utils.GetChainId(ctx)
		key, err := txn.ReadBlockKey(chainId, bh)
		if err != nil {
			return err
		}
		// Should return 0 for transaction count if no block found for the provided hash
		if key == nil {
			resp = 0
			return nil
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
	th := hash.Data32
	err = h.db.View(func(txn *core.ViewTxn) error {
		var key *dbt.TransactionKey
		chainId := utils.GetChainId(ctx)
		key, err = txn.ReadTxKey(chainId, th)
		if err != nil {
			return err
		}
		if key == nil {
			return &errs.KeyNotFoundError{}
		}
		resp, err = txn.ReadTx(chainId, *key)
		return err
	})
	return resp, err
}

func (h *BlockHandler) GetTransactionByBlockHashAndIndex(ctx context.Context, hash common.H256, index common.Uint64) (*response.Transaction, error) {
	var resp *response.Transaction
	var err error
	bh := hash.Data32
	err = h.db.View(func(txn *core.ViewTxn) error {
		chainId := utils.GetChainId(ctx)
		key, err := txn.ReadBlockKey(chainId, bh)
		if err != nil {
			return err
		}
		if key == nil {
			return &errs.KeyNotFoundError{}
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
	th := hash.Data32
	err = h.db.View(func(txn *core.ViewTxn) error {
		var key *dbt.TransactionKey
		chainId := utils.GetChainId(ctx)
		key, err = txn.ReadTxKey(chainId, th)
		if err != nil {
			return err
		}
		if key == nil {
			return &errs.KeyNotFoundError{}
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

func (h *BlockHandler) GetIndexerState(chainId uint64) ([]byte, error) {
	var resp []byte
	err := h.db.View(func(txn *core.ViewTxn) error {
		data, err := txn.ReadIndexerState(chainId)
		if err != nil {
			return err
		}
		resp = data
		return nil
	})
	return resp, err
}

func (h *BlockHandler) BlockHashToNumber(ctx context.Context, hash common.H256) (*uint64, error) {
	var resp uint64
	var err error
	bh := hash.Data32
	err = h.db.View(func(txn *core.ViewTxn) error {
		chainId := utils.GetChainId(ctx)
		key, err := txn.ReadBlockKey(chainId, bh)
		if err != nil {
			return err
		}
		if key == nil {
			return &errs.KeyNotFoundError{}
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
			return &errs.KeyNotFoundError{}
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

	gasUsed := big.NewInt(0)
	cumulativeGas := big.NewInt(0)

	for i, t := range block.Transactions {
		txnIndex := uint64(i)
		gasUsed.SetUint64(t.GasUsed)
		cumulativeGas.Add(cumulativeGas, gasUsed)
		err = writer.InsertTransaction(chainId, height, txnIndex, t.Hash, utils.IndexerTxnToDbTxn(t, primitives.QuantityFromBigInt(cumulativeGas)))
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

func (h *BlockHandler) SetIndexerState(chainId uint64, data []byte) error {
	return h.db.InsertIndexerState(chainId, data)
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

	if filter.To.BlockHeight == 0 && filter.To.TransactionIndex == dbkey.MaxTxIndex && filter.To.LogIndex == dbkey.MaxLogIndex {
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
	topics := make([][]primitives.Data32, len(filter.Topics.Content))
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
		limit := int(100_000) // Max limit
		// If the block range is higher than ScanRangeThreshold, then limit the maximum logs in the response to MaxScanIterators
		if to.BlockHeight-from.BlockHeight > uint64(h.Config.Core.ScanRangeThreshold) {
			limit = int(h.Config.Core.MaxScanIterators)
		}
		resp, lastKey, err = txn.ReadLogs(ctx, chainId, from, to, addresses, topics, limit)
		if err != nil {
			if err == core.ErrLimited {
				var err error
				if limit == int(h.Config.Core.MaxScanIterators) {
					err = &errs.LogResponseRangeLimitError{
						Err: fmt.Errorf("Log response size exceeded. You can make eth_getLogs requests with "+
							"up to a %d block range, or you can request any block range with a cap of %d logs in the response.",
							int(h.Config.Core.ScanRangeThreshold), int(h.Config.Core.MaxScanIterators)),
					}
				} else {
					err = &errs.LogResponseRangeLimitError{
						Err: fmt.Errorf("Log response size exceeded. Your requests can not exceed the maximum "+
							"capacity of %d logs in the response.", limit),
					}
				}
				return resp, lastKey, err
			}
			return resp, lastKey, err
		} else if limit < 0 || len(resp) <= limit {
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

// postProcessPrehistoryBlock updates the block hash and parent hash fields of the prehistory block according to the prehistory chainId
func postProcessPrehistoryBlock(preBlock *response.Block, blockHeight, chainId uint64) (*response.Block, error) {
	var err error
	preBlock.Hash.Content, err = encodeBlockHeight(preBlock.Hash.Content, blockHeight)
	if err != nil {
		return nil, err
	}
	if blockHeight != 0 {
		preBlock.ParentHash.Content, err = encodeBlockHeight(preBlock.ParentHash.Content, blockHeight-1)
		if err != nil {
			return nil, err
		}
	}
	return preBlock, nil
}

// encodeBlockHeight embeds the height to the provided byte slice
func encodeBlockHeight(hash []byte, height uint64) ([]byte, error) {
	if len(hash) < 32 {
		return nil, fmt.Errorf("hash retrieved from block [%d] is invalid", height)
	}

	bufHeightBE := make([]byte, 8)
	bufHeightLE := make([]byte, 8)
	binary.BigEndian.PutUint64(bufHeightBE, height)
	binary.LittleEndian.PutUint64(bufHeightLE, height)
	hashBaseSection1 := hash[8:16]
	hashSection1ToEncode := hash[0:8]
	hashBaseSection2 := hash[24:32]
	hashSection2ToEncode := hash[16:24]
	for i := 0; i < 8; i++ {
		hashSection1ToEncode[i] = bufHeightBE[i] ^ hashBaseSection1[i]
		hashSection2ToEncode[i] = bufHeightLE[i] ^ hashBaseSection2[i]
	}

	return hash, nil
}

// decodeBlockHeight decodes the embedded height from the provided
func decodeBlockHeight(hash []byte) (*uint64, error) {
	if len(hash) < 32 {
		return nil, errors.New("received hash is invalid")
	}

	bufHeightBE := make([]byte, 8)
	bufHeightLE := make([]byte, 8)
	hashBaseSection1 := hash[8:16]
	hashSection1ToDecode := hash[0:8]
	hashBaseSection2 := hash[24:32]
	hashSection2ToDecode := hash[16:24]
	for i := 0; i < 8; i++ {
		bufHeightBE[i] = hashSection1ToDecode[i] ^ hashBaseSection1[i]
		bufHeightLE[i] = hashSection2ToDecode[i] ^ hashBaseSection2[i]
	}
	heightBE := binary.BigEndian.Uint64(bufHeightBE)
	heightLE := binary.LittleEndian.Uint64(bufHeightLE)

	if heightLE != heightBE {
		return nil, nil
	}
	return &heightBE, nil
}
