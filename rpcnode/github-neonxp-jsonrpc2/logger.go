package github_neonxp_jsonrpc2

import "github.com/aurora-is-near/relayer2-base/log"

type Logger struct {
	*log.Logger
}

func NewNeonxpJsonRpc2Logger(log *log.Logger) Logger {
	return Logger{log}
}

func (l Logger) Logf(format string, args ...interface{}) {
	l.Printf(format, args)
}
