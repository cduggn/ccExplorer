package aws

import (
	"context"
	"github.com/cduggn/cloudcost/internal/pkg/service/aws"
	"github.com/spf13/cobra"
)

func CostForecast(cmd *cobra.Command, args []string) error {

	apiClient := aws.NewAPIClient()
	
	req := NewGetCostForecastRequestType(cmd)
	res, _ := apiClient.GetCostForecast(context.TODO(), apiClient.Client, req)

	aws.PrintGetCostForecastReport(res)

	//TODO add error handling

	return nil
}

func NewGetCostForecastRequestType(cmd *cobra.Command) aws.
GetCostForecastRequest {

	granularity, _ := cmd.Flags().GetString("granularity")

	predictionIntervalLevel, _ := cmd.Flags().GetInt32("predictionIntervalLevel")
	start := cmd.Flags().Lookup("start").Value.String()
	end := cmd.Flags().Lookup("end").Value.String()

	dimensions, _ := cmd.Flags().GetStringToString("forecast-filter-dimensions")
	tags, _ := cmd.Flags().GetStringToString("filter-tags")

	//TODO add validation for all flags

	return aws.GetCostForecastRequest{
		Granularity:             granularity,
		Metric:                  "UNBLENDED_COST",
		PredictionIntervalLevel: predictionIntervalLevel,
		Time: aws.Time{
			Start: start,
			End:   end,
		},
		Filter: populateFilter(dimensions, tags),
	}
}

func populateFilter(dimensions map[string]string, tags map[string]string) aws.Filter {
	return aws.Filter{
		Dimensions: createDimensionFilter(dimensions),
		Tags:       createTagFilter(tags),
	}
}

func createDimensionFilter(m map[string]string) []aws.Dimension {
	var dimensions []aws.Dimension
	for k, v := range m {
		dimensions = append(dimensions, aws.Dimension{
			Key:   k,
			Value: []string{v},
		})
	}
	return dimensions
}

func createTagFilter(m map[string]string) []aws.Tag {
	var tags []aws.Tag
	for k, v := range m {
		tags = append(tags, aws.Tag{
			Key:   k,
			Value: []string{v},
		})
	}
	return tags
}
