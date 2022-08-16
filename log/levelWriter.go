package log

import (
	"github.com/rs/zerolog"
	"io"
)

type LevelWriter struct {
	io.Writer
	ErrorWriter io.Writer
}

// NewLevelWriter returns a LevelWriter given two io.Writers for errors and others
func NewLevelWriter(stdWriter io.Writer, errWriter io.Writer) LevelWriter {
	return LevelWriter{
		stdWriter,
		errWriter,
	}
}

// WriteLevel writes logs to different writers for different log levels
func (lw *LevelWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	w := lw.Writer
	if level > zerolog.InfoLevel {
		w = lw.ErrorWriter
	}
	return w.Write(p)
}
