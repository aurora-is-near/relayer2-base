package main

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/aurora-is-near/relayer2-base/db/codec"
	dbt "github.com/aurora-is-near/relayer2-base/types/db"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/spf13/afero"
)

func NewUnarchiver(fs afero.Fs, codec codec.Codec) (Unarchiver, error) {
	ids := fileIDs()
	u := unarchiver{
		files:    make(map[fileID]*bufio.Reader, len(ids)),
		rawFiles: make(map[fileID]io.ReadCloser, len(ids)),
		codec:    codec,
	}
	for _, id := range ids {
		f, err := fs.Open(id.String())
		if err != nil {
			return nil, err
		}

		u.rawFiles[id] = f
		u.files[id] = bufio.NewReaderSize(f, 10*1024*1024)
	}
	return &u, nil
}

type unarchiver struct {
	files    map[fileID]*bufio.Reader
	rawFiles map[fileID]io.ReadCloser
	codec    codec.Codec
}

func (u *unarchiver) Close() error {
	for _, f := range u.rawFiles {
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (u *unarchiver) ReadBlock() (*dbt.Block, error) {
	return read[uint32, dbt.Block](u.files[blockData], u.codec)
}

func (u *unarchiver) ReadBlockHash() (*primitives.Data32, error) {
	return read[uint8, primitives.Data32](u.files[blockHash], u.codec)
}

func (u *unarchiver) ReadBlockHeight() (uint64, error) {
	h, err := read[uint8, uint64](u.files[blockHeight], u.codec)
	return *h, err
}

func (u *unarchiver) ReadLog() (*dbt.Log, error) {
	return read[uint32, dbt.Log](u.files[logData], u.codec)
}

func (u *unarchiver) ReadLogHeight() (uint64, error) {
	h, err := read[uint8, uint64](u.files[logHeight], u.codec)
	return *h, err
}

func (u *unarchiver) ReadLogIndex() (uint64, error) {
	i, err := read[uint8, uint64](u.files[logIndex], u.codec)
	return *i, err
}

func (u *unarchiver) ReadLogTxIndex() (uint64, error) {
	i, err := read[uint8, uint64](u.files[logTxIndex], u.codec)
	return *i, err
}

func (u *unarchiver) ReadTx() (*dbt.Transaction, error) {
	return read[uint32, dbt.Transaction](u.files[txData], u.codec)
}

func (u *unarchiver) ReadTxHash() (*primitives.Data32, error) {
	return read[uint8, primitives.Data32](u.files[txHash], u.codec)
}

func (u *unarchiver) ReadTxHeight() (uint64, error) {
	h, err := read[uint8, uint64](u.files[txHeight], u.codec)
	return *h, err
}

func (u *unarchiver) ReadTxIndex() (uint64, error) {
	i, err := read[uint8, uint64](u.files[txIndex], u.codec)
	return *i, err
}

func read[LT uint8 | uint16 | uint32, R any](r io.Reader, codec codec.Decoder) (*R, error) {
	var l LT
	err := binary.Read(r, binary.BigEndian, &l)
	if err != nil {
		return nil, err
	}
	data := make([]byte, l)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, err
	}
	res := new(R)
	err = codec.Unmarshal(data, res)
	return res, err
}
