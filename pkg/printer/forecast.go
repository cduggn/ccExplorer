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
}
