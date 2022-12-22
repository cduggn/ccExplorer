package aws

import (
	"context"
	"github.com/cduggn/cloudcost/internal/pkg/service/aws"
	"github.com/spf13/cobra"
)

func CostForecast(cmd *cobra.Command, args []string) error {

	apiClient := aws.NewAPIClient()
	dimensions := GetDimensionValues(apiClient)

	req := NewGetCostForecastRequestType(cmd, dimensions)
	res, _ := apiClient.GetCostForecast(context.TODO(), apiClient.Client, req)

	aws.PrintGetCostForecastReport(res)

	//TODO add error handling

	return nil
}

func NewGetCostForecastRequestType(cmd *cobra.Command, dimensions []string) aws.
	GetCostForecastRequest {

	granularity, _ := cmd.Flags().GetString("granularity")

	predictionIntervalLevel, _ := cmd.Flags().GetInt32("predictionIntervalLevel")
	start := cmd.Flags().Lookup("start").Value.String()
	end := cmd.Flags().Lookup("end").Value.String()
	filterName, _ := cmd.Flags().GetString(
		"forecast-filter")
	filterValue, _ := cmd.Flags().GetString("forecast-filter-value")
	filterType := cmd.Flags().Lookup("forecast-filter-type").Value.String()

	//TODO add validation for all flags

	return aws.GetCostForecastRequest{
		Granularity:             granularity,
		Metric:                  "UNBLENDED_COST",
		PredictionIntervalLevel: predictionIntervalLevel,
		Time: aws.Time{
			Start: start,
			End:   end,
		},
		Filter: generateFilter(filterType, filterName, filterValue),
	}
}

func generateFilter(filterType, filterName, filterValue string) aws.Filter {

	if filterType == "TAG" {
		return aws.Filter{
			Tags: []aws.Tag{
				{
					Key:   filterName,
					Value: []string{filterValue},
				},
			},
		}
	} else if filterType == "DIMENSION" {
		return aws.Filter{
			Dimensions: []aws.Dimension{
				{
					Key:   filterName,
					Value: []string{filterValue},
				},
			},
		}
	}

	return aws.Filter{}
}
