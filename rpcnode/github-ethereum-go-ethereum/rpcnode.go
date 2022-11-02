// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package github_ethereum_go_ethereum

import (
	"aurora-relayer-go-common/broker"
	eventbroker "aurora-relayer-go-common/rpcnode/github-ethereum-go-ethereum/events"
	"github.com/ethereum/go-ethereum/node"
	"github.com/jinzhu/copier"
	"net/http"
)

// GoEthereum is a container on which underlying go-ethereum services can be registered.
type GoEthereum struct {
	node.Node
	Broker broker.Broker
}

// New creates a new node with default config
func New() (*GoEthereum, error) {
	conf := GetConfig()
	return NewWithConf(conf)
}

// NewWithConf creates a new node with given config and the types broker if node supports websocket comm
func NewWithConf(conf *Config) (*GoEthereum, error) {
	ethConf := convertConfigurationToEthNode(conf)
	n, err := node.New(ethConf)
	if err != nil {
		return nil, err
	}

	// Start eventbroker if WS configured
	eb := eventbroker.NewEventBroker()
	if conf.WSHost != "" && conf.WSPort > 0 {
		go eb.Start()
	}

	return &GoEthereum{
		Node:   *n,
		Broker: eb,
	}, nil
}

// WithMiddleware places the given middleware at the beginning of HTTP handlers chain. A middleware is a function accepting
// http.Handler and returning http.Handler. Any request matching the path argument is first processed by this middleware.
// Middleware should either return response for the request or pass the request to next HTTP handler in the chain.
// Refer to http.HandlerFunc for middleware creation.
func (ge GoEthereum) WithMiddleware(name string, path string, middleware func(handler http.Handler) http.Handler) {
	h, err := ge.RPCHandler()
	if err != nil {
		panic(err)
	}
	ge.RegisterHandler(name, path, middleware(h))
}

func convertConfigurationToEthNode(confAurora *Config) *node.Config {
	confEth := &node.Config{}
	copier.CopyWithOption(confEth, confAurora, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	return confEth
}
