package endpoint

import (
	"aurora-relayer-go-common/db"
	"aurora-relayer-go-common/log"
)

const (
	ConfigPath = "Endpoint"
)

type Preprocessor func(name string, endpoint *Endpoint, args ...any) (bool, *any, error)

func Preprocess[T any](name string, endpoint *Endpoint, handler func() (T, error), args ...any) (T, error) {

	var zeroVal T
	var resp any
	var err error
	var stop bool

	for _, m := range endpoint.Preprocessors {
		stop, resp, err = m(name, endpoint, args)
		if stop {
			if err != nil {
				return zeroVal, err
			} else {
				return resp.(T), nil
			}
		}
	}

	resp, err = handler()
	if err != nil {
		return zeroVal, err
	}

	return resp.(T), nil
}

type Endpoint struct {
	DbHandler        *db.Handler
	Logger           *log.Log
	WithPreprocessor func(Preprocessor)
	Preprocessors    []Preprocessor
}

func New(dbh *db.Handler) *Endpoint {
	if dbh == nil {
		panic("DB Handler should be initialized")
	}
	logger := log.New()
	ep := Endpoint{
		DbHandler:     dbh,
		Logger:        logger,
		Preprocessors: []Preprocessor{},
	}

	ep.WithPreprocessor = func(p Preprocessor) {
		withPreprocessor(&ep, p)
	}

	return &ep
}

func withPreprocessor(e *Endpoint, p Preprocessor) {
	e.Preprocessors = append(e.Preprocessors, p)
}
