package get

import (
	"github.com/cduggn/ccexplorer/internal/commands/get/aws"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

var (
	costAndUsageCmd = &cobra.Command{
		Use:   "get",
		Short: "Fetch Cost and Usage information for cloud provider",
		Long:  paintHeader(),
	}
	awsCost = &cobra.Command{
		Use:   "aws",
		Short: "Return unblended cost summary",
		Long: `
Command: aws 
Description: Returns cost and usage summary for the specified time period.

Prerequisites:
- AWS credentials must be configured in ~/.aws/credentials
- AWS region must be configured in ~/.aws/config
- Cost Allocation Tags must exist in AWS console if you want to filter by tag ( 
Note cost allocation tags can take up to 24 hours to be applied )`,
		RunE: aws.CostAndUsageSummary,
	}
	forecast = &cobra.Command{
		Use: "forecast",
		Short: "Return cost, usage, " +
			"and resoucrce information including ARN",
		RunE: aws.CostForecast,
	}
)

func paintHeader() string {
	myFigure := figure.NewFigure("Cost And Usage", "thin", true)
	return myFigure.String()
}

func AWSCostAndUsageCommand() *cobra.Command {

	costAndUsageCmd.AddCommand(aws.CostAndUsageCommand(awsCost))
	awsCost.AddCommand(aws.ForecastCommand(forecast))
	return costAndUsageCmd
}
