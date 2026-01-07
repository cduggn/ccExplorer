package cli

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/ccexplorer/internal/flags"
	awsservice "github.com/cduggn/ccexplorer/internal/awsservice"
	"github.com/cduggn/ccexplorer/internal/ports"
	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/cduggn/ccexplorer/internal/utils"
	"github.com/cduggn/ccexplorer/internal/writer"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"time"
)

var (
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Cost and usage summary for AWS services",
		Long:  paintHeader(),
	}
	costUsageGroupBy                *flags.GroupByFlag
	costUsageGranularity            string
	costUsageStartDate              string
	costUsageEndDate                string
	costUsageWithoutDiscounts       bool
	costAndUsagePrintFormat         string
	costAndUsageMetric              string
	costUsageSortByDate             bool
	forecastStartDate               string
	forecastEndDate                 string
	forecastGranularity             string
	forecastPredictionIntervalLevel int32
	srv                             *service
)

type service struct {
	aws ports.AWSService
}

type CostCommandType struct {
	Cmd *cobra.Command
}

type ForecastCommandType struct {
	Cmd *cobra.Command
}

func Initialize() {
	var err error
	srv, err = configureServices()
	if err != nil {
		panic(err.Error())
	}
}

func configureServices() (*service, error) {
	awsService, err := awsservice.New()
	if err != nil {
		return &service{}, err
	}
	awsClient := &service{
		aws: awsService,
	}
	return awsClient, nil
}

func CostAndForecast() *cobra.Command {
	costCommand := CostCommandType{
		Cmd: &cobra.Command{
			Use:   "aws",
			Short: "Explore UNBLENDED cost summaries for AWS",
			Long: `
Command: aws 
Description: Cost and usage summary for AWS services.

Prerequisites:
- AWS credentials configured in ~/.aws/credentials and default region configured in ~/.aws/config. Alternatively, 
you can set the environment variables AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY and AWS_REGION.`,
			Example: CostAndUsageExamples,
		},
	}
	costCommand.Cmd.RunE = costCommand.RunE
	costCommand.DefineFlags()
	getCmd.AddCommand(costCommand.Cmd)

	forecastCommand := ForecastCommandType{
		Cmd: &cobra.Command{
			Use:     "forecast",
			Short:   "Return cost and usage forecasts for your account.",
			Example: ForecastExamples,
		},
	}
	forecastCommand.Cmd.RunE = forecastCommand.RunE

	forecastCommand.DefineFlags()
	costCommand.Cmd.AddCommand(forecastCommand.Cmd)
	return getCmd
}

func (c *CostCommandType) DefineFlags() {
	costUsageGroupBy = flags.NewGroupByFlag()
	c.Cmd.Flags().VarP(costUsageGroupBy, "groupBy", "g",
		"Group by DIMENSION and/or TAG ")
	// add required flag for groupBy
	_ = c.Cmd.MarkFlagRequired("groupBy")

	costUsageFilterBy := flags.NewFilterByFlag()
	c.Cmd.Flags().VarP(costUsageFilterBy, "filterBy", "f",
		"Filter by DIMENSION and/or TAG")

	// Optional flag to dictate the granularity of the data returned
	c.Cmd.Flags().StringVarP(&costUsageGranularity, "granularity", "m",
		"MONTHLY",
		"Valid values: DAILY, MONTHLY, "+
			"HOURLY. (default: MONTHLY)")

	c.Cmd.Flags().BoolVarP(&costUsageWithoutDiscounts, "excludeDiscounts", "l",
		false,
		"Excludes the following charge categories: Credit, Refund, Discount, BundledDiscount, DiscountedUsage, SavingsPlanCoveredUsage, SavingsPlanNegation . ( Exclusions not enabled by default)")

	c.Cmd.Flags().BoolVarP(&costUsageSortByDate, "sortByDate", "d",
		false,
		"Sort results by date in descending order("+
			"default is to sort by cost in descending order)")

	c.Cmd.Flags().StringVarP(&costUsageStartDate, "startDate", "s",
		utils.DefaultStartDate(utils.DayOfCurrentMonth, utils.SubtractDays),
		"Start date (defaults to the start of the previous month)")
	c.Cmd.Flags().StringVarP(&costUsageEndDate, "endDate", "e",
		utils.DefaultEndDate(utils.Format),
		"End date *(defaults to the present day)")

	c.Cmd.Flags().StringVarP(&costAndUsagePrintFormat, "printFormat", "p", "stdout",
		"Valid values: stdout, csv, chart, pinecone (default: stdout)")

	c.Cmd.Flags().StringVarP(&costAndUsageMetric, "metric", "i", "UnblendedCost",
		"Valid values: AmortizedCost, BlendedCost, NetAmortizedCost, "+
			"NetUnblendedCost, NormalizedUsageAmount, UnblendedCost, UsageQuantity (default: UnblendedCost)")

}

