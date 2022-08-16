package log

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"path/filepath"
)

type RollingFileWriter struct {
	io.Writer
}

// NewFileWriter returns a rolling file writer implementation
func NewFileWriter(filePath string) RollingFileWriter {
	return RollingFileWriter{
		&lumberjack.Logger{
			Filename:   filepath.FromSlash(filePath),
			MaxBackups: 5,
			MaxSize:    10,
		},
	}
}
