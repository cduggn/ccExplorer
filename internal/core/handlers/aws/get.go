package aws

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/ccexplorer/internal/core/domain"
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	flags "github.com/cduggn/ccexplorer/internal/core/handlers/aws/flags"
	"github.com/cduggn/ccexplorer/internal/core/ports"
	"github.com/cduggn/ccexplorer/internal/core/service/aws"
	"github.com/cduggn/ccexplorer/internal/core/usecases"
	"github.com/cduggn/ccexplorer/internal/core/util"
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
	costUsageGroupBy                flags.DimensionAndTagFlag
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

type PresetCommandType struct {
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
	awsService, err := aws.New()
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
			Example: domain.CostAndUsageExamples,
		},
	}
	costCommand.Cmd.RunE = costCommand.RunE
	costCommand.DefineFlags()
	getCmd.AddCommand(costCommand.Cmd)

	forecastCommand := ForecastCommandType{
		Cmd: &cobra.Command{
			Use:     "forecast",
			Short:   "Return cost and usage forecasts for your account.",
			Example: domain.ForecastExamples,
		},
	}
	forecastCommand.Cmd.RunE = forecastCommand.RunE

	forecastCommand.DefineFlags()
	costCommand.Cmd.AddCommand(forecastCommand.Cmd)
	return getCmd
}

func Presets() *cobra.Command {
	presetCommand := PresetCommandType{
		Cmd: &cobra.Command{
			Use:   "run-query",
			Short: "Predefined AWS Cost and Usage queries",
		},
	}
	presetCommand.Cmd.RunE = presetCommand.RunE
	return presetCommand.Cmd
}

func (c *CostCommandType) DefineFlags() {

	c.Cmd.Flags().VarP(&costUsageGroupBy, "groupBy", "g",
		"Group by DIMENSION and/or TAG ")
	// add required flag for groupBy
	_ = c.Cmd.MarkFlagRequired("groupBy")

	costUsageFilterBy := flags.NewFilterBy()
	c.Cmd.Flags().VarP(&costUsageFilterBy, "filterBy", "f",
		"Filter by DIMENSION and/or TAG")

	// Optional flag to dictate the granularity of the data returned
	c.Cmd.Flags().StringVarP(&costUsageGranularity, "granularity", "m",
		"MONTHLY",
		"Valid values: DAILY, MONTHLY, "+
			"HOURLY. (default: MONTHLY)")

	c.Cmd.Flags().BoolVarP(&costUsageWithoutDiscounts, "excludeDiscounts", "l",
		false,
		"Exclude credit, refunds, "+
			"and discounts (default is to include)")

	c.Cmd.Flags().BoolVarP(&costUsageSortByDate, "sortByDate", "d",
		false,
		"Sort results by date in descending order("+
			"default is to sort by cost in descending order)")

	c.Cmd.Flags().StringVarP(&costUsageStartDate, "startDate", "s",
		util.DefaultStartDate(util.DayOfCurrentMonth, util.SubtractDays),
		"Start date (defaults to the start of the previous month)")
	c.Cmd.Flags().StringVarP(&costUsageEndDate, "endDate", "e",
		util.DefaultEndDate(util.Format),
		"End date *(defaults to the present day)")

	c.Cmd.Flags().StringVarP(&costAndUsagePrintFormat, "printFormat", "p", "stdout",
		"Valid values: stdout, csv, chart, pinecone (default: stdout)")

	c.Cmd.Flags().StringVarP(&costAndUsageMetric, "metric", "i", "UnblendedCost",
		"Valid values: AmortizedCost, BlendedCost, NetAmortizedCost, "+
			"NetUnblendedCost, NormalizedUsageAmount, UnblendedCost, UsageQuantity (default: UnblendedCost)")

}

