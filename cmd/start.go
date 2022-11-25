package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultConfigFile = "config/testnet.yaml"
)

func StartCmd(f func(cmd *cobra.Command, args []string)) *cobra.Command {

	startCmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"s"},
		Short:   "Starts Aurora Relayer",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return bindConfiguration(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			f(cmd, args)
		},
	}
	startCmd.PersistentFlags().StringP("config", "c", "", "Path of the configuration file (default -> config/testnet.yaml)")
	return startCmd
}

func bindConfiguration(cmd *cobra.Command) error {
	configFile, _ := cmd.Flags().GetString("config")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigFile(defaultConfigFile)
	}

	viper.WatchConfig()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return nil
}
