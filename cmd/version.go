package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var versionInfo = ""

var buildInfo = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.String()
	}
	return "Build info not available"
}()

func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Command to display version and build information",
		Run: func(cmd *cobra.Command, args []string) {
			if versionInfo != "" {
				fmt.Println("Version:\t", versionInfo)
				println("")
			}
			println("Build Info:")
			println(buildInfo)
		},
	}
}
