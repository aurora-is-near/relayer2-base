package probe

import (
	"github.com/aurora-is-near/relayer2-base/cmdutils"
	"github.com/aurora-is-near/relayer2-base/log"
	"github.com/spf13/viper"
)

const (
	defaultAddress   = "127.0.0.1:28080"
	defaultNamespace = "Aurora"
	defaultSubsystem = "Relayer2"
	defaultEnabled   = false

	configPath = "probe"
)

type ServerConfig struct {
	Address   string `mapstructure:"address"`
	Namespace string `mapstructure:"namespace"`
	Subsystem string `mapstructure:"subsystem"`
}

type MetricConfig struct {
	Id          string    `mapstructure:"id"`
	Type        string    `mapstructure:"type"`
	Name        string    `mapstructure:"name"`
	Help        string    `mapstructure:"help"`
	Buckets     []float64 `mapstructure:"buckets"`
	LabelNames  []string  `mapstructure:"labelNames"`
	LabelValues []string  `mapstructure:"labelValues"`
}

type Config struct {
	Enabled       bool            `mapstructure:"enable"`
	ServerConfig  *ServerConfig   `mapstructure:"server"`
	MetricConfigs *[]MetricConfig `mapstructure:"metrics"`
}

func defaultConfig() *Config {
	return &Config{
		Enabled: defaultEnabled,
		ServerConfig: &ServerConfig{
			Address:   defaultAddress,
			Namespace: defaultNamespace,
			Subsystem: defaultSubsystem,
		},
		MetricConfigs: &[]MetricConfig{},
	}
}

func GetConfig() *Config {
	config := defaultConfig()
	sub := viper.Sub(configPath)
	if sub != nil {
		cmdutils.BindSubViper(sub, configPath)
		if err := sub.Unmarshal(&config); err != nil {
			log.Log().Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
				"falling back to defaults", configPath, viper.ConfigFileUsed())
		}
	}

	return config
}
