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
	"fmt"
)

var logIndex = core.NewIndex(5, false)

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
	a := dbp.HexUint(76500507)
	return &a, nil
}

func (h *BlockHandler) GetBlockByHash(ctx context.Context, hash utils.H256) (*dbresponses.Block, error) {

	bh := dbp.DataFromHex[dbp.Len32](hash.String())

	var resp *dbresponses.Block

	h.db.View(func(txn *cc.ViewTxn) error {
		dbKey, err := txn.ReadBlockKey(chainID, bh)
		if err == nil {
			resp, err = txn.ReadBlock(chainID, *dbKey, false)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	})

	fmt.Println(resp.ParentHash)

	return nil, nil
}

func (h *BlockHandler) GetBlockByNumber(ctx context.Context, number utils.Uint256) (*dbresponses.Block, error) {
	return nil, nil
}

func (h *BlockHandler) GetBlockTransactionCountByHash(ctx context.Context, hash utils.H256) (*dbp.HexUint, error) {
	return nil, nil
}

func (h *BlockHandler) GetBlockTransactionCountByNumber(ctx context.Context, number utils.Uint256) (*dbp.HexUint, error) {
	return nil, nil
}

func (h *BlockHandler) GetTransactionByHash(ctx context.Context, hash utils.H256) (*dbresponses.Transaction, error) {
	th := dbp.DataFromHex[dbp.Len32](hash.String())

	var resp *dbresponses.Transaction

	h.db.View(func(txn *cc.ViewTxn) error {
		txKey, err := txn.ReadTxKey(chainID, th)
		if err == nil && txKey != nil {
			resp, err = txn.ReadTx(chainID, *txKey)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	})

	fmt.Println(resp.Hash)

	return nil, nil
}

func (h *BlockHandler) GetTransactionByBlockHashAndIndex(ctx context.Context, hash utils.H256, index utils.Uint256) (*dbresponses.Transaction, error) {
	return nil, nil
}

func (h *BlockHandler) GetTransactionByBlockNumberAndIndex(ctx context.Context, number utils.Uint256, index utils.Uint256) (*dbresponses.Transaction, error) {
	return nil, nil
}

func (h *BlockHandler) GetLogs(ctx context.Context, filter utils.LogFilter) ([]*dbresponses.Log, error) {
	return nil, nil
}

func (h *BlockHandler) GetLogsForTransaction(ctx context.Context, tx *utils.Transaction) ([]*dbresponses.Log, error) {
	return nil, nil

}

func (h *BlockHandler) GetTransactionsForBlock(ctx context.Context, block *utils.Block) ([]*dbresponses.Transaction, error) {
	return nil, nil
}

func (h *BlockHandler) GetBlockHashesSinceNumber(ctx context.Context, number utils.Uint256) ([]*dbp.Data32, error) {
	return nil, nil
}

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

func (h *BlockHandler) BlockHashToNumber(ctx context.Context, hash utils.H256) (*dbp.HexUint, error) {
	return nil, nil
}

func (h *BlockHandler) BlockNumberToHash(ctx context.Context, number utils.Uint256) (*dbp.Data32, error) {
	return nil, nil
}

func (h *BlockHandler) CurrentBlockSequence(ctx context.Context) uint64 {
	return 0
}

func toBlockStore(block *utils.Block) *dbtypes.Block {
	b := dbtypes.Block{
		ParentHash:       dbp.DataFromHex[dbp.Len32](block.ParentHash.String()),
		Miner:            dbp.DataFromBytes[dbp.Len20](block.Miner.Bytes()),
		Timestamp:        uint64(block.Timestamp),
		GasLimit:         dbp.Quantity(dbp.DataFromBytes[dbp.Len32](block.GasLimit.Bytes())),
		GasUsed:          dbp.Quantity(dbp.DataFromBytes[dbp.Len32](block.GasUsed.Bytes())),
		LogsBloom:        dbp.DataFromBytes[dbp.Len256]([]byte(block.LogsBloom)),
		TransactionsRoot: dbp.DataFromBytes[dbp.Len32](block.TransactionsRoot.Bytes()),
		StateRoot:        dbp.DataFromBytes[dbp.Len32]([]byte(block.StateRoot)),
		ReceiptsRoot:     dbp.DataFromBytes[dbp.Len32](block.ReceiptsRoot.Bytes()),
		Size:             block.Size.Uint64(),
	}
	return &b
}

func toTxnStore(txn *utils.Transaction) *dbtypes.Transaction {

	toOrContract := dbp.DataFromHex[dbp.Len20]("0x0") // TODO: corrupted block both contract address and to address is null, what should we do?
	isContractDeployment := false
	if txn.ContractAddress != nil {
		if txn.To != nil {
			log.Log().Warn().Msgf("both contract address and to address is set for txn: [%v], to: [%s], contract: [%s]",
				txn.Hash, txn.To.String(), txn.ContractAddress.String())
		}
		isContractDeployment = true
		toOrContract = dbp.DataFromBytes[dbp.Len20](txn.ContractAddress.Bytes())
	} else if txn.To != nil {
		toOrContract = dbp.DataFromBytes[dbp.Len20](txn.To.Bytes())
	} else {
		log.Log().Warn().Msgf("both contract address and to address is null for txn: [%v]", txn.Hash)
	}

	var accessListEntries []dbtypes.AccessListEntry
	for _, al := range txn.AccessList {
		var storageKeys []dbp.Data32
		for _, sk := range al.StorageKeys {
			storageKey := dbp.DataFromBytes[dbp.Len32](sk.Bytes())
			storageKeys = append(storageKeys, storageKey)
		}
		sk := tinypack.CreateList[dbp.VarLen, dbp.Data32](storageKeys...)
		accessListEntry := dbtypes.AccessListEntry{
			Address:     dbp.DataFromBytes[dbp.Len20](al.Address.Bytes()),
			StorageKeys: tinypack.VarList[dbp.Data32]{sk},
		}
		accessListEntries = append(accessListEntries, accessListEntry)
	}
	ake := tinypack.CreateList[dbp.VarLen, dbtypes.AccessListEntry](accessListEntries...)

	t := dbtypes.Transaction{
		Type:                 uint64(txn.TxType),
		From:                 dbp.DataFromBytes[dbp.Len20](txn.From.Bytes()),
		IsContractDeployment: isContractDeployment,
		ToOrContract:         toOrContract,
		Nonce:                dbp.Quantity(dbp.DataFromBytes[dbp.Len32](txn.Nonce.Bytes())),
		GasPrice:             dbp.Quantity(dbp.DataFromBytes[dbp.Len32](txn.GasPrice.Bytes())),
		GasLimit:             dbp.Quantity(dbp.DataFromBytes[dbp.Len32](txn.GasLimit.Bytes())),
		GasUsed:              txn.GasUsed.Uint64(),
		Value:                dbp.Quantity(dbp.DataFromBytes[dbp.Len32](txn.Value.Bytes())),
		Input:                dbp.DataFromBytes[dbp.VarLen](txn.Input),
		NearHash:             dbp.DataFromBytes[dbp.Len32](txn.NearTransaction.Hash.Bytes()),
		NearReceiptHash:      dbp.DataFromBytes[dbp.Len32](txn.NearTransaction.ReceiptHash.Bytes()),
		Status:               txn.Status,
		V:                    txn.V,
		R:                    dbp.Quantity(dbp.DataFromBytes[dbp.Len32](txn.R.Bytes())),
		S:                    dbp.Quantity(dbp.DataFromBytes[dbp.Len32](txn.S.Bytes())),
		LogsBloom:            dbp.DataFromHex[dbp.Len256](txn.LogsBloom),
		AccessList:           tinypack.VarList[dbtypes.AccessListEntry]{ake},
		MaxPriorityFeePerGas: dbp.Quantity(dbp.DataFromBytes[dbp.Len32](txn.MaxPriorityFeePerGas.Bytes())),
		MaxFeePerGas:         dbp.Quantity(dbp.DataFromBytes[dbp.Len32](txn.MaxFeePerGas.Bytes())),
	}
	return &t
}

func toLogStore(log *utils.Log) *dbtypes.Log {

	var topics []dbp.Data32
	for _, t := range log.Topics {
		topic := dbp.DataFromBytes[dbp.Len32](t)
		topics = append(topics, topic)
	}
	t := tinypack.CreateList[dbp.VarLen, dbp.Data32](topics...)

	l := dbtypes.Log{
		Address: dbp.DataFromBytes[dbp.Len20](log.Address.Bytes()),
		Data:    dbp.DataFromBytes[dbp.VarLen](log.Data),
		Topics:  tinypack.VarList[dbp.Data32]{t},
	}
	return &l
}
