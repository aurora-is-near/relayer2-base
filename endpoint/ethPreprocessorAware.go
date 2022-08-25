package endpoint

import (
	"aurora-relayer-go-common/utils"
	"context"
)

type EthPreprocessorAware struct {
	Eth
}

func NewEthPreprocessorAware(eth *Eth) *EthPreprocessorAware {
	return &EthPreprocessorAware{*eth}
}

func (e *EthPreprocessorAware) Accounts(ctx context.Context) ([]string, error) {
	return Preprocess("eth_account", e.Eth.Endpoint, func() ([]string, error) {
		return e.Eth.Accounts(ctx)
	}, ctx)
}

func (e *EthPreprocessorAware) Coinbase(ctx context.Context) (*utils.Uint256, error) {
	return Preprocess("eth_coinbase", e.Eth.Endpoint, func() (*utils.Uint256, error) {
		return e.Eth.Coinbase(ctx)
	}, ctx)
}

func (e *EthPreprocessorAware) ProtocolVersion(ctx context.Context) (*utils.Uint256, error) {
	return Preprocess("eth_protocolVersion", e.Eth.Endpoint, func() (*utils.Uint256, error) {
		return e.Eth.ProtocolVersion(ctx)
	}, ctx)
}

func (e *EthPreprocessorAware) Hashrate(ctx context.Context) (string, error) {
	return Preprocess("eth_hashrate", e.Eth.Endpoint, func() (string, error) {
		return e.Eth.Hashrate(ctx)
	}, ctx)
}

func (e *EthPreprocessorAware) BlockNumber(ctx context.Context) (*utils.Uint256, error) {
	return Preprocess("eth_blockNumber", e.Eth.Endpoint, func() (*utils.Uint256, error) {
		return e.Eth.BlockNumber(ctx)
	}, ctx)
}

func (e *EthPreprocessorAware) GetBlockByHash(ctx context.Context, hash utils.H256, isFull bool) (*utils.BlockResponse, error) {
	return Preprocess("eth_getBlockByHash", e.Eth.Endpoint, func() (*utils.BlockResponse, error) {
		return e.Eth.GetBlockByHash(ctx, hash, isFull)
	}, ctx, hash, isFull)
}

func (e *EthPreprocessorAware) GetBlockByNumber(ctx context.Context, number utils.Uint256, isFull bool) (*utils.BlockResponse, error) {
	return Preprocess("eth_getBlockByNumber", e.Eth.Endpoint, func() (*utils.BlockResponse, error) {
		return e.Eth.GetBlockByNumber(ctx, number, isFull)
	}, ctx, number, isFull)
}

func (e *EthPreprocessorAware) GetBlockTransactionCountByHash(ctx context.Context, hash utils.H256) (utils.Uint256, error) {
	return Preprocess("eth_getBlockTransactionCountByHash", e.Endpoint, func() (utils.Uint256, error) {
		return e.Eth.GetBlockTransactionCountByHash(ctx, hash)
	}, ctx, hash)
}

func (e *EthPreprocessorAware) GetBlockTransactionCountByNumber(ctx context.Context, number utils.Uint256) (utils.Uint256, error) {
	return Preprocess("eth_getBlockTransactionCountByNumber", e.Endpoint, func() (utils.Uint256, error) {
		return e.Eth.GetBlockTransactionCountByNumber(ctx, number)
	}, ctx, number)
}

func (e *EthPreprocessorAware) GetTransactionByHash(ctx context.Context, hash utils.H256) (*utils.TransactionResponse, error) {
	return Preprocess("eth_GetTransactionByHash", e.Endpoint, func() (*utils.TransactionResponse, error) {
		return e.Eth.GetTransactionByHash(ctx, hash)
	}, ctx, hash)
}

