package github_ethereum_go_ethereum

import (
	"github.com/aurora-is-near/relayer2-base/log"
	gel "github.com/ethereum/go-ethereum/log"
)

type Logger struct {
	log *log.Logger
}

func NewGoEthLogger(log *log.Logger) Logger {
	return Logger{
		log: log,
	}
}

func (l Logger) New(ctx ...interface{}) gel.Logger {
	return Logger{
		log: &log.Logger{Logger: l.log.With().Fields(ctx).Logger()},
	}
}

func (l Logger) GetHandler() gel.Handler {
	return nil
}

func (l Logger) SetHandler(h gel.Handler) {

}

func (l Logger) Trace(msg string, ctx ...interface{}) {
	l.log.Trace().Fields(ctx).Msg(msg)
}

func (l Logger) Debug(msg string, ctx ...interface{}) {
	l.log.Debug().Fields(ctx).Msg(msg)
}

func (l Logger) Info(msg string, ctx ...interface{}) {
	l.log.Info().Fields(ctx).Msg(msg)
}

func (l Logger) Warn(msg string, ctx ...interface{}) {
	l.log.Warn().Fields(ctx).Msg(msg)
}

func (l Logger) Error(msg string, ctx ...interface{}) {
	l.log.Error().Fields(ctx).Msg(msg)
}

func (l Logger) Crit(msg string, ctx ...interface{}) {
	l.log.Panic().Fields(ctx).Msg(msg)
}
