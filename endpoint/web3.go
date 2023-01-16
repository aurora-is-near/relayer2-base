package endpoint

import (
	errs "aurora-relayer-go-common/types/errors"
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

// ClientVersion returns client version
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
// 	TODO: implement
func (e *Web3) ClientVersion(_ context.Context) (*string, error) {
	return utils.Constants.ClientVersion(), nil
}

// Sha3 returns Keccak-256 hash of the given data.
//
// 	If API is disabled, returns errors code '-32601' with message 'the method does not exist/is not available'.
// 	On failure, returns errors code '-32000' with custom message.
func (e *Web3) Sha3(_ context.Context, in string) (*string, error) {
	in = strings.TrimPrefix(in, "0x")
	dec := make([]byte, hex.DecodedLen(len(in)))
	_, err := hex.Decode(dec, []byte(in))
	if err != nil {
		e.Logger.Err(err).Msgf("could hex decode [%s]", in)
		return nil, &errs.GenericError{Err: err}
	}
	keccak256 := crypto.Keccak256(dec)
	hash := fmt.Sprintf("0x%s", hex.EncodeToString(keccak256))
	return &hash, nil
}
