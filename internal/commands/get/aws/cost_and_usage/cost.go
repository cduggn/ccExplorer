package cost_and_usage

import (
	"context"
	"github.com/cduggn/ccexplorer/internal/commands/get/aws"
	"github.com/cduggn/ccexplorer/internal/commands/get/aws/custom_flags"
	"github.com/cduggn/ccexplorer/pkg/printer"
	aws2 "github.com/cduggn/ccexplorer/pkg/service/aws"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

func CostAndUsageRunCmd(cmd *cobra.Command, args []string) error {
	userInput, err := handleCommandLineInput(cmd)
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

func handleCommandLineInput(cmd *cobra.Command) (CommandLineInput, error) {

	var err error

	// groupBY dimensions and tags
	groupByValues := cmd.Flags().Lookup("groupBy").Value
	groupBy, _ := groupByValues.(*custom_flags.DimensionAndTagFlag)

	// groupBy TAGs
	var groupByTag []string
	if len(groupBy.Tags) > 0 {
		groupByTag = groupBy.Tags
	}

	// filterBY dimensions and tags
	filterByValues := cmd.Flags().Lookup("filterBy").Value
	filterBy, _ := filterByValues.(*custom_flags.DimensionAndTagFilterFlag)

	// check if filter TAGs are set
	var isFilterByTag bool
	var tagFilterValue = ""
	// ( currently only supports one tag )
	if len(filterBy.Tags) > 1 {
		return CommandLineInput{}, aws.ValidationError{
			Message: "Currently only supports one TAG filter",
		}
	} else if len(filterBy.Tags) == 1 {
		isFilterByTag = true
		tagFilterValue = filterBy.Tags[0]
	}

	// check if filter DIMENSIONS are set
	var isFilterByDimension bool
	if len(filterBy.Dimensions) > 2 {
		return CommandLineInput{}, aws.ValidationError{
			Message: "Currently only supports two DIMENSION filter",
		}
	} else if len(filterBy.Dimensions) > 0 {
		isFilterByDimension = true
	}

	// get start time
	start := cmd.Flags().Lookup("startDate").Value.String()
	err = aws.ValidateStartDate(start)
	if err != nil {
		return CommandLineInput{}, err
	}

	// get end time
	end := cmd.Flags().Lookup("endDate").Value.String()
	err = aws.ValidateEndDate(end, start)
	if err != nil {
		return CommandLineInput{}, err
	}

	// check if exclude tag is set
	excludeDiscounts, _ := cmd.Flags().GetBool("excludeDiscounts")

	// get granularity
	granularity := cmd.Flags().Lookup("granularity").Value.String()
	granularity = strings.ToUpper(granularity)
	isValidGranularity := IsValidGranularity(granularity)
	if !isValidGranularity {
		return CommandLineInput{}, aws.ValidationError{
			Message: "Invalid granularity. Valid values are: DAILY, MONTHLY, HOURLY",
		}
	}

	// get the open_ai_api_key from the viper config
	openAIKey := viper.GetString("open_ai_api_key")

	printFormat := cmd.Flags().Lookup("printFormat").Value.
		String()
	printFormat = strings.ToLower(printFormat)
	isValidPrintFormat := IsValidPrintFormat(printFormat)
	if !isValidPrintFormat {
		return CommandLineInput{}, aws.ValidationError{
			Message: "Invalid print format. " +
				"Please use one of the following: stdout, csv, chart, gpt3",
		}
	}

	if printFormat == "gpt3" && openAIKey == "" {
		return CommandLineInput{}, aws.ValidationError{
			Message: "OpenAI API key not set. Please set the open_ai_api_key in the config file or environment variable",
		}
	}

	metric := cmd.Flags().Lookup("metric").Value.String()
	IsValid := IsValidMetric(metric)
	if !IsValid {
		return CommandLineInput{}, aws.ValidationError{
			Message: "Invalid metric. " +
				"Please use one of the following: AmortizedCost, BlendedCost, NetAmortizedCost, NetUnblendedCost, NormalizedUsageAmount, UnblendedCost, UsageQuantity",
		}
	}

	// check if exclude tag is set
	sortByDate, _ := cmd.Flags().GetBool("sortByDate")

	return CommandLineInput{
		GroupByValues:       groupBy,
		GroupByTag:          groupByTag,
		FilterByValues:      filterBy,
		IsFilterByTag:       isFilterByTag,
		TagFilterValue:      tagFilterValue,
		IsFilterByDimension: isFilterByDimension,
		Start:               start,
		End:                 end,
		ExcludeDiscounts:    excludeDiscounts,
		Interval:            granularity,
		PrintFormat:         printFormat,
		Metrics:             []string{metric},
		SortByDate:          sortByDate,
		OpenAIAPIKey:        openAIKey,
	}, nil

}

func synthesizeRequest(input CommandLineInput) aws2.CostAndUsageRequestType {

	return aws2.CostAndUsageRequestType{
		Granularity: input.Interval,
		GroupBy:     input.GroupByValues.Dimensions,
		Time: aws2.Time{
			Start: input.Start,
			End:   input.End,
		},
		IsFilterByTagEnabled:       input.IsFilterByTag,
		IsFilterByDimensionEnabled: input.IsFilterByDimension,
		GroupByTag:                 input.GroupByTag,
		TagFilterValue:             input.TagFilterValue,
		DimensionFilter:            input.FilterByValues.Dimensions,
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
