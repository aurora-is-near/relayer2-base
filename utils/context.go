package utils

import (
	"sync/atomic"

	"golang.org/x/net/context"
)

var (
	defaultChainId          uint64 = 1313161554
	prehistoryChainId       uint64 = defaultChainId
	defaultPrehistoryHeight uint64 = 37157758
	prehistoryHeight        uint64 = defaultPrehistoryHeight
)

type chainIdKey struct{}

// GetChainId returns the chainId in the following order;
//  1. returns chainId if exists in ctx
//  2. returns default chainId (1313161554 unless manually set)
func GetChainId(ctx context.Context) uint64 {
	if ctx != nil {
		if cid, ok := ctx.Value(chainIdKey{}).(*uint64); ok && cid != nil {
			return *cid
		}
	}
	return atomic.LoadUint64(&defaultChainId)
}

// PutChainId is a helper function to put chainId in the context, also see GetChainId
func PutChainId(ctx context.Context, chainId uint64) context.Context {
	return context.WithValue(ctx, chainIdKey{}, &chainId)
}

// SetDefaultChainId sets the default chainId for GetChainId.
func SetDefaultChainId(val uint64) {
	atomic.StoreUint64(&defaultChainId, val)
}

// GetPrehistoryChainId returns the chainId (1313161554 by default) config that the prehistory was indexed for.
func GetPrehistoryChainId() uint64 {
	return atomic.LoadUint64(&prehistoryChainId)
}

// SetPrehistoryChainId sets the prehistory chainId for GetPrehistoryChainId.
func SetPrehistoryChainId(val uint64) {
	atomic.StoreUint64(&prehistoryChainId, val)
}

// GetPrehistoryHeight returns the height of the prehistory (37157758 for mainnet by default)
func GetPrehistoryHeight() uint64 {
	return atomic.LoadUint64(&prehistoryHeight)
}

// SetPrehistoryHeight sets the prehistory height for GetPrehistoryHeight.
func SetPrehistoryHeight(val uint64) {
	atomic.StoreUint64(&prehistoryHeight, val)
}
