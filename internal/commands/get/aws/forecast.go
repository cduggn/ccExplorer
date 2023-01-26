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

	var dimensions []string
	for _, d := range req.Filter.Dimensions {
		dimensions = append(dimensions, d.Key)
	}
	printer.PrintGetCostForecastReport(res, dimensions)
	return nil
}

func NewGetCostForecastRequestType(cmd *cobra.Command) (aws2.GetCostForecastRequest, error) {

	filterByValues := cmd.Flags().Lookup("filterBy").Value
	filterBy, _ := filterByValues.(*ForecastFilterBy)

	granularity, _ := cmd.Flags().GetString("granularity")

	predictionIntervalLevel, _ := cmd.Flags().GetInt32("predictionIntervalLevel")

	start := cmd.Flags().Lookup("start").Value.String()

	end := cmd.Flags().Lookup("end").Value.String()

	return aws2.GetCostForecastRequest{
		Granularity:             granularity,
		Metric:                  "UNBLENDED_COST",
		PredictionIntervalLevel: predictionIntervalLevel,
		Time: aws2.Time{
			Start: start,
			End:   end,
		},
		Filter: aws2.ExtractForecastFilters(filterBy.Dimensions),
	}, nil
}
