package writers

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	aws2 "github.com/cduggn/ccexplorer/internal/core/domain/model"
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
	u aws2.CostAndUsageRequestType) aws2.CostAndUsageOutputType {
	return CurateCostAndUsageReport(r, u)
}

func CurateCostAndUsageReport(
	d *costexplorer.GetCostAndUsageOutput, query aws2.CostAndUsageRequestType) aws2.CostAndUsageOutputType {

	c := aws2.CostAndUsageOutputType{
		Services:     make(map[int]aws2.Service),
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

func ResultsToServicesMap(res []types.ResultByTime) map[int]aws2.Service {
	services := make(map[int]aws2.Service)
	count := 0
	for _, v := range res {
		for _, g := range v.Groups {
			keys := append(make([]string, 0), g.Keys...)
			service := aws2.Service{
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

func MetricsToService(m map[string]types.MetricValue) []aws2.Metrics {
	var metrics []aws2.Metrics
	for k, v := range m {
		metrics = append(metrics, aws2.Metrics{
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

func ConvertServiceToSlice(s aws2.Service, granularity string) [][]string {

	var r [][]string
	for _, v := range s.Metrics {
		t := []string{s.Keys[0], ReturnIfPresent(s.Keys), v.Name,
			granularity, s.Start, s.End,
			v.Amount, v.Unit}
		r = append(r, t)
	}
	return r
}

func ToPrintWriterType(s string) aws2.PrintWriterType {
	switch s {
	case "csv":
		return aws2.CSV
	case "stdout":
		return aws2.Stdout
	case "chart":
		return aws2.Chart
	case "gpt3":
		return aws2.OpenAPI
	default:
		return aws2.Stdout
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

func SortFunction(sortBy string) func(r map[int]aws2.Service) []aws2.Service {
	switch sortBy {
	case "date":
		return SortServicesByStartDate
	case "cost":
		return SortServicesByMetricAmount
	default:
		return SortServicesByMetricAmount
	}
}

func SortServicesByStartDate(r map[int]aws2.Service) []aws2.Service {
	// Create a slice of key-value pairs
	pairs := make([]struct {
		Key   int
		Value aws2.Service
	}, len(r))
	i := 0
	for k, v := range r {
		pairs[i] = struct {
			Key   int
			Value aws2.Service
		}{k, v}
		i++
	}

	sort.SliceStable(pairs, func(i, j int) bool {

		t1, _ := time.Parse("2006-01-02", pairs[i].Value.Start)
		t2, _ := time.Parse("2006-01-02", pairs[j].Value.Start)
		return t1.After(t2)
	})

	result := make([]aws2.Service, len(pairs))
	for i, pair := range pairs {
		result[i] = pair.Value
	}
	return result
}

func SortServicesByMetricAmount(r map[int]aws2.Service) []aws2.Service {
	// Create a slice of key-value pairs
	pairs := make([]struct {
		Key   int
		Value aws2.Service
	}, len(r))
	i := 0
	for k, v := range r {
		pairs[i] = struct {
			Key   int
			Value aws2.Service
		}{k, v}
		i++
	}

	// Sort the slice by the Value.Metrics[0].Amount field
	sort.SliceStable(pairs, func(i, j int) bool {
		return pairs[i].Value.Metrics[0].NumericAmount > pairs[j].Value.
			Metrics[0].NumericAmount
	})

	result := make([]aws2.Service, len(pairs))
	for i, pair := range pairs {
		result[i] = pair.Value
	}
	return result
}

func ConvertServiceMapToArray(s map[int]aws2.Service,
	granularity string) [][]string {
	var rows [][]string
	for _, v := range s {
		rows = append(rows, ConvertServiceToSlice(v, granularity)...)
	}
	return rows
}

func ConvertServiceSliceToArray(s []aws2.Service, granularity string) [][]string {
	var rows [][]string
	for _, v := range s {
		rows = append(rows, ConvertServiceToSlice(v, granularity)...)
	}
	return rows
}

func ConvertToStdoutType(s []aws2.Service,
	granularity string) aws2.CostAndUsageStdoutType {

	outputType := aws2.CostAndUsageStdoutType{
		Granularity: granularity,
	}

	var services []aws2.Service
	for _, v := range s {
		var metrics []aws2.Metrics
		for _, m := range v.Metrics {
			metrics = append(metrics, aws2.Metrics{
				Name:          m.Name,
				Amount:        m.Amount,
				Unit:          m.Unit,
				NumericAmount: m.NumericAmount,
			})
		}
		services = append(services, aws2.Service{
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

func ConvertToChartInputType(r aws2.CostAndUsageOutputType,
	s []aws2.Service) aws2.InputType {

	input := aws2.InputType{
		Granularity: r.Granularity,
		Start:       r.Start,
		End:         r.End,
		Dimensions:  r.Dimensions,
		Tags:        r.Tags,
	}

	var services []aws2.Service
	for _, service := range s {
		var metrics []aws2.Metrics
		for _, metric := range service.Metrics {
			metrics = append(metrics, aws2.Metrics{
				Name:          metric.Name,
				Amount:        metric.Amount,
				Unit:          metric.Unit,
				UsageQuantity: metric.UsageQuantity,
				NumericAmount: metric.NumericAmount,
			})
		}

		services = append(services, aws2.Service{
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

func ConvertToForecastStdoutType(r aws2.ForecastPrintData,
	filteredBy string) aws2.ForecastStdoutType {
	var forecast []aws2.ForecastResults
	for _, v := range r.Forecast.ForecastResultsByTime {
		forecast = append(forecast, aws2.ForecastResults{
			TimePeriod: aws2.DateInterval{
				Start: *v.TimePeriod.Start,
				End:   *v.TimePeriod.End,
			},
			MeanValue:                    *v.MeanValue,
			PredictionIntervalLowerBound: *v.PredictionIntervalLowerBound,
			PredictionIntervalUpperBound: *v.PredictionIntervalUpperBound,
		})
	}

	return aws2.ForecastStdoutType{
		Forecast: forecast,
		Total: aws2.Total{
			Amount: *r.Forecast.Total.Amount,
			Unit:   *r.Forecast.Total.Unit,
		},
		FilteredBy: filteredBy,
	}

}
