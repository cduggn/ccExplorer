package printer

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"strconv"
)

var tableDivider = table.Row{"-", "-", "-",
	"-", "-", "-", "-",
	"-",
	"-", ""}

func CurateCostAndUsageReport(output *costexplorer.GetCostAndUsageOutput,
	granularity string) CostAndUsageReport {

	c := CostAndUsageReport{
		Services:    make(map[int]Service),
		Granularity: granularity,
	}
	count := 0
	for _, v := range output.ResultsByTime {
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
		"Metric Name", "Numeric Amount", "String Amount", "Unit",
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
				// add to total
				//total.Add(&total, &v.NumericAmount)
				total += v.NumericAmount
			}
			tempRow := table.Row{index, m.Keys[0], ReturnIfPresent(m.Keys),
				v.Name, fmt.Sprintf("%f10", v.NumericAmount), v.Amount, v.Unit,
				r.Granularity,
				m.Start, m.End}
			t.AppendRow(tempRow)

		}
	}
	totalHeaderRow := table.Row{"", "", "", "", "", "", "", "", "", ""}
	totalRow := table.Row{"", "", "", "", "TOTAL COST", total, "", "", "", ""}
	t.AppendRow(totalHeaderRow)
	t.AppendRow(totalRow)
	t.Render()
}

func PrintGetCostForecastReport(r *costexplorer.GetCostForecastOutput) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Mean Value", "PredictionIntervalLowerBound",
		"PredictionIntervalUpperBound", "Start", "End", "Unit", "Total"})
	for _, v := range r.ForecastResultsByTime {

		tempRow := table.Row{*v.MeanValue, *v.PredictionIntervalUpperBound,
			*v.PredictionIntervalLowerBound, *v.TimePeriod.Start,
			*v.TimePeriod.End}
		t.AppendRow(tempRow)
	}

	totalRow := table.Row{"", "", "", "", "", *r.Total.Unit, *r.Total.Amount}
	t.AppendRow(totalRow)
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
