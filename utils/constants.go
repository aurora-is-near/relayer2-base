package utils

import "github.com/aurora-is-near/relayer2-base/types/common"

type constants struct {
	clientVersion  string
	zeroStrUint256 string
	zeroStrUint160 string
	zeroStrUint128 string
	response0x     string
	syncing        bool
	mining         bool
	full           bool
	gasLimit       uint64
	emptyArray     []string
	zeroUint256    common.Uint256
}

func (c *constants) ClientVersion() *string {
	return &c.clientVersion
}

func (c *constants) ZeroStrUint256() *string {
	return &c.zeroStrUint256
}

func (c *constants) ZeroStrUint160() *string {
	return &c.zeroStrUint160
}

func (c constants) ZeroStrUint128() *string {
	return &c.zeroStrUint128
}

func (c *constants) Response0x() *string {
	return &c.response0x
}

func (c *constants) Syncing() *bool {
	return &c.syncing
}

func (c *constants) Mining() *bool {
	return &c.mining
}

func (c *constants) Full() *bool {
	return &c.full
}

func (c *constants) GasLimit() *uint64 {
	return &c.gasLimit
}

func (c *constants) EmptyArray() *[]string {
	return &c.emptyArray
}

func (c *constants) ZeroUint256() *common.Uint256 {
	return &c.zeroUint256
}

var Constants constants

func init() {
	Constants.clientVersion = "Aurora"
	Constants.zeroStrUint256 = "0x0000000000000000000000000000000000000000000000000000000000000000"
	Constants.zeroStrUint160 = "0x000000000000000000000000000000000000"
	Constants.zeroStrUint128 = "0x00000000000000000000000000000000"
	Constants.response0x = "0x"
	Constants.syncing = false
	Constants.mining = false
	Constants.full = false
	Constants.gasLimit = 9007199254740991 // hex value 0x1fffffffffffff
	Constants.emptyArray = []string{}
	Constants.zeroUint256 = common.IntToUint256(0)
}
