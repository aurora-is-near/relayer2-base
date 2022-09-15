package badger

import (
	"aurora-relayer-go-common/db"
	"aurora-relayer-go-common/db/badger/core"
	"aurora-relayer-go-common/db/codec"
	"aurora-relayer-go-common/log"
	"aurora-relayer-go-common/utils"
	"context"
	"github.com/dgraph-io/badger/v3"
	"time"
)

var logIndex = core.NewIndex(5, false)

const maxPendingTxns = 100

type BlockHandler struct {
	db     *badger.DB
	codec  codec.Codec
	config *Config
}

func NewBlockHandler() (db.BlockHandler, error) {
	return NewBlockHandlerWithCodec(codec.NewCborCodec())
}

func NewBlockHandlerWithCodec(codec codec.Codec) (db.BlockHandler, error) {
	config := GetConfig()
	bdb, err := core.Open(config.BadgerConfig, config.GcIntervalSeconds)
	if err != nil {
		return nil, err
	}
	return &BlockHandler{
		db:     bdb,
		codec:  codec,
		config: config,
	}, nil
}

func (h *BlockHandler) Close() error {
	return core.Close()
}

func (h *BlockHandler) BlockNumber(ctx context.Context) (*utils.Uint256, error) {
	return fetch[utils.Uint256](ctx, h.codec, prefixCurrentBlockId.Key())
}

func (h *BlockHandler) GetBlockByHash(ctx context.Context, hash utils.H256) (*utils.Block, error) {
	number, err := h.BlockHashToNumber(ctx, hash)
	if err != nil {
		return nil, err
	}
	return h.GetBlockByNumber(ctx, *number)
}

func (h *BlockHandler) GetBlockByNumber(ctx context.Context, number utils.Uint256) (*utils.Block, error) {
	block, err := fetch[utils.Block](ctx, h.codec, prefixBlockByNumber.Key(number))
	if err != nil {
		return nil, err
	}
	block.Transactions, err = h.GetTransactionsForBlock(ctx, block)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (h *BlockHandler) GetBlockTransactionCountByHash(ctx context.Context, hash utils.H256) (int64, error) {
	count, err := fetch[int64](ctx, h.codec, prefixTransactionCountByBlockHash.Key(hash))
	if count == nil {
		return 0, err
	}
	return *count, err
}

func (h *BlockHandler) GetBlockTransactionCountByNumber(ctx context.Context, number utils.Uint256) (int64, error) {
	count, err := fetch[int64](ctx, h.codec, prefixTransactionCountByBlockNumber.Key(number))
	if count == nil {
		return 0, err
	}
	return *count, err
}

func (h *BlockHandler) GetTransactionByHash(ctx context.Context, hash utils.H256) (*utils.Transaction, error) {
	key, err := fetch[[]byte](ctx, h.codec, prefixTransactionByHash.Key(hash))
	if err != nil {
		return nil, err
	}
	return fetch[utils.Transaction](ctx, h.codec, *key)
}

func (h *BlockHandler) GetTransactionByBlockHashAndIndex(ctx context.Context, hash utils.H256, index utils.Uint256) (*utils.Transaction, error) {
	idx, err := index.ToUint32Key()
	if err != nil {
		return nil, err
	}
	return fetch[utils.Transaction](ctx, h.codec, prefixTransactionByBlockHashAndIndex.Key(hash, idx))
}

func (h *BlockHandler) GetTransactionByBlockNumberAndIndex(ctx context.Context, number utils.Uint256, index utils.Uint256) (*utils.Transaction, error) {
	idx, err := index.ToUint32Key()
	if err != nil {
		return nil, err
	}
	key, err := fetch[[]byte](ctx, h.codec, prefixTransactionByBlockNumberAndIndex.Key(number, idx))
	if err != nil {
		return nil, err
	}
	return fetch[utils.Transaction](ctx, h.codec, *key)
}

func (h *BlockHandler) GetLogs(ctx context.Context, filter utils.LogFilter) (*[]utils.LogResponse, error) {
	logResponses := []utils.LogResponse{}
	timeout := time.NewTimer(time.Second * time.Duration(h.config.IterationTimeoutSeconds))
	defer timeout.Stop()
	err := h.db.View(func(txn *badger.Txn) error {
		from, to := filter.FromBlock.KeyBytes(), filter.ToBlock.Add(1).KeyBytes()
		fieldFilters := [][][]byte{filter.Address}
		fieldFilters = append(fieldFilters, filter.Topics...)
		scan := logIndex.StartScan(
			&h.config.ScanConfig,
			txn,
			prefixLogTable.Bytes(),
			prefixLogIndexTable.Bytes(),
			fieldFilters,
			from,
			to,
		)
	loop:
		for {
			select {
			case <-timeout.C:
				break loop

			case itemEnc, hasItem := <-scan.Output():
				if !hasItem {
					break loop
				}
				var item utils.LogResponse
				if err := h.codec.Unmarshal(itemEnc.Value, &item); err != nil {
					return err
				}
				logResponses = append(logResponses, item)
				if uint(len(logResponses)) >= h.config.IterationMaxItems {
					break loop
				}
			}
		}
		if err := scan.Stop(); err != nil {
			log.Log().Err(err).Msg("error during log scan")
		}
		return nil
	})
	return &logResponses, err
}

func (h *BlockHandler) GetLogsForTransaction(ctx context.Context, tx *utils.Transaction) ([]*utils.LogResponse, error) {
	blockNum := utils.UintToUint256(tx.BlockHeight)
	txIdx, err := utils.UintToUint256(tx.TransactionIndex).ToUint32Key()
	if err != nil {
		return nil, err
	}
	return fetchPrefixedWithLimitAndTimeout[utils.LogResponse](ctx, h.codec, h.config.IterationMaxItems,
		h.config.IterationTimeoutSeconds, prefixLogTable.Key(blockNum, txIdx))

}

func (h *BlockHandler) GetTransactionsForBlock(ctx context.Context, block *utils.Block) ([]*utils.Transaction, error) {
	return fetchPrefixedWithLimitAndTimeout[utils.Transaction](ctx, h.codec, h.config.IterationMaxItems,
		h.config.IterationTimeoutSeconds, prefixTransactionByBlockHashAndIndex.Key(block.Hash))
}

func (h *BlockHandler) GetBlockHashesSinceNumber(ctx context.Context, number utils.Uint256) ([]utils.H256, error) {
	results := make([]utils.H256, 0)
	err := h.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		var res utils.Block
		for it.Seek(prefixBlockByNumber.Key(number.Add(1))); it.ValidForPrefix(prefixBlockByNumber.Bytes()); it.Next() {
			item, err := it.Item().ValueCopy(nil)
			if err != nil {
				return err
			}
			if err := h.codec.Unmarshal(item, &res); err != nil {
				return err
			}
			results = append(results, res.Hash)
		}
		return nil
	})
	return results, err
}

