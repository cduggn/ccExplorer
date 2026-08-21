package cli

import (
	"github.com/cduggn/ccexplorer/internal/config"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ccexplorer",
		Short: "A CLI tool to explore cloud costs and usage",
		Long:  paintRootHeader(),
	}
)

// RootCommand returns the root command annotated with the build version.
// AWS clients are not constructed here: they are resolved lazily by the
// subcommands that need them so that --help and --version work without
// credentials.
func RootCommand(version string) *cobra.Command {
	config.LoadConfigFunc(".")()
	rootCmd.Version = version
	return rootCmd
}

func init() {
	rootCmd.AddCommand(CostAndForecast())
	rootCmd.AddCommand(mcpCommand())
}

func paintRootHeader() string {
	myFigure := figure.NewFigure("ccExplorer", "thin", true)
	return myFigure.String()
}
