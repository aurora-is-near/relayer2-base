package github_ethereum_go_ethereum

import (
	"time"

	gel "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/aurora-is-near/relayer2-base/log"
)

const (
	DefaultHost       = "localhost" // Default host interface for the HTTP RPC server
	DefaultHTTPPort   = 8545        // Default TCP port for the HTTP RPC server
	DefaultWSPort     = 8546        // Default TCP port for the websocket RPC server
	DefaultPathPrefix = ""          // Default TCP port for the websocket RPC server
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
func DefaultConfig() *Config {
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
