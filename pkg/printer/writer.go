package printer

import (
	"github.com/jedib0t/go-pretty/v6/table"
)

var (
	costAndUsageHeader = table.Row{"Rank", "Dimension/Tag", "Dimension/Tag",
		"Metric Name", "Truncated USD Amount", "Amount",
		"Unit",
		"Granularity",
		"Start",
		"End"}
	forecastedHeader = table.Row{"Start", "End", "Mean Value",
		"Prediction Interval LowerBound",
		"Prediction Interval UpperBound", "Unit", "Total"}
	tableDivider = table.Row{"-", "-", "-",
		"-", "-", "-", "-",
		"-",
		"-", ""}
	costAndUsageTableFooter = func(t float64) table.Row {
		return table.
			Row{"", "",
			"",
			"",
			"TOTAL COST",
			t, "", "", "", ""}
	}
	forecasteTableFooter = func(filter string, unit string,
		amount string) table.Row {
		return table.Row{"FilteredBy", filter, "", "", "",
			unit,
			amount}
	}
)

func PrintFactory(printType PrintWriterType, variant string) Printer {
	switch printType {
	case Stdout:
		return &StdoutPrinter{variant}
	case CSV:
		return &CsvPrinter{variant}
	case Chart:
		return &ChartPrinter{variant}
	default:
		panic("Invalid print type")
	}
}

func (p *StdoutPrinter) Print(f interface{}, c interface{}) error {

	switch p.Variant {
	case "forecast":
		ForecastToStdout(f.(ForecastPrintData), c.([]string))
	case "costAndUsage":
		CostAndUsageToStdout(f.(func(r map[int]Service) []Service), c.(CostAndUsageOutputType))
	}
	return nil

}

func (p *CsvPrinter) Print(f interface{}, c interface{}) error {
	switch p.Variant {
	// no requirement for csv printing for forecast
	case "costAndUsage":
		err := CostAndUsageToCSV(f.(func(r map[int]Service) []Service),
			c.(CostAndUsageOutputType))
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
		err := CostAndUsageToChart(f.(func(r map[int]Service) []Service),
			c.(CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}