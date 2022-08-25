package log

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"io"
	"os"
	"sync"
)

const (
	defaultLogFilePath = "/tmp/relayer/log"
	defaultLogLevel    = "info"
	defaultLogToFile   = true
	defaultLogToStdOut = true
)

type config struct {
	LogToFile    bool
	LogToConsole bool
	Level        string
	FilePath     string
}

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
	config := defaultConfig()
	if err := viper.Sub("Logger"); err != nil {
		_ = viper.Unmarshal(&config)
	}
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

func defaultConfig() config {
	return config{
		LogToFile:    defaultLogToFile,
		LogToConsole: defaultLogToStdOut,
		Level:        defaultLogLevel,
		FilePath:     defaultLogFilePath,
	}
}
