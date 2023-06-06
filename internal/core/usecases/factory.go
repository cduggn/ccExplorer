package usecases

import (
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	"github.com/cduggn/ccexplorer/internal/core/util"
	"os"
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
	case model.Pinecone:
		return &PineconePrinter{variant}
	default:
		panic("Invalid print type")
	}
}

func (p *PineconePrinter) Write(f interface{}, c interface{}) error {
	switch p.Variant {
	case "costAndUsage":
		err := CostAndUsageToVectorMapper(c.(model.CostAndUsageOutputType))
		if err != nil {
			return err
		}
		/// working with CostAndUsageOutputType
	}
	return nil
}

func (p *StdoutPrinter) Write(f interface{}, c interface{}) error {
	switch p.Variant {
	case "forecast":
		ForecastToStdoutMapper(f.(model.ForecastPrintData),
			c.([]string))
	case "costAndUsage":
		fn := util.SortFunction(f.(string))
		err := CostAndUsageToStdoutMapper(fn,
			c.(model.CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *CsvPrinter) Write(f interface{}, c interface{}) error {
	switch p.Variant {
	case "costAndUsage":
		fn := util.SortFunction(f.(string))
		err := CostAndUsageToCSVMapper(fn,
			c.(model.CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ChartPrinter) Write(f interface{}, c interface{}) error {
	switch p.Variant {
	case "costAndUsage":
		fn := util.SortFunction(f.(string))
		err := CostAndUsageToChartMapper(fn,
			c.(model.CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}
