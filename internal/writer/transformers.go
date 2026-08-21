package writer

import (
	"fmt"
	"strings"

	costexplorertypes "github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/cduggn/ccexplorer/internal/http"
	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/cduggn/ccexplorer/internal/utils"
)

// CostUsageToTableTransformer transforms cost and usage data to table format
type CostUsageToTableTransformer struct {
	sortFunc func(map[int]types.Service) []types.Service
}

// NewCostUsageToTableTransformer creates a new transformer for cost and usage data
func NewCostUsageToTableTransformer(sortBy string) *CostUsageToTableTransformer {
	return &CostUsageToTableTransformer{
		sortFunc: utils.SortFunction(sortBy),
	}
}

// Transform implements the Transformer interface for cost and usage data
func (t *CostUsageToTableTransformer) Transform(input types.CostAndUsageOutputType) (*TableOutput, error) {
	sortedServices := t.sortFunc(input.Services)

	// First pass: calculate total for percentage calculation
	var total float64
	var metricName, unit string
	var hasTagColumn bool

	for _, service := range sortedServices {
		for _, metric := range service.Metrics {
			if metric.Unit == "USD" {
				total += metric.NumericAmount
			}
			if metricName == "" {
				metricName = metric.Name
				unit = metric.Unit
			}
		}
		// Check if tag/dimension column has values
		if !hasTagColumn && len(service.Keys) > 1 && service.Keys[1] != "" {
			hasTagColumn = true
		}
	}

	// Build headers based on whether tag column has data
	var headers []string
	if hasTagColumn {
		headers = []string{"#", "Service", "Tag/Dimension", "Cost", "%", "Start", "End"}
	} else {
		headers = []string{"#", "Service", "Cost", "%", "Start", "End"}
	}

	var rows [][]string
	var startDate, endDate string

	// Second pass: create rows with percentage
	for index, service := range sortedServices {
		// Track date range
		if index == 0 {
			startDate = service.Start
		}
		endDate = service.End

		serviceRows := utils.Transform(service.Metrics, func(metric types.Metrics) []string {
			// Calculate percentage of total
			var pct float64
			if total > 0 {
				pct = (metric.NumericAmount / total) * 100
			}

			if hasTagColumn {
				return []string{
					fmt.Sprintf("%d", index+1),
					service.Keys[0],
					utils.ReturnIfPresent(service.Keys),
					fmt.Sprintf("%.2f", metric.NumericAmount),
					fmt.Sprintf("%.1f%%", pct),
					service.Start,
					service.End,
				}
			}
			return []string{
				fmt.Sprintf("%d", index+1),
				service.Keys[0],
				fmt.Sprintf("%.2f", metric.NumericAmount),
				fmt.Sprintf("%.1f%%", pct),
				service.Start,
				service.End,
			}
		})

		rows = append(rows, serviceRows...)
	}

	totalFormatted := fmt.Sprintf("$%.2f", total)
	dateRange := fmt.Sprintf("%s to %s", startDate, endDate)
	sortedBy := "Cost"

	return NewTableOutput(headers, rows, totalFormatted, input.Granularity, dateRange, metricName, unit, sortedBy, hasTagColumn), nil
}

// CostUsageToCSVTransformer transforms cost and usage data to CSV format
type CostUsageToCSVTransformer struct {
	sortFunc func(map[int]types.Service) []types.Service
}

// NewCostUsageToCSVTransformer creates a new transformer for CSV output
func NewCostUsageToCSVTransformer(sortBy string) *CostUsageToCSVTransformer {
	return &CostUsageToCSVTransformer{
		sortFunc: utils.SortFunction(sortBy),
	}
}

// Transform implements the Transformer interface for CSV output
func (t *CostUsageToCSVTransformer) Transform(input types.CostAndUsageOutputType) (*CSVOutput, error) {
	headers := []string{
		"Dimension/Tag", "Dimension/Tag", "Metric",
		"Granularity", "Start", "End", "USD Amount", "Unit",
	}
	
	rows := utils.ConvertServiceMapToArray(input.Services, input.Granularity)
	
	return NewCSVOutput(headers, rows, "ccexplorer.csv"), nil
}

// CostUsageToChartTransformer transforms cost and usage data to chart format
type CostUsageToChartTransformer struct {
	sortFunc func(map[int]types.Service) []types.Service
	builder  Builder
}

// NewCostUsageToChartTransformer creates a new transformer for chart output
func NewCostUsageToChartTransformer(sortBy string) *CostUsageToChartTransformer {
	return &CostUsageToChartTransformer{
		sortFunc: utils.SortFunction(sortBy),
		builder:  Builder{},
	}
}

// Transform implements the Transformer interface for chart output
func (t *CostUsageToChartTransformer) Transform(input types.CostAndUsageOutputType) (*ChartOutput, error) {
	sortedServices := t.sortFunc(input.Services)
	chartInput := utils.ConvertToChartInputType(input, sortedServices)
	
	page, err := t.builder.NewCharts(chartInput)
	if err != nil {
		return nil, err
	}
	
	return NewChartOutput(page, "Cost and Usage Report", "ccexplorer_chart.html"), nil
}

// CostUsageToVectorTransformer transforms cost and usage data to vector format
type CostUsageToVectorTransformer struct{}

// NewCostUsageToVectorTransformer creates a new transformer for vector output
func NewCostUsageToVectorTransformer() *CostUsageToVectorTransformer {
	return &CostUsageToVectorTransformer{}
}

// Transform implements the Transformer interface for vector output
func (t *CostUsageToVectorTransformer) Transform(input types.CostAndUsageOutputType) (*VectorOutput, error) {
	// NewVectorStoreClient(builder, openAIAPIKey, indexURL, pineconeAPIKey)
	client := NewVectorStoreClient(http.NewRequestBuilder(),
		input.OpenAIAPIKey, input.PineconeIndex, input.PineconeAPIKey)
	items, err := client.CreateVectorStoreInput(input)
	if err != nil {
		return nil, err
	}

	return NewVectorOutput(items, input.PineconeIndex,
		input.PineconeAPIKey, input.OpenAIAPIKey), nil
}

// ForecastToTableTransformer transforms forecast data to table format
type ForecastToTableTransformer struct{}

// NewForecastToTableTransformer creates a new transformer for forecast table output
func NewForecastToTableTransformer() *ForecastToTableTransformer {
	return &ForecastToTableTransformer{}
}

// Transform implements the Transformer interface for forecast data
func (t *ForecastToTableTransformer) Transform(input types.ForecastPrintData) (*ForecastTableOutput, error) {
	headers := []string{
		"Start", "End", "Mean Value",
		"Prediction Interval LowerBound",
		"Prediction Interval UpperBound", "Unit", "Total",
	}
	
	// Use generic transformation for creating rows
	rows := utils.Transform(input.Forecast.ForecastResultsByTime, func(forecast costexplorertypes.ForecastResult) []string {
		return []string{
			*forecast.TimePeriod.Start,
			*forecast.TimePeriod.End,
			*forecast.MeanValue,
			*forecast.PredictionIntervalLowerBound,
			*forecast.PredictionIntervalUpperBound,
		}
	})
	
	total := types.Total{
		Amount: *input.Forecast.Total.Amount,
		Unit:   *input.Forecast.Total.Unit,
	}
	
	filterInfo := strings.Join(input.Filters, " | ")
	
	return NewForecastTableOutput(headers, rows, filterInfo, total), nil
}