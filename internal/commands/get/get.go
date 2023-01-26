package get

import (
	"github.com/cduggn/ccexplorer/internal/commands/get/aws"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

const (
	usageExample = `
  # Costs grouped by LINKED_ACCOUNT 
  ccexplorer get aws -g DIMENSION=LINKED_ACCOUNT
  
  # Costs grouped by CommittedThroughput operation and SERVICE
  ccexplorer get aws -g DIMENSION=OPERATION,DIMENSION=SERVICE -s 2022-10-10 -f OPERATION="CommittedThroughput" -l

  # Costs grouped by CommittedThroughput and LINKED_ACCOUNT
  ccexplorer get aws -g DIMENSION=OPERATION,DIMENSION=LINKED_ACCOUNT  -s 2022-10-10 -f OPERATION="CommittedThroughput" -l

  # DynamodDB costs grouped by OPERATION
  ccexplorer get aws -g DIMENSION=OPERATION,DIMENSION=SERVICE -s 2022-10-10 -f SERVICE="Amazon DynamoDB" -l

  # All service costs grouped by SERVICE
  ccexplorer get aws -g DIMENSION=SERVICE -s 2022-10-10

  # All service costs grouped by SERVICE and OPERATION
  ccexplorer get aws -g DIMENSION=SERVICE,DIMENSION=OPERATION -s 2022-10- -l

  # S3 costs grouped by OPERATION
  ccexplorer get aws -g DIMENSION=OPERATION,DIMENSION=SERVICE -s 2022-04-04  -f SERVICE="Amazon Simple Storage Service" -l

  # Costs grpuped by ApplicationName Cost Allocation Tag
  ccexplorer get aws -g TAG=ApplicationName,DIMENSION=OPERATION -s 2022-12-10 -l
`
)

var (
	costAndUsageCmd = &cobra.Command{
		Use:   "get",
		Short: "Cost and usage summary for AWS services",
		Long:  paintHeader(),
	}
	awsCost = &cobra.Command{
		Use:   "aws",
		Short: "Explore UNBLENDED cost summaries for AWS",
		Long: `
Command: aws 
Description: Cost and usage summary for AWS services.

Prerequisites:
- AWS credentials configured in ~/.aws/credentials and default region configured in ~/.aws/config. Alternatively, 
you can set the environment variables AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY and AWS_REGION.`,
		Example: usageExample,
		RunE:    aws.CostAndUsageSummary,
	}
	forecast = &cobra.Command{
		Use:   "forecast",
		Short: "Return cost and usage forecasts for your account.",
		RunE:  aws.CostForecast,
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
