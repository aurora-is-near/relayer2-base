package log

import (
	"io"
	"os"
	"relayer2-base/syncutils"

	"github.com/rs/zerolog"
)

var globalPtr syncutils.LockablePtr[Logger]

type Logger struct {
	zerolog.Logger
}

func (l *Logger) HandleConfigChange() {
	oldLvl := zerolog.GlobalLevel()
	cfg := GetConfig()
	if cfg.Level != oldLvl.String() {
		lvl, err := zerolog.ParseLevel(cfg.Level)
		if err != nil {
			lvl = oldLvl
		}
		zerolog.SetGlobalLevel(lvl)
	}
}

// Log returns the common library global logger
func Log() *Logger {
	global, unlock := globalPtr.LockIfNil()
	if unlock != nil {
		global = log()
		unlock(global)
	}
	return global
}

func log() *Logger {
	config := GetConfig()
	l, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		l = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(l)
	var writers []io.Writer
	if config.LogToConsole {
		writers = append(writers, NewLevelWriter(os.Stdout, os.Stderr))
	}
	if config.LogToFile {
		writers = append(writers, NewFileWriter(config.FilePath))
	}
	return &Logger{zerolog.New(io.MultiWriter(writers...)).With().Timestamp().Logger()}
}
