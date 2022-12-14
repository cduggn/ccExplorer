package aws

import (
	"fmt"
	"github.com/cduggn/cloudcost/internal/pkg/billing"
	"github.com/spf13/cobra"
)

var (
	excludeCredits bool
	groupBy        []string
	groupByTag     string
	granularity    string
	filterBy       string
	rates          []string
	startDate      string
	endDate        string
	report         *billing.CostAndUsageReport
)

func GetCommand(c *cobra.Command) *cobra.Command {
	// create new flagset for get command
	c.Flags().StringSliceVarP(&groupBy, "group-by-dimension", "d", []string{"SERVICE", "USAGE_TYPE"}, "Group by at most 2 dimension tags [ Dimensions: AZ, SERVICE, USAGE_TYPE ]")
	c.Flags().StringVarP(&groupByTag, "group-by-tag", "t", "", "Group by cost allocation tag")
	c.Flags().StringVarP(&filterBy, "filter-by", "f", "", "When grouping by tag, filter by tag value")

	c.Flags().StringVarP(&startDate, "start-date", "s", PastMonth(), "Start date for billing information. Defaults to the past 7 days")
	c.Flags().StringVarP(&endDate, "end-date", "e", Today(), "End date for billing information. Default is todays date.")

	// Granularity and rate flags
	c.Flags().StringVarP(&granularity, "granularity", "g", "MONTHLY", "Granularity of billing information to fetch. Monthly, Daily or Hourly")
	c.Flags().StringSliceVarP(&rates, "rates", "r", []string{"UNBLENDED_COST"}, "Cost and Usage rates to fetch [ Rates: BLENDED_COST, UNBLENDED_COST, AMORTIZED_COST, NET_AMORTIZED_COST, NET_UNBLENDED_COST, USAGE_QUANTITY ]. Defaults to UNBLENDED_COST")

	// Other flags
	c.Flags().BoolVarP(&excludeCredits, "exclude-credit", "c", false, "Exclude credit and refund information in the report. This is enabled by default")

	return c
}

func GetBillingSummary(cmd *cobra.Command, args []string) {
	req := NewCostAndUsageRequest(cmd)
	report = billing.GetAWSCostAndUsage(req)
	report.Print()
}

func NewCostAndUsageRequest(cmd *cobra.Command) billing.CostAndUsageRequest {

	dimensions, err := cmd.Flags().GetStringSlice("group-by-dimension")
	if err != nil {
		fmt.Println(err)
	}

	rates, err := cmd.Flags().GetStringSlice("rates")
	if err != nil {
		fmt.Println(err)
	}
	filterBy, _ := cmd.Flags().GetString("filter-by")
	excludeCredits, _ := cmd.Flags().GetBool("exclude-credit")

	return billing.CostAndUsageRequest{
		Granularity: cmd.Flags().Lookup("granularity").Value.String(),
		GroupBy:     dimensions,
		Tag:         cmd.Flags().Lookup("group-by-tag").Value.String(),
		Time: billing.Time{
			Start: cmd.Flags().Lookup("start-date").Value.String(),
			End:   cmd.Flags().Lookup("end-date").Value.String(),
		},
		IsFilterEnabled: isFilterEnabled(filterBy),
		TagFilterValue:  filterBy,
		Rates:           rates,
		ExcludeCredits:  excludeCredits,
	}

}

func isFilterEnabled(filterBy string) bool {
	if filterBy != "" {
		return true
	} else {
		return false
	}
}
