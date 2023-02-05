package printer

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	aws2 "github.com/cduggn/ccexplorer/pkg/service/aws"
	"github.com/go-echarts/go-echarts/v2/opts"
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

func CreateSubTitle(granularity string, start string, end string) string {
	return fmt.Sprintf("Response granularity: %s. Timeframe: %s-%s",
		granularity,
		start, end)
}

func CreateTitle(dimension string) string {
	return fmt.Sprintf("Pie chart for dimension: [ %s ]", dimension)
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
				m.Name, fmt.Sprintf("%f10",
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

func ToCostAndUsageOutputType(r *costexplorer.GetCostAndUsageOutput,
	u aws2.CostAndUsageRequestType) CostAndUsageOutputType {
	return CurateCostAndUsageReport(r, u)
}

func CurateCostAndUsageReport(
	d *costexplorer.GetCostAndUsageOutput, query aws2.CostAndUsageRequestType) CostAndUsageOutputType {

	c := CostAndUsageOutputType{
		Services:    make(map[int]Service),
		Granularity: query.Granularity,
		Dimensions:  query.GroupBy,
		Tags:        query.GroupByTag,
		Start:       query.Time.Start,
		End:         query.Time.End,
	}

	c.Services = ResultsToServicesMap(d.ResultsByTime)
	return c
}

func ResultsToServicesMap(res []types.ResultByTime) map[int]Service {
	services := make(map[int]Service)
	count := 0
	for _, v := range res {
		for _, g := range v.Groups {
			keys := append(make([]string, 0), g.Keys...)
			service := Service{
				Start: *v.TimePeriod.Start,
				End:   *v.TimePeriod.End,
				Keys:  keys,
			}

			service.Metrics = MetricsToService(g.Metrics)
			services[count] = service
			count++
		}
	}
	return services
}

func MetricsToService(m map[string]types.MetricValue) []Metrics {
	var metrics []Metrics
	for k, v := range m {
		metrics = append(metrics, Metrics{
			Name:          k,
			Amount:        *v.Amount,
			NumericAmount: ConvertToFloat(*v.Amount),
			Unit:          *v.Unit,
		})
	}
	return metrics
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

func ReturnFirst(s []string) string {
	if len(s) > 1 {
		return s[0]
	}
	return ""
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

func PopulatePieDate(s map[int]Service, key int) []opts.
	PieData {
	items := make([]opts.PieData, 0)

	services := SortServicesByMetricAmount(s)

	for index, v := range services {
		if index < 15 {
			items = append(items, opts.PieData{Name: v.Keys[key],
				Value: v.Metrics[0].NumericAmount})
		}

	}
	return items
}
