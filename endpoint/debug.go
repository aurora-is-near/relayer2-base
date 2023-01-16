package endpoint

import (
	"aurora-relayer-go-common/types/common"
	errs "aurora-relayer-go-common/types/errors"
	"aurora-relayer-go-common/types/response"
	"errors"

	"golang.org/x/net/context"
)

type Debug struct {
	*Endpoint
}

func NewDebug(endpoint *Endpoint) *Debug {
	return &Debug{endpoint}
}

// TraceTransaction attempts to run the transaction in the exact same manner as it was executed on the network.
// It replays any transaction that may have been executed prior to this one before it will finally attempt to
// execute the transaction that corresponds to the given hash.
//
//	On missing or invalid param returns errors code '-32602' with custom message.
//	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
func (d *Debug) TraceTransaction(_ context.Context, hash common.H256) (*response.CallFrame, error) {
	err := errors.New("`debug_traceTransaction` method is not supported. Please use the proxy feature to consume it from the Aurora Infrastructure")
	return nil, &errs.GenericError{Err: err}
}
