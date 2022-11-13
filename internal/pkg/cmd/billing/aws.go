package billing

import (
	"fmt"
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
		Use:   "cost-and-usage",
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
	groupByDimensionCmd = &cobra.Command{
		Use:   "dimension",
		Short: "Group by dimension",
		Long:  "Group by at most 2 dimension tags [ Dimensions: AZ, SERVICE, USAGE_TYPE ]",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Group by dimension")
		},
	}
	groupByTagCmd = &cobra.Command{
		Use:   "tag",
		Short: "Group by tag",
		Long:  "Group by cost allocation tag",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Group by cost allocation tag")
		},
	}
	groupByTagAndDimensionCmd = &cobra.Command{
		Use:   "tag-and-dimension",
		Short: "Group by tag and dimension",
		Long:  "Group by cost allocation tag and dimension",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Group by cost allocation tag and dimension")
		},
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
	getCmd.AddCommand(GroupByDimensionommand())
	getCmd.AddCommand(GroupByTagCmd())
	getCmd.AddCommand(GroupByTagAndDimensionCmd())
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

func GroupByTagCmd() *cobra.Command {

	groupByTagCmd.Flags().StringVarP(&granularity, "granularity", "g", "DAILY", "Granularity of billing information to fetch")
	ok := groupByTagCmd.MarkFlagRequired("granularity")
	if ok != nil {
		panic(ok)
	}

	groupByTagCmd.Flags().StringVarP(&startDate, "start-date", "s", "", "Start date for billing information")
	ok = groupByTagCmd.MarkFlagRequired("start-date")
	if ok != nil {
		panic(ok)
	}
	groupByTagCmd.Flags().StringVarP(&endDate, "end-date", "e", time.Now().Format("2006-01-02"), "End date for billing information")
	groupByTagCmd.Flags().StringVarP(&filterBy, "filter-by", "f", "", "When grouping by tag, filter by tag value")

	return groupByTagCmd
}

func GroupByDimensionommand() *cobra.Command {

	groupByDimensionCmd.Flags().StringVarP(&granularity, "granularity", "g", "DAILY", "Granularity of billing information to fetch")
	ok := groupByDimensionCmd.MarkFlagRequired("granularity")
	if ok != nil {
		panic(ok)
	}

	groupByDimensionCmd.Flags().StringVarP(&startDate, "start-date", "s", "", "Start date for billing information")
	ok = groupByDimensionCmd.MarkFlagRequired("start-date")
	if ok != nil {
		panic(ok)
	}
	groupByDimensionCmd.Flags().StringVarP(&endDate, "end-date", "e", time.Now().Format("2006-01-02"), "End date for billing information")
	
	return groupByDimensionCmd
}

func GroupByTagAndDimensionCmd() *cobra.Command {

	groupByTagAndDimensionCmd.Flags().StringVarP(&granularity, "granularity", "g", "DAILY", "Granularity of billing information to fetch")
	ok := groupByTagAndDimensionCmd.MarkFlagRequired("granularity")
	if ok != nil {
		panic(ok)
	}

	groupByTagAndDimensionCmd.Flags().StringVarP(&startDate, "start-date", "s", "", "Start date for billing information")
	ok = groupByTagAndDimensionCmd.MarkFlagRequired("start-date")
	if ok != nil {
		panic(ok)
	}
	groupByTagAndDimensionCmd.Flags().StringVarP(&endDate, "end-date", "e", time.Now().Format("2006-01-02"), "End date for billing information")
	groupByTagAndDimensionCmd.Flags().StringVarP(&filterBy, "filter-by", "f", "", "When grouping by tag, filter by tag value")

	return groupByTagAndDimensionCmd
}
