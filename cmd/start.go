package cmd

import "github.com/spf13/cobra"

func StartCmd(f func(cmd *cobra.Command, args []string)) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "starts Aurora Relayer",
		Run: func(cmd *cobra.Command, args []string) {
			f(cmd, args)
		},
	}
}
