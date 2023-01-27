package printer

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"strconv"
	"strings"
)

var tableDivider = table.Row{"-", "-", "-",
	"-", "-", "-", "-",
	"-",
	"-", ""}

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
func PrintCostAndUsageReport(s func(r map[int]Service) []Service,
	r CostAndUsageReport) {
	sortedServices := s(r.Services)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Rank", "Dimension/Tag", "Dimension/Tag",
		"Metric Name", "Numeric Amount", "String Amount",
		"Unit",
		"Granularity",
		"Start",
		"End"})
	var total float64
	for index, m := range sortedServices {

		for _, v := range m.Metrics {

			if index%10 == 0 {
				t.AppendRow(tableDivider)
			}
			if v.Unit == "USD" {
				total += v.NumericAmount
			}

			tempRow := table.Row{index, m.Keys[0], ReturnIfPresent(m.Keys),
				v.Name, fmt.Sprintf("%f10", v.NumericAmount), v.Amount,
				v.Unit,
				r.Granularity,
				m.Start, m.End}
			t.AppendRow(tempRow)

		}
	}
	totalHeaderRow := table.Row{"", "", "", "", "", "", "", "", "", ""}
	totalRow := table.Row{"", "", "", "", "TOTAL COST", total, "", "", "",
		""}
	t.AppendRow(totalHeaderRow)
	t.AppendRow(totalRow)
	t.Render()
}

func PrintGetCostForecastReport(r ForecastPrintData,
	dimensions []string) {
	filteredBy := strings.Join(dimensions, " | ")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"Start", "End", "Mean Value",
		"Prediction Interval LowerBound",
		"Prediction Interval UpperBound", "Unit", "Total"})

	for _, v := range r.Forecast.ForecastResultsByTime {

		tempRow := table.Row{*v.TimePeriod.Start,
			*v.TimePeriod.End, *v.MeanValue, *v.PredictionIntervalUpperBound,
			*v.PredictionIntervalLowerBound}
		t.AppendRow(tempRow)
	}

	t.AppendSeparator()
	t.AppendRow(table.Row{"FilteredBy", filteredBy, "", "", "",
		*r.Forecast.Total.Unit,
		*r.Forecast.Total.Amount})
	t.Render()
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
