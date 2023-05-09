package tar

import (
	"bytes"
	"fmt"

	"github.com/aurora-is-near/stream-backup/chunks"
	"github.com/aurora-is-near/stream-backup/messagebackup"
	"github.com/fxamacker/cbor/v2"

	"github.com/aurora-is-near/relayer2-base/db"
	"github.com/aurora-is-near/relayer2-base/db/codec"
	"github.com/aurora-is-near/relayer2-base/log"
	"github.com/aurora-is-near/relayer2-base/types/indexer"
)

type Indexer struct {
	Config *Config
	dbh    db.Handler
	reader chunks.Chunks
	logger *log.Logger
	mode   cbor.DecMode
}

func New(dbh db.Handler) (*Indexer, error) {
	logger := log.Log()
	config := GetConfig()

	if !config.IndexFromBackup {
		return &Indexer{Config: config}, nil
	}

	i := &Indexer{
		Config: config,
		dbh:    dbh,
		logger: logger,
		mode:   codec.CborDecoder(),
		reader: chunks.Chunks{
			Dir:             config.Dir,
			ChunkNamePrefix: config.NamePrefix,
		},
	}
	return i, nil
}

func (i *Indexer) Start() {
	if !i.Config.IndexFromBackup {
		return
	}

	if err := i.reader.Open(); err != nil {
		i.logger.Fatal().Err(err).Msg("failed to open backup file reader")
	}
	defer i.reader.CloseReader()

	if err := i.reader.SeekReader(i.Config.From); err != nil {
		i.logger.Fatal().Err(err).Msg("failed to position file reader")
	}

	for {
		seq, data, err := i.reader.ReadNext()
		if err != nil {
			if err.Error() == "not found" {
				break
			}
			i.logger.Fatal().Err(err).Msgf("reader failed to read")
		}

		var m messagebackup.MessageBackup
		if err = m.UnmarshalVT(data); err != nil {
			i.logger.Fatal().Err(err).Msgf("failed to decode backup [%d]", seq)
		}

		block, err := DecodeAugmentedCBOR[indexer.Block](m.Data, i.mode)
		if err != nil {
			i.logger.Fatal().Err(err).Msgf("failed decode block from backup seq [%d]", m.Sequence)
		}

		if seq%1000 == 0 {
			i.logger.Info().Msgf("inserting backup block: [%d] - seq: [%d]", block.Height, seq)
		}

		err = i.dbh.InsertBlock(block)
		if err != nil {
			i.logger.Fatal().Err(err).Msgf("failed to insert block [%d]", block.Height)
		}
		if i.Config.To != 0 && i.Config.To == seq {
			break
		}
	}
	i.logger.Info().Msgf("backup indexer finished")
}

func (i *Indexer) Close() {
}

func DecodeAugmentedCBOR[T any](input []byte, mode cbor.DecMode) (*T, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("input is too short")
	}

	reader := bytes.NewReader(input[1:])
	var decoder *cbor.Decoder
	if mode != nil {
		decoder = mode.NewDecoder(reader)
	} else {
		decoder = cbor.NewDecoder(reader)
	}

	a := new(cbor.RawMessage)
	if err := decoder.Decode(a); err != nil {
		return nil, err
	}
	b := new(T)
	if err := decoder.Decode(b); err != nil {
		return nil, err
	}

	return b, nil
}
