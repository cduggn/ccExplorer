package billing

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

var (
	groupBy     []string
	groupByTag  string
	granularity string
	filterBy    string
	billingCmd  = &cobra.Command{
		Use:   "get",
		Short: "Fetch Cost and Usage information for cloud provider",
		Long:  paintHeader(),
	}
	getCmd = &cobra.Command{
		Use:   "aws",
		Short: "Cost and Usage information for AWS service",
		Long: `
		GetBill = DESCRIPTION
		Fetches billing information for the time interval provided using the AWS Cost Explorer API
		
		Prerequisites:
		- AWS credentials must be configured in ~/.aws/credentials
		- AWS region must be configured in ~/.aws/config
		- Cost Allocation Tags if you want to filter by tag ( Note cost allocation tags can take up to 24 hours to be applied )`,
		Run: GetBillingSummary,
	}
	startDate string
	endDate   string
)

func paintHeader() string {
	myFigure := figure.NewFigure("billing", "thin", true)
	return myFigure.String()
}

func CostAndUsageCommand() *cobra.Command {
	billingCmd.AddCommand(GetCommand())

	return billingCmd
}

func GetCommand() *cobra.Command {
	getCmd.Flags().StringSliceVarP(&groupBy, "group-by-dimension", "d", []string{"SERVICE", "USAGE_TYPE"}, "Group by at most 2 dimension tags [ Dimensions: AZ, SERVICE, USAGE_TYPE ]")
	getCmd.Flags().StringVarP(&groupByTag, "group-by-tag", "t", "", "Group by cost allocation tag")
	getCmd.Flags().StringVarP(&granularity, "granularity", "g", "DAILY", "Granularity of billing information to fetch")

	getCmd.Flags().StringVarP(&startDate, "start-date", "s", Format(SubtractDays(Time(), 7)), "Start date for billing information. Defaults to the past 7 days")
	getCmd.Flags().StringVarP(&endDate, "end-date", "e", Format(Time()), "End date for billing information. Default is todays date.")
	getCmd.Flags().StringVarP(&filterBy, "filter-by", "f", "", "When grouping by tag, filter by tag value")

	return getCmd
}
