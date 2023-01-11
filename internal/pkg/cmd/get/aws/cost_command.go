package aws

import "github.com/spf13/cobra"

var (
	costUsageGroupBy    []string
	costUsageGroupByTag string
	//costUsageFilterByDimension          string
	costUsageGranularity      string
	costUsageFilterBy         string
	costUsageStartDate        string
	costUsageEndDate          string
	costUsageWithoutDiscounts bool
	//costUsageFilterByDimensionValue     string
	costUsageFilterByDimension map[string]string
)

func CostAndUsageCommand(c *cobra.Command) *cobra.Command {

	// Optional Flags used to manage start and end dates for billing
	//information retrieval

	// Mandatory tags used to specify how data will be grouped.
	//This also dictates the type of data that will be returned.
	c.Flags().StringSliceVarP(&costUsageGroupBy, "groupByDimension", "d",
		[]string{},
		"Group by at most 2 dimension tags [ Dimensions: AZ, SERVICE, "+
			"USAGE_TYPE ]")
	c.Flags().StringVarP(&costUsageGroupByTag, "groupByTag", "t", "",
		"Group by cost allocation tag. Example: ApplicationName, Environment, BucketName")

	// Optional flag used to filter data by tag value,
	//this is only relevant when the data is grouped by tag
	c.Flags().StringVarP(&costUsageFilterBy, "filterByTagName", "f", "",
		"Results can be filtered by custom cost allocation tags. "+
			"The groupByTag flag must be set with an active cost allocation"+
			" tag. Once the tag is set, the filterByTagName flag can be used")

	c.Flags().StringToStringVarP(&costUsageFilterByDimension,
		"filterByDimensionNameValue",
		"u",
		nil, "Filter by dimension . "+
			"Example: -U SERVICE='Amazon Simple Storage Service'")

	// Optional flag to dictate the granularity of the data returned
	c.Flags().StringVarP(&costUsageGranularity, "granularity", "g",
		"MONTHLY",
		"Sets the Amazon Web Services cost granularity to MONTHLY or DAILY , or HOURLY . If Granularity isn't set, the response object doesn't include the Granularity , either MONTHLY or DAILY , or HOURLY")

	c.Flags().BoolVarP(&costUsageWithoutDiscounts, "excludeDiscounts", "c",
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
