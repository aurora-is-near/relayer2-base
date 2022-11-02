package response

import "aurora-relayer-go-common/types/primitives"

// https://openethereum.github.io/JSONRPC
type AccessListEntry struct {
	Address     primitives.Data20   `json:"address"`
	StorageKeys []primitives.Data32 `json:"storageKeys"`
}
