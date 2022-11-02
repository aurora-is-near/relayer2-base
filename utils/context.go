package utils

import (
	"aurora-relayer-go-common/log"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

const (
	chainIdConfigPath = "endpoint.chainID"
)

var (
	defaultChainId uint64 = 1313161554
	chainId        *uint64
)

type chainIdKey struct{}

// GetChainId returns the chainId in the following order;
// 	1. returns chainId if exists in ctx
//	2. returns chainId if exists in relayer.yml
//	3. returns defaultChainId=1313161554
func GetChainId(ctx context.Context) uint64 {
	if cid, ok := ctx.Value(chainIdKey{}).(*uint64); ok && cid != nil {
		return *cid
	} else {
		if chainId == nil {
			chainId = getChainId()
		}
		return *chainId
	}
}

// PutChainId is a helper function to put chainId in the context, also see GetChainId
func PutChainId(ctx context.Context, chainId uint64) context.Context {
	return context.WithValue(ctx, chainIdKey{}, &chainId)
}

func getChainId() *uint64 {
	var cid *uint64
	sub := viper.GetUint64(chainIdConfigPath)
	if sub != 0 {
		cid = &sub
	} else {
		log.Log().Warn().Msgf("failed to parse configuration [%s] from [%s], "+
			"falling back to defaults", chainIdConfigPath, viper.ConfigFileUsed())
		cid = &defaultChainId
	}
	return cid
}
