package get

import (
	"github.com/cduggn/ccexplorer/internal/commands/get/aws/cost_and_usage"
	"github.com/cduggn/ccexplorer/internal/commands/get/aws/forecast"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

var (
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Cost and usage summary for AWS services",
		Long:  paintHeader(),
	}
	costAndUsageCmd = &cobra.Command{
		Use:   "aws",
		Short: "Explore UNBLENDED cost summaries for AWS",
		Long: `
Command: aws 
Description: Cost and usage summary for AWS services.

Prerequisites:
- AWS credentials configured in ~/.aws/credentials and default region configured in ~/.aws/config. Alternatively, 
you can set the environment variables AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY and AWS_REGION.`,
		Example: costAndUsageExamples,
		RunE:    cost_and_usage.CostAndUsageRunCmd,
	}
	forecastCmd = &cobra.Command{
		Use:     "forecast",
		Short:   "Return cost and usage forecasts for your account.",
		Example: forecastExamples,
		RunE:    forecast.CostForecastRunCmd,
	}
)

func paintHeader() string {
	myFigure := figure.NewFigure("CostAndUsage", "thin", true)
	return myFigure.String()
}

func AWSCostAndUsageCommand() *cobra.Command {
	getCmd.AddCommand(cost_and_usage.CostAndUsageCommand(costAndUsageCmd))
	costAndUsageCmd.AddCommand(forecast.ForecastCommand(forecastCmd))
	return getCmd
}
