package cmd

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cloudcost",
		Short: "A CLI tool to get AWS Costs, Usage and Forecasts",
		Long:  paintRootHeader(),
	}
)

func init() {

	rootCmd.AddCommand(AWSCostAndUsageCommand())
}

func paintRootHeader() string {
	myFigure := figure.NewFigure("CloudCost", "thin", true)
	return myFigure.String()
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(126)
	}
}
