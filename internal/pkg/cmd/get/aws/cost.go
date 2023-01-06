package aws

import (
	"context"
	"github.com/cduggn/cloudcost/internal/pkg/service/aws"
	"github.com/cduggn/cloudcost/internal/pkg/service/display"
	"github.com/spf13/cobra"
)

func CostAndUsageSummary(cmd *cobra.Command, args []string) error {

	req, err := NewCostAndUsageRequest(cmd)
	if err != nil {
		return err
	}

	awsClient := aws.NewAPIClient()
	usage, err := awsClient.GetCostAndUsage(context.Background(), awsClient.Client, req)
	if err != nil {
		return err
	}

	report := aws.CostAndUsageReport{
		Services: make(map[int]aws.Service),
	}
	report.Granularity = req.Granularity
	display.CurateCostAndUsageReport(usage, &report)
	display.PrintCostAndUsageReport(&report)

	return nil
}

func NewCostAndUsageRequest(cmd *cobra.Command) (aws.CostAndUsageRequestType, error) {

	dimensions, err := cmd.Flags().GetStringSlice("dimensions")
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}
	err = ValidateDimension(dimensions)
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}

	tag := cmd.Flags().Lookup("tags").Value.String()
	err = ValidateTag(tag, dimensions)
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}

	filter, _ := cmd.Flags().GetString("filter-by")
	err = ValidateFilterBy(filter, tag)
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}

	start := cmd.Flags().Lookup("start").Value.String()
	err = ValidateStartDate(start)
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}

	end := cmd.Flags().Lookup("end").Value.String()
	err = ValidateEndDate(end, start)
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}

	excludeDiscounts, _ := cmd.Flags().GetBool("exclude-discounts")
	interval := cmd.Flags().Lookup("granularity").Value.String()

	return aws.CostAndUsageRequestType{
		Granularity: interval,
		GroupBy:     dimensions,
		Tag:         tag,
		Time: aws.Time{
			Start: start,
			End:   end,
		},
		IsFilterEnabled: isFilterEnabled(filter),
		TagFilterValue:  filter,
		//Rates:            rates,
		ExcludeDiscounts: excludeDiscounts,
	}, nil

}
