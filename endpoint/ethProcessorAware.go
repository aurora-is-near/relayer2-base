package endpoint

import (
	"aurora-relayer-go-common/types/common"
	"aurora-relayer-go-common/types/engine"
	"aurora-relayer-go-common/types/primitives"
	"aurora-relayer-go-common/types/request"
	"aurora-relayer-go-common/types/response"
	"context"
)

type EthProcessorAware struct {
	*Eth
}

func NewEthProcessorAware(eth *Eth) *EthProcessorAware {
	return &EthProcessorAware{eth}
}

func (e *EthProcessorAware) Accounts(ctx context.Context) (*[]string, error) {
	return Process(ctx, "eth_account", e.Endpoint, func(ctx context.Context) (*[]string, error) {
		return e.Eth.Accounts(ctx)
	})
}

func (e *EthProcessorAware) Coinbase(ctx context.Context) (*string, error) {
	return Process(ctx, "eth_coinbase", e.Endpoint, func(ctx context.Context) (*string, error) {
		return e.Eth.Coinbase(ctx)
	})
}

func (e *EthProcessorAware) ProtocolVersion(ctx context.Context) (*common.Uint256, error) {
	return Process(ctx, "eth_protocolVersion", e.Endpoint, func(ctx context.Context) (*common.Uint256, error) {
		return e.Eth.ProtocolVersion(ctx)
	})
}

func (e *EthProcessorAware) Hashrate(ctx context.Context) (*common.Uint256, error) {
	return Process(ctx, "eth_hashrate", e.Endpoint, func(ctx context.Context) (*common.Uint256, error) {
		return e.Eth.Hashrate(ctx)
	})
}

func (e *EthProcessorAware) Mining(ctx context.Context) (*bool, error) {
	return Process(ctx, "eth_mining", e.Endpoint, func(ctx context.Context) (*bool, error) {
		return e.Eth.Mining(ctx)
	})
}

func (e *EthProcessorAware) Syncing(ctx context.Context) (*bool, error) {
	return Process(ctx, "eth_syncing", e.Endpoint, func(ctx context.Context) (*bool, error) {
		return e.Eth.Syncing(ctx)
	})
}

func (e *EthProcessorAware) BlockNumber(ctx context.Context) (*primitives.HexUint, error) {
	return Process(ctx, "eth_blockNumber", e.Endpoint, func(ctx context.Context) (*primitives.HexUint, error) {
		return e.Eth.BlockNumber(ctx)
	})
}

func (e *EthProcessorAware) GetBlockByHash(ctx context.Context, hash common.H256, isFull *bool) (*response.Block, error) {
	return Process(ctx, "eth_getBlockByHash", e.Endpoint, func(ctx context.Context) (*response.Block, error) {
		return e.Eth.GetBlockByHash(ctx, hash, isFull)
	}, hash, isFull)
}

func (e *EthProcessorAware) GetBlockByNumber(ctx context.Context, number common.BN64, isFull *bool) (*response.Block, error) {
	return Process(ctx, "eth_getBlockByNumber", e.Endpoint, func(ctx context.Context) (*response.Block, error) {
		return e.Eth.GetBlockByNumber(ctx, number, isFull)
	}, number, isFull)
}

func (e *EthProcessorAware) GetBlockTransactionCountByHash(ctx context.Context, hash common.H256) (*primitives.HexUint, error) {
	return Process(ctx, "eth_getBlockTransactionCountByHash", e.Endpoint, func(ctx context.Context) (*primitives.HexUint, error) {
		return e.Eth.GetBlockTransactionCountByHash(ctx, hash)
	}, hash)
}

func (e *EthProcessorAware) GetBlockTransactionCountByNumber(ctx context.Context, number *common.BN64) (*primitives.HexUint, error) {
	return Process(ctx, "eth_getBlockTransactionCountByNumber", e.Endpoint, func(ctx context.Context) (*primitives.HexUint, error) {
		return e.Eth.GetBlockTransactionCountByNumber(ctx, number)
	}, number)
}

func (e *EthProcessorAware) GetTransactionByHash(ctx context.Context, hash common.H256) (*response.Transaction, error) {
	return Process(ctx, "eth_getTransactionByHash", e.Endpoint, func(ctx context.Context) (*response.Transaction, error) {
		return e.Eth.GetTransactionByHash(ctx, hash)
	}, hash)
}

func (e *EthProcessorAware) GetTransactionByBlockHashAndIndex(ctx context.Context, hash common.H256, index common.Uint64) (*response.Transaction, error) {
	return Process(ctx, "eth_getTransactionByBlockHashAndIndex", e.Endpoint, func(ctx context.Context) (*response.Transaction, error) {
		return e.Eth.GetTransactionByBlockHashAndIndex(ctx, hash, index)
	}, hash, index)
}

func (e *EthProcessorAware) GetTransactionByBlockNumberAndIndex(ctx context.Context, number common.BN64, index common.Uint64) (*response.Transaction, error) {
	return Process(ctx, "eth_getTransactionByBlockNumberAndIndex", e.Endpoint, func(ctx context.Context) (*response.Transaction, error) {
		return e.Eth.GetTransactionByBlockNumberAndIndex(ctx, number, index)
	}, number, index)
}

