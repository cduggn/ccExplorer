package printer

import (
	"github.com/cduggn/ccexplorer/pkg/printer/writers/stdout"
	"strings"
)

func ForecastToStdout(r ForecastPrintData,
	dimensions []string) {

	filteredBy := strings.Join(dimensions, " | ")

	output := ConvertToForecastStdoutType(r, filteredBy)

	w, err := stdout.NewStdoutWriter("forecast")
	if err != nil {
		return
	}
	w.Writer(output)
	//t := CreateTable(forecastedHeader)
	//rows := ForecastToRows(r)
	//t.AppendRows(rows)

	//footer := forecastedTableFooter(filteredBy, *r.Forecast.Total.Unit, *r.Forecast.Total.Amount)
	//t.AppendRow(footer)
	//t.Render()
	//
	//stdout.Writer(output, "forecast")

}
