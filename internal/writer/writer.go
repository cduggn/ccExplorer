// Package writer provides output formatting for AWS cost data.
//
// The package uses a strategy pattern with composable transformers and renderers.
// The main components are:
//   - Transformers: convert domain data to output-specific formats
//   - Renderers: write formatted data to destinations (stdout, files, vector DB)
//   - CompositeWriter: combines a transformer and renderer
//
// For new code, use the factory functions in writers.go directly.
// The PrinterAdapter type provides backward compatibility with legacy code.
package writer

import (
	"github.com/cduggn/ccexplorer/internal/types"
	"os"
)

var (
	OutputDir = "./writer"
)

type Builder struct {
}

func init() {
	if _, err := os.Stat(OutputDir); os.IsNotExist(err) {
		err := os.Mkdir(OutputDir, 0755)
		if err != nil {
			panic("Unable writer directory")
		}
	}
}

// PrinterAdapter adapts the new CompositeWriter to the legacy Printer interface
type PrinterAdapter struct {
	printType types.PrintWriterType
	variant   string
}

func NewPrintWriter(printType types.PrintWriterType, variant string) Printer {
	return &PrinterAdapter{
		printType: printType,
		variant:   variant,
	}
}

// Write implements the Printer interface by delegating to appropriate CompositeWriter
func (p *PrinterAdapter) Write(f interface{}, c interface{}) error {
	switch p.printType {
	case types.Stdout:
		return p.writeStdout(f, c)
	case types.CSV:
		return p.writeCSV(f, c)
	case types.Chart:
		return p.writeChart(f, c)
	case types.Pinecone:
		return p.writeVector(f, c)
	default:
		panic("Invalid print type")
	}
}

func (p *PrinterAdapter) writeStdout(f interface{}, c interface{}) error {
	switch p.variant {
	case "forecast":
		writer := NewForecastTableWriter()
		return writer.Write(f.(types.ForecastPrintData))
	case "costAndUsage":
		sortBy := f.(string)
		writer := NewCostUsageTableWriter(sortBy)
		return writer.Write(c.(types.CostAndUsageOutputType))
	default:
		panic("Invalid stdout variant: " + p.variant)
	}
}

func (p *PrinterAdapter) writeCSV(f interface{}, c interface{}) error {
	if p.variant != "costAndUsage" {
		panic("Invalid CSV variant: " + p.variant)
	}
	sortBy := f.(string)
	writer := NewCostUsageCSVWriter(sortBy)
	return writer.Write(c.(types.CostAndUsageOutputType))
}

func (p *PrinterAdapter) writeChart(f interface{}, c interface{}) error {
	if p.variant != "costAndUsage" {
		panic("Invalid chart variant: " + p.variant)
	}
	sortBy := f.(string)
	writer := NewCostUsageChartWriter(sortBy)
	return writer.Write(c.(types.CostAndUsageOutputType))
}

func (p *PrinterAdapter) writeVector(f interface{}, c interface{}) error {
	if p.variant != "costAndUsage" {
		panic("Invalid vector variant: " + p.variant)
	}
	writer := NewCostUsageVectorWriter()
	return writer.Write(c.(types.CostAndUsageOutputType))
}
