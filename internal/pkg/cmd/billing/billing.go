package billing

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"time"
)

var (
	groupBy     []string
	groupByTag  string
	granularity string
	filterBy    string
	billingCmd  = &cobra.Command{
		Use:   "aws",
		Short: "Fetch Cost and Usage information for default account and region",
		Long:  paintHeader(),
	}
	getCmd = &cobra.Command{
		Use:   "get-cost-and-usage",
		Short: "Cost and Usage information",
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
	ok := getCmd.MarkFlagRequired("granularity")
	if ok != nil {
		panic(ok)
	}

	getCmd.Flags().StringVarP(&startDate, "start-date", "s", "", "Start date for billing information")
	ok = getCmd.MarkFlagRequired("start-date")
	if ok != nil {
		panic(ok)
	}
	getCmd.Flags().StringVarP(&endDate, "end-date", "e", time.Now().Format("2006-01-02"), "End date for billing information")
	getCmd.Flags().StringVarP(&filterBy, "filter-by", "f", "", "When grouping by tag, filter by tag value")
	return getCmd
}
