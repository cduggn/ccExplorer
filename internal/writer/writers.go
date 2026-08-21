package writer

import (
	"github.com/cduggn/ccexplorer/internal/types"
)

// CostUsageWriter is the interface for writing cost and usage data
type CostUsageWriter interface {
	Write(data types.CostAndUsageOutputType) error
}

// ForecastWriter is the interface for writing forecast data
type ForecastWriter interface {
	Write(data types.ForecastPrintData) error
}

// Type aliases for specific writer combinations
type costUsageTableWriter = CompositeWriter[types.CostAndUsageOutputType, *TableOutput]
type costUsageCSVWriter = CompositeWriter[types.CostAndUsageOutputType, *CSVOutput]
type costUsageChartWriter = CompositeWriter[types.CostAndUsageOutputType, *ChartOutput]
type costUsageVectorWriter = CompositeWriter[types.CostAndUsageOutputType, *VectorOutput]
type forecastTableWriter = CompositeWriter[types.ForecastPrintData, *ForecastTableOutput]

// NewCostUsageWriter creates a writer for cost and usage data based on print format
func NewCostUsageWriter(printFormat string, sortBy string) CostUsageWriter {
	switch printFormat {
	case "csv":
		return newCostUsageCSVWriter(sortBy)
	case "chart":
		return newCostUsageChartWriter(sortBy)
	case "pinecone":
		return newCostUsageVectorWriter()
	default: // stdout
		return newCostUsageTableWriter(sortBy)
	}
}

// NewForecastWriter creates a writer for forecast data
func NewForecastWriter() ForecastWriter {
	return newForecastTableWriter()
}

// Factory functions for creating specific writer types

func newCostUsageTableWriter(sortBy string) *costUsageTableWriter {
	transformer := NewCostUsageToTableTransformer(sortBy)
	renderer := NewStdoutTableRenderer()
	return NewCompositeWriter[types.CostAndUsageOutputType, *TableOutput](transformer, renderer)
}

func newCostUsageCSVWriter(sortBy string) *costUsageCSVWriter {
	transformer := NewCostUsageToCSVTransformer(sortBy)
	renderer := NewCSVRenderer()
	return NewCompositeWriter[types.CostAndUsageOutputType, *CSVOutput](transformer, renderer)
}

func newCostUsageChartWriter(sortBy string) *costUsageChartWriter {
	transformer := NewCostUsageToChartTransformer(sortBy)
	renderer := NewChartRenderer()
	return NewCompositeWriter[types.CostAndUsageOutputType, *ChartOutput](transformer, renderer)
}

func newCostUsageVectorWriter() *costUsageVectorWriter {
	transformer := NewCostUsageToVectorTransformer()
	renderer := NewVectorRenderer()
	return NewCompositeWriter[types.CostAndUsageOutputType, *VectorOutput](transformer, renderer)
}

func newForecastTableWriter() *forecastTableWriter {
	transformer := NewForecastToTableTransformer()
	renderer := NewForecastTableRenderer()
	return NewCompositeWriter[types.ForecastPrintData, *ForecastTableOutput](transformer, renderer)
}
