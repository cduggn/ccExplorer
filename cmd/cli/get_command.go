package cli

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/ccexplorer/internal/flags"
	awsservice "github.com/cduggn/ccexplorer/internal/awsservice"
	gcpservice "github.com/cduggn/ccexplorer/internal/gcpservice"
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
		Short: "Cost and usage summary for AWS and GCP services",
		Long:  paintHeader(),
	}
	// AWS-specific variables
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
	
	// GCP-specific variables
	gcpProjectID                    string
	gcpProjectIDs                   []string
	gcpBillingAccount              string
	gcpOrganizationID              string
	gcpFolderIDs                   []string
	gcpServices                    []string
	gcpRegions                     []string
	gcpSKUs                        []string
	gcpLabels                      map[string]string
	gcpStartDate                   string
	gcpEndDate                     string
	gcpGranularity                 string
	gcpCurrency                    string
	gcpGroupBy                     []string
	gcpIncludeCredits              bool
	gcpIncludeDiscounts            bool
	gcpPrintFormat                 string
	gcpSortByDate                  bool
	gcpCostThreshold               float64
	
	srv                             *service
)

type service struct {
	aws ports.AWSService
	gcp ports.GCPService
}

type CostCommandType struct {
	Cmd *cobra.Command
}

type ForecastCommandType struct {
	Cmd *cobra.Command
}

// GCP command types following the same pattern as AWS
type GCPCommandType struct {
	Cmd *cobra.Command
}

type GCPForecastCommandType struct {
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
	// Initialize AWS service
	awsService, err := awsservice.New()
	if err != nil {
		return &service{}, err
	}
	
	// Initialize GCP service
	gcpConfig := gcpservice.NewConfigFromEnv()
	gcpService, err := gcpservice.NewService(gcpConfig)
	if err != nil {
		// Log GCP service initialization error but continue with AWS only
		// This allows the tool to work even if GCP credentials are not configured
		fmt.Printf("Warning: GCP service initialization failed: %v\n", err)
		fmt.Println("AWS functionality will be available, but GCP commands will not work.")
		gcpService = nil
	}
	
	client := &service{
		aws: awsService,
		gcp: gcpService,
	}
	return client, nil
}

