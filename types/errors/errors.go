package errors

import (
	"fmt"
)

const (
	// Use the same error code with Aurora Plus infra
	TxsStatus = 3

	// -32000 to -32768 are reserved for pre-defined errors by JSON-RPC standard.
	Generic               = -32000
	LogRangeLimitExceeded = -32005
	InvalidRequest        = -32600
	MethodNotFound        = -32601
	InvalidParams         = -32602
	Internal              = -32603
	ParseError            = -32700

	// -32000 to -32999 space is reserved for Aurora Relayer application specific errors.
	KeyNotFound = -32900
)

type Error interface {
	Error() string  // returns the message
	ErrorCode() int // returns the code
}

type DataError interface {
	Error() string     // returns the message
	ErrorData() string // returns the error data
}

type MethodNotFoundError struct{ Method string }

func (e *MethodNotFoundError) ErrorCode() int { return MethodNotFound }

func (e *MethodNotFoundError) Error() string {
	return fmt.Sprintf("the method %s does not exist/is not available", e.Method)
}

// received rpc request not valid
type InvalidRequestError struct{ Message string }

func (e *InvalidRequestError) ErrorCode() int { return InvalidRequest }

func (e *InvalidRequestError) Error() string { return e.Message }

// received rpc request parameters not valid
type InvalidParamsError struct{ Message string }

func (e *InvalidParamsError) ErrorCode() int { return InvalidParams }

func (e *InvalidParamsError) Error() string { return e.Message }

// provided event subscription method does not exist
type SubscriptionNotFoundError struct{ Subscription string }

func (e *SubscriptionNotFoundError) ErrorCode() int { return -32601 }

func (e *SubscriptionNotFoundError) Error() string {
	return fmt.Sprintf("no %q subscription found", e.Subscription)
}

type GenericError struct{ Err error }

func (e *GenericError) ErrorCode() int {
	err, ok := e.Err.(Error)
	if ok {
		return err.ErrorCode()
	}
	return Generic
}

func (e *GenericError) Error() string { return e.Err.Error() }

type KeyNotFoundError struct{}

func (e *KeyNotFoundError) ErrorCode() int { return KeyNotFound }

func (e *KeyNotFoundError) Error() string {
	return "record not found in DB"
}

type LogResponseRangeLimitError struct{ Err error }

func (e *LogResponseRangeLimitError) ErrorCode() int { return LogRangeLimitExceeded }

func (e *LogResponseRangeLimitError) Error() string { return e.Err.Error() }

type TxsStatusError struct{ Message string }

func (e *TxsStatusError) ErrorCode() int { return TxsStatus }

func (e *TxsStatusError) Error() string { return e.Message }

type TxsRevertError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func (e *TxsRevertError) ErrorCode() int { return e.Code }

func (e *TxsRevertError) Error() string { return e.Message }

func (e *TxsRevertError) ErrorData() string { return e.Data }

type InternalError struct{ Message string }

func (e *InternalError) ErrorCode() int { return Internal }

func (e *InternalError) Error() string {
	return e.Message
}
