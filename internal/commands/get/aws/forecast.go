package aws

import (
	"context"
	"github.com/cduggn/ccexplorer/pkg/printer"
	aws2 "github.com/cduggn/ccexplorer/pkg/service/aws"
	"github.com/spf13/cobra"
)

func CostForecast(cmd *cobra.Command, args []string) error {

	apiClient := aws2.NewAPIClient()

	req, err := NewGetCostForecastRequestType(cmd)
	if err != nil {
		return err
	}
	res, err := apiClient.GetCostForecast(context.TODO(), apiClient.Client, req)
	if err != nil {
		return err
	}
	printer.PrintGetCostForecastReport(res)

	//TODO add error handling

	return nil
}

func NewGetCostForecastRequestType(cmd *cobra.Command) (aws2.GetCostForecastRequest, error) {

	// TODO add error handling
	granularity, _ := cmd.Flags().GetString("granularity")

	predictionIntervalLevel, _ := cmd.Flags().GetInt32("predictionIntervalLevel")
	start := cmd.Flags().Lookup("start").Value.String()
	end := cmd.Flags().Lookup("end").Value.String()

	dimensions, _ := cmd.Flags().GetStringToString("forecast-filter-dimensions")
	err := validateForecastDimensionKey(dimensions)
	if err != nil {
		return aws2.GetCostForecastRequest{}, err
	}
	tags, _ := cmd.Flags().GetStringToString("filter-tags")

	//TODO add validation for all flags

	return aws2.GetCostForecastRequest{
		Granularity:             granularity,
		Metric:                  "UNBLENDED_COST",
		PredictionIntervalLevel: predictionIntervalLevel,
		Time: aws2.Time{
			Start: start,
			End:   end,
		},
		Filter: populateFilter(dimensions, tags),
	}, nil
}

func populateFilter(dimensions map[string]string, tags map[string]string) aws2.Filter {

	if len(dimensions) == 0 && len(tags) == 0 {
		return aws2.Filter{}
	}

	return aws2.Filter{
		Dimensions: createDimensionFilter(dimensions),
		Tags:       createTagFilter(tags),
	}
}

func createDimensionFilter(m map[string]string) []aws2.Dimension {

	if len(m) == 0 {
		return nil
	}

	var dimensions []aws2.Dimension
	for k, v := range m {
		dimensions = append(dimensions, aws2.Dimension{
			Key:   k,
			Value: []string{v},
		})
	}
	return dimensions
}

func createTagFilter(m map[string]string) []aws2.Tag {

	if len(m) == 0 {
		return nil
	}

	var tags []aws2.Tag
	for k, v := range m {
		tags = append(tags, aws2.Tag{
			Key:   k,
			Value: []string{v},
		})
	}
	return tags
}
