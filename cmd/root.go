package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfigName = "relayer"
	envPrefix         = "AURORA_RELAYER"
)

func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "relayer",
		Short: "Aurora Relayer Implementation",
		Long:  `A JSON-RPC server implementation compatible with Ethereum's Web3 API for Aurora Engine instances deployed on the NEAR Protocol`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return bind(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.OutOrStdout()
		},
	}

	return rootCmd
}

func bind(_ *cobra.Command) error {
	viper.SetConfigName(defaultConfigName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/home/altug/repo/new-relayer-2/aurora-relayer-go")
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return nil
}

// TODO flag and env binding
// func bindFlags(cmd *cobra.Command, v *viper.Viper) {
// 	cmd.Flags().VisitAll(func(f *pflag.Flag) {
// 		// Environment variables can't have dashes in them, so bind them to their equivalent
// 		if strings.Contains(f.Name, "-") {
// 			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
// 			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
// 		}
//
// 		// Apply the viper config value to the flag when the flag is not set and viper has a value
// 		if !f.Changed && v.IsSet(f.Name) {
// 			val := v.Get(f.Name)
// 			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
// 		}
// 	})
// }