func (f *ForecastCommandType) DefineFlags() {

	forecastFilterBy := flags.NewForecastFilterBy()
	f.Cmd.Flags().VarP(&forecastFilterBy, "filterBy", "f",
		"Filter by DIMENSION  (default: none)")

	f.Cmd.Flags().StringVarP(&forecastStartDate, "start", "s",
		util.Format(time.Now()), "Start date (defaults to the present day)")

	f.Cmd.Flags().StringVarP(&forecastEndDate, "end", "e",
		util.LastDayOfMonth(),
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

func (c *CostCommandType) InputHandler(validatorFn func(input model.CommandLineInput) error) (model.CommandLineInput, error) {

	groupByTag, groupByDimension := c.ExtractGroupBySelections()

	filterSelection, err := c.ExtractFilterBySelection()
	if err != nil {
		return model.CommandLineInput{}, err
	}

	start, end, err := c.ExtractStartAndEndDates()
	if err != nil {
		return model.CommandLineInput{}, err
	}

	printOptions := c.ExtractPrintPreferences()

	input := model.CommandLineInput{
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
		return model.CommandLineInput{}, err
	}

	return input, nil
}

func (c *CostCommandType) SynthesizeRequest(input model.CommandLineInput) model.
	CostAndUsageRequestType {

	return model.CostAndUsageRequestType{
		Granularity: input.Interval,
		GroupBy:     input.GroupByDimension,
		Time: model.Time{
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

func (c *CostCommandType) Execute(req model.CostAndUsageRequestType) error {

	costAndUsageResponse, err := srv.aws.GetCostAndUsage(
		context.Background(), req)
	if err != nil {
		return err
	}

	report := util.ToCostAndUsageOutputType(costAndUsageResponse, req)

	w := usecases.NewPrintWriter(util.ToPrintWriterType(req.PrintFormat),
		"costAndUsage")

	err = w.Write(util.SortByFn(req.SortByDate), report)
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

	p := usecases.NewPrintWriter(util.ToPrintWriterType("stdout"),
		"forecast")
	err = p.Write(printData, filters)
	if err != nil {
		return err
	}

	return nil
}

func (f *ForecastCommandType) InputHandler() model.ForecastCommandLineInput {
	filterByValues := f.Cmd.Flags().Lookup("filterBy").Value
	granularity, _ := f.Cmd.Flags().GetString("granularity")
	predictionIntervalLevel, _ := f.Cmd.Flags().GetInt32(
		"predictionIntervalLevel")

	filterFlag := filterByValues.(*flags.DimensionFilterByFlag)
	dimensions := aws.ExtractForecastFilters(filterFlag.Dimensions)

	return model.ForecastCommandLineInput{
		FilterByValues:          dimensions,
		Granularity:             granularity,
		PredictionIntervalLevel: predictionIntervalLevel,
		Start:                   f.Cmd.Flags().Lookup("start").Value.String(),
		End:                     f.Cmd.Flags().Lookup("end").Value.String(),
	}
}

func (f *ForecastCommandType) SynthesizeRequest(input model.
	ForecastCommandLineInput) (model.
	GetCostForecastRequest, error) {

	return model.GetCostForecastRequest{
		Granularity:             input.Granularity,
		Metric:                  "UNBLENDED_COST",
		PredictionIntervalLevel: input.PredictionIntervalLevel,
		Time: model.Time{
			Start: input.Start,
			End:   input.End,
		},
		Filter: input.FilterByValues,
	}, nil
}

func (f *ForecastCommandType) Execute(r model.GetCostForecastRequest) (
	*costexplorer.GetCostForecastOutput, error) {

	res, err := srv.aws.GetCostForecast(context.TODO(), r)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *PresetCommandType) RunE(cmd *cobra.Command, args []string) error {
	var presets = domain.PresetList()
	optionsNameList := presetQueryList(presets)

	var prompt = displayPresetPrompt(optionsNameList)
	selection := 0
	err := survey.AskOne(prompt, &selection)
	if err != nil {
		return model.PresetError{
			Msg: fmt.Sprintf("Error during query selection %v\n",
				err),
		}

	}

	selectedOption := selectedPreset(presets, selection)
	apiRequest, err := p.SynthesizeRequest(selectedOption)
	if err != nil {
		return model.PresetError{
			Msg: fmt.Sprintf("Error synthesizing query %v\n",
				err),
		}
	}

	displaySynthesizedPresetQuery(selectedOption)
	err = p.Execute(apiRequest)
	if err != nil {
		return model.PresetError{
			Msg: fmt.Sprintf("Error executing query %v\n",
				err),
		}
	}
	return nil
}

func (p *PresetCommandType) SynthesizeRequest(m model.PresetParams) (model.
	CostAndUsageRequestType,
	error) {
	return model.CostAndUsageRequestType{
		GroupBy:                    m.Dimension,
		DimensionFilter:            m.Filter,
		IsFilterByTagEnabled:       m.FilterByTag,
		IsFilterByDimensionEnabled: m.FilterByDimension,
		Time: model.Time{
			Start: util.DefaultStartDate(util.DayOfCurrentMonth,
				util.SubtractDays),
			End: util.DefaultEndDate(util.Format),
		},
		Granularity:      m.Granularity,
		ExcludeDiscounts: m.ExcludeDiscounts,
		PrintFormat:      m.PrintFormat,
		SortByDate:       false,
		Metrics:          m.Metric,
	}, nil
}

func (p *PresetCommandType) Execute(q model.CostAndUsageRequestType) error {
	err := executePreset(q)
	if err != nil {
		err := model.PresetError{
			Msg: fmt.Sprintf("Error executing preset query %v\n", err),
		}
		return err
	}
	return nil
}

func executePreset(q model.CostAndUsageRequestType) error {
	// relies on the cost command to execute query
	costCommand := CostCommandType{}

	err := costCommand.Execute(q)
	if err != nil {
		return err
	}
	return nil
}

func displayPresetPrompt(o []string) *survey.Select {
	return &survey.Select{
		Message: "Choose a query to execute:",
		Options: o,
	}
}
func presetQueryList(p []model.PresetParams) []string {
	queries := make([]string, len(p))
	for i, preset := range p {
		queries[i] = preset.Alias
	}
	return queries
}

func selectedPreset(p []model.PresetParams, s int) model.PresetParams {
	return p[s]
}

func prepareResponseForRendering(res *costexplorer.
	GetCostForecastOutput) model.ForecastPrintData {
	return model.ForecastPrintData{
		Forecast: res,
	}
}

func filterList(r model.GetCostForecastRequest) []string {
	var dimensions []string
	for _, d := range r.Filter.Dimensions {
		dimensions = append(dimensions, d.Key)
	}
	return dimensions
}

func displaySynthesizedPresetQuery(p model.PresetParams) {
	fmt.Println("")
	fmt.Printf("Synthesized Query: %v \n", p.CommandSyntax)
	fmt.Println("")
}
