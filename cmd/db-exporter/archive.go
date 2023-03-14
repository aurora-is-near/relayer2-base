package main

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/aurora-is-near/relayer2-base/db/codec"
	dbt "github.com/aurora-is-near/relayer2-base/types/db"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func NewArchiver(fs afero.Fs, codec codec.Codec) (Archiver, error) {
	ids := fileIDs()
	u := archiver{
		files:    make(map[fileID]*bufio.Writer, len(ids)),
		rawFiles: make(map[fileID]io.WriteCloser, len(ids)),
		codec:    codec,
	}
	for _, id := range ids {
		f, err := fs.Create(id.String())
		if err != nil {
			return nil, errors.Wrap(err, "failed to create file for exporting")
		}
		u.rawFiles[id] = f
		u.files[id] = bufio.NewWriterSize(f, 10*1024*1024)
	}
	return &u, nil
}

type fileID uint64

const (
	blockData   fileID = 0
	blockHash   fileID = 1
	blockHeight fileID = 2
	txData      fileID = 3
	txHash      fileID = 4
	txIndex     fileID = 5
	txHeight    fileID = 6
	logData     fileID = 7
	logIndex    fileID = 8
	logTxIndex  fileID = 9
	logHeight   fileID = 10
)

func fileIDs() []fileID {
	return []fileID{
		blockData,
		blockHash,
		blockHeight,
		txData,
		txHash,
		txIndex,
		txHeight,
		logData,
		logIndex,
		logTxIndex,
		logHeight,
	}
}

func (id fileID) String() string {
	return [...]string{
		"blockData",
		"blockHash",
		"blockHeight",
		"txData",
		"txHash",
		"txIndex",
		"txHeight",
		"logData",
		"logIndex",
		"logTxIndex",
		"logHeight",
	}[id]
}

type archiver struct {
	fs       afero.Fs
	files    map[fileID]*bufio.Writer
	rawFiles map[fileID]io.WriteCloser
	codec    codec.Codec
}

func (a *archiver) Close() error {
	for _, f := range a.files {
		if err := f.Flush(); err != nil {
			return errors.Wrap(err, "failed to flush data to disk")
		}
	}
	for _, f := range a.rawFiles {
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (a *archiver) WriteBlock(bl *dbt.Block) error {
	return a.writeWithLongLen(bl, blockData)
}

func (a *archiver) WriteBlockHash(hash *primitives.Data32) error {
	return a.writeWithShortLen(hash, blockHash)
}

func (a *archiver) WriteBlockHeight(height uint64) error {
	return a.writeWithShortLen(&height, blockHeight)
}

func (a *archiver) WriteLog(log *dbt.Log) error {
	return a.writeWithLongLen(log, logData)
}

func (a *archiver) WriteLogHeight(height uint64) error {
	return a.writeWithShortLen(&height, logHeight)
}

func (a *archiver) WriteLogIndex(index uint64) error {
	return a.writeWithShortLen(&index, logIndex)
}

func (a *archiver) WriteLogTxIndex(txIndex uint64) error {
	return a.writeWithShortLen(&txIndex, logTxIndex)
}

func (a *archiver) WriteTx(tx *dbt.Transaction) error {
	return a.writeWithLongLen(tx, txData)
}

func (a *archiver) WriteTxHash(hash *primitives.Data32) error {
	return a.writeWithShortLen(hash, txHash)
}

func (a *archiver) WriteTxHeight(height uint64) error {
	return a.writeWithShortLen(&height, txHeight)
}

func (a *archiver) WriteTxIndex(index uint64) error {
	return a.writeWithShortLen(&index, txIndex)
}

func (a *archiver) writeWithShortLen(data any, file fileID) error {
	b, err := a.codec.Marshal(data)
	if err != nil {
		return err
	}
	return write[uint8](b, a.files[file])
}

func (a *archiver) writeWithLongLen(data any, file fileID) error {
	b, err := a.codec.Marshal(data)
	if err != nil {
		return err
	}
	return write[uint32](b, a.files[file])
}

func write[LenT uint8 | uint16 | uint32](data []byte, w io.Writer) error {
	err := binary.Write(w, binary.BigEndian, LenT(len(data)))
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
