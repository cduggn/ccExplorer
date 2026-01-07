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

// NewCostUsageTableWriter creates a writer for cost usage table output
func NewCostUsageTableWriter(sortBy string) *CostUsageTableWriter {
	transformer := NewCostUsageToTableTransformer(sortBy)
	renderer := NewStdoutTableRenderer("costAndUsage")
	return NewCompositeWriter[types.CostAndUsageOutputType, *TableOutput](transformer, renderer)
}

// NewCostUsageCSVWriter creates a writer for cost usage CSV output
func NewCostUsageCSVWriter(sortBy string) *CostUsageCSVWriter {
	transformer := NewCostUsageToCSVTransformer(sortBy)
	renderer := NewCSVRenderer()
	return NewCompositeWriter[types.CostAndUsageOutputType, *CSVOutput](transformer, renderer)
}

// NewCostUsageChartWriter creates a writer for cost usage chart output
func NewCostUsageChartWriter(sortBy string) *CostUsageChartWriter {
	transformer := NewCostUsageToChartTransformer(sortBy)
	renderer := NewChartRenderer()
	return NewCompositeWriter[types.CostAndUsageOutputType, *ChartOutput](transformer, renderer)
}

// NewCostUsageVectorWriter creates a writer for cost usage vector output
func NewCostUsageVectorWriter() *CostUsageVectorWriter {
	transformer := NewCostUsageToVectorTransformer()
	renderer := NewVectorRenderer()
	return NewCompositeWriter[types.CostAndUsageOutputType, *VectorOutput](transformer, renderer)
}

// NewForecastTableWriter creates a writer for forecast table output
func NewForecastTableWriter() *ForecastTableWriter {
	transformer := NewForecastToTableTransformer()
	renderer := NewForecastTableRenderer()
	return NewCompositeWriter[types.ForecastPrintData, *ForecastTableOutput](transformer, renderer)
}