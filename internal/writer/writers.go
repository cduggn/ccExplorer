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

// Legacy compatibility types - these wrap the new generic writers to maintain the old interface
type GenericStdoutPrinter struct {
	variant string
}

type GenericCsvPrinter struct {
	variant string
}

type GenericChartPrinter struct {
	variant string
}

type GenericPineconePrinter struct {
	variant string
}

// NewGenericStdoutPrinter creates a new stdout printer with backward compatibility
func NewGenericStdoutPrinter(variant string) *GenericStdoutPrinter {
	return &GenericStdoutPrinter{variant: variant}
}

// Write implements the legacy Printer interface for stdout
func (p *GenericStdoutPrinter) Write(f interface{}, c interface{}) error {
	switch p.variant {
	case "forecast":
		writer := NewForecastTableWriter()
		return writer.Write(f.(types.ForecastPrintData))
	case "costAndUsage":
		sortBy := f.(string)
		writer := NewCostUsageTableWriter(sortBy)
		return writer.Write(c.(types.CostAndUsageOutputType))
	}
	return nil
}

// NewGenericCsvPrinter creates a new CSV printer with backward compatibility
func NewGenericCsvPrinter(variant string) *GenericCsvPrinter {
	return &GenericCsvPrinter{variant: variant}
}

// Write implements the legacy Printer interface for CSV
func (p *GenericCsvPrinter) Write(f interface{}, c interface{}) error {
	switch p.variant {
	case "costAndUsage":
		sortBy := f.(string)
		writer := NewCostUsageCSVWriter(sortBy)
		return writer.Write(c.(types.CostAndUsageOutputType))
	}
	return nil
}

// NewGenericChartPrinter creates a new chart printer with backward compatibility
func NewGenericChartPrinter(variant string) *GenericChartPrinter {
	return &GenericChartPrinter{variant: variant}
}

// Write implements the legacy Printer interface for charts
func (p *GenericChartPrinter) Write(f interface{}, c interface{}) error {
	switch p.variant {
	case "costAndUsage":
		sortBy := f.(string)
		writer := NewCostUsageChartWriter(sortBy)
		return writer.Write(c.(types.CostAndUsageOutputType))
	}
	return nil
}

// NewGenericPineconePrinter creates a new Pinecone printer with backward compatibility
func NewGenericPineconePrinter(variant string) *GenericPineconePrinter {
	return &GenericPineconePrinter{variant: variant}
}

// Write implements the legacy Printer interface for vector databases
func (p *GenericPineconePrinter) Write(f interface{}, c interface{}) error {
	switch p.variant {
	case "costAndUsage":
		writer := NewCostUsageVectorWriter()
		return writer.Write(c.(types.CostAndUsageOutputType))
	}
	return nil
}