package github_neonxp_jsonrpc2

import "aurora-relayer-go-common/log"

type Logger struct {
	*log.Logger
}

func NewNeonxpJsonRpc2Logger(log *log.Logger) Logger {
	return Logger{log}
}

func (l Logger) Logf(format string, args ...interface{}) {
	l.Printf(format, args)
}
