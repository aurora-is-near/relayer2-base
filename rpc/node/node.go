package node

import (
	"context"
	"strings"

	"github.com/aurora-is-near/relayer2-base/broker"
	"github.com/aurora-is-near/relayer2-base/log"
	"github.com/aurora-is-near/relayer2-base/rpc"
	"github.com/aurora-is-near/relayer2-base/rpc/node/events"
)

type RpcNode struct {
	rpc.RpcServer
	Broker broker.Broker
}

func New() (*RpcNode, error) {
	config := GetConfig()
	return NewWithConf(config)
}

func NewWithConf(config *Config) (*RpcNode, error) {
	logger := log.Log()
	transports := []rpc.TransportOption{}

	// If httpEndpoint is not empty, then a HttpServer should be initialized (no matters if it is http only or http and ws)
	if config.httpEndpoint() != "" {
		httpCfg := rpc.HttpConfig{
			HttpEndpoint:       config.httpEndpoint(),
			HttpCors:           config.HttpCors,
			HttpCompress:       config.HttpCompress,
			HttpTimeout:        config.HttpTimeout,
			WsEndpoint:         config.wsEndpoint(),
			WsHandshakeTimeout: config.WsHandshakeTimeout,
			WsOnly:             false,
		}
		transports = append(transports, rpc.WithTransport(&rpc.HttpServer{Config: httpCfg, Logger: logger}))
	}

	// If wsEndpoint is not empty and different from httpEndpoint, then another HttpServer should be initialized to
	// handle ws connections. Please note that WsOnly field is used to understand to case while handling the incoming request
	if config.wsEndpoint() != "" && !strings.EqualFold(config.httpEndpoint(), config.wsEndpoint()) {
		httpCfg := rpc.HttpConfig{
			HttpEndpoint:       config.wsEndpoint(),
			HttpCors:           []string{},
			HttpCompress:       false,
			HttpTimeout:        config.HttpTimeout,
			WsEndpoint:         config.wsEndpoint(),
			WsHandshakeTimeout: config.WsHandshakeTimeout,
			WsOnly:             true,
		}
		transports = append(transports, rpc.WithTransport(&rpc.HttpServer{Config: httpCfg, Logger: logger}))
	}
	if len(transports) == 0 {
		logger.Fatal().Msg("rpc server configuration error, no transport configured")
	}

	srv := rpc.New(logger, config.MaxBatchRequests, transports...)
	node := &RpcNode{RpcServer: *srv}

	// Start eventbroker if WS configured
	if config.wsEndpoint() != "" {
		eb := events.NewEventBroker()
		go eb.Start()
		node.Broker = eb
	}
	return node, nil
}

// Start starts RPC server as a seperate go routine
func (n *RpcNode) Start() {
	go func() {
		err := n.Run(context.Background())
		if err != nil {
			log.Log().Fatal().Err(err).Msg("can not start rpc server...")
		}
	}()
}
