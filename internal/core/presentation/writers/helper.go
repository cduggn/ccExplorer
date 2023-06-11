package writers

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	model "github.com/cduggn/ccexplorer/internal/core/domain/model"
	"github.com/cduggn/ccexplorer/internal/core/util"
	"sort"
	"time"
)

func CreateSubTitle(granularity string, start string, end string) string {
	return fmt.Sprintf("Response granularity: %s. Timeframe: %s-%s",
		granularity,
		start, end)
}

func ToCostAndUsageOutputType(r *costexplorer.GetCostAndUsageOutput,
	u model.CostAndUsageRequestType) model.CostAndUsageOutputType {
	return CurateCostAndUsageReport(r, u)
}

func CurateCostAndUsageReport(
	d *costexplorer.GetCostAndUsageOutput, query model.CostAndUsageRequestType) model.CostAndUsageOutputType {

	c := model.CostAndUsageOutputType{
		Services:     make(map[int]model.Service),
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

func ResultsToServicesMap(res []types.ResultByTime) map[int]model.Service {
	services := make(map[int]model.Service)
	count := 0
	for _, v := range res {
		for _, g := range v.Groups {
			keys := append(make([]string, 0), g.Keys...)
			service := model.Service{
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

func MetricsToService(m map[string]types.MetricValue) []model.Metrics {
	var metrics []model.Metrics
	for k, v := range m {
		metrics = append(metrics, model.Metrics{
			Name:          k,
			Amount:        *v.Amount,
			NumericAmount: util.ConvertToFloat(*v.Amount),
			Unit:          *v.Unit,
		})
	}
	return metrics
}

func ConvertServiceToSlice(s model.Service, granularity string) [][]string {
	var r [][]string
	for _, v := range s.Metrics {
		t := []string{s.Keys[0], util.ReturnIfPresent(s.Keys), v.Name,
			granularity, s.Start, s.End,
			v.Amount, v.Unit}
		r = append(r, t)
	}
	return r
}

func SortFunction(sortBy string) func(r map[int]model.Service) []model.Service {
	switch sortBy {
	case "date":
		return SortServicesByStartDate
	case "cost":
		return SortServicesByMetricAmount
	default:
		return SortServicesByMetricAmount
	}
}

func SortServicesByStartDate(r map[int]model.Service) []model.Service {
	// Create a slice of key-value pairs
	pairs := make([]struct {
		Key   int
		Value model.Service
	}, len(r))
	i := 0
	for k, v := range r {
		pairs[i] = struct {
			Key   int
			Value model.Service
		}{k, v}
		i++
	}

	sort.SliceStable(pairs, func(i, j int) bool {

		t1, _ := time.Parse("2006-01-02", pairs[i].Value.Start)
		t2, _ := time.Parse("2006-01-02", pairs[j].Value.Start)
		return t1.After(t2)
	})

	result := make([]model.Service, len(pairs))
	for i, pair := range pairs {
		result[i] = pair.Value
	}
	return result
}

func SortServicesByMetricAmount(r map[int]model.Service) []model.Service {
	// Create a slice of key-value pairs
	pairs := make([]struct {
		Key   int
		Value model.Service
	}, len(r))
	i := 0
	for k, v := range r {
		pairs[i] = struct {
			Key   int
			Value model.Service
		}{k, v}
		i++
	}

	// Sort the slice by the Value.Metrics[0].Amount field
	sort.SliceStable(pairs, func(i, j int) bool {
		return pairs[i].Value.Metrics[0].NumericAmount > pairs[j].Value.
			Metrics[0].NumericAmount
	})

	result := make([]model.Service, len(pairs))
	for i, pair := range pairs {
		result[i] = pair.Value
	}
	return result
}

func ConvertServiceMapToArray(s map[int]model.Service,
	granularity string) [][]string {
	var rows [][]string
	for _, v := range s {
		rows = append(rows, ConvertServiceToSlice(v, granularity)...)
	}
	return rows
}

func ConvertServiceSliceToArray(s []model.Service, granularity string) [][]string {
	var rows [][]string
	for _, v := range s {
		rows = append(rows, ConvertServiceToSlice(v, granularity)...)
	}
	return rows
}

func ConvertToStdoutType(s []model.Service,
	granularity string) model.CostAndUsageStdoutType {

	outputType := model.CostAndUsageStdoutType{
		Granularity: granularity,
	}

	var services []model.Service
	for _, v := range s {
		var metrics []model.Metrics
		for _, m := range v.Metrics {
			metrics = append(metrics, model.Metrics{
				Name:          m.Name,
				Amount:        m.Amount,
				Unit:          m.Unit,
				NumericAmount: m.NumericAmount,
			})
		}
		services = append(services, model.Service{
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

func ConvertToChartInputType(r model.CostAndUsageOutputType,
	s []model.Service) model.InputType {

	input := model.InputType{
		Granularity: r.Granularity,
		Start:       r.Start,
		End:         r.End,
		Dimensions:  r.Dimensions,
		Tags:        r.Tags,
	}

	var services []model.Service
	for _, service := range s {
		var metrics []model.Metrics
		for _, metric := range service.Metrics {
			metrics = append(metrics, model.Metrics{
				Name:          metric.Name,
				Amount:        metric.Amount,
				Unit:          metric.Unit,
				UsageQuantity: metric.UsageQuantity,
				NumericAmount: metric.NumericAmount,
			})
		}

		services = append(services, model.Service{
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

func ConvertToForecastStdoutType(r model.ForecastPrintData,
	filteredBy string) model.ForecastStdoutType {
	var forecast []model.ForecastResults
	for _, v := range r.Forecast.ForecastResultsByTime {
		forecast = append(forecast, model.ForecastResults{
			TimePeriod: model.DateInterval{
				Start: *v.TimePeriod.Start,
				End:   *v.TimePeriod.End,
			},
			MeanValue:                    *v.MeanValue,
			PredictionIntervalLowerBound: *v.PredictionIntervalLowerBound,
			PredictionIntervalUpperBound: *v.PredictionIntervalUpperBound,
		})
	}

	return model.ForecastStdoutType{
		Forecast: forecast,
		Total: model.Total{
			Amount: *r.Forecast.Total.Amount,
			Unit:   *r.Forecast.Total.Unit,
		},
		FilteredBy: filteredBy,
	}

}
