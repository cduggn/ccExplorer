package writer

import (
	"github.com/cduggn/ccexplorer/internal/types"
)

// Type aliases for specific writer combinations
type CostUsageTableWriter = CompositeWriter[types.CostAndUsageOutputType, *TableOutput]
type CostUsageCSVWriter = CompositeWriter[types.CostAndUsageOutputType, *CSVOutput]
type CostUsageChartWriter = CompositeWriter[types.CostAndUsageOutputType, *ChartOutput]
type CostUsageVectorWriter = CompositeWriter[types.CostAndUsageOutputType, *VectorOutput]
type ForecastTableWriter = CompositeWriter[types.ForecastPrintData, *ForecastTableOutput]

// Factory functions for creating specific writer types
// All factory functions accept options for configuration

// NewCostUsageTableWriter creates a writer for cost usage table output
func NewCostUsageTableWriter(opts ...Option) *CostUsageTableWriter {
	config := NewConfig(opts...)
	transformer := NewCostUsageToTableTransformer(config.SortBy)
	renderer := NewStdoutTableRenderer("costAndUsage")
	return NewCompositeWriter[types.CostAndUsageOutputType, *TableOutput](transformer, renderer)
}

// NewCostUsageCSVWriter creates a writer for cost usage CSV output
func NewCostUsageCSVWriter(opts ...Option) *CostUsageCSVWriter {
	config := NewConfig(opts...)
	transformer := NewCostUsageToCSVTransformer(config.SortBy)
	renderer := NewCSVRenderer(config)
	return NewCompositeWriter[types.CostAndUsageOutputType, *CSVOutput](transformer, renderer)
}

// NewCostUsageChartWriter creates a writer for cost usage chart output
func NewCostUsageChartWriter(opts ...Option) *CostUsageChartWriter {
	config := NewConfig(opts...)
	transformer := NewCostUsageToChartTransformer(config.SortBy)
	renderer := NewChartRenderer(config)
	return NewCompositeWriter[types.CostAndUsageOutputType, *ChartOutput](transformer, renderer)
}

// NewCostUsageVectorWriter creates a writer for cost usage vector output
func NewCostUsageVectorWriter(opts ...Option) *CostUsageVectorWriter {
	config := NewConfig(opts...)
	transformer := NewCostUsageToVectorTransformer()
	renderer := NewVectorRenderer(config)
	return NewCompositeWriter[types.CostAndUsageOutputType, *VectorOutput](transformer, renderer)
}

// NewForecastTableWriter creates a writer for forecast table output
func NewForecastTableWriter(opts ...Option) *ForecastTableWriter {
	transformer := NewForecastToTableTransformer()
	renderer := NewForecastTableRenderer()
	return NewCompositeWriter[types.ForecastPrintData, *ForecastTableOutput](transformer, renderer)
}