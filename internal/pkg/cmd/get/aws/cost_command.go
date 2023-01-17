package aws

import "github.com/spf13/cobra"

var (
	costUsageGroupBy          GroupBy
	costUsageGranularity      string
	costUsageStartDate        string
	costUsageEndDate          string
	costUsageWithoutDiscounts bool
)

func CostAndUsageCommand(c *cobra.Command) *cobra.Command {

	c.Flags().VarP(&costUsageGroupBy, "groupBy", "g",
		"Group by DIMENSION and/or TAG . "+
			"Example: --groupBy dimension=SERVICE --groupBy tag=Name"+
			"Example: --groupBy dimension=SERVICE,TAG=Name")

	costUsageFilterBy := NewFilterBy()
	c.Flags().VarP(&costUsageFilterBy, "filterBy", "f",
		"Filter by DIMENSION and/or TAG . "+
			"Example: --filterBy dimension=SERVICE --filterBy tag=Name"+
			"Example: --filterBy dimension=SERVICE,TAG=Name")

	// Optional flag to dictate the granularity of the data returned
	c.Flags().StringVarP(&costUsageGranularity, "granularity", "m",
		"MONTHLY",
		"Sets the Amazon Web Services cost granularity to MONTHLY or DAILY , or HOURLY . If Granularity isn't set, the response object doesn't include the Granularity , either MONTHLY or DAILY , or HOURLY")

	c.Flags().BoolVarP(&costUsageWithoutDiscounts, "excludeDiscounts", "l",
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
