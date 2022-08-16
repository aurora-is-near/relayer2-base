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
	"aurora-relayer-go-common/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
)

// GoEthereum is a container on which underlying go-ethereum services can be registered.
type GoEthereum struct {
	node.Node
}

const (
	configPath = "RpcNode.GoEthereum"
)

// New creates a new node with default config
func New() (*GoEthereum, error) {
	logger := log.New()
	conf := DefaultConfig()
	conf.Logger = NewGoEthLogger(logger)
	sub := viper.Sub(configPath)
	if sub != nil {
		if err := sub.Unmarshal(&conf); err != nil {
			logger.Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
				"falling back to defaults", configPath, viper.ConfigFileUsed())
		}
	}
	return NewWithConf(conf)
}

// NewWithConf creates a new node with given config
func NewWithConf(conf *Config) (*GoEthereum, error) {
	ethConf := convertConfigurationToEthNode(conf)
	n, err := node.New(ethConf)
	if err != nil {
		return nil, err
	}
	return &GoEthereum{*n}, nil
}

// // Start starts Eth Node.
// func (n *GoEthereum) Start() error {
// 	return n.Start()
// }

// // Close stops the Eth Node and clean up the underlying resources
// func (n *GoEthereum) Close() error {
// 	return n.Close()
// }

// // RegisterEndpointAPIs takes the implemented RPC methods, converts them to go-ethereum-rpc-API type and registers them
// func (n *GoEthereum) RegisterEndpointAPIs(apis []endpoint.API) {
// 	a := make([]rpc.API, len(apis))
// 	for i, e := range apis {
// 		a[i].Namespace = e.Namespace
// 		a[i].Service = e.EndpointInstance
// 		a[i].Authenticated = false
// 	}
// 	n.RegisterAPIs(a)
// }

func convertConfigurationToEthNode(confAurora *Config) *node.Config {
	confEth := &node.Config{}
	copier.CopyWithOption(confEth, confAurora, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	return confEth
}
