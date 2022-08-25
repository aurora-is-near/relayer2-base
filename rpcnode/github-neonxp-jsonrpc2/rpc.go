package github_neonxp_jsonrpc2

import (
	"aurora-relayer-go-common/utils"
	"context"
	"encoding/json"
	"fmt"

	"go.neonxp.dev/jsonrpc2/rpc"
)

func CreateNoParamHandler[RS any](handler func(context.Context) (RS, error)) rpc.HandlerFunc {
	return func(ctx context.Context, in json.RawMessage) (json.RawMessage, error) {
		params := make([]json.RawMessage, 0)
		if err := json.Unmarshal(in, &params); err != nil {
			return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("error while parsing the rpc: %s", err.Error()))
		}
		if len(params) != 0 {
			return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("two many arguments, want at most 0, got %d", len(params)))
		}
		resp, err := handler(ctx)
		if err != nil {
			return nil, convertToRpcError(err)
		}
		return json.Marshal(resp)
	}
}

func CreateOneParamHandler[RQ, RS any](handler func(context.Context, RQ) (RS, error)) rpc.HandlerFunc {
	return func(ctx context.Context, in json.RawMessage) (json.RawMessage, error) {
		params := make([]json.RawMessage, 0, 1)
		if err := json.Unmarshal(in, &params); err != nil {
			return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("error while parsing the rpc: %s", err.Error()))
		}
		lenParams := len(params)
		if lenParams > 1 {
			return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("two many arguments, want at most 1, got %d", len(params)))
		} else if lenParams < 1 {
			return nil, createRpcError(utils.InvalidParams, "missing value for required argument 0")
		}
		one := new(RQ)
		if err := json.Unmarshal(params[0], one); err != nil {
			return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("invalid argument 0: %s", err.Error()))
		}

		resp, err := handler(ctx, *one)
		if err != nil {
			return nil, convertToRpcError(err)
		}
		return json.Marshal(resp)
	}
}

func CreateOneParamHandlerOptional[RQ, RS any](handler func(context.Context, *RQ) (RS, error)) rpc.HandlerFunc {
	return func(ctx context.Context, in json.RawMessage) (json.RawMessage, error) {
		params := make([]json.RawMessage, 0, 1)
		if err := json.Unmarshal(in, &params); err != nil {
			return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("error while parsing the rpc: %s", err.Error()))
		}

		lenParams := len(params)
		if lenParams > 1 {
			return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("two many arguments, want at most 1, got %d", len(params)))
		}
		one := new(RQ)
		switch lenParams {
		case 0:
			one = nil
		case 1:
			if err := json.Unmarshal(params[0], one); err != nil {
				return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("invalid argument 0: %s", err.Error()))
			}
		}

		resp, err := handler(ctx, one)
		if err != nil {
			return nil, convertToRpcError(err)
		}
		return json.Marshal(resp)
	}
}

func CreateTwoParamHandler[RQone, RQtwo, RS any](handler func(context.Context, RQone, RQtwo) (RS, error)) rpc.HandlerFunc {
	return func(ctx context.Context, in json.RawMessage) (json.RawMessage, error) {
		params := make([]json.RawMessage, 0, 2)
		if err := json.Unmarshal(in, &params); err != nil {
			return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("error while parsing the rpc: %s", err.Error()))
		}

		lenParams := len(params)
		if lenParams > 2 {
			return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("two many arguments, want at most 2, got %d", len(params)))
		} else if lenParams < 2 {
			return nil, createRpcError(utils.InvalidParams, "missing value for required argument")
		}
		one := new(RQone)
		two := new(RQtwo)
		if err := json.Unmarshal(params[0], one); err != nil {
			return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("invalid argument 0: %s", err.Error()))
		}
		if err := json.Unmarshal(params[1], two); err != nil {
			return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("invalid argument 1: %s", err.Error()))
		}

		resp, err := handler(ctx, *one, *two)
		if err != nil {
			return nil, convertToRpcError(err)
		}
		return json.Marshal(resp)
	}
}

func CreateTwoParamHandlerOneOptional[RQone, RQtwo, RS any](handler func(context.Context, RQone, *RQtwo) (RS, error)) rpc.HandlerFunc {
	return func(ctx context.Context, in json.RawMessage) (json.RawMessage, error) {
		params := make([]json.RawMessage, 0, 2)
		if err := json.Unmarshal(in, &params); err != nil {
			return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("error while parsing the rpc: %s", err.Error()))
		}

		lenParams := len(params)
		if lenParams > 2 {
			return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("two many arguments, want at most 2, got %d", len(params)))
		} else if lenParams == 0 {
			return nil, createRpcError(utils.InvalidParams, "missing value for required argument 0")
		}
		one := new(RQone)
		two := new(RQtwo)
		switch lenParams {
		case 1:
			if err := json.Unmarshal(params[0], one); err != nil {
				return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("invalid argument 0: %s", err.Error()))
			}
			two = nil
		case 2:
			if err := json.Unmarshal(params[0], one); err != nil {
				return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("invalid argument 0: %s", err.Error()))
			}
			if err := json.Unmarshal(params[1], two); err != nil {
				return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("invalid argument 1: %s", err.Error()))
			}
		}

		resp, err := handler(ctx, *one, two)
		if err != nil {
			return nil, convertToRpcError(err)
		}
		return json.Marshal(resp)
	}
}

func CreateTwoParamHandlerTwoOptional[RQone, RQtwo, RS any](handler func(context.Context, *RQone, *RQtwo) (RS, error)) rpc.HandlerFunc {
	return func(ctx context.Context, in json.RawMessage) (json.RawMessage, error) {
		params := make([]json.RawMessage, 0, 2)
		if err := json.Unmarshal(in, &params); err != nil {
			return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("error while parsing the rpc: %s", err.Error()))
		}

		lenParams := len(params)
		if lenParams > 2 {
			return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("two many arguments, want at most 2, got %d", len(params)))
		}
		one := new(RQone)
		two := new(RQtwo)
		switch lenParams {
		case 0:
			one = nil
			two = nil
		case 1:
			if err := json.Unmarshal(params[0], one); err != nil {
				return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("invalid argument 0: %s", err.Error()))
			}
			two = nil
		case 2:
			if err := json.Unmarshal(params[0], one); err != nil {
				return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("invalid argument 0: %s", err.Error()))
			}
			if err := json.Unmarshal(params[1], two); err != nil {
				return nil, createRpcError(utils.InvalidParams, fmt.Sprintf("invalid argument 1: %s", err.Error()))
			}
		}

		resp, err := handler(ctx, one, two)
		if err != nil {
			return nil, convertToRpcError(err)
		}
		return json.Marshal(resp)
	}
}

func createRpcError(code int, msg string) rpc.Error {
	return rpc.Error{
		Code:    code,
		Message: msg,
	}
}

func convertToRpcError(err error) rpc.Error {
	i, ok := err.(utils.Error)
	if ok {
		return rpc.Error{
			Code:    i.ErrorCode(),
			Message: i.Error(),
		}
	} else {
		return rpc.Error{
			Code:    utils.Generic,
			Message: err.Error(),
		}
	}
}
