package printer

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/cduggn/ccexplorer/pkg/printer/writers/chart"
	"github.com/cduggn/ccexplorer/pkg/printer/writers/stdout"
	aws2 "github.com/cduggn/ccexplorer/pkg/service/aws"
	"sort"
	"strconv"
	"time"
)

func CreateSubTitle(granularity string, start string, end string) string {
	return fmt.Sprintf("Response granularity: %s. Timeframe: %s-%s",
		granularity,
		start, end)
}

func ToCostAndUsageOutputType(r *costexplorer.GetCostAndUsageOutput,
	u aws2.CostAndUsageRequestType) CostAndUsageOutputType {
	return CurateCostAndUsageReport(r, u)
}

func CurateCostAndUsageReport(
	d *costexplorer.GetCostAndUsageOutput, query aws2.CostAndUsageRequestType) CostAndUsageOutputType {

	c := CostAndUsageOutputType{
		Services:     make(map[int]Service),
		Granularity:  query.Granularity,
		Dimensions:   query.GroupBy,
		Tags:         query.GroupByTag,
		Start:        query.Time.Start,
		End:          query.Time.End,
		OpenAIAPIKey: query.OpenAIAPIKey,
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
	case "gpt3":
		return OpenAPI
	default:
		return Stdout
	}
}

//func CreateOutputDir(outputDir string) (string, error) {
//
//	dir, err := os.Getwd()
//	if err != nil {
//		return "", err
//	}
//	dir = dir + outputDir
//	if _, err := os.Stat(dir); os.IsNotExist(err) {
//		err = os.Mkdir(dir, 0755)
//		if err != nil {
//			return "", err
//		}
//	}
//	return dir, nil
//}

func SortFunction(sortBy string) func(r map[int]Service) []Service {
	switch sortBy {
	case "date":
		return SortServicesByStartDate
	case "cost":
		return SortServicesByMetricAmount
	default:
		return SortServicesByMetricAmount
	}
}

func SortServicesByStartDate(r map[int]Service) []Service {
	// Create a slice of key-value pairs
	pairs := make([]struct {
		Key   int
		Value Service
	}, len(r))
	i := 0
	for k, v := range r {
		pairs[i] = struct {
			Key   int
			Value Service
		}{k, v}
		i++
	}

	sort.SliceStable(pairs, func(i, j int) bool {

		t1, _ := time.Parse("2006-01-02", pairs[i].Value.Start)
		t2, _ := time.Parse("2006-01-02", pairs[j].Value.Start)
		return t1.After(t2)
	})

	result := make([]Service, len(pairs))
	for i, pair := range pairs {
		result[i] = pair.Value
	}
	return result
}

func SortServicesByMetricAmount(r map[int]Service) []Service {
	// Create a slice of key-value pairs
	pairs := make([]struct {
		Key   int
		Value Service
	}, len(r))
	i := 0
	for k, v := range r {
		pairs[i] = struct {
			Key   int
			Value Service
		}{k, v}
		i++
	}

	// Sort the slice by the Value.Metrics[0].Amount field
	sort.SliceStable(pairs, func(i, j int) bool {
		return pairs[i].Value.Metrics[0].NumericAmount > pairs[j].Value.
			Metrics[0].NumericAmount
	})

	result := make([]Service, len(pairs))
	for i, pair := range pairs {
		result[i] = pair.Value
	}
	return result
}

func ConvertServiceMapToArray(s map[int]Service,
	granularity string) [][]string {
	var rows [][]string
	for _, v := range s {
		rows = append(rows, ConvertServiceToSlice(v, granularity)...)
	}
	return rows
}

func ConvertServiceSliceToArray(s []Service, granularity string) [][]string {
	var rows [][]string
	for _, v := range s {
		rows = append(rows, ConvertServiceToSlice(v, granularity)...)
	}
	return rows
}

func ConvertToStdoutType(s []Service,
	granularity string) stdout.CostAndUsageStdoutType {

	outputType := stdout.CostAndUsageStdoutType{
		Granularity: granularity,
	}

	var services []stdout.Service
	for _, v := range s {
		var metrics []stdout.Metrics
		for _, m := range v.Metrics {
			metrics = append(metrics, stdout.Metrics{
				Name:          m.Name,
				Amount:        m.Amount,
				Unit:          m.Unit,
				NumericAmount: m.NumericAmount,
			})
		}
		services = append(services, stdout.Service{
			Name:    v.Keys[0],
			Keys:    v.Keys,
			Start:   v.Start,
			End:     v.End,
			Metrics: metrics,
		})
	}
	outputType.Services = services

	return outputType
}

func ConvertToChartInputType(r CostAndUsageOutputType,
	s []Service) chart.InputType {

	input := chart.InputType{
		Granularity: r.Granularity,
		Start:       r.Start,
		End:         r.End,
		Dimensions:  r.Dimensions,
		Tags:        r.Tags,
	}

	var services []chart.Service
	for _, service := range s {
		var metrics []chart.Metrics
		for _, metric := range service.Metrics {
			metrics = append(metrics, chart.Metrics{
				Name:          metric.Name,
				Amount:        metric.Amount,
				Unit:          metric.Unit,
				UsageQuantity: metric.UsageQuantity,
				NumericAmount: metric.NumericAmount,
			})
		}

		services = append(services, chart.Service{
			Name:    service.Name,
			Keys:    service.Keys,
			Start:   service.Start,
			End:     service.End,
			Metrics: metrics,
		})
	}

	input.Services = services

	return input

}

func ConvertToForecastStdoutType(r ForecastPrintData,
	filteredBy string) stdout.ForecastStdoutType {
	var forecast []stdout.ForecastResults
	for _, v := range r.Forecast.ForecastResultsByTime {
		forecast = append(forecast, stdout.ForecastResults{
			TimePeriod: stdout.DateInterval{
				Start: *v.TimePeriod.Start,
				End:   *v.TimePeriod.End,
			},
			MeanValue:                    *v.MeanValue,
			PredictionIntervalLowerBound: *v.PredictionIntervalLowerBound,
			PredictionIntervalUpperBound: *v.PredictionIntervalUpperBound,
		})
	}

	return stdout.ForecastStdoutType{
		Forecast: forecast,
		Total: stdout.Total{
			Amount: *r.Forecast.Total.Amount,
			Unit:   *r.Forecast.Total.Unit,
		},
		FilteredBy: filteredBy,
	}

}
