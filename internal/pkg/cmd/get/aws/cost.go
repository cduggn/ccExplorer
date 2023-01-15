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

	groupByValue := cmd.Flags().Lookup("groupBy").Value
	groupBy, _ := groupByValue.(*GroupBy)

	var err error

	//tagMap, _ := cmd.Flags().GetStringToString("groupByTag")
	var tag string = ""
	//if len(tagMap) > 0 {
	//	tag, err = ValidateGroupByTag(tagMap)
	//	if err != nil {
	//		return aws.CostAndUsageRequestType{}, err
	//	}
	//}

	if len(groupBy.Tags) > 0 {
		tag = groupBy.Tags[0]
	}

	tagFilterValue, _ := cmd.Flags().GetString("filterByTag")
	err = ValidateTagFilterValue(tagFilterValue, tag)
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}

	dimensionFilterMap, _ := cmd.Flags().GetStringToString(
		"filterByDimension")

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
		IsFilterByTagEnabled:       isFilterEnabled(tagFilterValue),
		IsFilterByDimensionEnabled: isFilterDimensionEnabled(dimensionFilterMap),
		Tag:                        tag,
		TagFilterValue:             tagFilterValue,
		DimensionFilter:            dimensionFilterMap,
		ExcludeDiscounts:           excludeDiscounts,
	}, nil

}
