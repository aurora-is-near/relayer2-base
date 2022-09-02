package log

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"sync"
)

type Logger struct {
	zerolog.Logger
}

var global *Logger
var lock = &sync.Mutex{}

// Log returns the common library global logger
func Log() *Logger {
	if global == nil {
		lock.Lock()
		defer lock.Unlock()
		if global == nil {
			global = log()
		}
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
