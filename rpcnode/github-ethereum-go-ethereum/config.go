// Copyright 2014 The go-ethereum Authors
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
	"relayer2-base/cmd"
	"relayer2-base/log"
	"time"

	gel "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
)

const (
	DefaultHost       = "localhost" // Default host interface for the HTTP RPC server
	DefaultHTTPPort   = 8545        // Default TCP port for the HTTP RPC server
	DefaultWSPort     = 8546        // Default TCP port for the websocket RPC server
	DefaultPathPrefix = ""          // Default TCP port for the websocket RPC server

	configPath = "rpcNode.geth"
)

// defaultHTTPTimeouts represents the default timeout values used by the RPC server if further
// configuration is not provided.
var defaultHTTPTimeouts = rpc.HTTPTimeouts{
	ReadTimeout:       300 * time.Second,
	ReadHeaderTimeout: 120 * time.Second,
	WriteTimeout:      120 * time.Second,
	IdleTimeout:       120 * time.Second,
}

// Config represents a small collection of configuration values to fine tune the
// P2P network layer of a protocol stack. These values can be further extended by
// all registered services.
type Config struct {

	// HTTPHost is the host interface on which to start the HTTP RPC server. If this
	// field is empty, no HTTP API endpoint will be started.
	HTTPHost string

	// HTTPPort is the TCP port number on which to start the HTTP RPC server. The
	// default zero value is/ valid and will pick a port number randomly (useful
	// for ephemeral nodes).
	HTTPPort int `mapstructure:",omitempty"`

	// HTTPCors is the Cross-Origin Resource Sharing header to send to requesting
	// clients. Please be aware that CORS is a browser enforced security, it's fully
	// useless for custom HTTP clients.
	HTTPCors []string `mapstructure:",omitempty"`

	// HTTPVirtualHosts is the list of virtual hostnames which are allowed on incoming requests.
	// This is by default {'localhost'}. Using this prevents attacks like
	// DNS rebinding, which bypasses SOP by simply masquerading as being within the same
	// origin. These attacks do not utilize CORS, since they are not cross-domain.
	// By explicitly checking the Host-header, the server will not allow requests
	// made against the server with a malicious host domain.
	// Requests using ip address directly are not affected
	HTTPVirtualHosts []string `mapstructure:",omitempty"`

	// HTTPModules is a list of API modules to expose via the HTTP RPC interface.
	// If the module list is empty, all RPC API endpoints designated public will be
	// exposed.
	HTTPModules []string

	// HTTPTimeouts allows for customization of the timeout values used by the HTTP RPC
	// interface.
	HTTPTimeouts rpc.HTTPTimeouts

	// HTTPPathPrefix specifies a path prefix on which http-rpc is to be served.
	HTTPPathPrefix string `mapstructure:",omitempty"`

	// WSHost is the host interface on which to start the websocket RPC server. If
	// this field is empty, no websocket API endpoint will be started.
	WSHost string

	// WSPort is the TCP port number on which to start the websocket RPC server. The
	// default zero value is/ valid and will pick a port number randomly (useful for
	// ephemeral nodes).
	WSPort int `mapstructure:",omitempty"`

	// WSPathPrefix specifies a path prefix on which ws-rpc is to be served.
	WSPathPrefix string `mapstructure:",omitempty"`

	// WSOrigins is the list of domain to accept websocket requests from. Please be
	// aware that the server can only act upon the HTTP request the client sends and
	// cannot verify the validity of the request header.
	WSOrigins []string `mapstructure:",omitempty"`

	// WSModules is a list of API modules to expose via the websocket RPC interface.
	// If the module list is empty, all RPC API endpoints designated public will be
	// exposed.
	WSModules []string

	// Logger is a custom log to use with the p2p.Server.
	Logger gel.Logger `mapstructure:",omitempty"`
}

// defaultConfig is a helper to initialize Go-Ethereum RPC node with following defaults
// HTTPHost: DefaultHost
// HTTPPort: DefaultHTTPPort
// HTTPPathPrefix: DefaultPathPrefix
// HTTPCors: []
// HTTPModules: ["net", "web3", "eth"]
// HTTPVirtualHosts: []
// HTTPTimeouts: defaultHTTPTimeouts
// -> WS parameters are optional. Add the following WS parameters to make them mandatory.
// WSHost: DefaultHost
// WSPort: DefaultWSPort
// WSModules: ["net", "web3", "eth"]
// WSPathPrefix: DefaultPathPrefix
// WSOrigins: []
func defaultConfig() *Config {
	return &Config{
		HTTPHost:         DefaultHost,
		HTTPPort:         DefaultHTTPPort,
		HTTPPathPrefix:   DefaultPathPrefix,
		HTTPModules:      []string{"net", "web3", "eth", "parity"},
		HTTPVirtualHosts: []string{},
		HTTPCors:         []string{},
		HTTPTimeouts:     defaultHTTPTimeouts,
		Logger:           NewGoEthLogger(log.Log()),
	}
}

func GetConfig() *Config {
	config := defaultConfig()
	sub := viper.Sub(configPath)
	if sub != nil {
		cmd.BindSubViper(sub, configPath)
		if err := sub.Unmarshal(&config); err != nil {
			log.Log().Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
				"falling back to defaults", configPath, viper.ConfigFileUsed())
		}
	}
	return config
}
