package aws

import "github.com/spf13/cobra"

var (
	costUsageGroupBy          []string
	costUsageGroupByTag       string
	costUsageGranularity      string
	costUsageFilterBy         string
	costUsageStartDate        string
	costUsageEndDate          string
	costUsageWithoutDiscounts bool
	costAndUsageVerboseMode   bool
)

func CostAndUsageCommand(c *cobra.Command) *cobra.Command {

	// Optional Flags used to manage start and end dates for billing
	//information retrieval
	c.Flags().StringVarP(&costUsageStartDate, "start", "s",
		DefaultStartDate(DayOfCurrentMonth, SubtractDays),
		"Defaults to the start of the current month")
	c.Flags().StringVarP(&costUsageEndDate, "end", "e", DefaultEndDate(Format),
		"Defaults to the present day")

	// Mandatory tags used to specify how data will be grouped.
	//This also dictates the type of data that will be returned.
	c.Flags().StringSliceVarP(&costUsageGroupBy, "dimensions", "d",
		[]string{},
		"Group by at most 2 dimension tags [ Dimensions: AZ, SERVICE, "+
			"USAGE_TYPE ]")
	c.Flags().StringVarP(&costUsageGroupByTag, "tags", "t", "",
		"Group by cost allocation tag")

	// Optional flag used to filter data by tag value,
	//this is only relevant when the data is grouped by tag
	c.Flags().StringVarP(&costUsageFilterBy, "filter-by", "f", "",
		"When grouping by tag, filter by tag value")

	// Optional flag to dictate the granularity of the data returned
	c.Flags().StringVarP(&costUsageGranularity, "granularity", "g", "MONTHLY",
		"Granularity of billing information to fetch. Monthly, Daily or Hourly")

	c.Flags().BoolVarP(&costUsageWithoutDiscounts, "exclude-discounts", "c", false,
		"Exclude credit, refund, "+
			"and discount information in the report summary. "+
			"Disabled by default.")

	return c
}
