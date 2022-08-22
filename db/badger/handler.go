package badger

import (
	"aurora-relayer-go-common/db"
	"aurora-relayer-go-common/log"
	"aurora-relayer-go-common/utils"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/fxamacker/cbor/v2"
	"github.com/spf13/viper"
	"os"
	"path"
	"time"
)

const (
	defaultDataPath = "/tmp/badger/data"
	configPath      = "DB.Badger"

	blockByNumberPrefix string = "/block/number/"
	logTablePrefix      string = "/logs/"
	logIndexTablePrefix string = "/logs-index/"
)

type Encoder interface {
	Marshal(v interface{}) ([]byte, error)
}

type Decoder interface {
	Unmarshal(data []byte, v interface{}) error
}

type Handler struct {
	db      *badger.DB
	encoder Encoder
	decoder Decoder
}

func New() (db.Handler, error) {
	opts := badger.DefaultOptions(defaultDataPath)
	logger := log.New()
	opts.Logger = NewBadgerLogger(logger)
	sub := viper.Sub(configPath)
	if sub != nil {
		if err := sub.Unmarshal(&opts); err != nil {
			logger.Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
				"falling back to defaults", configPath, viper.ConfigFileUsed())
		}
	}

	bdb, err := NewWithOpts(opts)

	if err != nil {
		snapshotBaseName := path.Base(opts.Dir) + "_" + time.Now().Format("2006-01-02T15-04-05.000000000")
		snapshotPath := path.Join(path.Dir(opts.Dir), snapshotBaseName)
		logger.Warn().Err(err).Msgf("saving old database snapshot at [%s]", snapshotPath)
		if err := os.Rename(opts.Dir, snapshotPath); err != nil {
			logger.Error().Err(err).Msg("failed to save old snapshot")
			return nil, err
		}
		bdb, err = NewWithOpts(opts)
	}

	return bdb, err
}

func NewWithOpts(opt badger.Options) (db.Handler, error) {
	opt.Logger.Infof("opening database with path [%s]", opt.Dir)
	bdb, err := badger.Open(opt)
	if err != nil {
		opt.Logger.Errorf("failed to open database with path [%s]", opt.Dir)
		return nil, err
	}
	enc, dec := newEncoderAndDecoder()
	return &Handler{db: bdb,
		encoder: enc,
		decoder: dec,
	}, err
}

func (b Handler) Close() error {
	return b.db.Close()
}

func (b Handler) BlockNumber() (*utils.Uint256, error) {
	return fetchFromDB[utils.Uint256](b.db, currentBlockIdKey())
}

func (b Handler) GetBlockByHash(hash utils.H256) (*utils.Block, error) {
	block, err := fetchFromDB[utils.Block](b.db, blockByHashKey(hash))
	if err != nil {
		return nil, err
	}
	block.Transactions, err = fetchPrefixedFromDB[utils.Transaction](b.db, txByBlockHashPrefixKey(hash))
	if err != nil {
		return nil, err
	}
	return block, err
}

func (b Handler) GetBlockByNumber(number utils.Uint256) (*utils.Block, error) {
	hash, err := fetchFromDB[utils.H256](b.db, blockByNumberKey(number))
	if err != nil {
		return nil, err
	}
	return b.GetBlockByHash(*hash)
}

func (b Handler) GetBlockTransactionCountByHash(hash utils.H256) (*int64, error) {
	return fetchFromDB[int64](b.db, txCountByHashKey(hash))
}

func (b Handler) GetBlockTransactionCountByNumber(number utils.Uint256) (*int64, error) {
	return fetchFromDB[int64](b.db, txCountByNumberKey(number))
}

func (b Handler) GetTransactionByHash(hash utils.H256) (*utils.Transaction, error) {
	key, err := fetchFromDB[[]byte](b.db, txByHashKey(hash))
	if err != nil {
		return nil, err
	}
	return fetchFromDB[utils.Transaction](b.db, *key)
}

