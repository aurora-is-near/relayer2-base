package endpoint

import (
	"aurora-relayer-go-common/db"
	"aurora-relayer-go-common/log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

const (
	ConfigPath = "endpoint"
)

func Process[T any](ctx context.Context, name string, endpoint *Endpoint, handler func(ctx context.Context) (T, error), args ...any) (T, error) {

	var zeroVal T
	var resp any
	var err error
	var stop bool
	var childCtx context.Context

	for _, p := range endpoint.Processors {
		childCtx, stop, resp, err = p.Pre(ctx, name, endpoint, args)
		defer p.Post(childCtx, name, &resp, &err)
		if stop {
			if err != nil {
				return zeroVal, err
			} else {
				return resp.(T), nil
			}
		}
	}

	resp, err = handler(childCtx)
	if err != nil {
		return zeroVal, err
	}

	return resp.(T), nil
}

type Endpoint struct {
	DbHandler     *db.Handler
	Logger        *log.Logger
	Config        *Config
	WithProcessor func(Processor)
	Processors    []Processor
}

func New(dbh db.Handler) *Endpoint {
	if dbh == nil {
		panic("DB Handler should be initialized")
	}
	ep := Endpoint{
		DbHandler:  &dbh,
		Logger:     log.Log(),
		Config:     GetConfig(),
		Processors: []Processor{},
	}

	ep.WithProcessor = func(p Processor) {
		withProcessor(&ep, p)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		handleConfigChange(&ep)
	})

	return &ep
}

func withProcessor(e *Endpoint, p Processor) {
	e.Processors = append(e.Processors, p)
}

func handleConfigChange(e *Endpoint) {
	e.Config = GetConfig()
}
