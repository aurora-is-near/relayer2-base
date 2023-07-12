package cmdutils

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	EnvPrefix = "AURORA_RELAYER"
)

// BindSubViper is a work around because of a Viper limitation (spf13/viper#507)
// This work around allows environment variables to be used with sub configs as well.
func BindSubViper(sub *viper.Viper, subConfigPath string) {
	sub.AutomaticEnv()
	sub.SetEnvPrefix(EnvPrefix + "_" + subConfigPath)
	sub.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
