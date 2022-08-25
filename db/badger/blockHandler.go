package badger

import (
	"aurora-relayer-go-common/db"
	"aurora-relayer-go-common/db/badger/core"
	"aurora-relayer-go-common/log"
	"aurora-relayer-go-common/utils"
	"github.com/dgraph-io/badger/v3"
	"time"
)

var logIndex = core.NewIndex(5, false)

type BlockHandler struct {
	db     *badger.DB
	codec  db.Codec
	config *Config
}

func NewBlockHandler() (db.BlockHandler, error) {
	config := GetConfig()
	codec := db.NewCborCodec()
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

func (h *BlockHandler) BlockNumber() (*utils.Uint256, error) {
	return fetch[utils.Uint256](h.codec, prefixCurrentBlockId.Key())
}

func (h *BlockHandler) GetBlockByHash(hash utils.H256) (*utils.Block, error) {
	number, err := h.BlockHashToNumber(hash)
	if err != nil {
		return nil, err
	}
	return h.GetBlockByNumber(*number)
}

func (h *BlockHandler) GetBlockByNumber(number utils.Uint256) (*utils.Block, error) {
	block, err := fetch[utils.Block](h.codec, prefixBlockByNumber.Key(number))
	if err != nil {
		return nil, err
	}
	block.Transactions, err = h.GetTransactionsForBlock(block)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (h *BlockHandler) GetBlockTransactionCountByHash(hash utils.H256) (int64, error) {
	count, err := fetch[int64](h.codec, prefixTransactionCountByBlockHash.Key(hash))
	if count == nil {
		return 0, err
	}
	return *count, err
}

func (h *BlockHandler) GetBlockTransactionCountByNumber(number utils.Uint256) (int64, error) {
	count, err := fetch[int64](h.codec, prefixTransactionCountByBlockNumber.Key(number))
	if count == nil {
		return 0, err
	}
	return *count, err
}

func (h *BlockHandler) GetTransactionByHash(hash utils.H256) (*utils.Transaction, error) {
	key, err := fetch[[]byte](h.codec, prefixTransactionByHash.Key(hash))
	if err != nil {
		return nil, err
	}
	return fetch[utils.Transaction](h.codec, *key)
}

func (h *BlockHandler) GetTransactionByBlockHashAndIndex(hash utils.H256, index utils.Uint256) (*utils.Transaction, error) {
	idx, err := index.ToUint32Key()
	if err != nil {
		return nil, err
	}
	return fetch[utils.Transaction](h.codec, prefixTransactionByBlockHashAndIndex.Key(hash, idx))
}

func (h *BlockHandler) GetTransactionByBlockNumberAndIndex(number utils.Uint256, index utils.Uint256) (*utils.Transaction, error) {
	idx, err := index.ToUint32Key()
	if err != nil {
		return nil, err
	}
	key, err := fetch[[]byte](h.codec, prefixTransactionByBlockNumberAndIndex.Key(number, idx))
	if err != nil {
		return nil, err
	}
	return fetch[utils.Transaction](h.codec, *key)
}

func (h *BlockHandler) GetLogs(filter utils.LogFilter) (*[]utils.LogResponse, error) {
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

func (h *BlockHandler) GetLogsForTransaction(tx *utils.Transaction) ([]*utils.LogResponse, error) {
	blockNum := utils.UintToUint256(tx.BlockHeight)
	txIdx, err := utils.UintToUint256(tx.TransactionIndex).ToUint32Key()
	if err != nil {
		return nil, err
	}
	return fetchPrefixedWithLimitAndTimeout[utils.LogResponse](h.codec, h.config.IterationMaxItems,
		h.config.IterationTimeoutSeconds, prefixLogTable.Key(blockNum, txIdx))

}

func (h *BlockHandler) GetTransactionsForBlock(block *utils.Block) ([]*utils.Transaction, error) {
	return fetchPrefixedWithLimitAndTimeout[utils.Transaction](h.codec, h.config.IterationMaxItems,
		h.config.IterationTimeoutSeconds, prefixTransactionByBlockHashAndIndex.Key(block.Hash))
}

func (h *BlockHandler) GetBlockHashesSinceNumber(number utils.Uint256) ([]utils.H256, error) {
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
	// TODO do these in one transaction

	savedBlock := *block
	savedBlock.Transactions = nil
	if err := insert(h.codec, prefixBlockByNumber.Key(block.Sequence), savedBlock); err != nil {
		return err
	}
	if err := insert(h.codec, prefixCurrentBlockId.Key(), block.Sequence); err != nil {
		return err
	}
	if err := insert(h.codec, prefixBlockByHash.Key(block.Hash), block.Sequence); err != nil {
		return err
	}
	if err := insert(h.codec, prefixTransactionCountByBlockHash.Key(block.Hash), block.TxCount()); err != nil {
		return err
	}
	if err := insert(h.codec, prefixTransactionCountByBlockNumber.Key(block.Sequence), block.TxCount()); err != nil {
		return err
	}

	for idx, tx := range block.Transactions {
		if err := h.InsertTransaction(tx, idx, block); err != nil {
			return err
		}
	}
	return nil
}

func (h *BlockHandler) InsertTransaction(tx *utils.Transaction, index int, block *utils.Block) error {
	savedTx := *tx
	savedTx.Logs = nil
	idx, err := utils.IntToUint256(index).ToUint32Key()
	if err != nil {
		return err
	}
	mainKey := prefixTransactionByBlockHashAndIndex.Key(block.Hash, idx)
	if err := insert(h.codec, mainKey, savedTx); err != nil {
		return err
	}
	if err := insert(h.codec, prefixTransactionByHash.Key(tx.Hash), mainKey); err != nil {
		return err
	}
	if err := insert(h.codec, prefixTransactionByBlockNumberAndIndex.Key(block.Sequence, idx), mainKey); err != nil {
		return err
	}

	for lIdx, l := range tx.Logs {
		if err := h.InsertLog(l, lIdx, tx, index, block); err != nil {
			return err
		}
	}
	return nil
}

func (h *BlockHandler) InsertLog(log *utils.Log, idx int, tx *utils.Transaction, txIdx int, block *utils.Block) error {
	data := utils.LogResponse{
		Removed:          false,
		LogIndex:         utils.IntToUint256(idx),
		TransactionIndex: utils.IntToUint256(txIdx),
		TransactionHash:  tx.Hash,
		BlockHash:        block.Hash,
		BlockNumber:      block.Sequence,
		Address:          log.Address,
		Data:             log.Data,
		Topics:           log.Topics,
	}
	dataEnc, err := h.codec.Marshal(data)
	if err != nil {
		return err
	}
	batch := h.db.NewWriteBatch()

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
	err = batch.Set(
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
			return batch.Set(key, value)
		},
	)
	if err != nil {
		return err
	}
	if err := batch.Flush(); err != nil {
		return err
	}
	return nil
}

func (h *BlockHandler) BlockHashToNumber(hash utils.H256) (*utils.Uint256, error) {
	return fetch[utils.Uint256](h.codec, prefixBlockByHash.Key(hash))
}

func (h *BlockHandler) BlockNumberToHash(number utils.Uint256) (*utils.H256, error) {
	block, err := h.GetBlockByNumber(number)
	if err != nil {
		return nil, err
	}
	return &block.Hash, nil
}
