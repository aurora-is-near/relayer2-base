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

type GenericError struct{ Err error }

func (e *GenericError) ErrorCode() int { return -32000 }

func (e *GenericError) Error() string { return e.Err.Error() }
