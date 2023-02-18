package printer

import "strings"

func ForecastToStdout(r ForecastPrintData,
	dimensions []string) {
	filteredBy := strings.Join(dimensions, " | ")

	t := CreateTable(forecastedHeader)
	rows := ForecastToRows(r)
	t.AppendRows(rows)

	footer := forecasteTableFooter(filteredBy, *r.Forecast.Total.Unit,
		*r.Forecast.Total.Amount)
	t.AppendRow(footer)
	t.Render()
}
