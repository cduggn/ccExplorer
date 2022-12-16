package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/cloudcost/internal/pkg/logger"
	"github.com/cduggn/cloudcost/internal/pkg/storage"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

var (
	conn *storage.CostDataStorage
)

func init() {
	newConnection()
}

func newConnection() {
	conn = &storage.CostDataStorage{}
	err := conn.New("./cloudcost.db")
	if err != nil {
		logger.Error(err.Error())
	}
}

func (c *CostAndUsageReport) Print() {
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

			conn.Insert(storage.CostDataInsert{
				Dimension:   m.Keys[0],
				Dimension2:  "",
				Tag:         "",
				MetricName:  "",
				Amount:      v.Amount,
				Unit:        v.Unit,
				Granularity: c.Granularity,
				StartDate:   m.Start,
				EndDate:     m.End,
			})
		}
	}
	totalHeaderRow := table.Row{"", "", "", "", "", "", "", ""}
	totalRow := table.Row{"", "", "TOTAL COST", total, "", "", "", ""}
	t.AppendRow(totalHeaderRow)
	t.AppendRow(totalRow)
	t.Render()
}

func isEmpty(s []string) string {
	if len(s) == 1 {
		return ""
	} else {
		return s[1]
	}

}

func (c *CostAndUsageReport) CurateReport(output *costexplorer.GetCostAndUsageOutput) {
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
