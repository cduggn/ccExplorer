package aws

import (
	"context"
	"github.com/cduggn/ccexplorer/pkg/printer"
	aws2 "github.com/cduggn/ccexplorer/pkg/service/aws"
	"github.com/spf13/cobra"
)

func CostAndUsageSummary(cmd *cobra.Command, args []string) error {

	req, err := NewCostAndUsageRequest(cmd)
	if err != nil {
		return err
	}
	err = ExecuteCostCommand(req)
	if err != nil {
		return err
	}
	return nil
}

func ExecuteCostCommand(q aws2.CostAndUsageRequestType) error {
	awsClient := aws2.NewAPIClient()
	usage, err := awsClient.GetCostAndUsage(context.Background(),
		awsClient.Client, q)
	if err != nil {
		return err
	}

	report := printer.CurateCostAndUsageReport(usage, q.Granularity)
	printer.PrintCostAndUsageReport(printer.SortServicesByMetricAmount, report)
	return nil
}

func NewCostAndUsageRequest(cmd *cobra.Command) (aws2.CostAndUsageRequestType, error) {

	var err error

	// groupBY dimensions and tags
	groupByValues := cmd.Flags().Lookup("groupBy").Value
	groupBy, _ := groupByValues.(*GroupBy)

	// groupBy TAGs
	var groupByTag []string
	if len(groupBy.Tags) > 0 {
		groupByTag = groupBy.Tags
	}

	// filterBY dimensions and tags
	filterByValues := cmd.Flags().Lookup("filterBy").Value
	filterBy, _ := filterByValues.(*FilterBy)

	// check if filter TAGs are set
	var isFilterByTag bool
	var tagFilterValue = ""
	// ( currently only supports one tag )
	if len(filterBy.Tags) > 1 {
		return aws2.CostAndUsageRequestType{}, ValidationError{
			msg: "Currently only supports one TAG filter",
		}
	} else if len(filterBy.Tags) == 1 {
		isFilterByTag = true
		tagFilterValue = filterBy.Tags[0]
	}

	// check if filter DIMENSIONS are set
	var isFilterByDimension bool
	if len(filterBy.Dimensions) > 2 {
		return aws2.CostAndUsageRequestType{}, ValidationError{
			msg: "Currently only supports two DIMENSION filter",
		}
	} else if len(filterBy.Dimensions) > 0 {
		isFilterByDimension = true
	}

	// get start time
	start := cmd.Flags().Lookup("startDate").Value.String()
	err = ValidateStartDate(start)
	if err != nil {
		return aws2.CostAndUsageRequestType{}, err
	}

	// get end time
	end := cmd.Flags().Lookup("endDate").Value.String()
	err = ValidateEndDate(end, start)
	if err != nil {
		return aws2.CostAndUsageRequestType{}, err
	}

	// check if exclude tag is set
	excludeDiscounts, _ := cmd.Flags().GetBool("excludeDiscounts")

	// get granularity
	interval := cmd.Flags().Lookup("granularity").Value.String()

	return aws2.CostAndUsageRequestType{
		Granularity: interval,
		GroupBy:     groupBy.Dimensions,
		Time: aws2.Time{
			Start: start,
			End:   end,
		},
		IsFilterByTagEnabled:       isFilterByTag,
		IsFilterByDimensionEnabled: isFilterByDimension,
		GroupByTag:                 groupByTag,
		TagFilterValue:             tagFilterValue,
		DimensionFilter:            filterBy.Dimensions,
		ExcludeDiscounts:           excludeDiscounts,
	}, nil

}
