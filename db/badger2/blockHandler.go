package badger

import (
	"aurora-relayer-go-common/db"
	"aurora-relayer-go-common/db/badger/core"
	cc "aurora-relayer-go-common/db/badger2/core"
	"aurora-relayer-go-common/db/badger2/core/dbcore"
	dbp "aurora-relayer-go-common/db/badger2/core/dbprimitives"
	"aurora-relayer-go-common/db/badger2/core/dbresponses"
	"aurora-relayer-go-common/db/badger2/core/dbtypes"
	"aurora-relayer-go-common/db/codec"
	"aurora-relayer-go-common/log"
	"aurora-relayer-go-common/tinypack"
	"aurora-relayer-go-common/utils"
	"context"
	"github.com/pkg/errors"
)

var (
	logIndex         = core.NewIndex(5, false)
	keyNotFoundError = errors.New("key not found")
	txnParseError    = errors.New("failed to parse transaction")
)

const (
	chainID        = 1313161555
	maxPendingTxns = 100
)

type BlockHandler struct {
	db     *cc.DB
	config *Config
}

// type BlockHandler struct {
// 	db     *badger.DB
// 	codec  codec.Codec
// 	config *Config
// }

func NewBlockHandler() (db.BlockHandler, error) {
	return NewBlockHandlerWithCodec(codec.NewCborCodec())
}

func NewBlockHandlerWithCodec(codec codec.Codec) (db.BlockHandler, error) {
	config := GetConfig()

	asd := &cc.DB{
		CoreOpts: &dbcore.DBCoreOpts{
			Dir:               config.BadgerConfig.Dir,
			GCIntervalSeconds: uint(config.GcIntervalSeconds),
			InMemory:          false,
		},
		Encoder:     tinypack.DefaultEncoder(),
		Decoder:     tinypack.DefaultDecoder(),
		LogScanOpts: &config.ScanConfig,
	}

	asd.Open(config.BadgerConfig.Logger)

	return &BlockHandler{
		db:     asd,
		config: config,
	}, nil
}

func (h *BlockHandler) Close() error {
	return core.Close()
}

func (h *BlockHandler) BlockNumber(ctx context.Context) (*dbp.HexUint, error) {
	var bn dbp.HexUint
	err := h.db.View(func(txn *cc.ViewTxn) error {
		key, err := txn.ReadLatestBlockKey(chainID)
		if err != nil {
			return err
		}
		if key == nil {
			return keyNotFoundError
		}
		bn = dbp.HexUint(key.Height)
		return nil
	})
	return &bn, err
}

func (h *BlockHandler) GetBlockByHash(ctx context.Context, hash utils.H256) (*dbresponses.Block, error) {
	var resp *dbresponses.Block
	var err error
	bh := dbp.DataFromHex[dbp.Len32](hash.String())
	err = h.db.View(func(txn *cc.ViewTxn) error {
		key, err := txn.ReadBlockKey(chainID, bh)
		if err != nil {
			return err
		}
		if key == nil {
			return keyNotFoundError
		}
		resp, err = txn.ReadBlock(chainID, *key, true) // TODO fix fullTransaction flag
		return err
	})
	return resp, err
}

func (h *BlockHandler) GetBlockByNumber(ctx context.Context, number utils.Uint256) (*dbresponses.Block, error) {
	var resp *dbresponses.Block
	var err error
	err = h.db.View(func(txn *cc.ViewTxn) error {
		key := dbtypes.BlockKey{Height: number.Uint64()}
		resp, err = txn.ReadBlock(chainID, key, true) // TODO fix fullTransaction flag
		return err
	})
	return resp, err
}

func (h *BlockHandler) GetBlockTransactionCountByHash(ctx context.Context, hash utils.H256) (*dbp.HexUint, error) {
	var resp dbp.HexUint
	var err error
	bh := dbp.DataFromHex[dbp.Len32](hash.String())
	err = h.db.View(func(txn *cc.ViewTxn) error {
		key, err := txn.ReadBlockKey(chainID, bh)
		if err != nil {
			return err
		}
		if key == nil {
			return keyNotFoundError
		}
		resp, err = txn.ReadBlockTxCount(chainID, *key)
		return err
	})
	return &resp, err
}

func (h *BlockHandler) GetBlockTransactionCountByNumber(ctx context.Context, number utils.Uint256) (*dbp.HexUint, error) {
	var resp dbp.HexUint
	var err error
	err = h.db.View(func(txn *cc.ViewTxn) error {
		key := dbtypes.BlockKey{Height: number.Uint64()}
		resp, err = txn.ReadBlockTxCount(chainID, key)
		return err
	})
	return &resp, err
}