func CostAndForecast() *cobra.Command {
	// AWS Commands
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

	// GCP Commands
	gcpCommand := GCPCommandType{
		Cmd: &cobra.Command{
			Use:   "gcp",
			Short: "Explore cost summaries for Google Cloud Platform",
			Long: `
Command: gcp 
Description: Cost and usage summary for GCP services.

Prerequisites:
- GCP credentials configured via service account or Application Default Credentials (ADC)
- Set GOOGLE_APPLICATION_CREDENTIALS environment variable to service account JSON file path
- Or set GCP_PROJECT_ID environment variable for your project
- Billing API must be enabled for the project`,
			Example: GCPCostExamples,
		},
	}
	gcpCommand.Cmd.RunE = gcpCommand.RunE
	gcpCommand.DefineFlags()
	getCmd.AddCommand(gcpCommand.Cmd)

	gcpForecastCommand := GCPForecastCommandType{
		Cmd: &cobra.Command{
			Use:     "forecast",
			Short:   "Return cost forecasts for GCP services using historical data analysis.",
			Example: GCPForecastExamples,
		},
	}
	gcpForecastCommand.Cmd.RunE = gcpForecastCommand.RunE
	gcpForecastCommand.DefineFlags()
	gcpCommand.Cmd.AddCommand(gcpForecastCommand.Cmd)

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

// ===== GCP Command Implementations =====

// DefineFlags defines all flags for GCP cost analysis commands
func (g *GCPCommandType) DefineFlags() {
	// Required project ID flag
	g.Cmd.Flags().StringVarP(&gcpProjectID, "project", "p", "", 
		"GCP Project ID (required)")
	g.Cmd.MarkFlagRequired("project")

	// Multi-project support
	g.Cmd.Flags().StringSliceVar(&gcpProjectIDs, "projects", []string{}, 
		"Multiple GCP Project IDs for batch analysis")

	// Billing and organization flags
	g.Cmd.Flags().StringVarP(&gcpBillingAccount, "billingAccount", "b", "", 
		"GCP Billing Account ID")
	g.Cmd.Flags().StringVar(&gcpOrganizationID, "organization", "", 
		"GCP Organization ID for enterprise queries")
	g.Cmd.Flags().StringSliceVar(&gcpFolderIDs, "folders", []string{}, 
		"GCP Folder IDs for hierarchical analysis")

	// Service and resource filtering
	g.Cmd.Flags().StringSliceVarP(&gcpServices, "services", "s", []string{}, 
		"Filter by GCP service names (e.g., 'Compute Engine', 'Cloud Storage')")
	g.Cmd.Flags().StringSliceVarP(&gcpRegions, "regions", "r", []string{}, 
		"Filter by GCP regions (e.g., 'us-central1', 'europe-west1')")
	g.Cmd.Flags().StringSliceVar(&gcpSKUs, "skus", []string{}, 
		"Filter by specific GCP SKUs")

	// Time range flags
	g.Cmd.Flags().StringVar(&gcpStartDate, "startDate", 
		utils.DefaultStartDate(utils.DayOfCurrentMonth, utils.SubtractDays),
		"Start date (defaults to the start of the previous month)")
	g.Cmd.Flags().StringVar(&gcpEndDate, "endDate", 
		utils.DefaultEndDate(utils.Format),
		"End date (defaults to the present day)")

	// Granularity and grouping
	g.Cmd.Flags().StringVarP(&gcpGranularity, "granularity", "g", "MONTHLY",
		"Time granularity: DAILY, MONTHLY, HOURLY (default: MONTHLY)")
	g.Cmd.Flags().StringSliceVar(&gcpGroupBy, "groupBy", []string{},
		"Group results by: service, project, region, currency")

	// Currency and cost options
	g.Cmd.Flags().StringVarP(&gcpCurrency, "currency", "c", "USD",
		"Currency for cost reporting (default: USD)")
	g.Cmd.Flags().Float64Var(&gcpCostThreshold, "costThreshold", 0.01,
		"Minimum cost threshold for filtering (default: 0.01)")

	// Credits and discounts
	g.Cmd.Flags().BoolVar(&gcpIncludeCredits, "includeCredits", true,
		"Include promotional credits in analysis (default: true)")
	g.Cmd.Flags().BoolVar(&gcpIncludeDiscounts, "includeDiscounts", true,
		"Include discounts in analysis (default: true)")

	// Output options
	g.Cmd.Flags().StringVar(&gcpPrintFormat, "printFormat", "stdout",
		"Output format: stdout, csv, chart, pinecone (default: stdout)")
	g.Cmd.Flags().BoolVar(&gcpSortByDate, "sortByDate", false,
		"Sort results by date instead of cost (default: false)")
}

// RunE executes the GCP cost analysis command
func (g *GCPCommandType) RunE(cmd *cobra.Command, args []string) error {
	// Check if GCP service is available
	if srv.gcp == nil {
		return fmt.Errorf("GCP service not available. Please check your GCP credentials and configuration")
	}

	userInput, err := g.InputHandler()
	if err != nil {
		return err
	}

	req := g.SynthesizeRequest(userInput)
	
	// Validate request
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	err = g.Execute(req)
	if err != nil {
		return err
	}
	return nil
}

// InputHandler processes command-line input for GCP commands
func (g *GCPCommandType) InputHandler() (types.GCPCommandLineInput, error) {
	input := types.GCPCommandLineInput{
		ProjectID:          gcpProjectID,
		ProjectIDs:         gcpProjectIDs,
		BillingAccount:     gcpBillingAccount,
		OrganizationID:     gcpOrganizationID,
		FolderIDs:          gcpFolderIDs,
		Services:           gcpServices,
		Regions:            gcpRegions,
		SKUs:               gcpSKUs,
		Labels:             gcpLabels, // Will be populated if label flags are added
		Start:              gcpStartDate,
		End:                gcpEndDate,
		Granularity:        gcpGranularity,
		Currency:           gcpCurrency,
		GroupBy:            gcpGroupBy,
		IncludeCredits:     gcpIncludeCredits,
		IncludeDiscounts:   gcpIncludeDiscounts,
		PrintFormat:        gcpPrintFormat,
		SortByDate:         gcpSortByDate,
		CostThreshold:      gcpCostThreshold,
		// OpenAI and Pinecone keys would be populated from environment or flags
	}

	// Validate input
	if input.ProjectID == "" && len(input.ProjectIDs) == 0 {
		return input, fmt.Errorf("at least one project ID must be specified")
	}

	return input, nil
}

// SynthesizeRequest converts command-line input to GCP billing request
func (g *GCPCommandType) SynthesizeRequest(input types.GCPCommandLineInput) types.GCPBillingRequest {
	return types.GCPBillingRequest{
		ProjectID:          input.ProjectID,
		ProjectIDs:         input.ProjectIDs,
		BillingAccount:     input.BillingAccount,
		OrganizationID:     input.OrganizationID,
		FolderIDs:          input.FolderIDs,
		Services:           input.Services,
		Regions:            input.Regions,
		SKUs:               input.SKUs,
		Labels:             input.Labels,
		Time: types.GCPTimeRange{
			Start: input.Start,
			End:   input.End,
		},
		Granularity:        input.Granularity,
		Currency:           input.Currency,
		GroupBy:            input.GroupBy,
		IncludeCredits:     input.IncludeCredits,
		IncludeDiscounts:   input.IncludeDiscounts,
		CostThreshold:      input.CostThreshold,
		PrintFormat:        input.PrintFormat,
		SortByDate:         input.SortByDate,
		OpenAIAPIKey:       input.OpenAIAPIKey,
		PineconeIndex:      input.PineconeIndex,
		PineconeAPIKey:     input.PineconeAPIKey,
	}
}

// Execute runs the GCP billing data query and formats output
func (g *GCPCommandType) Execute(req types.GCPBillingRequest) error {
	// Get billing data from GCP
	response, err := srv.gcp.GetBillingData(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to get GCP billing data: %w", err)
	}

	// Transform to output format (this would need to be implemented)
	// For now, we'll create a placeholder transformation
	outputData := g.transformGCPResponse(response, req)

	// Create writer and output results
	w := writer.NewPrintWriter(utils.ToPrintWriterType(req.PrintFormat), "gcpCostAndUsage")
	
	// Create a sorting function based on user preference
	sortFn := func(data interface{}) interface{} {
		// This would implement GCP-specific sorting logic
		return data
	}

	err = w.Write(sortFn, outputData)
	if err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

// transformGCPResponse transforms GCP billing response to output format
// This is a placeholder - would need full implementation based on writer expectations
func (g *GCPCommandType) transformGCPResponse(response *types.GCPBillingResponse, req types.GCPBillingRequest) interface{} {
	// This would transform the GCP response to match the expected output format
	// For now, return the response as-is
	return response
}

// ===== GCP Forecast Command Implementation =====

// DefineFlags defines flags for GCP forecast commands
func (gf *GCPForecastCommandType) DefineFlags() {
	// Reuse many of the same flags as the main GCP command
	gf.Cmd.Flags().StringVarP(&gcpProjectID, "project", "p", "", 
		"GCP Project ID for forecasting (required)")
	gf.Cmd.MarkFlagRequired("project")

	gf.Cmd.Flags().StringSliceVar(&gcpServices, "services", []string{}, 
		"Filter forecast by GCP service names")
	gf.Cmd.Flags().StringSliceVar(&gcpRegions, "regions", []string{}, 
		"Filter forecast by GCP regions")

	// Forecast-specific flags
	gf.Cmd.Flags().StringVar(&gcpStartDate, "start", 
		utils.DefaultStartDate(utils.DayOfCurrentMonth, utils.SubtractDays),
		"Historical data start date")
	gf.Cmd.Flags().StringVar(&gcpEndDate, "end", 
		utils.DefaultEndDate(utils.Format),
		"Historical data end date")

	gf.Cmd.Flags().StringVar(&gcpGranularity, "granularity", "MONTHLY",
		"Forecast granularity: DAILY, MONTHLY (default: MONTHLY)")
	gf.Cmd.Flags().StringVar(&gcpCurrency, "currency", "USD",
		"Currency for forecast (default: USD)")
}

// RunE executes the GCP forecast command
func (gf *GCPForecastCommandType) RunE(cmd *cobra.Command, args []string) error {
	if srv.gcp == nil {
		return fmt.Errorf("GCP service not available. Please check your GCP credentials and configuration")
	}

	userInput := gf.InputHandler()
	req := gf.SynthesizeRequest(userInput)

	response, err := gf.Execute(req)
	if err != nil {
		return err
	}

	// Output forecast results
	w := writer.NewPrintWriter(utils.ToPrintWriterType("stdout"), "gcpForecast")
	err = w.Write(response, nil)
	if err != nil {
		return fmt.Errorf("failed to write forecast output: %w", err)
	}

	return nil
}

// InputHandler processes input for GCP forecast commands
func (gf *GCPForecastCommandType) InputHandler() types.GCPCommandLineInput {
	return types.GCPCommandLineInput{
		ProjectID:   gcpProjectID,
		Services:    gcpServices,
		Regions:     gcpRegions,
		Start:       gcpStartDate,
		End:         gcpEndDate,
		Granularity: gcpGranularity,
		Currency:    gcpCurrency,
	}
}

// SynthesizeRequest creates a GCP forecast request
func (gf *GCPForecastCommandType) SynthesizeRequest(input types.GCPCommandLineInput) types.GCPForecastRequest {
	return types.GCPForecastRequest{
		ProjectID:   input.ProjectID,
		Services:    input.Services,
		Regions:     input.Regions,
		TimeRange: types.GCPTimeRange{
			Start: input.Start,
			End:   input.End,
		},
		ForecastPeriod: types.GCPTimeRange{
			Start: input.End, // Forecast starts where historical data ends
			End:   utils.LastDayOfMonth(), // Default to end of current month
		},
		Granularity:     input.Granularity,
		Currency:        input.Currency,
		ConfidenceLevel: 0.95, // Default 95% confidence level
		IncludeCredits:  true,
		IncludeDiscounts: true,
		ModelType:       "linear_regression", // Default model type
	}
}

// Execute runs the GCP forecast analysis
func (gf *GCPForecastCommandType) Execute(req types.GCPForecastRequest) (*types.GCPForecastResponse, error) {
	response, err := srv.gcp.GetCostForecast(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to get GCP cost forecast: %w", err)
	}
	return response, nil
}

// ===== Command Examples =====

const GCPCostExamples = `
  # Get cost summary for a single project
  ccexplorer get gcp --project my-project-id

  # Get cost summary for multiple projects
  ccexplorer get gcp --project my-project-id --projects project1,project2,project3

  # Filter by specific services
  ccexplorer get gcp --project my-project-id --services "Compute Engine,Cloud Storage"

  # Filter by regions
  ccexplorer get gcp --project my-project-id --regions us-central1,europe-west1

  # Get daily cost breakdown for last 30 days
  ccexplorer get gcp --project my-project-id --granularity DAILY --startDate 2024-06-01 --endDate 2024-06-30

  # Group results by service and region
  ccexplorer get gcp --project my-project-id --groupBy service,region

  # Export to CSV
  ccexplorer get gcp --project my-project-id --printFormat csv

  # Filter by cost threshold (only show costs > $10)
  ccexplorer get gcp --project my-project-id --costThreshold 10.00

  # Exclude credits and discounts from analysis
  ccexplorer get gcp --project my-project-id --includeCredits=false --includeDiscounts=false

  # Organization-wide cost analysis
  ccexplorer get gcp --organization 123456789 --granularity MONTHLY
`

const GCPForecastExamples = `
  # Get cost forecast for next month
  ccexplorer get gcp forecast --project my-project-id

  # Get forecast filtered by specific services
  ccexplorer get gcp forecast --project my-project-id --services "Compute Engine,BigQuery"

  # Get daily forecast with custom date range
  ccexplorer get gcp forecast --project my-project-id --granularity DAILY --start 2024-01-01 --end 2024-03-31

  # Regional forecast analysis
  ccexplorer get gcp forecast --project my-project-id --regions us-central1,europe-west1
`
