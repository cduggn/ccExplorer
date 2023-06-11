package presentation

import (
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	"github.com/cduggn/ccexplorer/internal/core/presentation/writers"
	"os"
)

var (
	OutputDir = "./output"
)

func init() {
	if _, err := os.Stat(OutputDir); os.IsNotExist(err) {
		err := os.Mkdir(OutputDir, 0755)
		if err != nil {
			panic("Unable output directory")
		}
	}
}

func NewPrintWriter(printType model.PrintWriterType, variant string) Printer {
	switch printType {
	case model.Stdout:
		return &StdoutPrinter{variant}
	case model.CSV:
		return &CsvPrinter{variant}
	case model.Chart:
		return &ChartPrinter{variant}
	case model.OpenAPI:
		return &OpenAIPrinter{variant}
	default:
		panic("Invalid print type")
	}
}

func (p *StdoutPrinter) Print(f interface{}, c interface{}) error {
	switch p.Variant {
	case "forecast":
		writers.ForecastToStdout(f.(model.ForecastPrintData), c.([]string))
	case "costAndUsage":
		fn := writers.SortFunction(f.(string))
		err := writers.CostAndUsageToStdout(fn, c.(model.CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *CsvPrinter) Print(f interface{}, c interface{}) error {
	switch p.Variant {
	case "costAndUsage":
		fn := writers.SortFunction(f.(string))
		err := writers.CostAndUsageToCSV(fn, c.(model.CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ChartPrinter) Print(f interface{}, c interface{}) error {
	switch p.Variant {
	case "costAndUsage":
		fn := writers.SortFunction(f.(string))
		err := writers.CostAndUsageToChart(fn, c.(model.CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *OpenAIPrinter) Print(f interface{}, c interface{}) error {
	switch p.Variant {
	case "costAndUsage":
		fn := writers.SortFunction(f.(string))
		err := writers.CostAndUsageToOpenAI(fn, c.(model.CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}
