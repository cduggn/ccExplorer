package billing

import (
	"fmt"
	billingapi "github.com/cduggn/cloudcost/internal/pkg/billing"
	"github.com/spf13/cobra"
)

func GetBillingSummary(cmd *cobra.Command, args []string) {
	req := NewCostAndUsageRequest(cmd)
	report := billingapi.GetAWSCostAndUsage(req)
	report.Print()
}

func NewCostAndUsageRequest(cmd *cobra.Command) billingapi.CostAndUsageRequest {

	dimensions, err := cmd.Flags().GetStringSlice("group-by-dimension")
	if err != nil {
		fmt.Println(err)
	}

	filterBy, _ := cmd.Flags().GetString("filter-by")

	return billingapi.CostAndUsageRequest{
		Granularity: cmd.Flags().Lookup("granularity").Value.String(),
		GroupBy:     dimensions,
		Tag:         cmd.Flags().Lookup("group-by-tag").Value.String(),
		Time: billingapi.Time{
			Start: cmd.Flags().Lookup("start-date").Value.String(),
			End:   cmd.Flags().Lookup("end-date").Value.String(),
		},
		IsFilterEnabled: isFilterEnabled(filterBy),
		TagFilterValue:  filterBy,
	}

}

func isFilterEnabled(filterBy string) bool {
	if filterBy != "" {
		return true
	} else {
		return false
	}
}
