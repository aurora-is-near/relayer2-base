package cmd

import "github.com/spf13/cobra"

const (
	version = "v0.1.0"
)

func VersionCmd(f func(cmd *cobra.Command, args []string)) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "command to manage/display config",
		Run: func(cmd *cobra.Command, args []string) {
			f(cmd, args)
			println("lib: ", version)
		},
	}
}
