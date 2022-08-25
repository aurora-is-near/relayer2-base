package preprocessor

import (
	"aurora-relayer-go-common/endpoint"
	"aurora-relayer-go-common/log"
	"aurora-relayer-go-common/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type EnableDisable struct {
	disabledEndpoints map[string]bool
}

func NewEnableDisable() endpoint.Preprocessor {

	conf, err := readConfig()
	if err != nil {
		log.Log().Warn().Err(err).Msgf("failed to parse configuration [%s] from [%s], "+
			"falling back to defaults", endpoint.ConfigPath, viper.ConfigFileUsed())
	}

	ed := EnableDisable{disabledEndpoints: setDisabledEndpoints(conf.DisabledEndpoints)}
	viper.OnConfigChange(func(e fsnotify.Event) {
		handleConfigChange(&ed)
	})

	return func(name string, _ *endpoint.Endpoint, _ ...any) (bool, *any, error) {
		if ed.disabledEndpoints[name] {
			return true, nil, &utils.MethodNotFoundError{Method: name}
		}
		return false, nil, nil
	}
}

func setDisabledEndpoints(conf []string) map[string]bool {
	de := make(map[string]bool, len(conf))
	for _, e := range conf {
		de[e] = true
	}
	return de
}

func handleConfigChange(e *EnableDisable) {

	newConf, err := readConfig()
	if err != nil {
		log.Log().Warn().Err(err).Msgf("failed to parse new configuration [%s] from [%s], "+
			"falling back to the old config", endpoint.ConfigPath, viper.ConfigFileUsed())
	}

	if len(newConf.DisabledEndpoints) != len(e.disabledEndpoints) {
		e.disabledEndpoints = setDisabledEndpoints(newConf.DisabledEndpoints)
	} else {
		for _, de := range newConf.DisabledEndpoints {
			if !e.disabledEndpoints[de] && len(de) != 0 {
				e.disabledEndpoints = setDisabledEndpoints(newConf.DisabledEndpoints)
				break
			}
		}
	}
}

func readConfig() (endpoint.Config, error) {
	conf := endpoint.DefaultConfig()
	sub := viper.Sub(endpoint.ConfigPath)
	if sub != nil {
		if err := sub.Unmarshal(&conf); err != nil {
			return conf, err
		}
	}
	return conf, nil
}
