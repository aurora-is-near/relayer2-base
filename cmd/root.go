package cmd

import (
	"strings"

	"github.com/aurora-is-near/relayer2-base/cmdutils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	viper.SetEnvPrefix(cmdutils.EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
