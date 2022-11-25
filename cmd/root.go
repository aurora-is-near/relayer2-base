package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	envPrefix = "AURORA_RELAYER"
)

func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "relayer",
		Short: "Aurora Relayer Implementation",
		Long:  `A JSON-RPC server implementation compatible with Ethereum's Web3 API for Aurora Engine instances deployed on the NEAR Protocol`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			bindRootViper()
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.OutOrStdout()
		},
	}

	return rootCmd
}

func bindRootViper() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// BindSubViper is a work around because of a Viper limitation (spf13/viper#507)
// This work around allows environment variables to be used with sub configs as well.
func BindSubViper(sub *viper.Viper, subConfigPath string) {
	sub.AutomaticEnv()
	sub.SetEnvPrefix(envPrefix + "_" + subConfigPath)
	sub.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
