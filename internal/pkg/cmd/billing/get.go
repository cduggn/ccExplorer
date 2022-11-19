package billing

import (
	"fmt"
	"github.com/cduggn/cloudcost/internal/pkg/billing"
	"github.com/spf13/cobra"
)

var report *billing.CostAndUsageReport

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
