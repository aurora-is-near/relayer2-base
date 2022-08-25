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

func (e *EthPreprocessorAware) GetBlockTransactionCountByHash(ctx context.Context, hash utils.H256) (*utils.Uint256, error) {
	return Preprocess("eth_getBlockTransactionCountByHash", e.Endpoint, func() (*utils.Uint256, error) {
		return e.Eth.GetBlockTransactionCountByHash(ctx, hash)
	}, ctx, hash)
}

func (e *EthPreprocessorAware) GetBlockTransactionCountByNumber(ctx context.Context, number utils.Uint256) (*utils.Uint256, error) {
	return Preprocess("eth_getBlockTransactionCountByNumber", e.Endpoint, func() (*utils.Uint256, error) {
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

func (e *EthPreprocessorAware) GetLogs(ctx context.Context, rawFilter utils.FilterOptions) (*[]utils.LogResponse, error) {
	return Preprocess("eth_getLogs", e.Endpoint, func() (*[]utils.LogResponse, error) {
		return e.Eth.GetLogs(ctx, rawFilter)
	}, ctx, rawFilter)
}

// TODO NewFilter, NewBlockFilter, UninstallFilter, GetFilterChanges
