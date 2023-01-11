package aws

import (
	"context"
	"github.com/cduggn/ccexplorer/internal/pkg/service/aws"
	"github.com/cduggn/ccexplorer/internal/pkg/service/display"
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

	report := display.CurateCostAndUsageReport(usage, req.Granularity)
	display.PrintCostAndUsageReport(display.SortServicesByMetricAmount, report)

	return nil
}

func NewCostAndUsageRequest(cmd *cobra.Command) (aws.CostAndUsageRequestType, error) {

	dimensions, err := cmd.Flags().GetStringSlice("groupByDimension")
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}
	err = ValidateDimension(dimensions)
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}

	tag := cmd.Flags().Lookup("groupByTag").Value.String()
	err = ValidateTag(tag, dimensions)
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}

	filterByTag, _ := cmd.Flags().GetString("filterByTagName")
	err = ValidateFilterBy(filterByTag, tag)
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}

	dimensionFilterMap, _ := cmd.Flags().GetStringToString(
		"filterByDimensionNameValue")

	start := cmd.Flags().Lookup("startDate").Value.String()
	err = ValidateStartDate(start)
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}

	end := cmd.Flags().Lookup("endDate").Value.String()
	err = ValidateEndDate(end, start)
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}

	excludeDiscounts, _ := cmd.Flags().GetBool("excludeDiscounts")
	interval := cmd.Flags().Lookup("granularity").Value.String()

	return aws.CostAndUsageRequestType{
		Granularity: interval,
		GroupBy:     dimensions,
		Time: aws.Time{
			Start: start,
			End:   end,
		},
		IsFilterByTagEnabled:       isFilterEnabled(filterByTag),
		IsFilterByDimensionEnabled: isFilterDimensionEnabled(dimensionFilterMap),
		Tag:                        tag,
		TagFilterValue:             filterByTag,
		DimensionFilter:            dimensionFilterMap,
		ExcludeDiscounts:           excludeDiscounts,
	}, nil

}
