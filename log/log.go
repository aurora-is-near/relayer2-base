package log

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"io"
	"os"
)

const (
	DefaultLogFilePath = "/tmp/relayer/log"
	DefaultLogLevel    = "info"
	DefaultLogToFile   = true
	DefaultLogToStdOut = true
)

type Log struct {
	zerolog.Logger
}

type Config struct {
	LogToFile    bool
	LogToConsole bool
	Level        string
	FilePath     string
}

// New returns common library logger with default config
// Common library logger is a zerolog implementation
func New() *Log {
	logConfig := DefaultConfig()
	_ = viper.Sub("Log").Unmarshal(&logConfig)
	return NewWithConf(logConfig)
}

// NewWithConf returns common library logger with the specified config
// Common library logger is a zerolog implementation
func NewWithConf(config Config) *Log {
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
	logger := zerolog.New(io.MultiWriter(writers...)).With().Timestamp().Logger()

	return &Log{
		logger,
	}
}

// DefaultConfig returns the default configuration of common library logger
func DefaultConfig() Config {
	return Config{
		LogToFile:    DefaultLogToFile,
		LogToConsole: DefaultLogToStdOut,
		Level:        DefaultLogLevel,
		FilePath:     DefaultLogFilePath,
	}
}
