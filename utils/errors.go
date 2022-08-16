package utils

import "fmt"

type Error interface {
	Error() string  // returns the message
	ErrorCode() int // returns the code
}

type MethodNotFoundError struct{ Method string }

func (e *MethodNotFoundError) ErrorCode() int { return -32601 }

func (e *MethodNotFoundError) Error() string {
	return fmt.Sprintf("the method %s does not exist/is not available", e.Method)
}

type InvalidParamsError struct{ Message string }

func (e *InvalidParamsError) ErrorCode() int { return -32602 }

func (e *InvalidParamsError) Error() string { return e.Message }

type TxNotFoundError struct {
	Number Uint256
	Hash   H256
}

func (e *TxNotFoundError) ErrorCode() int { return 1 }

func (e *TxNotFoundError) Error() string {
	if e.Number != "" {
		return fmt.Sprintf("tx with number [%s] not found", e.Number)
	}
	if e.Hash != "" {
		return fmt.Sprintf("tx with hash [%s] not found", e.Hash)
	}
	return "tx not found"
}
