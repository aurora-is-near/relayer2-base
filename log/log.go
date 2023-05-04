package log

import (
	"io"
	"os"

	"github.com/rs/zerolog"

	"github.com/aurora-is-near/relayer2-base/syncutils"
)

var globalPtr syncutils.LockablePtr[Logger]

type Logger struct {
	zerolog.Logger
}

func (l *Logger) HandleConfigChange(config *Config) {
	oldLvl := zerolog.GlobalLevel()
	if config.Level != oldLvl.String() {
		lvl, err := zerolog.ParseLevel(config.Level)
		if err != nil {
			lvl = oldLvl
		}
		zerolog.SetGlobalLevel(lvl)
	}
}

func Initialize(config *Config) *Logger {
	l := log(config)
	if _, unlock := globalPtr.LockIfNil(); unlock != nil {
		unlock(l)
	} else if _, unlock := globalPtr.LockIfNotNil(); unlock != nil {
		unlock(l)
	}
	return l
}

// Log returns the common library global logger
func Log() *Logger {
	global, unlock := globalPtr.LockIfNil()
	if unlock != nil {
		global = log(DefaultConfig())
		unlock(global)
	}
	return global
}

func log(config *Config) *Logger {
	lvl, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)
	var writers []io.Writer
	if config.LogToConsole {
		writers = append(writers, NewLevelWriter(os.Stdout, os.Stderr))
	}
	if config.LogToFile {
		writers = append(writers, NewFileWriter(config.FilePath))
	}
	return &Logger{zerolog.New(io.MultiWriter(writers...)).With().Timestamp().Logger()}
}