func (h *BlockHandler) GetTransactionByHash(ctx context.Context, hash utils.H256) (*dbresponses.Transaction, error) {
	var resp *dbresponses.Transaction
	var err error
	th := dbp.DataFromHex[dbp.Len32](hash.String())
	err = h.db.View(func(txn *cc.ViewTxn) error {
		key, err := txn.ReadTxKey(chainID, th)
		if err != nil {
			return err
		}
		if key == nil {
			return keyNotFoundError
		}
		resp, err = txn.ReadTx(chainID, *key)
		return err
	})
	return resp, err
}

func (h *BlockHandler) GetTransactionByBlockHashAndIndex(ctx context.Context, hash utils.H256, index utils.Uint256) (*dbresponses.Transaction, error) {
	var resp dbresponses.Transaction
	var err error
	bh := dbp.DataFromHex[dbp.Len32](hash.String())
	err = h.db.View(func(txn *cc.ViewTxn) error {
		key, err := txn.ReadBlockKey(chainID, bh)
		if err != nil {
			return err
		}
		if key == nil {
			return keyNotFoundError
		}
		block, err := txn.ReadBlock(chainID, *key, true)
		if err != nil {
			return err
		}
		var ok bool
		resp, ok = block.Transactions[index.Uint64()].(dbresponses.Transaction)
		if !ok {
			return txnParseError
		}
		return nil
	})
	return &resp, err
}

func (h *BlockHandler) GetTransactionByBlockNumberAndIndex(ctx context.Context, number utils.Uint256, index utils.Uint256) (*dbresponses.Transaction, error) {
	var resp dbresponses.Transaction
	var err error
	err = h.db.View(func(txn *cc.ViewTxn) error {
		key := dbtypes.BlockKey{Height: number.Uint64()}
		block, err := txn.ReadBlock(chainID, key, true)
		if err != nil {
			return err
		}
		var ok bool
		resp, ok = block.Transactions[index.Uint64()].(dbresponses.Transaction)
		if !ok {
			return txnParseError
		}
		return nil
	})
	return &resp, err
}

// TODO implement
func (h *BlockHandler) GetLogs(ctx context.Context, filter utils.LogFilter) ([]*dbresponses.Log, error) {
	return nil, nil
}

// TODO implement if necessary
// func (h *BlockHandler) GetLogsForTransaction(ctx context.Context, tx *utils.Transaction) ([]*dbresponses.Log, error) {
// 	return nil, nil
//
// }
//
// func (h *BlockHandler) GetTransactionsForBlock(ctx context.Context, block *utils.Block) ([]*dbresponses.Transaction, error) {
//
// 	return nil, nil
// }
//
// func (h *BlockHandler) GetBlockHashesSinceNumber(ctx context.Context, number utils.Uint256) ([]*dbp.Data32, error) {
// 	return nil, nil
// }
//
// func (h *BlockHandler) BlockHashToNumber(ctx context.Context, hash utils.H256) (*dbp.HexUint, error) {
// 	return nil, nil
// }
//
// func (h *BlockHandler) BlockNumberToHash(ctx context.Context, number utils.Uint256) (*dbp.Data32, error) {
// 	return nil, nil
// }
//
// func (h *BlockHandler) CurrentBlockSequence(ctx context.Context) uint64 {
// 	return 0
// }

func (h *BlockHandler) InsertBlock(block *utils.Block) error {

	h.db.OpenWriter()
	defer h.db.CloseWriter()

	hash := dbp.DataFromBytes[dbp.Len32](block.Hash.Bytes())
	e := h.db.InsertBlock(chainID, block.Height, hash, toBlockStore(block))
	if e != nil {
		return e
	}

	for i, t := range block.Transactions {
		txnHash := dbp.DataFromBytes[dbp.Len32](t.Hash.Bytes())
		txnIndex := uint64(i)
		e = h.db.InsertTransaction(chainID, block.Height, txnIndex, txnHash, toTxnStore(t))
		if e != nil {
			return e
		}
		for j, l := range t.Logs {
			e = h.db.InsertLog(chainID, block.Height, txnIndex, uint64(j), toLogStore(l))
			if e != nil {
				return e
			}
		}
	}

	h.db.FlushWriter()
	return nil
}

