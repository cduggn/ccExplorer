package cost_and_usage

import (
	"github.com/cduggn/ccexplorer/internal/commands/get/aws/custom_flags"
	"github.com/cduggn/ccexplorer/pkg/helpers"
	"github.com/spf13/cobra"
)

var (
	costUsageGroupBy          custom_flags.DimensionAndTagFlag
	costUsageGranularity      string
	costUsageStartDate        string
	costUsageEndDate          string
	costUsageWithoutDiscounts bool
	costAndUsagePrintFormat   string
)

func CostAndUsageCommand(c *cobra.Command) *cobra.Command {

	c.Flags().VarP(&costUsageGroupBy, "groupBy", "g",
		"Group by DIMENSION and/or TAG ")

	costUsageFilterBy := custom_flags.NewFilterBy()
	c.Flags().VarP(&costUsageFilterBy, "filterBy", "f",
		"Filter by DIMENSION and/or TAG")

	// Optional flag to dictate the granularity of the data returned
	c.Flags().StringVarP(&costUsageGranularity, "granularity", "m",
		"MONTHLY",
		"Valid values: DAILY, MONTHLY, "+
			"HOURLY. (default: MONTHLY)")

	c.Flags().BoolVarP(&costUsageWithoutDiscounts, "excludeDiscounts", "l",
		false,
		"Exclude credit, refunds, "+
			"and discounts (default is to include)")

	c.Flags().StringVarP(&costUsageStartDate, "startDate", "s",
		helpers.DefaultStartDate(helpers.DayOfCurrentMonth, helpers.SubtractDays),
		"Start date (defaults to the start of the current month)")
	c.Flags().StringVarP(&costUsageEndDate, "endDate", "e",
		helpers.DefaultEndDate(helpers.Format),
		"End date *(defaults to the present day)")

	c.Flags().StringVarP(&costAndUsagePrintFormat, "printFormat", "p", "stdout",
		"Valid values: stdout, csv, chart (default: stdout)")

	return c
}
