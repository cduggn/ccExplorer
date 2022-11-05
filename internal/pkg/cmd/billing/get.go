package billing

import (
	"fmt"
	"github.com/cduggn/cloudcost/internal/pkg/billing"
	"github.com/spf13/cobra"
)

func GetBillingSummary(cmd *cobra.Command, args []string) {
	req := NewCostAndUsageRequest(cmd)
	report := billing.GetAWSCostAndUsage(req)
	report.Print()
}

func NewCostAndUsageRequest(cmd *cobra.Command) billing.CostAndUsageRequest {

	dimensions, err := cmd.Flags().GetStringSlice("group-by-dimension")
	if err != nil {
		fmt.Println(err)
	}

	return billing.CostAndUsageRequest{
		Granularity: cmd.Flags().Lookup("granularity").Value.String(),
		GroupBy:     dimensions,
		Tag:         cmd.Flags().Lookup("group-by-tag").Value.String(),
		Time: billing.Time{
			Start: cmd.Flags().Lookup("start-date").Value.String(),
			End:   cmd.Flags().Lookup("end-date").Value.String(),
		},
	}

}
