package endpoint

import (
	"context"
	"encoding/hex"
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

func (ep Web3) Sha3(_ context.Context, in *string) (string, error) {
	if err := ep.IsEndpointAllowed("web3_sha3"); err != nil {
		return "", err
	}
	*in = strings.TrimPrefix(*in, "0x")
	dec := make([]byte, hex.DecodedLen(len(*in)))
	_, err := hex.Decode(dec, []byte(*in))
	if err != nil {
		ep.Logger.Err(err).Msgf("could hex decode [%s]", in)
		return "", err
	}
	hash := crypto.Keccak256(dec)
	return "0x" + hex.EncodeToString(hash), nil
}
