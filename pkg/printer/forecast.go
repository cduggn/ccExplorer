package printer

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"strings"
)

var (
	forecastedHeader = table.Row{"Start", "End", "Mean Value",
		"Prediction Interval LowerBound",
		"Prediction Interval UpperBound", "Unit", "Total"}
	forecastedTableFooter = func(filter string, unit string,
		amount string) table.Row {
		return table.Row{"FilteredBy", filter, "", "", "",
			unit,
			amount}
	}
)

func ForecastToStdout(r ForecastPrintData,
	dimensions []string) {
	filteredBy := strings.Join(dimensions, " | ")

	t := CreateTable(forecastedHeader)
	rows := ForecastToRows(r)
	t.AppendRows(rows)

	footer := forecastedTableFooter(filteredBy, *r.Forecast.Total.Unit,
		*r.Forecast.Total.Amount)
	t.AppendRow(footer)
	t.Render()
}
