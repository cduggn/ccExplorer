package aws

import (
	"github.com/spf13/cobra"
)

func AWSCostWithDiscountsCommand(c *cobra.Command) *cobra.Command {
	// create new flagset for get command
	c.Flags().StringSliceVarP(&groupBy, "group-by-dimension", "d", []string{"SERVICE", "USAGE_TYPE"}, "Group by at most 2 dimension tags [ Dimensions: AZ, SERVICE, USAGE_TYPE ]")
	c.Flags().StringVarP(&groupByTag, "group-by-tag", "t", "", "Group by cost allocation tag")
	c.Flags().StringVarP(&filterBy, "filter-by", "f", "", "When grouping by tag, filter by tag value")

	c.Flags().StringVarP(&startDate, "start-date", "s", PastMonth(), "Start date for billing information. Defaults to the past 7 days")
	c.Flags().StringVarP(&endDate, "end-date", "e", Today(), "End date for billing information. Default is todays date.")

	// Granularity and rate flags
	c.Flags().StringVarP(&granularity, "granularity", "g", "MONTHLY", "Granularity of billing information to fetch. Monthly, Daily or Hourly")

	return c
}
