package db

import (
	"github.com/aurora-is-near/relayer2-base/types/primitives"
)

type Log struct {
	Address primitives.Data20
	Data    primitives.VarData
	Topics  []primitives.Data32
}
