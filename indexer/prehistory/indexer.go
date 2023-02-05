package prehistory

import (
	"context"
	"errors"
	"relayer2-base/db"
	"relayer2-base/log"
	"relayer2-base/types/indexer"
	"relayer2-base/types/primitives"
	"relayer2-base/utils"
	"sync"
	"time"

	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	blankHash = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

type Indexer struct {
	dbh     db.Handler
	config  *Config
	logger  *log.Logger
	lock    sync.Mutex
	started bool
	stopCh  chan bool
	reader  PreHistoryReader
}

type PreHistoryReader struct {
	dbPool    *pgxpool.Pool
	startTime time.Time
}

type queryResultMapping struct {
	ts             uint64
	txsAndRcptRoot primitives.Data32
}

// New creates the prehistory indexer, the db.Handler should not be nil and
// configuration file prehistory part should be properly set
func New(dbh db.Handler) (*Indexer, error) {
	if dbh == nil {
		return nil, errors.New("db handler is not initialized")
	}

	logger := log.Log()
	config := GetConfig()
	if !config.IndexFromPrehistory {
		return nil, nil
	}

	if config.To > config.PrehistoryHeight {
		err := fmt.Errorf("invalid Prehistory Indexer config, to: [%d] must be lower than PrehistoryHeight: [%d]", config.To, config.PrehistoryHeight)
		return nil, err
	}

	if config.To <= config.From {
		err := fmt.Errorf("invalid Prehistory Indexer config, to: [%d] must be greater than from: [%d]", config.To, config.From)
		return nil, err
	}

	i := &Indexer{
		dbh:    dbh,
		logger: logger,
		config: config,
		stopCh: make(chan bool),
		reader: PreHistoryReader{},
	}
	return i, nil
}

// Start starts the prehistory indexing as a goroutine based on the config file settings
func (i *Indexer) Start() {
	i.lock.Lock()
	defer i.lock.Unlock()
	if !i.started {
		i.started = true
		go i.index()
	}
}

// Close gracefully stops the prehistory indexer
func (i *Indexer) Close() {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.logger.Info().Msgf("Prehistory indexer reveived close signal")
	i.stopCh <- true
}

// Start starts the prehistory indexing as a goroutine based on the config file settings
func (i *Indexer) index() {
	var err error
	i.reader.dbPool, err = pgxpool.New(context.Background(), i.config.ArchiveURL)
	if err != nil {
		i.logger.Error().Msgf("Unable to connect to prehistory database %s: %v\n", i.config.ArchiveURL, err)
		return
	}
	defer i.reader.dbPool.Close()

	// Declare and initialize required variables
	emptyBytes := make([]byte, 2)
	epmtyTxsAndRcptRoot := primitives.Data32FromHex("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421") // use keccak(rlp(''))
	defaultTxsAndRcptRoot := primitives.Data32FromBytes(emptyBytes)                                                       // use "0x0..........0"
	d20 := primitives.Data20FromBytes(emptyBytes)
	d32 := primitives.Data32FromBytes(emptyBytes)
	d256 := primitives.Data256FromBytes(emptyBytes)
	quantity := primitives.QuantityFromBytes(emptyBytes)
	parentHash := primitives.Data32FromBytes(emptyBytes)
	blockHash := primitives.Data32FromBytes(emptyBytes)
	chainId := utils.GetChainId(context.Background())
	from := i.config.From
	to := i.config.To
	step := i.config.BatchSize

	i.reader.startTime = time.Now()
	i.logger.Info().Msgf("prehistory indexing started fromBlock: [%d] toBlock: [%d]", from, to)

	// Start retrieveing data from the near postgre DB for the provided range
loop:
	for bId := from; bId < to; bId += step {
		i.logger.Info().Msgf("inserting prehistory blocks: [%d] to [%d]", bId, (bId + step - 1))
		mapQR := make(map[uint64]queryResultMapping, 1000)

		rows, err := i.reader.dbPool.Query(context.Background(), "select block_height, block_timestamp from blocks where block_height>=$1 AND block_height<$2 limit $3", bId, (bId + step), step)
		if err != nil {
			if err != pgx.ErrNoRows {
				i.logger.Error().Msgf("prehistory query failed: %v\n", err)
				return
			}
		}
		// Process the incoming data set if the query returned rows
		for rows.Next() {
			values, err := rows.Values()
			if err != nil {
				i.logger.Error().Msgf("error while iterating prehistory query result: %v\n", err)
			}

			// block_height value
			v0, err := values[0].(pgtype.Numeric).Int64Value()
			if err != nil || !v0.Valid {
				i.logger.Error().Msgf("prehistory block height value failed: %v\n", err)
			}
			bh := uint64(v0.Int64)

			// timestemp value
			v1, err := values[1].(pgtype.Numeric).Int64Value()
			if err != nil || !v1.Valid {
				i.logger.Error().Msgf("prehistory timestamp value failed: %v\n", err)
			}

			mapQR[bh] = queryResultMapping{
				ts:             uint64(v1.Int64) / 1000000000, // convert nano seconds fromat to seconds format
				txsAndRcptRoot: epmtyTxsAndRcptRoot,
			}
		}

		for j := uint64(0); j < step; j++ {
			cBlock := bId + j
			if cBlock >= to {
				i.logger.Info().Msgf("reached to the last block [%d]", cBlock)
				break
			}

			// calculate hash and parent hash (if needed)
			blockHash.Content = utils.ComputeBlockHash(cBlock, chainId)
			if cBlock != 0 && parentHash.Hex() == blankHash {
				parentHash.Content = utils.ComputeBlockHash(cBlock-1, chainId)
			}
			// txs and receipt root hash values are either "0x0...0" or keccak(rlp('')) based on if they are before or after the Genesis
			var txsAndRcptRoot primitives.Data32
			if e, ok := mapQR[cBlock]; ok {
				txsAndRcptRoot = e.txsAndRcptRoot
			} else {
				txsAndRcptRoot = defaultTxsAndRcptRoot
			}

			nBlock := indexer.Block{
				ChainId:          chainId,
				Height:           cBlock,
				Sequence:         cBlock,
				GasLimit:         quantity,
				GasUsed:          quantity,
				Timestamp:        indexer.Timestamp(mapQR[cBlock].ts),
				Hash:             blockHash,
				ParentHash:       parentHash,
				TransactionsRoot: txsAndRcptRoot,
				ReceiptsRoot:     txsAndRcptRoot,
				StateRoot:        d32,
				Miner:            d20,
				LogsBloom:        d256,
			}

			err = i.dbh.InsertBlock(&nBlock)
			if err != nil {
				i.logger.Error().Msgf("failed to insert block [%d], with err: %v\n", nBlock.Height, err)
			}
			parentHash = blockHash
		}
		select {
		case <-i.stopCh:
			break loop
		default:
		}
	}
	i.logger.Info().Msgf("Prehistory indexer ended")
	i.logger.Info().Msgf("Prehistory indexer took %s", time.Since(i.reader.startTime).String())
}