func (e *EthProcessorAware) GetTransactionReceipt(ctx context.Context, hash common.H256) (*response.TransactionReceipt, error) {
	return Process(ctx, "eth_getTransactionReceipt", e.Endpoint, func(ctx context.Context) (*response.TransactionReceipt, error) {
		return e.Eth.GetTransactionReceipt(ctx, hash)
	}, hash)
}

func (e *EthProcessorAware) GetLogs(ctx context.Context, rawFilter request.Filter) (*[]*response.Log, error) {
	return Process(ctx, "eth_getLogs", e.Endpoint, func(ctx context.Context) (*[]*response.Log, error) {
		return e.Eth.GetLogs(ctx, rawFilter)
	}, rawFilter)
}

func (e *EthProcessorAware) GetFilterLogs(ctx context.Context, filterId common.Uint256) (*[]*response.Log, error) {
	return Process(ctx, "eth_getFilterLogs", e.Endpoint, func(ctx context.Context) (*[]*response.Log, error) {
		return e.Eth.GetFilterLogs(ctx, filterId)
	}, filterId)
}

func (e *EthProcessorAware) GetUncleCountByBlockHash(ctx context.Context, hash common.H256) (*common.Uint256, error) {
	return Process(ctx, "eth_getUncleCountByBlockHash", e.Endpoint, func(ctx context.Context) (*common.Uint256, error) {
		return e.Eth.GetUncleCountByBlockHash(ctx, hash)
	}, hash)
}

func (e *EthProcessorAware) GetUncleCountByBlockNumber(ctx context.Context, number *common.BN64) (*common.Uint256, error) {
	return Process(ctx, "eth_getUncleCountByBlockNumber", e.Endpoint, func(ctx context.Context) (*common.Uint256, error) {
		return e.Eth.GetUncleCountByBlockNumber(ctx, number)
	}, number)
}

func (e *EthProcessorAware) GetUncleByBlockHashAndIndex(ctx context.Context, hash *common.H256, index *common.Uint64) (*string, error) {
	return Process(ctx, "eth_getUncleCountByHashAndIndex", e.Endpoint, func(ctx context.Context) (*string, error) {
		return e.Eth.GetUncleByBlockHashAndIndex(ctx, hash, index)
	}, hash, index)
}

func (e *EthProcessorAware) GetUncleByBlockNumberAndIndex(ctx context.Context, number *common.BN64, index *common.Uint64) (*string, error) {
	return Process(ctx, "eth_getUncleByBlockNumberAndIndex", e.Endpoint, func(ctx context.Context) (*string, error) {
		return e.Eth.GetUncleByBlockNumberAndIndex(ctx, number, index)
	}, number, index)
}

func (e *EthProcessorAware) NewFilter(ctx context.Context, filter request.Filter) (*common.Uint256, error) {
	return Process(ctx, "eth_newFilter", e.Endpoint, func(ctx context.Context) (*common.Uint256, error) {
		return e.Eth.NewFilter(ctx, filter)
	}, filter)
}

func (e *EthProcessorAware) NewBlockFilter(ctx context.Context) (*common.Uint256, error) {
	return Process(ctx, "eth_newBlockFilter", e.Endpoint, func(ctx context.Context) (*common.Uint256, error) {
		return e.Eth.NewBlockFilter(ctx)
	})
}

func (e *EthProcessorAware) NewPendingTransactionFilter(ctx context.Context) (*string, error) {
	return Process(ctx, "eth_newPendingTransactionFilter", e.Endpoint, func(ctx context.Context) (*string, error) {
		return e.Eth.NewPendingTransactionFilter(ctx)
	})
}

func (e *EthProcessorAware) UninstallFilter(ctx context.Context, filterId common.Uint256) (*bool, error) {
	return Process(ctx, "eth_uninstallFilter", e.Endpoint, func(ctx context.Context) (*bool, error) {
		return e.Eth.UninstallFilter(ctx, filterId)
	}, filterId)
}

func (e *EthProcessorAware) GetFilterChanges(ctx context.Context, filterId common.Uint256) (*[]interface{}, error) {
	return Process(ctx, "eth_getFilterChanges", e.Endpoint, func(ctx context.Context) (*[]interface{}, error) {
		return e.Eth.GetFilterChanges(ctx, filterId)
	}, filterId)
}

func (e *EthProcessorAware) GetCompilers(ctx context.Context) (*[]string, error) {
	return Process(ctx, "eth_getCompilers", e.Endpoint, func(ctx context.Context) (*[]string, error) {
		return e.Eth.GetCompilers(ctx)
	})
}

func (e *EthProcessorAware) PendingTransactions(ctx context.Context) (*[]string, error) {
	return Process(ctx, "eth_pendingTransactions", e.Endpoint, func(ctx context.Context) (*[]string, error) {
		return e.Eth.PendingTransactions(ctx)
	})
}

func (e *EthProcessorAware) EstimateGas(ctx context.Context, txs engine.TransactionForCall, number *common.BN64) (*common.Uint256, error) {
	return Process(ctx, "eth_estimateGas", e.Endpoint, func(ctx context.Context) (*common.Uint256, error) {
		return e.Eth.EstimateGas(ctx, txs, number)
	}, txs, number)
}

func (e *EthProcessorAware) GasPrice(ctx context.Context) (*common.Uint256, error) {
	return Process(ctx, "eth_gasPrice", e.Endpoint, func(ctx context.Context) (*common.Uint256, error) {
		return e.Eth.GasPrice(ctx)
	})
}
