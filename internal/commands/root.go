package commands

import (
	"github.com/cduggn/ccexplorer/internal/commands/aws-presets"
	"github.com/cduggn/ccexplorer/internal/commands/get"
	"github.com/cduggn/ccexplorer/pkg/logger"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ccexplorer",
		Short: "A CLI tool to explore cloud costs and usage",
		Long:  paintRootHeader(),
	}
)

func init() {
	rootCmd.AddCommand(get.AWSCostAndUsageCommand())
	rootCmd.AddCommand(aws_presets.AddAWSPresetCommands())
}

func Execute() {
	_, err := logger.New()
	if err != nil {
		panic(err.Error())
	}

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(126)
	}
}

func paintRootHeader() string {
	myFigure := figure.NewFigure("ccExplorer", "thin", true)
	return myFigure.String()
}
