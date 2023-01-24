package endpoint

import (
	"encoding/json"
	"relayer2-base/db"
	"relayer2-base/log"

	"golang.org/x/net/context"
)

func Process[T any](ctx context.Context, name string, endpoint *Endpoint, handler func(ctx context.Context) (*T, error), args ...any) (*T, error) {

	var resp any
	var err error
	var stop bool
	var childCtx context.Context

	for _, p := range endpoint.Processors {
		childCtx, stop, err = p.Pre(ctx, name, endpoint, &resp, args...)
		defer p.Post(childCtx, name, &resp, &err)
		if stop {
			if err != nil {
				return nil, err
			} else {
				if r, ok := resp.(T); ok {
					return &r, nil
				} else {
					var buff []byte
					buff, err = json.Marshal(resp)
					if err != nil {
						return nil, err
					}
					err = json.Unmarshal(buff, &r)
					if err != nil {
						return nil, err
					}
					return &r, nil
				}
			}
		}
	}

	// we could just 'return handler(childCtx)' but in that case 'resp' would not be passed to 'defer' correctly
	var tmpResp *T
	tmpResp, err = handler(childCtx)
	if err != nil {
		resp = nil
		return nil, err
	}
	resp = tmpResp

	return resp.(*T), err
}

type Endpoint struct {
	DbHandler     db.Handler
	Logger        *log.Logger
	Config        *Config
	WithProcessor func(Processor)
	Processors    []Processor
}

func New(dbh db.Handler) *Endpoint {
	if dbh == nil {
		log.Log().Fatal().Msg("DB Handler should be initialized")
	}
	ep := Endpoint{
		DbHandler:  dbh,
		Logger:     log.Log(),
		Config:     GetConfig(),
		Processors: []Processor{},
	}

	ep.WithProcessor = func(p Processor) {
		withProcessor(&ep, p)
	}

	return &ep
}

func withProcessor(e *Endpoint, p Processor) {
	e.Processors = append(e.Processors, p)
}

func (e *Endpoint) HandleConfigChange() {
	e.Config = GetConfig()
}
