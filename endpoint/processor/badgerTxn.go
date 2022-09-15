package processor

import (
	b "aurora-relayer-go-common/db/badger"
	"aurora-relayer-go-common/db/badger/core"
	"aurora-relayer-go-common/endpoint"
	"github.com/dgraph-io/badger/v3"
	"golang.org/x/net/context"
)

type BadgerTxn struct {
	db *badger.DB
}

func NewBadgerTxn() endpoint.Processor {
	config := b.GetConfig()
	db, err := core.Open(config.BadgerConfig, config.GcIntervalSeconds)
	if err != nil {
		panic("failed to open badger db")
	}
	processor := BadgerTxn{
		db: db,
	}
	return &processor
}

func (p *BadgerTxn) openTxn(update bool) *badger.Txn {
	return p.db.NewTransaction(update)
}

func (p *BadgerTxn) Pre(ctx context.Context, _ string, _ *endpoint.Endpoint, _ ...any) (context.Context, bool, *any, error) {
	txn := p.openTxn(false)
	childCtx := b.PutTxn(ctx, txn)
	return childCtx, false, nil, nil
}

func (p *BadgerTxn) Post(ctx context.Context, _ string, r *any, err *error) (context.Context, *any, *error) {
	txn := b.GetTxn(ctx)
	if txn != nil {
		if err != nil {
			txn.Commit()
		} else {
			txn.Discard()
		}
	}
	return ctx, r, err
}