func toBlockStore(block *utils.Block) *dbtypes.Block {
	b := dbtypes.Block{
		ParentHash:       dbp.Data32FromHex(block.ParentHash.String()),
		Miner:            dbp.Data20FromBytes(block.Miner.Bytes()),
		Timestamp:        uint64(block.Timestamp),
		GasLimit:         dbp.QuantityFromBytes(block.GasLimit.Bytes()),
		GasUsed:          dbp.QuantityFromBytes(block.GasUsed.Bytes()),
		LogsBloom:        dbp.Data256FromBytes([]byte(block.LogsBloom)),
		TransactionsRoot: dbp.Data32FromBytes(block.TransactionsRoot.Bytes()),
		StateRoot:        dbp.Data32FromBytes([]byte(block.StateRoot)),
		ReceiptsRoot:     dbp.Data32FromBytes(block.ReceiptsRoot.Bytes()),
		Size:             block.Size.Uint64(),
	}
	return &b
}

func toTxnStore(txn *utils.Transaction) *dbtypes.Transaction {

	toOrContract := dbp.Data20FromHex("0x0") // TODO: temporarily assigning zero addr for corrupted blocks; both contract address and to address is null, what should we do?
	isContractDeployment := false
	if txn.ContractAddress != nil {
		if txn.To != nil {
			log.Log().Warn().Msgf("both contract address and to address is set for txn: [%v], to: [%s], contract: [%s]",
				txn.Hash, txn.To.String(), txn.ContractAddress.String())
		}
		isContractDeployment = true
		toOrContract = dbp.Data20FromBytes(txn.ContractAddress.Bytes())
	} else if txn.To != nil {
		toOrContract = dbp.Data20FromBytes(txn.To.Bytes())
	} else {
		log.Log().Warn().Msgf("both contract address and to address is null for txn: [%v]", txn.Hash)
	}

	var accessListEntries []dbtypes.AccessListEntry
	for _, al := range txn.AccessList {
		var storageKeys []dbp.Data32
		for _, sk := range al.StorageKeys {
			storageKey := dbp.Data32FromBytes(sk.Bytes())
			storageKeys = append(storageKeys, storageKey)
		}
		sk := tinypack.CreateList[dbp.VarLen, dbp.Data32](storageKeys...)
		accessListEntry := dbtypes.AccessListEntry{
			Address:     dbp.Data20FromBytes(al.Address.Bytes()),
			StorageKeys: tinypack.VarList[dbp.Data32]{sk},
		}
		accessListEntries = append(accessListEntries, accessListEntry)
	}
	ake := tinypack.CreateList[dbp.VarLen, dbtypes.AccessListEntry](accessListEntries...)

	t := dbtypes.Transaction{
		Type:                 uint64(txn.TxType),
		From:                 dbp.Data20FromBytes(txn.From.Bytes()),
		IsContractDeployment: isContractDeployment,
		ToOrContract:         toOrContract,
		Nonce:                dbp.QuantityFromBytes(txn.Nonce.Bytes()),
		GasPrice:             dbp.QuantityFromBytes(txn.GasPrice.Bytes()),
		GasLimit:             dbp.QuantityFromBytes(txn.GasLimit.Bytes()),
		GasUsed:              txn.GasUsed.Uint64(),
		Value:                dbp.QuantityFromBytes(txn.Value.Bytes()),
		Input:                dbp.VarDataFromBytes(txn.Input),
		NearHash:             dbp.Data32FromBytes(txn.NearTransaction.Hash.Bytes()),
		NearReceiptHash:      dbp.Data32FromBytes(txn.NearTransaction.ReceiptHash.Bytes()),
		Status:               txn.Status,
		V:                    txn.V,
		R:                    dbp.QuantityFromBytes(txn.R.Bytes()),
		S:                    dbp.QuantityFromBytes(txn.S.Bytes()),
		LogsBloom:            dbp.Data256FromHex(txn.LogsBloom),
		AccessList:           tinypack.VarList[dbtypes.AccessListEntry]{ake},
		MaxPriorityFeePerGas: dbp.QuantityFromBytes(txn.MaxPriorityFeePerGas.Bytes()),
		MaxFeePerGas:         dbp.QuantityFromBytes(txn.MaxFeePerGas.Bytes()),
	}
	return &t
}

func toLogStore(log *utils.Log) *dbtypes.Log {

	var topics []dbp.Data32
	for _, t := range log.Topics {
		topic := dbp.Data32FromBytes(t)
		topics = append(topics, topic)
	}
	t := tinypack.CreateList[dbp.VarLen, dbp.Data32](topics...)

	l := dbtypes.Log{
		Address: dbp.Data20FromBytes(log.Address.Bytes()),
		Data:    dbp.VarDataFromBytes(log.Data),
		Topics:  tinypack.VarList[dbp.Data32]{t},
	}
	return &l
}
