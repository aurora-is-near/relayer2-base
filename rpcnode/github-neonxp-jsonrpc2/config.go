package github_neonxp_jsonrpc2

const (
	defaultHttpHost = "localhost" // Default host interface for the HTTP RPC server
	defaultHttpPort = 8545        // Default TCP port for the HTTP RPC server
)

type Config struct {
	HttpPort int16  `mapstructure:"httpPort"`
	HttpHost string `mapstructure:"httpHost"`
}

func DefaultConfig() *Config {
	return &Config{
		HttpPort: defaultHttpPort,
		HttpHost: defaultHttpHost,
	}
}
