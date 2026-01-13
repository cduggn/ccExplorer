package writer

import (
	"fmt"
	"strings"

	costexplorertypes "github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
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
	
	headers := []string{
		"Rank", "Dimension/Tag", "Dimension/Tag",
		"Metric Name", "Amount", "Rounded",
		"Unit", "Granularity", "Start", "End",
	}
	
	var rows [][]string
	var total float64
	
	// Use generic transformation for creating rows
	for index, service := range sortedServices {
		serviceRows := utils.Transform(service.Metrics, func(metric types.Metrics) []string {
			if metric.Unit == "USD" {
				total += metric.NumericAmount
			}
			
			return []string{
				fmt.Sprintf("%d", index+1),
				service.Keys[0],
				utils.ReturnIfPresent(service.Keys),
				metric.Name,
				metric.Amount,
				fmt.Sprintf("%.2f", metric.NumericAmount),
				metric.Unit,
				input.Granularity,
				service.Start,
				service.End,
			}
		})
		
		// Add periodic divider rows
		if index%10 == 0 && len(rows) > 0 {
			rows = append(rows, make([]string, len(headers)))
		}
		rows = append(rows, serviceRows...)
	}
	
	totalFormatted := fmt.Sprintf("$%.2f", total)
	
	return NewTableOutput(headers, rows, totalFormatted), nil
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
	client := NewVectorStoreClient(nil, input.PineconeIndex, input.PineconeAPIKey, input.OpenAIAPIKey)
	items, err := client.CreateVectorStoreInput(input)
	if err != nil {
		return nil, err
	}
	
	return NewVectorOutput(items, input.PineconeIndex), nil
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