func (f *ForecastCommandType) DefineFlags() {

	forecastFilterBy := flags.NewDimensionFilterFlag()
	f.Cmd.Flags().VarP(forecastFilterBy, "filterBy", "f",
		"Filter by DIMENSION  (default: none)")

	f.Cmd.Flags().StringVarP(&forecastStartDate, "start", "s",
		utils.Format(time.Now()), "Start date (defaults to the present day)")

	f.Cmd.Flags().StringVarP(&forecastEndDate, "end", "e",
		utils.LastDayOfMonth(),
		"End date (defaults to one month from the start date)")

	// Optional flag to dictate the granularity of the data returned
	f.Cmd.Flags().StringVarP(&forecastGranularity, "granularity", "g", "MONTHLY",
		"Valid values: DAILY, MONTHLY, HOURLY (default: MONTHLY)")

	f.Cmd.Flags().Int32VarP(&forecastPredictionIntervalLevel, "predictionIntervalLevel",
		"p", 95, "Prediction interval level (default: 95)")
}

func paintHeader() string {
	myFigure := figure.NewFigure("CostAndUsage", "thin", true)
	return myFigure.String()
}

func (c *CostCommandType) RunE(cmd *cobra.Command,
	args []string) error {
	userInput, err := c.InputHandler(ValidateInput)
	if err != nil {
		return err
	}

	req := c.SynthesizeRequest(userInput)
	err = c.Execute(req)
	if err != nil {
		return err
	}
	return nil
}

func (c *CostCommandType) InputHandler(validatorFn func(input types.CommandLineInput) error) (types.CommandLineInput, error) {

	groupByTag, groupByDimension := c.ExtractGroupBySelections()

	filterSelection, err := c.ExtractFilterBySelection()
	if err != nil {
		return types.CommandLineInput{}, err
	}

	start, end, err := c.ExtractStartAndEndDates()
	if err != nil {
		return types.CommandLineInput{}, err
	}

	printOptions := c.ExtractPrintPreferences()

	input := types.CommandLineInput{
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
		PineconeAPIKey:      printOptions.PineconeAPIKey,
		PineconeIndex:       printOptions.PineconeIndex,
	}

	err = validatorFn(input)
	if err != nil {
		return types.CommandLineInput{}, err
	}

	return input, nil
}

func (c *CostCommandType) SynthesizeRequest(input types.CommandLineInput) types.CostAndUsageRequestType {

	return types.CostAndUsageRequestType{
		Granularity: input.Interval,
		GroupBy:     input.GroupByDimension,
		Time: types.Time{
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
		PineconeAPIKey:             input.PineconeAPIKey,
		PineconeIndex:              input.PineconeIndex,
	}
}

func (c *CostCommandType) Execute(req types.CostAndUsageRequestType) error {

	costAndUsageResponse, err := srv.aws.GetCostAndUsage(
		context.Background(), req)
	if err != nil {
		return err
	}

	report := utils.ToCostAndUsageOutputType(costAndUsageResponse, req)

	w := writer.NewPrintWriter(utils.ToPrintWriterType(req.PrintFormat),
		"costAndUsage")

	err = w.Write(utils.SortByFn(req.SortByDate), report)
	if err != nil {
		return err
	}
	return nil
}

func (f *ForecastCommandType) RunE(cmd *cobra.Command, args []string) error {

	userInput := f.InputHandler()
	req, err := f.SynthesizeRequest(userInput)
	if err != nil {
		return err
	}

	res, err := f.Execute(req)
	if err != nil {
		return err
	}

	printData := prepareResponseForRendering(res)
	filters := filterList(req)
	printData.Filters = filters

	p := writer.NewPrintWriter(utils.ToPrintWriterType("stdout"),
		"forecast")
	err = p.Write(printData, filters)
	if err != nil {
		return err
	}

	return nil
}

func (f *ForecastCommandType) InputHandler() types.ForecastCommandLineInput {
	filterByValues := f.Cmd.Flags().Lookup("filterBy").Value
	granularity, _ := f.Cmd.Flags().GetString("granularity")
	predictionIntervalLevel, _ := f.Cmd.Flags().GetInt32(
		"predictionIntervalLevel")

	filterFlag := filterByValues.(*flags.DimensionFilterFlag)
	filterData := filterFlag.Value()
	dimensions := awsservice.ExtractForecastFilters(filterData)

	return types.ForecastCommandLineInput{
		FilterByValues:          dimensions,
		Granularity:             granularity,
		PredictionIntervalLevel: predictionIntervalLevel,
		Start:                   f.Cmd.Flags().Lookup("start").Value.String(),
		End:                     f.Cmd.Flags().Lookup("end").Value.String(),
	}
}

func (f *ForecastCommandType) SynthesizeRequest(input types.ForecastCommandLineInput) (types.GetCostForecastRequest, error) {

	return types.GetCostForecastRequest{
		Granularity:             input.Granularity,
		Metric:                  "UNBLENDED_COST",
		PredictionIntervalLevel: input.PredictionIntervalLevel,
		Time: types.Time{
			Start: input.Start,
			End:   input.End,
		},
		Filter: input.FilterByValues,
	}, nil
}

func (f *ForecastCommandType) Execute(r types.GetCostForecastRequest) (
	*costexplorer.GetCostForecastOutput, error) {

	res, err := srv.aws.GetCostForecast(context.TODO(), r)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func prepareResponseForRendering(res *costexplorer.
	GetCostForecastOutput) types.ForecastPrintData {
	return types.ForecastPrintData{
		Forecast: res,
	}
}

func filterList(r types.GetCostForecastRequest) []string {
	var dimensions []string
	for _, d := range r.Filter.Dimensions {
		dimensions = append(dimensions, d.Key)
	}
	return dimensions
}
