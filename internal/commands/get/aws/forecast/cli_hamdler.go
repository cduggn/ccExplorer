package forecast

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/ccexplorer/internal/commands/get/aws/custom_flags"
	"github.com/cduggn/ccexplorer/pkg/domain/aws"
	aws2 "github.com/cduggn/ccexplorer/pkg/domain/model"
	"github.com/cduggn/ccexplorer/pkg/presentation"
	formats2 "github.com/cduggn/ccexplorer/pkg/presentation/writers"
	"github.com/spf13/cobra"
)

func CostForecastRunCmd(cmd *cobra.Command, args []string) error {

	userInput := handleCommandLineInput(cmd)
	req, err := synthesizeRequest(userInput)
	if err != nil {
		return err
	}

	res, err := execute(req)
	if err != nil {
		return err
	}

	printData := prepareResponseForRendering(res)
	filters := filterList(req)
	printData.Filters = filters

	p := presentation.NewPrintWriter(formats2.ToPrintWriterType("stdout"), "forecast")
	err = p.Write(printData, filters)
	if err != nil {
		return err
	}

	return nil
}

func prepareResponseForRendering(res *costexplorer.
GetCostForecastOutput) formats2.ForecastPrintData {
	return formats2.ForecastPrintData{
		Forecast: res,
	}
}

func execute(r aws2.GetCostForecastRequest) (*costexplorer.GetCostForecastOutput, error) {
	apiClient := aws.NewAPIClient()
	res, err := apiClient.GetCostForecast(context.TODO(), apiClient.Client, r)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func handleCommandLineInput(cmd *cobra.Command) CommandLineInput {

	filterByValues := cmd.Flags().Lookup("filterBy").Value
	granularity, _ := cmd.Flags().GetString("granularity")
	predictionIntervalLevel, _ := cmd.Flags().GetInt32("predictionIntervalLevel")

	filterFlag := filterByValues.(*custom_flags.DimensionFilterByFlag)
	dimensions := aws.ExtractForecastFilters(filterFlag.Dimensions)

	return CommandLineInput{
		FilterByValues:          dimensions,
		Granularity:             granularity,
		PredictionIntervalLevel: predictionIntervalLevel,
		Start:                   cmd.Flags().Lookup("start").Value.String(),
		End:                     cmd.Flags().Lookup("end").Value.String(),
	}
}

func filterList(r aws2.GetCostForecastRequest) []string {
	var dimensions []string
	for _, d := range r.Filter.Dimensions {
		dimensions = append(dimensions, d.Key)
	}
	return dimensions
}

func synthesizeRequest(input CommandLineInput) (aws2.GetCostForecastRequest,
	error) {

	return aws2.GetCostForecastRequest{
		Granularity:             input.Granularity,
		Metric:                  "UNBLENDED_COST",
		PredictionIntervalLevel: input.PredictionIntervalLevel,
		Time: aws2.Time{
			Start: input.Start,
			End:   input.End,
		},
		Filter: input.FilterByValues,
	}, nil
}