func (h *BlockHandler) InsertBlock(block *utils.Block) error {

	savedBlock := *block
	savedBlock.Transactions = nil
	blockNum := utils.UintToUint256(block.Height)

	writer := h.db.NewWriteBatch()
	writer.SetMaxPendingTxns(maxPendingTxns)
	defer writer.Cancel()

	if err := insertBatch(writer, h.codec, prefixCurrentBlockSequence.Key(), block.Sequence); err != nil {
		return err
	}
	if err := insertBatch(writer, h.codec, prefixCurrentBlockId.Key(), blockNum); err != nil {
		return err
	}
	if err := insertBatch(writer, h.codec, prefixBlockByNumber.Key(blockNum), savedBlock); err != nil {
		return err
	}
	if err := insertBatch(writer, h.codec, prefixBlockByHash.Key(block.Hash), blockNum); err != nil {
		return err
	}
	if err := insertBatch(writer, h.codec, prefixTransactionCountByBlockHash.Key(block.Hash), block.TxCount()); err != nil {
		return err
	}
	if err := insertBatch(writer, h.codec, prefixTransactionCountByBlockNumber.Key(blockNum), block.TxCount()); err != nil {
		return err
	}

	for idx, tx := range block.Transactions {
		if err := h.insertTransaction(writer, tx, idx, block); err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

func (h *BlockHandler) BlockHashToNumber(ctx context.Context, hash utils.H256) (*utils.Uint256, error) {
	return fetch[utils.Uint256](ctx, h.codec, prefixBlockByHash.Key(hash))
}

func (h *BlockHandler) BlockNumberToHash(ctx context.Context, number utils.Uint256) (*utils.H256, error) {
	block, err := h.GetBlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	return &block.Hash, nil
}

func (h *BlockHandler) CurrentBlockSequence(ctx context.Context) uint64 {
	i, err := fetch[uint64](ctx, h.codec, prefixCurrentBlockSequence.Key())
	if err != nil {
		return 0
	}
	return *i
}

func (h *BlockHandler) insertTransaction(writer *badger.WriteBatch, tx *utils.Transaction, index int, block *utils.Block) error {
	savedTx := *tx
	savedTx.Logs = nil
	idx, err := utils.IntToUint256(index).ToUint32Key()
	if err != nil {
		return err
	}
	mainKey := prefixTransactionByBlockHashAndIndex.Key(block.Hash, idx)
	if err := insertBatch(writer, h.codec, mainKey, savedTx); err != nil {
		return err
	}
	if err := insertBatch(writer, h.codec, prefixTransactionByHash.Key(tx.Hash), mainKey); err != nil {
		return err
	}
	if err := insertBatch(writer, h.codec, prefixTransactionByBlockNumberAndIndex.Key(utils.UintToUint256(block.Height), idx), mainKey); err != nil {
		return err
	}

	for lIdx, l := range tx.Logs {
		if err := h.insertLog(writer, l, lIdx, tx, index, block); err != nil {
			return err
		}
	}
	return nil
}

func (h *BlockHandler) insertLog(writer *badger.WriteBatch, log *utils.Log, idx int, tx *utils.Transaction, txIdx int, block *utils.Block) error {
	data := utils.LogResponse{
		Removed:          false,
		LogIndex:         utils.IntToUint256(idx),
		TransactionIndex: utils.IntToUint256(txIdx),
		TransactionHash:  tx.Hash,
		BlockHash:        block.Hash,
		BlockNumber:      utils.UintToUint256(block.Height),
		Address:          log.Address,
		Data:             log.Data,
		Topics:           log.Topics,
	}
	dataEnc, err := h.codec.Marshal(data)
	if err != nil {
		return err
	}

	logSmallKey, err := utils.IntToUint256(idx).ToUint32Key()
	if err != nil {
		return err
	}
	txSmallKey, err := utils.IntToUint256(txIdx).ToUint32Key()
	if err != nil {
		return err
	}
	key := getLogInsertKey(data.BlockNumber, txSmallKey, logSmallKey)
	// Populate original table
	err = writer.Set(
		prefixLogTable.AppendBytes(key),
		dataEnc,
	)
	if err != nil {
		return err
	}

	topicsB := make([][]byte, 0, len(data.Topics))
	topicsB = append(topicsB, data.Address.Bytes())
	for _, t := range data.Topics {
		topicsB = append(topicsB, t)
	}

	// Populate index table
	err = logIndex.Insert(
		prefixLogIndexTable.Bytes(),
		topicsB,
		key,
		func(key, value []byte) error {
			return writer.Set(key, value)
		},
	)
	if err != nil {
		return err
	}
	return nil
}
