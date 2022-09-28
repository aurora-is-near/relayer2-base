package processor

import (
	b "aurora-relayer-go-common/db/badger"
	"aurora-relayer-go-common/db/badger/core"
	"aurora-relayer-go-common/endpoint"
	"aurora-relayer-go-common/log"
	"github.com/dgraph-io/badger/v3"
	"golang.org/x/net/context"
)

type BadgerTxn struct {
	db         *badger.DB
	excludeTxn map[string]bool
	updateTxn  map[string]bool
}

func NewBadgerTxn() endpoint.Processor {
	config := b.GetConfig()
	db, err := core.Open(config.BadgerConfig, config.GcIntervalSeconds)
	if err != nil {
		panic("failed to open badger db")
	}

	excludeTxn := make(map[string]bool, len(config.ExcludeTxn))
	for _, endpointName := range config.ExcludeTxn {
		excludeTxn[endpointName] = true
	}

	updateTxn := make(map[string]bool, len(config.UpdateTxn))
	for _, endpointName := range config.UpdateTxn {
		updateTxn[endpointName] = true
	}

	processor := BadgerTxn{
		db:         db,
		excludeTxn: excludeTxn,
		updateTxn:  updateTxn,
	}
	return &processor
}

func (p *BadgerTxn) openTxn(update bool) *badger.Txn {
	return p.db.NewTransaction(update)
}

func (p *BadgerTxn) Pre(ctx context.Context, name string, _ *endpoint.Endpoint, _ ...any) (context.Context, bool, *any, error) {
	if p.excludeTxn[name] == true {
		return ctx, false, nil, nil
	}
	txn := p.openTxn(p.updateTxn[name])
	childCtx := b.PutTxn(ctx, txn)
	return childCtx, false, nil, nil
}

func (p *BadgerTxn) Post(ctx context.Context, name string, r *any, err *error) (context.Context, *any, *error) {
	if p.excludeTxn[name] == true {
		return ctx, r, err
	}
	txn := b.GetTxn(ctx)
	if txn != nil {
		if err != nil {
			err := txn.Commit()
			if err != nil {
				log.Log().Err(err).Msgf("failed to commit transaction, endpoint: [%s]", name)
				txn.Discard()
			}
		} else {
			txn.Discard()
		}
	}
	return ctx, r, err
}
