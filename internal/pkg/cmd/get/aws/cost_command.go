package aws

import "github.com/spf13/cobra"

var (
	costUsageGroupBy           GroupBy
	costUsageGranularity       string
	costUsageFilterByTag       string
	costUsageStartDate         string
	costUsageEndDate           string
	costUsageWithoutDiscounts  bool
	costUsageFilterByDimension map[string]string
)

func CostAndUsageCommand(c *cobra.Command) *cobra.Command {

	c.Flags().VarP(&costUsageGroupBy, "groupBy", "b",
		"Group by DIMENSION and/or TAG . "+
			"Example: --groupBy dimension=SERVICE --groupBy tag=Name"+
			"Example: --groupBy dimension=SERVICE,TAG=Name")

	c.Flags().StringVarP(&costUsageFilterByTag, "filterByTag", "t", "",
		"Results can be filtered by custom cost allocation tags. "+
			"groupByTag must also be used in conjection with this flag.")

	c.Flags().StringToStringVarP(&costUsageFilterByDimension,
		"filterByDimension",
		"d",
		nil, "Filter by dimension . "+
			"Example: -u SERVICE='Amazon Simple Storage Service'")

	// Optional flag to dictate the granularity of the data returned
	c.Flags().StringVarP(&costUsageGranularity, "granularity", "g",
		"MONTHLY",
		"Sets the Amazon Web Services cost granularity to MONTHLY or DAILY , or HOURLY . If Granularity isn't set, the response object doesn't include the Granularity , either MONTHLY or DAILY , or HOURLY")

	c.Flags().BoolVarP(&costUsageWithoutDiscounts, "excludeDiscounts", "x",
		false,
		"Excludes credit, refund, "+
			"and discount information in the report summary. "+
			"Disabled by default.")

	c.Flags().StringVarP(&costUsageStartDate, "startDate", "s",
		DefaultStartDate(DayOfCurrentMonth, SubtractDays),
		"Defaults to the start of the current month")
	c.Flags().StringVarP(&costUsageEndDate, "endDate", "e",
		DefaultEndDate(Format),
		"Defaults to the present day")

	return c
}
