package endpoint

import (
	"aurora-relayer-common/db"
	"aurora-relayer-common/log"
	"aurora-relayer-common/utils"
	"github.com/spf13/viper"
)

const (
	configPath = "Endpoint"
)

type Endpoint struct {
	DbHandler         *db.Handler
	Logger            *log.Log
	disabledEndpoints map[string]bool
}

func New(dbh *db.Handler) *Endpoint {
	if dbh == nil {
		panic("DB Handler should be initialized")
	}

	logger := log.New()
	conf := DefaultConfig()
	sub := viper.Sub(configPath)
	if sub != nil {
		if err := sub.Unmarshal(&conf); err != nil {
			logger.Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
				"falling back to defaults", configPath, viper.ConfigFileUsed())
		}
	}

	de := make(map[string]bool, 0)
	for _, e := range conf.DisabledEndpoints {
		de[e] = true
	}

	return &Endpoint{
		DbHandler:         dbh,
		Logger:            logger,
		disabledEndpoints: de,
	}
}

func (e *Endpoint) IsEndpointAllowed(name string) error {
	if e.disabledEndpoints[name] {
		return &utils.MethodNotFoundError{Method: name}
	}
	return nil
}