func (b Handler) GetTransactionByBlockHashAndIndex(bh utils.H256, idx int64) (*utils.Transaction, error) {
	return fetchFromDB[utils.Transaction](b.db, txByBlockHashAndIdxKey(bh, idx))
}

func (b Handler) GetTransactionByBlockNumberAndIndex(bn utils.Uint256, idx int64) (*utils.Transaction, error) {
	key, err := fetchFromDB[[]byte](b.db, txByBlockNumAndIdxKey(bn, idx))
	if err != nil {
		return nil, err
	}
	return fetchFromDB[utils.Transaction](b.db, *key)
}

func (b Handler) GetLogs(filter *utils.LogFilter) (*[]utils.LogResponse, error) {
	index := NewIndex(4, false)
	logResponses := []utils.LogResponse{}
	err := b.db.View(func(txn *badger.Txn) error {
		from, to := getLogRangeKeys(filter)
		scan := index.StartScan(
			&ScanOpts{
				MaxJumps:         1000,
				MaxRangeScanners: 4,
				MaxValueFetchers: 4,
				KeysOnly:         false,
			},
			txn,
			[]byte(logTablePrefix),
			[]byte(logIndexTablePrefix),
			filter.Topics,
			from,
			to,
		)
		for {
			itemEnc, hasItem := <-scan.Output()
			if !hasItem {
				break
			} else if len(logResponses) == 10000 {
				break
			}
			var logResponse utils.LogResponse
			if err := b.decoder.Unmarshal(itemEnc.Value, &logResponse); err != nil {
				return err
			}
			if filter.Address != nil && !filter.Address[logResponse.Address] {
				continue
			}
			topicStrings := make([]string, 0, len(logResponse.Topics))
			for _, t := range logResponse.Topics {
				topicStrings = append(topicStrings, string(t))
			}
			logResponses = append(logResponses, logResponse)
		}
		return scan.Stop()

	})
	return &logResponses, err
}

func (b Handler) GetBlockHashesSinceNumber(number utils.Uint256) ([]utils.H256, error) {
	results := make([]utils.H256, 0)
	err := b.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		it.Seek(blockByNumberKey(number))

		var res utils.H256
		for ; it.ValidForPrefix([]byte(blockByNumberPrefix)); it.Next() {
			item, err := it.Item().ValueCopy(nil)
			if err != nil {
				return err
			}
			if err := b.decoder.Unmarshal(item, &res); err != nil {
				return err
			}
			results = append(results, res)
		}
		return nil
	})
	return results, err
}

func (b Handler) InsertBlock(block utils.Block) error {
	// TODO do these in one transaction

	savedBlock := block
	savedBlock.Transactions = nil
	if err := insertToDB(b.db, blockByHashKey(block.Hash), savedBlock); err != nil {
		return err
	}
	if err := insertToDB(b.db, currentBlockIdKey(), block.Sequence); err != nil {
		return err
	}
	if err := insertToDB(b.db, blockByNumberKey(block.Sequence), block.Hash); err != nil {
		return err
	}
	if err := insertToDB(b.db, txCountByHashKey(block.Hash), block.TxCount()); err != nil {
		return err
	}
	if err := insertToDB(b.db, txCountByNumberKey(block.Sequence), block.TxCount()); err != nil {
		return err
	}

	for idx, tx := range block.Transactions {
		if err := b.InsertTransaction(tx, idx, &block); err != nil {
			return err
		}
	}
	return nil
}

func (b Handler) InsertTransaction(tx utils.Transaction, idx int, block *utils.Block) error {
	savedTx := tx
	savedTx.Logs = nil
	mainKey := txByBlockHashAndIdxKey(block.Hash, idx)
	if err := insertToDB(b.db, mainKey, savedTx); err != nil {
		return err
	}
	if err := insertToDB(b.db, txByHashKey(tx.Hash), mainKey); err != nil {
		return err
	}
	if err := insertToDB(b.db, txByBlockNumAndIdxKey(block.Sequence, idx), mainKey); err != nil {
		return err
	}

	for lIdx, l := range tx.Logs {
		if err := b.InsertLog(l, lIdx, &tx, idx, block); err != nil {
			return err
		}
	}
	return nil
}

