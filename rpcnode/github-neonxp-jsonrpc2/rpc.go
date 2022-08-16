package github_neonxp_jsonrpc2

import (
	"context"
	"encoding/json"
	"fmt"

	jsonrpc "go.neonxp.dev/jsonrpc2/rpc"
)

func CreateNoParamHandler[RS any](handler func(context.Context) (RS, error)) jsonrpc.HandlerFunc {
	return jsonrpc.HS(handler)
}

func CreateOneParamHandler[RQ, RS any](handler func(context.Context, *RQ) (RS, error)) jsonrpc.HandlerFunc {
	return func(ctx context.Context, in json.RawMessage) (json.RawMessage, error) {
		params := make([]json.RawMessage, 0, 1)
		if err := json.Unmarshal(in, &params); err != nil {
			return nil, jsonrpc.ErrorFromCode(jsonrpc.ErrCodeParseError)
		}
		if len(params) != 1 {
			return nil, jsonrpc.Error{
				Code:    jsonrpc.ErrUser,
				Message: fmt.Sprintf("one parameters required, got: %d", len(params)),
			}
		}
		one := new(RQ)
		if err := json.Unmarshal(params[0], one); err != nil {
			return nil, jsonrpc.ErrorFromCode(jsonrpc.ErrCodeParseError)
		}
		resp, err := handler(ctx, one)
		if err != nil {
			return nil, jsonrpc.Error{
				Code:    jsonrpc.ErrUser,
				Message: err.Error(),
			}
		}
		return json.Marshal(resp)
	}
}

func CreateTwoParamHandler[RQone, RQtwo, RS any](handler func(context.Context, *RQone, *RQtwo) (RS, error)) jsonrpc.HandlerFunc {
	return func(ctx context.Context, in json.RawMessage) (json.RawMessage, error) {
		params := make([]json.RawMessage, 0, 2)
		if err := json.Unmarshal(in, &params); err != nil {
			return nil, jsonrpc.ErrorFromCode(jsonrpc.ErrCodeParseError)
		}
		if len(params) != 2 {
			return nil, jsonrpc.Error{
				Code:    jsonrpc.ErrUser,
				Message: fmt.Sprintf("two parameters required, got: %d", len(params)),
			}
		}
		one := new(RQone)
		if err := json.Unmarshal(params[0], one); err != nil {
			return nil, jsonrpc.ErrorFromCode(jsonrpc.ErrCodeParseError)
		}
		two := new(RQtwo)
		if err := json.Unmarshal(params[1], two); err != nil {
			return nil, jsonrpc.ErrorFromCode(jsonrpc.ErrCodeParseError)
		}
		resp, err := handler(ctx, one, two)
		if err != nil {
			return nil, jsonrpc.Error{
				Code:    jsonrpc.ErrUser,
				Message: err.Error(),
			}
		}
		return json.Marshal(resp)
	}
}
