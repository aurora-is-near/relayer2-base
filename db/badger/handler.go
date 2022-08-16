package badger

import (
	"aurora-relayer-go-common/db"
	"aurora-relayer-go-common/log"
	"aurora-relayer-go-common/utils"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/spf13/viper"
)

const (
	defaultPath = "/tmp/badger/data"
	configPath  = "DB.Badger"
)

type Handler struct {
	db *badger.DB
}

func New() (db.Handler, error) {
	opts := badger.DefaultOptions(defaultPath)
	logger := log.New()
	opts.Logger = NewBadgerLogger(logger)
	sub := viper.Sub(configPath)
	if sub != nil {
		if err := sub.Unmarshal(&opts); err != nil {
			logger.Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
				"falling back to defaults", configPath, viper.ConfigFileUsed())
		}
	}
	return NewWithOpts(opts)
}

func NewWithOpts(opt badger.Options) (db.Handler, error) {
	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}
	return &Handler{db: db}, err
}

func (b Handler) Close() error {
	return b.db.Close()
}

func (b Handler) BlockNumber() (*uint64, error) {
	return FetchFromDB[uint64](b.db, currentBlockIdKey())
}

func (b Handler) GetBlockByHash(hash utils.H256) (*utils.Block, error) {
	block, err := FetchFromDB[utils.Block](b.db, blockByHashKey(hash))
	if err != nil {
		return nil, err
	}
	block.Transactions, err = FetchPrefixedFromDB[utils.Transaction](b.db, txByBlockHashPrefixKey(hash))
	if err != nil {
		return nil, err
	}
	return block, err
}

func (b Handler) GetBlockByNumber(number utils.Uint256) (*utils.Block, error) {
	hash, err := FetchFromDB[utils.H256](b.db, blockByNumberKey(number))
	if err != nil {
		return nil, err
	}
	return b.GetBlockByHash(*hash)
}

func (b Handler) GetBlockTransactionCountByHash(hash utils.H256) (*uint64, error) {
	return FetchFromDB[uint64](b.db, txCountByHashKey(hash))
}

func (b Handler) GetBlockTransactionCountByNumber(number utils.Uint256) (*uint64, error) {
	return FetchFromDB[uint64](b.db, txCountByNumberKey(number))
}

func (b Handler) GetTransactionByHash(hash utils.H256) (*utils.Transaction, error) {
	blockHash, err := FetchFromDB[utils.H256](b.db, txHashToBlockHashIdxKey(hash))
	if err != nil {
		return nil, err
	}
	block, err := b.GetBlockByHash(*blockHash)
	if err != nil {
		return nil, err
	}
	for _, tx := range block.Transactions {
		if tx.Hash == hash {
			return &tx, nil
		}
	}
	return nil, &utils.TxNotFoundError{Hash: hash}
}

func (b Handler) GetTransactionByBlockHashAndIndex(bh utils.H256, idx int64) (*utils.Transaction, error) {
	return FetchFromDB[utils.Transaction](b.db, txByBlockHashKeyAndIdx(bh, idx))
}

func (b Handler) GetTransactionByBlockNumberAndIndex(bn utils.Uint256, idx int64) (*utils.Transaction, error) {
	hash, err := FetchFromDB[utils.H256](b.db, blockByNumberKey(bn))
	if err != nil {
		return nil, err
	}
	return FetchFromDB[utils.Transaction](b.db, txByBlockHashKeyAndIdx(*hash, idx))
}

func (b Handler) GetLogs(addr utils.Address, bn utils.Uint256, topic ...[]string) (*utils.Log, error) {
	// TODO implement me
	panic("implement me")
}

func (b Handler) InsertBlock(block utils.Block) error {
	// TODO do these in one transaction
	savedBlock := block
	savedBlock.Transactions = nil
	if err := InsertToDB(b.db, blockByHashKey(block.Hash), savedBlock); err != nil {
		return err
	}
	if err := InsertToDB(b.db, currentBlockIdKey(), block.Sequence); err != nil {
		return err
	}
	if err := InsertToDB(b.db, blockByNumberKey(fmt.Sprint(block.Sequence)), block.Hash); err != nil {
		return err
	}
	if err := InsertToDB(b.db, txCountByHashKey(block.Hash), block.TxCount()); err != nil {
		return err
	}
	if err := InsertToDB(b.db, txCountByNumberKey(fmt.Sprint(block.Height)), block.TxCount()); err != nil {
		return err
	}

	for idx, tx := range block.Transactions {
		if err := b.InsertTransaction(&block, idx, &tx); err != nil {
			return err
		}
	}
	return nil
}

func (b Handler) InsertTransaction(block *utils.Block, idx int, tx *utils.Transaction) error {
	if err := InsertToDB(b.db, txByBlockHashKeyAndIdx(block.Hash, idx), tx); err != nil {
		return err
	}
	return b.InsertLog(tx)
}

func (b Handler) InsertLog(tx *utils.Transaction) error {
	for logIdx, log := range tx.Logs {
		_ = logIdx
		_ = log
		// TODO log insertion
		panic("implement me")
	}
	return nil
}

func currentBlockIdKey() []byte {
	return []byte("/current-block-id")
}

func blockByHashKey[T ~string | ~[]byte](hash T) []byte {
	return []byte(fmt.Sprintf("/block/hash/%s", hash))
}

func blockByNumberKey[T ~string | ~[]byte](num T) []byte {
	return []byte(fmt.Sprintf("/block/number/%s", num))
}

func txByBlockHashKeyAndIdx[T ~string, I utils.Integer](hash T, idx I) []byte {
	return []byte(fmt.Sprintf("/tx-by-block-hash/%s/%08x", hash, idx))
}

func txByBlockHashPrefixKey[T ~string](hash T) []byte {
	return []byte(fmt.Sprintf("/tx-by-block-hash/%s/", hash))
}

func txCountByHashKey[T ~string | ~[]byte](hash T) []byte {
	return []byte(fmt.Sprintf("/transaction-count/hash/%s", hash))
}

func txCountByNumberKey[T ~string | ~[]byte](hash T) []byte {
	return []byte(fmt.Sprintf("/transaction-count/number/%s", hash))
}

func txHashToBlockHashIdxKey[T ~string](hash T) []byte {
	return []byte(fmt.Sprintf("/block/tx-hash/%s", hash))
}
