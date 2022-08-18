package github_neonxp_jsonrpc2

import (
	"aurora-relayer-go-common/log"
	"go.neonxp.dev/jsonrpc2/rpc"
)

type Config struct {
	HTTPConfig HTTPConfig // `yaml:"http"`
	Logger     rpc.Logger // `yaml:"logger, omitempty"`
}

type HTTPConfig struct {
	Port int16 // `yaml:"port"`
}

func DefaultConfig() *Config {
	return &Config{
		Logger: NewNeonxpJsonRpc2Logger(log.New()),
	}
}