func (e *EthPreprocessorAware) GetTransactionByBlockHashAndIndex(ctx context.Context, hash utils.H256, index utils.Uint256) (*utils.TransactionResponse, error) {
	return Preprocess("eth_getTransactionByBlockHashAndIndex", e.Endpoint, func() (*utils.TransactionResponse, error) {
		return e.Eth.GetTransactionByBlockHashAndIndex(ctx, hash, index)
	}, ctx, hash, index)
}

func (e *EthPreprocessorAware) GetTransactionByBlockNumberAndIndex(ctx context.Context, number, index utils.Uint256) (*utils.TransactionResponse, error) {
	return Preprocess("eth_getTransactionByBlockNumberAndIndex", e.Endpoint, func() (*utils.TransactionResponse, error) {
		return e.Eth.GetTransactionByBlockNumberAndIndex(ctx, number, index)
	}, ctx, number, index)
}

func (e *EthPreprocessorAware) GetTransactionReceipt(ctx context.Context, hash utils.H256) (*utils.TransactionReceiptResponse, error) {
	return Preprocess("eth_getTransactionReceipt", e.Endpoint, func() (*utils.TransactionReceiptResponse, error) {
		return e.Eth.GetTransactionReceipt(ctx, hash)
	}, ctx, hash)
}

func (e *EthPreprocessorAware) GetLogs(ctx context.Context, rawFilter utils.FilterOptions) (*[]utils.LogResponse, error) {
	return Preprocess("eth_getLogs", e.Endpoint, func() (*[]utils.LogResponse, error) {
		return e.Eth.GetLogs(ctx, rawFilter)
	}, ctx, rawFilter)
}

func (e *EthPreprocessorAware) GetFilterLogs(ctx context.Context, filterId utils.Uint256) (*[]interface{}, error) {
	return Preprocess("eth_getFilterLogs", e.Endpoint, func() (*[]interface{}, error) {
		return e.Eth.GetFilterLogs(ctx, filterId)
	}, ctx, filterId)
}

func (e *EthPreprocessorAware) GetUncleCountByBlockHash(ctx context.Context, hash utils.H256) (*utils.Uint256, error) {
	return Preprocess("eth_GetUncleCountByBlockHash", e.Endpoint, func() (*utils.Uint256, error) {
		return e.Eth.GetUncleCountByBlockHash(ctx, hash)
	}, ctx, hash)
}

func (e *EthPreprocessorAware) GetUncleCountByBlockNumber(ctx context.Context, number utils.Uint256) (*utils.Uint256, error) {
	return Preprocess("eth_GetUncleCountByBlockNumber", e.Endpoint, func() (*utils.Uint256, error) {
		return e.Eth.GetUncleCountByBlockNumber(ctx, number)
	}, ctx, number)
}

func (e *EthPreprocessorAware) NewFilter(ctx context.Context, filterOptions utils.FilterOptions) (*utils.Uint256, error) {
	return Preprocess("eth_NewFilter", e.Endpoint, func() (*utils.Uint256, error) {
		return e.Eth.NewFilter(ctx, filterOptions)
	}, ctx, filterOptions)
}

func (e *EthPreprocessorAware) NewBlockFilter(ctx context.Context) (*utils.Uint256, error) {
	return Preprocess("eth_NewBlockFilter", e.Endpoint, func() (*utils.Uint256, error) {
		return e.Eth.NewBlockFilter(ctx)
	}, ctx)
}

func (e *EthPreprocessorAware) UninstallFilter(ctx context.Context, filterId utils.Uint256) (bool, error) {
	return Preprocess("eth_UninstallFilter", e.Endpoint, func() (bool, error) {
		return e.Eth.UninstallFilter(ctx, filterId)
	}, ctx, filterId)
}

func (e *EthPreprocessorAware) GetFilterChanges(ctx context.Context, filterId utils.Uint256) (*[]interface{}, error) {
	return Preprocess("eth_GetFilterChanges", e.Endpoint, func() (*[]interface{}, error) {
		return e.Eth.GetFilterChanges(ctx, filterId)
	}, ctx, filterId)
}
