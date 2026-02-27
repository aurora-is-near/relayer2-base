package types

import (
	"time"

	"github.com/aurora-is-near/relayer2-base/db/badger/core/dbkey"
	"github.com/aurora-is-near/relayer2-base/types/db"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
)

type Filter struct {
	FromBlock *uint64
	FromTxn   *uint64
	FromLog   *uint64
	ToBlock   *uint64
	ToTxn     *uint64
	ToLog     *uint64
	Addresses []primitives.Data20
	Topics    [][]primitives.Data32
}

func (f *Filter) ToLogFilter() *db.LogFilter {

	var fb, tb, ft, tt, fl, tl uint64
	tt = dbkey.MaxTxIndex
	tl = dbkey.MaxLogIndex

	if f.FromBlock != nil {
		fb = *f.FromBlock
	}
	if f.ToBlock != nil {
		tb = *f.ToBlock
	}
	if f.FromTxn != nil {
		ft = *f.FromTxn
	}
	if f.ToTxn != nil {
		tt = *f.ToTxn
	}
	if f.FromLog != nil {
		fl = *f.FromLog
	}
	if f.ToLog != nil {
		tl = *f.ToLog
	}

	return &db.LogFilter{
		CreatedAt: uint64(time.Now().UnixNano()),
		Metadata:  primitives.VarDataFromBytes(nil),
		From:      db.LogKey{BlockHeight: fb, TransactionIndex: ft, LogIndex: fl},
		To:        db.LogKey{BlockHeight: tb, TransactionIndex: tt, LogIndex: tl},
		Addresses: f.Addresses,
		Topics:    f.Topics,
	}
}

func (f *Filter) ToBlockFilter() *db.BlockFilter {
	var fb, tb uint64
	if f.FromBlock != nil {
		fb = *f.FromBlock
	}
	if f.ToBlock != nil {
		tb = *f.ToBlock
	}
	return &db.BlockFilter{
		CreatedAt: uint64(time.Now().UnixNano()),
		Metadata:  primitives.VarDataFromBytes(nil),
		From:      db.BlockKey{Height: fb},
		To:        db.BlockKey{Height: tb},
	}
}

func (f *Filter) ToTxnFilter() *db.TransactionFilter {
	var fb, tb, ft, tt uint64
	if f.FromBlock != nil {
		fb = *f.FromBlock
	}
	if f.ToBlock != nil {
		tb = *f.ToBlock
	}
	if f.FromTxn != nil {
		ft = *f.FromTxn
	}
	if f.ToTxn != nil {
		tt = *f.ToTxn
	}

	return &db.TransactionFilter{
		CreatedAt: uint64(time.Now().UnixNano()),
		Metadata:  primitives.VarDataFromBytes(nil),
		From:      db.TransactionKey{BlockHeight: fb, TransactionIndex: ft},
		To:        db.TransactionKey{BlockHeight: tb, TransactionIndex: tt},
	}
}
