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
	//aws.RightSizingRecommendationS3(context.Background(), awsClient.Client)

	usage, err := awsClient.GetCostAndUsage(context.Background(), awsClient.Client, req)
	if err != nil {
		return err
	}

	report := display.CurateCostAndUsageReport(usage, req.Granularity)
	display.PrintCostAndUsageReport(display.SortServicesByMetricAmount, report)

	return nil
}

func NewCostAndUsageRequest(cmd *cobra.Command) (aws.CostAndUsageRequestType, error) {

	var err error
	groupByValue := cmd.Flags().Lookup("groupBy").Value
	groupBy, _ := groupByValue.(*GroupBy)

	var tag string = ""
	if len(groupBy.Tags) > 0 {
		tag = groupBy.Tags[0]
	}

	filterByValue := cmd.Flags().Lookup("filterBy").Value
	filterBy, _ := filterByValue.(*FilterBy)

	var isFilterByTag bool
	var tagFilter string = ""
	if len(filterBy.Tags) > 0 {
		isFilterByTag = true
		tagFilter = filterBy.Tags[0]
	}

	var isFilterByDimension bool
	if len(filterBy.Dimensions) > 0 {
		isFilterByDimension = true
	}

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
		GroupBy:     groupBy.Dimensions,
		Time: aws.Time{
			Start: start,
			End:   end,
		},
		IsFilterByTagEnabled:       isFilterByTag,
		IsFilterByDimensionEnabled: isFilterByDimension,
		Tag:                        tag,
		TagFilterValue:             tagFilter,
		DimensionFilter:            filterBy.Dimensions,
		ExcludeDiscounts:           excludeDiscounts,
	}, nil

}
