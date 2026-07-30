// Package cmd is used for command line argument passing
package cmd

import (
	"fmt"

	goversion "github.com/caarlos0/go-version"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cobra.EnableCommandSorting = false

	rootCmd := cobra.Command{
		Use:          "flix",
		Short:        "A tool for quickly finding links to movies or tv shows, downloading and more.",
		SilenceUsage: true,
		Version:      buildVersion().GitVersion,
	}

	rootCmd.AddCommand(
		WatchCommand(),
		InfoCommand(),
		versionCmd(),
	)

	return &rootCmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "Show detailed version information",
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(buildVersion().String())
		},
	}
}

// buildVersion constructs and returns the application's version information.
func buildVersion() goversion.Info {
	return goversion.GetVersionInfo(
		goversion.WithAppDetails("flix", "Tool to work with movies and tv shows.", ""),
	)
}
