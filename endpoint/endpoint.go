package endpoint

import (
	"aurora-relayer-go-common/db"
	"aurora-relayer-go-common/log"
	"aurora-relayer-go-common/utils"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	configPath = "Endpoint"
)

type Endpoint struct {
	DbHandler         *db.Handler
	Logger            *log.Log
	disabledEndpoints map[string]bool
	IsEndpointAllowed func(name string) error
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
	ep := &Endpoint{
		DbHandler:         dbh,
		Logger:            logger,
		disabledEndpoints: setDisabledEndpoints(conf),
	}
	ep.IsEndpointAllowed = func(name string) error {
		return isEndpointAllowed(name, ep)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		handleConfigChange(ep)
	})
	return ep
}

func isEndpointAllowed(name string, ep *Endpoint) error {
	if ep.disabledEndpoints[name] {
		return &utils.MethodNotFoundError{Method: name}
	}
	return nil
}

func handleConfigChange(e *Endpoint) error {
	newConf := Config{}
	sub := viper.Sub(configPath)
	if sub != nil {
		if err := sub.Unmarshal(&newConf); err != nil {
			e.Logger.Warn().Err(err).Msgf("failed to parse new configuration [%s] from [%s], "+
				"falling back to the old config", configPath, viper.ConfigFileUsed())
			return err
		}
		if len(newConf.DisabledEndpoints) != len(e.disabledEndpoints) {
			e.disabledEndpoints = setDisabledEndpoints(newConf)
		} else {
			for _, de := range newConf.DisabledEndpoints {
				if !e.disabledEndpoints[de] && len(de) != 0 {
					e.disabledEndpoints = setDisabledEndpoints(newConf)
					break
				}
			}
		}
	}
	return nil
}

func setDisabledEndpoints(conf Config) map[string]bool {
	de := make(map[string]bool, len(conf.DisabledEndpoints))
	for _, e := range conf.DisabledEndpoints {
		de[e] = true
	}
	return de
}
