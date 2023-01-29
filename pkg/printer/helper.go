package printer

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"strconv"
)

func CreateTable(header table.Row) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(header)
	return t
}

func CostUsageToRows(s []Service, granularity string) CostAndUsage {
	var rows []table.Row
	var total float64
	for index, v := range s {
		for _, m := range v.Metrics {

			if index%10 == 0 {
				rows = append(rows, tableDivider)
			}
			if m.Unit == "USD" {
				total += m.NumericAmount
			}

			tempRow := table.Row{index, v.Keys[0], ReturnIfPresent(v.Keys),
				v.Name, fmt.Sprintf("%f10",
					m.NumericAmount), m.Amount,
				m.Unit,
				granularity,
				v.Start, v.End}

			rows = append(rows, tempRow)
		}
	}

	return CostAndUsage{Rows: rows, Total: total}
}

func ForecastToRows(r ForecastPrintData) []table.Row {

	var rows []table.Row
	for _, v := range r.Forecast.ForecastResultsByTime {
		tempRow := table.Row{*v.TimePeriod.Start,
			*v.TimePeriod.End, *v.MeanValue, *v.PredictionIntervalUpperBound,
			*v.PredictionIntervalLowerBound}

		rows = append(rows, tempRow)
	}
	return rows
}

func CurateCostAndUsageReport(
	d CostAndUsageReportPrintData) CostAndUsageReport {

	c := CostAndUsageReport{
		Services:    make(map[int]Service),
		Granularity: d.Granularity,
	}
	count := 0
	for _, v := range d.Report.ResultsByTime {
		c.Start = *v.TimePeriod.Start
		c.End = *v.TimePeriod.End
		for _, g := range v.Groups {
			keys := make([]string, 0)
			service := Service{
				Start: c.Start,
				End:   c.End,
			}
			keys = append(keys, g.Keys...)

			for key, m := range g.Metrics {
				metrics := Metrics{
					Name:          key,
					Amount:        *m.Amount,
					NumericAmount: ConvertToFloat(*m.Amount),
					Unit:          *m.Unit,
				}
				service.Metrics = append(service.Metrics, metrics)
			}
			service.Keys = keys
			c.Services[count] = service
			count++
		}

	}
	return c
}

func ConvertToFloat(amount string) float64 {
	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func ReturnIfPresent(s []string) string {
	if len(s) == 1 {
		return ""
	} else {
		return s[1]
	}

}

func ConvertServiceToSlice(s Service, granularity string) [][]string {

	var r [][]string
	for _, v := range s.Metrics {
		t := []string{s.Keys[0], ReturnIfPresent(s.Keys), v.Name,
			granularity, s.Start, s.End,
			v.Amount, v.Unit}
		r = append(r, t)
	}
	return r
}

func ToPrintWriterType(s string) PrintWriterType {
	switch s {
	case "csv":
		return CSV
	case "stdout":
		return Stdout
	case "chart":
		return Chart
	default:
		return Stdout
	}
}
