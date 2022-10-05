package endpoint

import (
	"aurora-relayer-go-common/utils"
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

func (e *EthProcessorAware) ChainId(ctx context.Context) (*utils.Uint256, error) {
	return Process(ctx, "eth_chainId", e.Endpoint, func(ctx context.Context) (*utils.Uint256, error) {
		return e.Eth.ChainId(ctx)
	})
}

func (e *EthProcessorAware) ProtocolVersion(ctx context.Context) (*utils.Uint256, error) {
	return Process(ctx, "eth_protocolVersion", e.Endpoint, func(ctx context.Context) (*utils.Uint256, error) {
		return e.Eth.ProtocolVersion(ctx)
	})
}

func (e *EthProcessorAware) Hashrate(ctx context.Context) (*utils.Uint256, error) {
	return Process(ctx, "eth_hashrate", e.Endpoint, func(ctx context.Context) (*utils.Uint256, error) {
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

func (e *EthProcessorAware) BlockNumber(ctx context.Context) (*utils.Uint256, error) {
	return Process(ctx, "eth_blockNumber", e.Endpoint, func(ctx context.Context) (*utils.Uint256, error) {
		return e.Eth.BlockNumber(ctx)
	})
}

func (e *EthProcessorAware) GetBlockByHash(ctx context.Context, hash utils.H256, isFull bool) (*utils.BlockResponse, error) {
	return Process(ctx, "eth_getBlockByHash", e.Endpoint, func(ctx context.Context) (*utils.BlockResponse, error) {
		return e.Eth.GetBlockByHash(ctx, hash, isFull)
	}, hash, isFull)
}

func (e *EthProcessorAware) GetBlockByNumber(ctx context.Context, number utils.Uint256, isFull bool) (*utils.BlockResponse, error) {
	return Process(ctx, "eth_getBlockByNumber", e.Endpoint, func(ctx context.Context) (*utils.BlockResponse, error) {
		return e.Eth.GetBlockByNumber(ctx, number, isFull)
	}, number, isFull)
}

func (e *EthProcessorAware) GetBlockTransactionCountByHash(ctx context.Context, hash utils.H256) (*utils.Uint256, error) {
	return Process(ctx, "eth_getBlockTransactionCountByHash", e.Endpoint, func(ctx context.Context) (*utils.Uint256, error) {
		return e.Eth.GetBlockTransactionCountByHash(ctx, hash)
	}, hash)
}

func (e *EthProcessorAware) GetBlockTransactionCountByNumber(ctx context.Context, number utils.Uint256) (*utils.Uint256, error) {
	return Process(ctx, "eth_getBlockTransactionCountByNumber", e.Endpoint, func(ctx context.Context) (*utils.Uint256, error) {
		return e.Eth.GetBlockTransactionCountByNumber(ctx, number)
	}, number)
}

func (e *EthProcessorAware) GetTransactionByHash(ctx context.Context, hash utils.H256) (*utils.TransactionResponse, error) {
	return Process(ctx, "eth_getTransactionByHash", e.Endpoint, func(ctx context.Context) (*utils.TransactionResponse, error) {
		return e.Eth.GetTransactionByHash(ctx, hash)
	}, hash)
}

func (e *EthProcessorAware) GetTransactionByBlockHashAndIndex(ctx context.Context, hash utils.H256, index utils.Uint256) (*utils.TransactionResponse, error) {
	return Process(ctx, "eth_getTransactionByBlockHashAndIndex", e.Endpoint, func(ctx context.Context) (*utils.TransactionResponse, error) {
		return e.Eth.GetTransactionByBlockHashAndIndex(ctx, hash, index)
	}, hash, index)
}

func (e *EthProcessorAware) GetTransactionByBlockNumberAndIndex(ctx context.Context, number, index utils.Uint256) (*utils.TransactionResponse, error) {
	return Process(ctx, "eth_getTransactionByBlockNumberAndIndex", e.Endpoint, func(ctx context.Context) (*utils.TransactionResponse, error) {
		return e.Eth.GetTransactionByBlockNumberAndIndex(ctx, number, index)
	}, number, index)
}

func (e *EthProcessorAware) GetTransactionReceipt(ctx context.Context, hash utils.H256) (*utils.TransactionReceiptResponse, error) {
	return Process(ctx, "eth_getTransactionReceipt", e.Endpoint, func(ctx context.Context) (*utils.TransactionReceiptResponse, error) {
		return e.Eth.GetTransactionReceipt(ctx, hash)
	}, hash)
}

func (e *EthProcessorAware) GetLogs(ctx context.Context, rawFilter *utils.FilterOptions) (*[]utils.LogResponse, error) {
	return Process(ctx, "eth_getLogs", e.Endpoint, func(ctx context.Context) (*[]utils.LogResponse, error) {
		return e.Eth.GetLogs(ctx, rawFilter)
	}, rawFilter)
}

func (e *EthProcessorAware) GetFilterLogs(ctx context.Context, filterId utils.Uint256) (*[]interface{}, error) {
	return Process(ctx, "eth_getFilterLogs", e.Endpoint, func(ctx context.Context) (*[]interface{}, error) {
		return e.Eth.GetFilterLogs(ctx, filterId)
	}, filterId)
}

func (e *EthProcessorAware) GetUncleCountByBlockHash(ctx context.Context, hash utils.H256) (*utils.Uint256, error) {
	return Process(ctx, "eth_getUncleCountByBlockHash", e.Endpoint, func(ctx context.Context) (*utils.Uint256, error) {
		return e.Eth.GetUncleCountByBlockHash(ctx, hash)
	}, hash)
}

func (e *EthProcessorAware) GetUncleCountByBlockNumber(ctx context.Context, number utils.Uint256) (*utils.Uint256, error) {
	return Process(ctx, "eth_getUncleCountByBlockNumber", e.Endpoint, func(ctx context.Context) (*utils.Uint256, error) {
		return e.Eth.GetUncleCountByBlockNumber(ctx, number)
	}, number)
}

func (e *EthProcessorAware) NewFilter(ctx context.Context, filterOptions *utils.FilterOptions) (*utils.Uint256, error) {
	return Process(ctx, "eth_newFilter", e.Endpoint, func(ctx context.Context) (*utils.Uint256, error) {
		return e.Eth.NewFilter(ctx, filterOptions)
	}, filterOptions)
}

func (e *EthProcessorAware) NewBlockFilter(ctx context.Context) (*utils.Uint256, error) {
	return Process(ctx, "eth_newBlockFilter", e.Endpoint, func(ctx context.Context) (*utils.Uint256, error) {
		return e.Eth.NewBlockFilter(ctx)
	})
}

func (e *EthProcessorAware) UninstallFilter(ctx context.Context, filterId utils.Uint256) (*bool, error) {
	return Process(ctx, "eth_uninstallFilter", e.Endpoint, func(ctx context.Context) (*bool, error) {
		return e.Eth.UninstallFilter(ctx, filterId)
	}, filterId)
}

func (e *EthProcessorAware) GetFilterChanges(ctx context.Context, filterId utils.Uint256) (*[]interface{}, error) {
	return Process(ctx, "eth_getFilterChanges", e.Endpoint, func(ctx context.Context) (*[]interface{}, error) {
		return e.Eth.GetFilterChanges(ctx, filterId)
	}, filterId)
}
