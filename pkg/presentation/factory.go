package presentation

import (
	formats2 "github.com/cduggn/ccexplorer/pkg/presentation/writers"
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

func NewPrintWriter(printType formats2.PrintWriterType, variant string) Writer {
	switch printType {
	case formats2.Stdout:
		return &StdoutPrinter{variant}
	case formats2.CSV:
		return &CsvPrinter{variant}
	case formats2.Chart:
		return &ChartPrinter{variant}
	case formats2.OpenAPI:
		return &OpenAIPrinter{variant}
	default:
		panic("Invalid print type")
	}
}

func (p *StdoutPrinter) Write(f interface{}, c interface{}) error {
	switch p.Variant {
	case "forecast":
		formats2.ForecastToStdout(f.(formats2.ForecastPrintData), c.([]string))
	case "costAndUsage":
		fn := formats2.SortFunction(f.(string))
		err := formats2.CostAndUsageToStdout(fn, c.(formats2.CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *CsvPrinter) Write(f interface{}, c interface{}) error {
	switch p.Variant {
	case "costAndUsage":
		fn := formats2.SortFunction(f.(string))
		err := formats2.CostAndUsageToCSV(fn, c.(formats2.CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ChartPrinter) Write(f interface{}, c interface{}) error {
	switch p.Variant {
	case "costAndUsage":
		fn := formats2.SortFunction(f.(string))
		err := formats2.CostAndUsageToChart(fn, c.(formats2.CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *OpenAIPrinter) Write(f interface{}, c interface{}) error {
	switch p.Variant {
	case "costAndUsage":
		fn := formats2.SortFunction(f.(string))
		err := formats2.CostAndUsageToOpenAI(fn, c.(formats2.CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}
