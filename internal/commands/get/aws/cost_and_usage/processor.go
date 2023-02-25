package cost_and_usage

import (
	"context"
	"github.com/cduggn/ccexplorer/pkg/printer"
	aws2 "github.com/cduggn/ccexplorer/pkg/service/aws"
	"github.com/spf13/cobra"
)

func CostAndUsageRunCmd(cmd *cobra.Command, args []string) error {
	userInput, err := handleCommandLineInput(ValidateInput, cmd)
	if err != nil {
		return err
	}

	req := synthesizeRequest(userInput)

	err = execute(req)
	if err != nil {
		return err
	}

	return nil
}

func handleCommandLineInput(validatorFn func(input CommandLineInput) error,
	cmd *cobra.Command) (CommandLineInput, error) {

	groupByTag, groupByDimension := ExtractGroupBySelections(cmd)

	filterSelection, err := ExtractFilterBySelection(cmd)
	if err != nil {
		return CommandLineInput{}, err
	}

	start, end, err := ExtractStartAndEndDates(cmd)
	if err != nil {
		return CommandLineInput{}, err
	}

	printOptions := ExtractPrintPreferences(cmd)

	input := CommandLineInput{
		GroupByDimension:    groupByDimension,
		GroupByTag:          groupByTag,
		FilterByValues:      filterSelection.Dimensions,
		IsFilterByTag:       filterSelection.IsFilterByTag,
		TagFilterValue:      filterSelection.Tags,
		IsFilterByDimension: filterSelection.IsFilterByDimension,
		Start:               start,
		End:                 end,
		ExcludeDiscounts:    printOptions.ExcludeDiscounts,
		Interval:            printOptions.Granularity,
		PrintFormat:         printOptions.Format,
		Metrics:             []string{printOptions.Metric},
		SortByDate:          printOptions.IsSortByDate,
		OpenAIAPIKey:        printOptions.OpenAIKey,
	}

	err = validatorFn(input)
	if err != nil {
		return CommandLineInput{}, err
	}

	return input, nil
}

func synthesizeRequest(input CommandLineInput) aws2.CostAndUsageRequestType {

	return aws2.CostAndUsageRequestType{
		Granularity: input.Interval,
		GroupBy:     input.GroupByDimension,
		Time: aws2.Time{
			Start: input.Start,
			End:   input.End,
		},
		IsFilterByTagEnabled:       input.IsFilterByTag,
		IsFilterByDimensionEnabled: input.IsFilterByDimension,
		GroupByTag:                 input.GroupByTag,
		TagFilterValue:             input.TagFilterValue,
		DimensionFilter:            input.FilterByValues,
		ExcludeDiscounts:           input.ExcludeDiscounts,
		PrintFormat:                input.PrintFormat,
		Metrics:                    input.Metrics,
		SortByDate:                 input.SortByDate,
		OpenAIAPIKey:               input.OpenAIAPIKey,
	}
}

func ExecutePreset(q aws2.CostAndUsageRequestType) error {
	err := execute(q)
	if err != nil {
		return err
	}
	return nil
}

func execute(req aws2.CostAndUsageRequestType) error {
	awsClient := aws2.NewAPIClient()
	costAndUsageResponse, err := awsClient.GetCostAndUsage(context.
		Background(), awsClient.Client, req)
	if err != nil {
		return err
	}

	report := printer.ToCostAndUsageOutputType(costAndUsageResponse, req)
	p := printer.PrintFactory(printer.ToPrintWriterType(req.PrintFormat),
		"costAndUsage")
	err = p.Print(SortByFn(req.SortByDate), report)
	if err != nil {
		return err
	}
	return nil
}
