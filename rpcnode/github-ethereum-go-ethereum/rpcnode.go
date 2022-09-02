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
	"github.com/ethereum/go-ethereum/node"
	"github.com/jinzhu/copier"
)

// GoEthereum is a container on which underlying go-ethereum services can be registered.
type GoEthereum struct {
	node.Node
}

// New creates a new node with default config
func New() (*GoEthereum, error) {
	conf := GetConfig()
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

func convertConfigurationToEthNode(confAurora *Config) *node.Config {
	confEth := &node.Config{}
	copier.CopyWithOption(confEth, confAurora, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	return confEth
}
