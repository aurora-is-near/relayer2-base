package errors

import "fmt"

const (
	// Use the same error code with Aurora Plus infra
	TxsStatus = 3

	// -32768 to -32000 are reserved for pre-defined errors by JSON-RPC standart.
	Generic        = -32000
	MethodNotFound = -32601
	InvalidParams  = -32602

	// -32999 to -32900 space is reserved for Aurora Relayer application specific errors.
	KeyNotFound = -32900
)

type Error interface {
	Error() string  // returns the message
	ErrorCode() int // returns the code
}

type MethodNotFoundError struct{ Method string }

func (e *MethodNotFoundError) ErrorCode() int { return MethodNotFound }

func (e *MethodNotFoundError) Error() string {
	return fmt.Sprintf("the method %s does not exist/is not available", e.Method)
}

type InvalidParamsError struct{ Message string }

func (e *InvalidParamsError) ErrorCode() int { return InvalidParams }

func (e *InvalidParamsError) Error() string { return e.Message }

type GenericError struct{ Err error }

func (e *GenericError) ErrorCode() int { return Generic }

func (e *GenericError) Error() string { return e.Err.Error() }

type KeyNotFoundError struct{}

func (e *KeyNotFoundError) ErrorCode() int { return KeyNotFound }

func (e *KeyNotFoundError) Error() string {
	return "record not found in DB"
}

type TxsStatusError struct{ Message string }

func (e *TxsStatusError) ErrorCode() int { return TxsStatus }

func (e *TxsStatusError) Error() string { return e.Message }

type TxsRevertError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *TxsRevertError) ErrorCode() int { return TxsStatus }

func (e *TxsRevertError) Error() string { return e.Message }

func (e *TxsRevertError) ErrorData() interface{} { return e.Data }
