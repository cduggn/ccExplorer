package billing

import (
	"fmt"
	"github.com/cduggn/cloudcost/internal/pkg/billing"
	"github.com/spf13/cobra"
)

func GetBillingDetailCommand(cmd *cobra.Command, args []string) {

	service := cmd.Flags().Lookup("service").Value.String()
	fmt.Println("Service: ", service)
	billable := billing.FetchCloudCost()
	bill := billable.ForService(service)
	billable.Print(bill)
}
