package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"os"

	"github.com/aurora-is-near/relayer2-base/db/codec"
	dbt "github.com/aurora-is-near/relayer2-base/types/db"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
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
		f, err := fs.OpenFile(id.String(), os.O_TRUNC|os.O_WRONLY, 0666)
		if errors.Is(err, afero.ErrFileNotFound) {
			f, err = fs.Create(id.String())
		}
		if err != nil {
			return nil, err
		}
		u.rawFiles[id] = f
		u.files[id] = bufio.NewWriterSize(f, 10*1024*1024)
	}
	return &u, nil
}

type fileID uint64

const (
	blockData fileID = iota
	blockHash
	blockHeight
	txData
	txHash
	txIndex
	txHeight
	logData
	logIndex
	logTxIndex
	logHeight
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
			return err
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
	return a.writeWithLen(bl, blockData)
}

func (a *archiver) WriteBlockHash(hash *primitives.Data32) error {
	return a.writeWithLen(hash, blockHash)
}

func (a *archiver) WriteBlockHeight(height uint64) error {
	return a.writeWithLen(&height, blockHeight)
}

func (a *archiver) WriteLog(log *dbt.Log) error {
	return a.writeWithLen(log, logData)
}

func (a *archiver) WriteLogHeight(height uint64) error {
	return a.writeWithLen(&height, logHeight)
}

func (a *archiver) WriteLogIndex(index uint64) error {
	return a.writeWithLen(&index, logIndex)
}

func (a *archiver) WriteLogTxIndex(txIndex uint64) error {
	return a.writeWithLen(&txIndex, logTxIndex)
}

func (a *archiver) WriteTx(tx *dbt.Transaction) error {
	return a.writeWithLen(tx, txData)
}

func (a *archiver) WriteTxHash(hash *primitives.Data32) error {
	return a.writeWithLen(hash, txHash)
}

func (a *archiver) WriteTxHeight(height uint64) error {
	return a.writeWithLen(&height, txHeight)
}

func (a *archiver) WriteTxIndex(index uint64) error {
	return a.writeWithLen(&index, txIndex)
}

func (a *archiver) writeWithLen(data any, file fileID) error {
	b, err := a.codec.Marshal(data)
	if err != nil {
		return err
	}
	err = binary.Write(a.files[file], binary.BigEndian, uint16(len(b)))
	if err != nil {
		return err
	}
	_, err = a.files[file].Write(b)
	return err
}
