package endpoint

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/aurora-is-near/relayer2-base/log"
	"github.com/aurora-is-near/relayer2-base/rpc"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

const (
	cacheMaxPriorityFeePerGas = 10 * time.Second
	maxPriorityFeePerGasBody  = `{"id":1,"jsonrpc":"2.0","method":"eth_maxPriorityFeePerGas"}`
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type AuroraClient interface {
	MaxPriorityFeePerGas() (*primitives.Quantity, error)
}

type RPCResponse struct {
	Result *primitives.Quantity `json:"result"`
}

type AuroraRPC struct {
	maxPriorityFeePerGasCache      *primitives.Quantity
	maxPriorityFeePerGasValidUntil time.Time
	url                            string
	mu                             sync.Mutex
}

func NewAuroraRPC(url string) *AuroraRPC {
	return &AuroraRPC{url: url}
}

func (a *AuroraRPC) MaxPriorityFeePerGas() (*primitives.Quantity, error) {
	now := time.Now()

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.maxPriorityFeePerGasCache != nil && a.maxPriorityFeePerGasValidUntil.Before(now) {
		return a.maxPriorityFeePerGasCache, nil
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(a.url)
	req.Header.SetContentType(rpc.DefaultContentType)
	req.Header.SetMethod(http.MethodPost)

	req.SetBody([]byte(maxPriorityFeePerGasBody))
	err := fasthttp.Do(req, resp)
	if err != nil {
		log.Log().Warn().Err(err).Msg("failed to query max priority fee per gas")
		return nil, errors.New("failed to query max priority fee per gas")
	}

	var val RPCResponse
	if err := json.Unmarshal(resp.Body(), &val); err != nil {
		log.Log().Warn().Err(err).Msg("failed to read max priority fee per gas")
		return nil, errors.New("failed to read max priority fee per gas")
	}

	a.maxPriorityFeePerGasCache = val.Result
	a.maxPriorityFeePerGasValidUntil = now.Add(cacheMaxPriorityFeePerGas)

	return a.maxPriorityFeePerGasCache, nil
}
