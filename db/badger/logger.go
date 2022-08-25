package badger

import (
	"aurora-relayer-go-common/log"
	"fmt"
)

type Logger struct {
	log *log.Logger
}

func NewBadgerLogger(log *log.Logger) Logger {
	return Logger{
		log: log,
	}
}

func (l Logger) Errorf(f string, v ...interface{}) {
	l.log.Error().Msg(fmt.Sprintf(f, v))
}

func (l Logger) Warningf(f string, v ...interface{}) {
	l.log.Warn().Msg(fmt.Sprintf(f, v))
}

func (l Logger) Infof(f string, v ...interface{}) {
	l.log.Info().Msg(fmt.Sprintf(f, v))
}

func (l Logger) Debugf(f string, v ...interface{}) {
	l.log.Debug().Msg(fmt.Sprintf(f, v))
}
