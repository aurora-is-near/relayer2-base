package endpoint

import (
	"aurora-relayer-go-common/utils"
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

type Web3 struct {
	*Endpoint
}

func NewWeb3(endpoint *Endpoint) *Web3 {
	return &Web3{endpoint}
}

func (ep *Web3) ClientVersion(_ context.Context) (string, error) {
	if err := ep.IsEndpointAllowed("web3_clientVersion"); err != nil {
		return "", err
	}
	return "Aurora Relayer", nil
}

// A sample method to show the usage of single mandatory parameter
func (ep *Web3) Sha3(_ context.Context, in string) (string, error) {
	if err := ep.IsEndpointAllowed("web3_sha3"); err != nil {
		return "", err
	}
	in = strings.TrimPrefix(in, "0x")
	dec := make([]byte, hex.DecodedLen(len(in)))
	_, err := hex.Decode(dec, []byte(in))
	if err != nil {
		ep.Logger.Err(err).Msgf("could hex decode [%s]", in)
		return "", err
	}
	hash := crypto.Keccak256(dec)
	return "0x" + hex.EncodeToString(hash), nil
}

// A sample method to show the usage of single optional parameter
func (ep *Web3) Sha31(_ context.Context, arg1 *string) (string, error) {
	if err := ep.IsEndpointAllowed("web3_sha3"); err != nil {
		return "", err
	}

	in := "123456"
	if arg1 != nil {
		in = string(*arg1)
	}
	in = strings.TrimPrefix(in, "0x")
	dec := make([]byte, hex.DecodedLen(len(in)))
	_, err := hex.Decode(dec, []byte(in))
	if err != nil {
		ep.Logger.Err(err).Msgf("could hex decode [%s]", in)
		return "", err
	}
	hash := crypto.Keccak256(dec)
	return "0x" + hex.EncodeToString(hash), nil
}

// A sample method to show the usage of two mandatory arguments
func (ep *Web3) GetBlockByNumber(_ context.Context, param1 string, param2 bool) (string, error) {
	if err := ep.IsEndpointAllowed("web3_getBlockByNumber"); err != nil {
		return "", err
	}

	return fmt.Sprintf("First param is %s, and Second param is %t", param1, param2), nil
}

// A sample method to show the usage of single optional parameter
func (ep *Web3) GetBlockByNumber1(_ context.Context, param1 string, param2 *bool) (string, error) {
	if err := ep.IsEndpointAllowed("web3_getBlockByNumber"); err != nil {
		return "", err
	}
	block := param1
	hydratedTxs := false
	if param2 != nil {
		hydratedTxs = *param2
	}

	return fmt.Sprintf("First param is %s, and Second param is %t", block, hydratedTxs), nil
}

// A sample method to show the usage of two optional parameters
func (ep *Web3) GetBlockByNumber2(_ context.Context, param1, param2 *any) (string, error) {
	if err := ep.IsEndpointAllowed("web3_getBlockByNumber"); err != nil {
		return "", err
	}
	block := "LATEST"
	hydratedTxs := false
	boolCounter := 0
	if param1 != nil {
		paramFirst := *param1
		switch t1 := paramFirst.(type) {
		case string:
			block = t1
		case bool:
			hydratedTxs = t1
			boolCounter++
		default:
			return "", &utils.InvalidParamsError{Message: "invalid argument: incorrect first argument type"}
		}

		if param2 != nil && boolCounter == 0 {
			paramSecond := *param2
			switch t2 := paramSecond.(type) {
			case bool:
				hydratedTxs = t2
			default:
				return "", &utils.InvalidParamsError{Message: "invalid argument: incorrect second argument type"}
			}
		} else if param2 != nil && boolCounter > 0 {
			return "", &utils.InvalidParamsError{Message: "invalid argument: incorrect second argument type"}
		}
	}

	return fmt.Sprintf("First param is %s, and Second param is %t", block, hydratedTxs), nil
}
