package endpoint

import (
	"aurora-relayer-go-common/utils"
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	clientVersion = "Aurora Relayer"
)

type Web3 struct {
	*Endpoint
}

func NewWeb3(endpoint *Endpoint) *Web3 {
	return &Web3{endpoint}
}

// ClientVersion returns client version
//
// 	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
// 	TODO: implement
func (e *Web3) ClientVersion(_ context.Context) (*string, error) {
	return &clientVersion, nil
}

// Sha3 returns Keccak-256 hash of the given data.
//
// 	If API is disabled, returns error code '-32601' with message 'the method does not exist/is not available'.
// 	On failure, returns error code '-32000' with custom message.
func (e *Web3) Sha3(_ context.Context, in string) (*string, error) {
	in = strings.TrimPrefix(in, "0x")
	dec := make([]byte, hex.DecodedLen(len(in)))
	_, err := hex.Decode(dec, []byte(in))
	if err != nil {
		e.Logger.Err(err).Msgf("could hex decode [%s]", in)
		return nil, &utils.GenericError{Err: err}
	}
	keccak256 := crypto.Keccak256(dec)
	hash := fmt.Sprintf("0x%s", hex.EncodeToString(keccak256))
	return &hash, nil
}

// A sample method to show the usage of single optional parameter
// TODO: delete
func (e *Web3) Sha31(_ context.Context, arg1 *string) (string, error) {

	in := "123456"
	if arg1 != nil {
		in = string(*arg1)
	}
	in = strings.TrimPrefix(in, "0x")
	dec := make([]byte, hex.DecodedLen(len(in)))
	_, err := hex.Decode(dec, []byte(in))
	if err != nil {
		e.Logger.Err(err).Msgf("could hex decode [%s]", in)
		return "", err
	}
	hash := crypto.Keccak256(dec)
	return "0x" + hex.EncodeToString(hash), nil
}

// A sample method to show the usage of two mandatory arguments
// TODO: delete
func (e *Web3) GetBlockByNumber(_ context.Context, param1 string, param2 bool) (string, error) {
	return fmt.Sprintf("First param is %s, and Second param is %t", param1, param2), nil
}

// A sample method to show the usage of single optional parameter
// TODO: delete
func (e *Web3) GetBlockByNumber1(_ context.Context, param1 string, param2 *bool) (string, error) {

	block := param1
	hydratedTxs := false
	if param2 != nil {
		hydratedTxs = *param2
	}

	return fmt.Sprintf("First param is %s, and Second param is %t", block, hydratedTxs), nil
}

// A sample method to show the usage of two optional parameters
// TODO: delete
func (e *Web3) GetBlockByNumber2(_ context.Context, param1, param2 *any) (string, error) {

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
