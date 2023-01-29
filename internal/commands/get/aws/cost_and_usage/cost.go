package cost_and_usage

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/ccexplorer/internal/commands/get/aws"
	"github.com/cduggn/ccexplorer/internal/commands/get/aws/custom_flags"
	"github.com/cduggn/ccexplorer/pkg/printer"
	aws2 "github.com/cduggn/ccexplorer/pkg/service/aws"
	"github.com/spf13/cobra"
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

	printFormat := cmd.Flags().Lookup("printFormat").Value.
		String()
	printFormat = strings.ToLower(printFormat)
	isValidPrintFormat := IsValidPrintFormat(printFormat)
	if !isValidPrintFormat {
		return CommandLineInput{}, aws.ValidationError{
			Message: "Invalid print format. " +
				"Please use one of the following: stdout, csv, chart",
		}
	}

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
	}
}

func ExecutePreset(q aws2.CostAndUsageRequestType) error {
	err := execute(q)
	if err != nil {
		return err
	}
	return nil
}

func execute(q aws2.CostAndUsageRequestType) error {
	awsClient := aws2.NewAPIClient()
	costAndUsage, err := awsClient.GetCostAndUsage(context.Background(), awsClient.Client, q)
	if err != nil {
		return err
	}

	printData := prepareResponseForRendering(costAndUsage)
	printData.Granularity = q.Granularity

	report := printer.CurateCostAndUsageReport(printData)

	p := printer.PrintFactory(printer.ToPrintWriterType(q.PrintFormat),
		"costAndUsage")
	p.Print(printer.SortServicesByMetricAmount, report)

	return nil
}

func prepareResponseForRendering(r *costexplorer.GetCostAndUsageOutput) printer.CostAndUsageReportPrintData {
	return printer.CostAndUsageReportPrintData{
		Report: r,
	}
}
