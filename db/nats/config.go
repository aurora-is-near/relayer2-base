package nats

import (
	"aurora-relayer-go-common/log"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

const (
	defaultBucket = "filter_bucket"
	defaultUrl    = "nats://localhost:4222/"

	configPath = "db.nats"
)

type Config struct {
	Bucket     string       `mapstructure:"bucket"`
	NatsConfig nats.Options `mapstructure:"options"`
}

func defaultConfig() *Config {
	options := nats.GetDefaultOptions()
	options.Url = defaultUrl
	return &Config{
		Bucket:     defaultBucket,
		NatsConfig: options,
	}
}

func GetConfig() *Config {
	config := defaultConfig()
	sub := viper.Sub(configPath)
	if sub != nil {
		if err := sub.Unmarshal(&config); err != nil {
			log.Log().Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
				"falling back to defaults", configPath, viper.ConfigFileUsed())
		}
	}
	return config
}
