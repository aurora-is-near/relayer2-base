package response

import "github.com/aurora-is-near/relayer2-base/types/primitives"

type AuthorizationEntry struct {
	ChainID primitives.HexUint  `json:"chainId"`
	Address primitives.Data20   `json:"address"`
	Nonce   primitives.HexUint  `json:"nonce"`
	YParity primitives.HexUint  `json:"yParity"`
	R       primitives.Quantity `json:"r"`
	S       primitives.Quantity `json:"s"`
}
