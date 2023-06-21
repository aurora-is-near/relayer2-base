package rpc

import (
	"bytes"
	"fmt"
)

func createResponse(idRepr []byte, resultRepr []byte) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, `{"jsonrpc":"2.0","id":%s,"result":%s}`, idRepr, resultRepr)
	return buf.Bytes()
}

func createEventResponse(subscription, resultRepr []byte) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(
		&buf,
		`{"jsonrpc":"2.0","method":"eth_subscription","params":{"subscription":"%s","result":%s}}`,
		subscription,
		resultRepr,
	)
	return buf.Bytes()
}

func createErrorResponse(idRepr []byte, code int64, message string) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, `{"jsonrpc":"2.0","id":%s,"error":{"code":%d,"message":"%s"}}`, idRepr, code, message)
	return buf.Bytes()
}

func createDataErrorResponse(idRepr []byte, code int64, message string, data string) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, `{"jsonrpc":"2.0","id":%s,"error":{"code":%d, "data": "%s", "message":"%s"}}`, idRepr, code, data, message)
	return buf.Bytes()
}