func (b Handler) InsertLog(log utils.Log, idx int, tx *utils.Transaction, txIdx int, block *utils.Block) error {
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
	dataEnc, err := b.encoder.Marshal(data)
	if err != nil {
		return err
	}
	index := NewIndex(4, false)
	batch := b.db.NewWriteBatch()

	key := getLogInsertKey(data.BlockNumber, uint64(txIdx), uint64(idx))
	// Populate original table
	err = batch.Set(
		[]byte(logTablePrefix+string(key)),
		dataEnc,
	)
	if err != nil {
		return err
	}

	topicsB := make([][]byte, 0, len(data.Topics))
	for _, t := range data.Topics {
		topicsB = append(topicsB, t)
	}

	// Populate index table
	err = index.Insert(
		[]byte(logIndexTablePrefix),
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

func getLogRangeKeys(filter *utils.LogFilter) (from, to []byte) {
	fNum := filter.FromBlock.Bytes()
	fromBuf := new(bytes.Buffer)
	if err := binary.Write(fromBuf, binary.BigEndian, fNum); err != nil {
		panic(err)
	}

	tNum := filter.ToBlock.Bytes()
	toBuf := new(bytes.Buffer)
	if err := binary.Write(toBuf, binary.BigEndian, tNum); err != nil {
		panic(err)
	}

	return fromBuf.Bytes(), toBuf.Bytes()
}

func getLogInsertKey(blockNum utils.Uint256, txIdx, logIdx uint64) []byte {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, blockNum.Bytes()); err != nil {
		panic(err)
	}
	if err := binary.Write(buf, binary.BigEndian, txIdx); err != nil {
		panic(err)
	}
	if err := binary.Write(buf, binary.BigEndian, logIdx); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func currentBlockIdKey() []byte {
	return []byte("/current-block-id")
}

func blockByHashKey(hash fmt.Stringer) []byte {
	return []byte(fmt.Sprintf("/block/hash/%s", hash))
}

func blockByNumberKey(num utils.Uint256) []byte {
	return []byte(fmt.Sprintf("%s%s", blockByNumberPrefix, num))
}

func txByBlockHashAndIdxKey[I utils.Integer](hash fmt.Stringer, idx I) []byte {
	return []byte(fmt.Sprintf("/tx-by-block-hash/%s/%08x", hash, idx))
}

func txByBlockNumAndIdxKey[I utils.Integer](num utils.Uint256, idx I) []byte {
	return []byte(fmt.Sprintf("/tx-by-block-num/%s/%08x", num, idx))
}

func txByBlockHashPrefixKey(hash fmt.Stringer) []byte {
	return []byte(fmt.Sprintf("/tx-by-block-hash/%s/", hash))
}

func txCountByHashKey(hash fmt.Stringer) []byte {
	return []byte(fmt.Sprintf("/transaction-count/hash/%s", hash))
}

func txCountByNumberKey(num utils.Uint256) []byte {
	return []byte(fmt.Sprintf("/transaction-count/number/%s", num))
}

func txByHashKey(hash fmt.Stringer) []byte {
	return []byte(fmt.Sprintf("/tx-by-hash/%s", hash))
}

func newEncoderAndDecoder() (Encoder, Decoder) {
	enc, err := cbor.EncOptions{
		BigIntConvert: cbor.BigIntConvertShortest,
	}.EncMode()
	if err != nil {
		panic(err)
	}

	dec, err := cbor.DecOptions{}.DecMode()
	if err != nil {
		panic(err)
	}

	return enc, dec
}
