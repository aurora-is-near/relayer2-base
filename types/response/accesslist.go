package response

import "github.com/aurora-is-near/relayer2-base/types/primitives"

// https://openethereum.github.io/JSONRPC
type AccessListEntry struct {
	Address     primitives.Data20   `json:"address"`
	StorageKeys []primitives.Data32 `json:"storageKeys"`
}
