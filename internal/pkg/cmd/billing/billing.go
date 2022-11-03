package billing

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

func paintHeader() string {
	myFigure := figure.NewFigure("billing", "thin", true)
	return myFigure.String()
}

func BillingCmd() *cobra.Command {

	var billingCmd = &cobra.Command{
		Use:   "billing",
		Short: "Fetch Billing information for default account and region",
		Long:  paintHeader(),
	}
	billingCmd.AddCommand(GetBill())

	return billingCmd
}

func GetBill() *cobra.Command {

	var getCommand = &cobra.Command{
		Use:   "get",
		Short: "Bill information",
		Long: `
		GetBill = DESCRIPTION
		Fetches billing information for the current month using the AWS Cost Explorer API
		
		Prerequisites:
		- AWS credentials must be configured in ~/.aws/credentials
		- AWS region must be configured in ~/.aws/config`,
		Run: GetBillingDetailCommand,
	}

	//getCommand.AddCommand(GetBill())

	service := ""
	getCommand.Flags().StringVarP(&service, "service", "s", "", "AWS Service to get billing information for")
	err := getCommand.MarkFlagRequired("service")
	if err != nil {
		panic(err)
	}

	return getCommand
}
