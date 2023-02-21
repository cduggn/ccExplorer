package printer

import (
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

func PrintFactory(printType PrintWriterType, variant string) Printer {
	switch printType {
	case Stdout:
		return &StdoutPrinter{variant}
	case CSV:
		return &CsvPrinter{variant}
	case Chart:
		return &ChartPrinter{variant}
	case OpenAPI:
		return &OpenAIPrinter{variant}
	default:
		panic("Invalid print type")
	}
}

func (p *StdoutPrinter) Print(f interface{}, c interface{}) error {

	switch p.Variant {
	case "forecast":
		ForecastToStdout(f.(ForecastPrintData), c.([]string))
	case "costAndUsage":
		fn := SortFunction(f.(string))
		CostAndUsageToStdout(fn, c.(CostAndUsageOutputType))
	}
	return nil

}

func (p *CsvPrinter) Print(f interface{}, c interface{}) error {
	switch p.Variant {
	// no requirement for csv printing for forecast
	case "costAndUsage":
		fn := SortFunction(f.(string))
		err := CostAndUsageToCSV(fn, c.(CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ChartPrinter) Print(f interface{}, c interface{}) error {
	switch p.Variant {
	// no requirement for csv printing for forecast
	case "costAndUsage":
		fn := SortFunction(f.(string))
		err := CostAndUsageToChart(fn, c.(CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *OpenAIPrinter) Print(f interface{}, c interface{}) error {

	// no requirement for csv printing for forecast
	switch p.Variant {
	case "costAndUsage":
		fn := SortFunction(f.(string))
		err := CostAndUsageToOpenAI(fn, c.(CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}
