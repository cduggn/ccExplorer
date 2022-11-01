package billing

import (
	"github.com/cduggn/cloudcost/internal/pkg/billing"
	"github.com/spf13/cobra"
)

func GetBillingDetailCommand(cmd *cobra.Command, args []string) {
	service := cmd.Flags().Lookup("service").Value.String()
	billable := billing.FetchCloudCost()
	bill := billable.ForService(service)
	billable.Print(bill)
}
