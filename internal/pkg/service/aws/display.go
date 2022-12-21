package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

func PrintGetCostForecastReport(r *costexplorer.GetCostForecastOutput) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Mean Value", "PredictionIntervalLowerBound",
		"PredictionIntervalUpperBound", "Start", "End", "Unit", "Total"})
	for _, v := range r.ForecastResultsByTime {

		tempRow := table.Row{*v.MeanValue, v.PredictionIntervalUpperBound,
			v.PredictionIntervalLowerBound, *v.TimePeriod.Start,
			*v.TimePeriod.End}
		t.AppendRow(tempRow)
	}

	totalRow := table.Row{"", "", "", "", "", *r.Total.Unit, *r.Total.Amount}
	t.AppendRow(totalRow)
	t.Render()
}

func (c *CostAndUsageReport) CurateCostAndUsageReport(output *costexplorer.
	GetCostAndUsageOutput) {
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
					Name:   key,
					Amount: *m.Amount,
					Unit:   *m.Unit,
				}
				service.Metrics = append(service.Metrics, metrics)
			}
			service.Keys = keys
			c.Services[count] = service
			count++
		}

	}
}

func (c *CostAndUsageReport) PrintCostAndUsageReport() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Dimension/Tag", "Dimension/Tag", "Metric Name", "Amount", "Unit", "Granularity", "Start", "End"})
	var total float64
	for _, m := range c.Services {
		for _, v := range m.Metrics {
			if v.Unit == "USD" {
				total += ConvertToFloat(v.Amount)
			}
			tempRow := table.Row{m.Keys[0], isEmpty(m.Keys), v.Name, v.Amount, v.Unit, c.Granularity, m.Start, m.End}
			t.AppendRow(tempRow)

			//_, err := conn.Insert(storage.CostDataInsert{
			//	Dimension:   m.Keys[0],
			//	Dimension2:  "",
			//	Tag:         "",
			//	MetricName:  "",
			//	Amount:      v.Amount,
			//	Unit:        v.Unit,
			//	Granularity: c.Granularity,
			//	StartDate:   m.Start,
			//	EndDate:     m.End,
			//})
			//if err != nil {
			//	logger.Error(err.Error())
			//}
		}
	}
	totalHeaderRow := table.Row{"", "", "", "", "", "", "", ""}
	totalRow := table.Row{"", "", "TOTAL COST", total, "", "", "", ""}
	t.AppendRow(totalHeaderRow)
	t.AppendRow(totalRow)
	t.Render()
}
