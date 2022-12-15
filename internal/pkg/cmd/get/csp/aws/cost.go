package aws

import (
	"github.com/cduggn/cloudcost/internal/pkg/csp/aws"
	"github.com/spf13/cobra"
)

var (
	groupBy     []string
	groupByTag  string
	granularity string
	filterBy    string
	rates       []string
	startDate   string
	endDate     string
	report      *aws.CostAndUsageReport
)

func AWSCostCommand(c *cobra.Command) *cobra.Command {

	// Optional Flags used to manage start and end dates for billing
	//information retrieval
	c.Flags().StringVarP(&startDate, "start", "s",
		DefaultStartDate(DayOfCurrentMonth, SubtractDays),
		"Defaults to the start of the current month")
	c.Flags().StringVarP(&endDate, "end", "e", DefaultEndDate(Format),
		"Defaults to the present day")

	// Mandatory tags used to specify how data will be grouped.
	//This also dictates the type of data that will be returned.
	c.Flags().StringSliceVarP(&groupBy, "dimensions", "d",
		[]string{},
		"Group by at most 2 dimension tags [ Dimensions: AZ, SERVICE, "+
			"USAGE_TYPE ]")
	c.Flags().StringVarP(&groupByTag, "tags", "t", "",
		"Group by cost allocation tag")

	// Optional flag used to filter data by tag value,
	//this is only relevant when the data is grouped by tag
	c.Flags().StringVarP(&filterBy, "filter-by", "f", "",
		"When grouping by tag, filter by tag value")

	// Optional flag to dictate the granularity of the data returned
	c.Flags().StringVarP(&granularity, "granularity", "g", "MONTHLY",
		"Granularity of billing information to fetch. Monthly, Daily or Hourly")

	return c
}