package github_neonxp_jsonrpc2

import "aurora-relayer-go-common/log"

type Logger struct {
	*log.Log
}

func NewNeonxpJsonRpc2Logger(log *log.Log) Logger {
	return Logger{log}
}

func (l Logger) Logf(format string, args ...interface{}) {
	l.Printf(format, args)
}
