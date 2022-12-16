package aws

import "github.com/spf13/cobra"

var (
	rgroupBy    []string
	rgroupByTag string
	rfilterBy   string
)

func CostAndResourcesCommand(c *cobra.Command) *cobra.Command {

	// Optional Flags used to manage start and end dates for billing
	//information retrieval
	// Mandatory tags used to specify how data will be grouped.
	//This also dictates the type of data that will be returned.
	c.Flags().StringSliceVarP(&rgroupBy, "dimensions", "d",
		[]string{},
		"Group by at most 2 dimension tags [ Dimensions: AZ, SERVICE, "+
			"USAGE_TYPE ]")
	c.Flags().StringVarP(&rgroupByTag, "tags", "t", "",
		"Group by cost allocation tag")

	// Optional flag used to filter data by tag value,
	//this is only relevant when the data is grouped by tag
	c.Flags().StringVarP(&rfilterBy, "filter-by", "f", "",
		"When grouping by tag, filter by tag value")

	return c
}
