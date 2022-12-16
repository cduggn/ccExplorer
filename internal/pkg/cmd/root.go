package cmd

import (
	"github.com/cduggn/cloudcost/internal/pkg/cmd/get"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ccxplorer",
		Short: "A CLI tool to explore cloud costs and usage",
		Long:  paintRootHeader(),
	}
)

func init() {
	rootCmd.AddCommand(get.AWSCostAndUsageCommand())
}

func paintRootHeader() string {
	myFigure := figure.NewFigure("ccXplorer", "thin", true)
	return myFigure.String()
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(126)
	}
}
