package utils

import (
	"net"

	"github.com/aurora-is-near/relayer2-base/log"

	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

const (
	chainIdConfigPath           = "endpoint.chainID"
	prehistoryChainIdConfigPath = "prehistoryIndexer.prehistoryChainID"
	prehsitoryHeightConfigPath  = "prehistoryIndexer.prehistoryHeight"
)

var (
	defaultChainId          uint64 = 1313161554
	chainId                 *uint64
	prehistoryChainId       *uint64
	defaultPrehistoryHeight uint64 = 37157758
	prehistoryHeight        *uint64
)

type chainIdKey struct{}

// GetChainId returns the chainId in the following order;
// 	1. returns chainId if exists in ctx
//	2. returns chainId if exists in relayer configuration .yml file
//	3. returns defaultChainId=1313161554
func GetChainId(ctx context.Context) uint64 {
	if ctx != nil {
		if cid, ok := ctx.Value(chainIdKey{}).(*uint64); ok && cid != nil {
			return *cid
		}
	}

	if chainId == nil {
		chainId = getChainId()
	}
	return *chainId
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

// GetPrehistoryChainId returns the chainId config that the prehistory was indexed for;
//	1. returns prehistoryChainId if exists in relayer configuration .yml file
//	2. returns defaultChainId=1313161554
func GetPrehistoryChainId() uint64 {
	if prehistoryChainId == nil {
		prehistoryChainId = getPrehistoryChainId()
	}
	return *prehistoryChainId
}

func getPrehistoryChainId() *uint64 {
	var pcid *uint64
	sub := viper.GetUint64(prehistoryChainIdConfigPath)
	if sub != 0 {
		pcid = &sub
	} else {
		log.Log().Warn().Msgf("failed to parse configuration [%s] from [%s], "+
			"falling back to defaults", prehistoryChainIdConfigPath, viper.ConfigFileUsed())
		pcid = &defaultChainId
	}
	return pcid
}

// GetPrehistoryHeight returns the height of the prehistory given in relayer configuration .yml file
//	1. returns prehistoryHeight if exists in relayer configuration .yml file
//	2. returns defaultPrehistoryHeight=37157758 that is for mainnet
func GetPrehistoryHeight() uint64 {
	if prehistoryHeight == nil {
		prehistoryHeight = getPrehistoryHeight()
	}
	return *prehistoryHeight
}

func getPrehistoryHeight() *uint64 {
	var ph *uint64
	sub := viper.GetUint64(prehsitoryHeightConfigPath)
	if sub != 0 {
		ph = &sub
	} else {
		log.Log().Warn().Msgf("failed to parse configuration [%s] from [%s], "+
			"falling back to defaults", prehsitoryHeightConfigPath, viper.ConfigFileUsed())
		ph = &defaultPrehistoryHeight
	}
	return ph
}

type clientIpKey struct{}

// PutClientIpKey is a helper function to put clientIp in the context so that handlers can use it
func PutClientIpKey(ctx context.Context, ip net.IP) context.Context {
	return context.WithValue(ctx, clientIpKey{}, ip)
}

// ClientIpFromContext returns the clientIp value stored in ctx, if any.
func ClientIpFromContext(ctx context.Context) (*net.IP, bool) {
	ip, ok := ctx.Value(clientIpKey{}).(net.IP)
	return &ip, ok
}
