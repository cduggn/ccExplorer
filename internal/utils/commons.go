package utils_new

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/cduggn/ccexplorer/internal/pinecone"
	types2 "github.com/cduggn/ccexplorer/internal/types"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

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

func ToPrintWriterType(s string) types2.PrintWriterType {
	switch s {
	case "csv":
		return types2.CSV
	case "stdout":
		return types2.Stdout
	case "chart":
		return types2.Chart
	case "pinecone":
		return types2.Pinecone
	default:
		return types2.Stdout
	}
}

func NewFile(dir string, file string) (*os.File, error) {
	filePath := BuildOutputFilePath(dir, file)
	return os.Create(filePath)
}

func BuildOutputFilePath(dir string, fileName string) string {
	return dir + "/" + fileName
}

func DefaultEndDate(f func(date time.Time) string) string {
	return f(time.Now())
}

func Format(date time.Time) string {
	return date.Format("2006-01-02")
}

// DefaultStartDate function which returns  the first day of the previous month

func DefaultStartDate(dayOfCurrentMonth func(time time.Time) int,
	subtractDays func(time time.Time, days int) string) string {
	today := time.Now()
	firstDayOfCurrentMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	firstDayOfPreviousMonth := firstDayOfCurrentMonth.AddDate(0, -1, 0)
	dayOfMonth := dayOfCurrentMonth(firstDayOfPreviousMonth)

	return subtractDays(firstDayOfPreviousMonth, dayOfMonth-1)
}

//func DefaultStartDate(d func(time time.Time) int, s func(time time.Time, days int) string) string {
//	today := time.Now()
//	dayOfMonth := d(today)
//
//	if dayOfMonth == 1 {
//		return s(today, 1)
//	}
//	return s(today, dayOfMonth-1) // subtract 1 to get the first day of the month
//}

func DayOfCurrentMonth(time time.Time) int {
	return time.Day()
}

func SubtractDays(today time.Time, days int) string {
	return today.AddDate(0, 0, -days).Format("2006-01-02")
}

func LastDayOfMonth() string {
	return time.Now().AddDate(0, 1, -1).Format("2006-01-02")
}

func SortByFn(sortByDate bool) string {
	if sortByDate {
		return "date"
	}
	return "cost"
}

func SplitCommaSeparatedString(value string) []string {
	var args []string
	if strings.Contains(value, ",") {
		args = strings.Split(value, ",")
	} else {
		args = []string{value}
	}
	return args
}

func SplitNameValuePair(value string) ([]string, error) {
	parts := strings.Split(value, "=")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid group by flag: %s", value)
	}
	return parts, nil
}

func MaxSupportedRows(rows [][]string, maxRows int) int {
	if len(rows) > maxRows {
		return maxRows
	}
	return len(rows)
}

func ToCostAndUsageOutputType(r *costexplorer.GetCostAndUsageOutput,
	u types2.CostAndUsageRequestType) types2.CostAndUsageOutputType {
	return CurateCostAndUsageReport(r, u)
}

func CurateCostAndUsageReport(
	d *costexplorer.GetCostAndUsageOutput, query types2.CostAndUsageRequestType) types2.CostAndUsageOutputType {

	c := types2.CostAndUsageOutputType{
		Services:       make(map[int]types2.Service),
		Granularity:    query.Granularity,
		Dimensions:     query.GroupBy,
		Tags:           query.GroupByTag,
		Start:          query.Time.Start,
		End:            query.Time.End,
		OpenAIAPIKey:   query.OpenAIAPIKey,
		PineconeAPIKey: query.PineconeAPIKey,
		PineconeIndex:  query.PineconeIndex,
	}

	c.Services = ResultsToServicesMap(d.ResultsByTime)
	return c
}

func ResultsToServicesMap(res []types.ResultByTime) map[int]types2.Service {
	services := make(map[int]types2.Service)
	count := 0
	for _, v := range res {
		for _, g := range v.Groups {
			keys := append(make([]string, 0), g.Keys...)
			service := types2.Service{
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

func MetricsToService(m map[string]types.MetricValue) []types2.Metrics {
	var metrics []types2.Metrics
	for k, v := range m {
		metrics = append(metrics, types2.Metrics{
			Name:          k,
			Amount:        *v.Amount,
			NumericAmount: ConvertToFloat(*v.Amount),
			Unit:          *v.Unit,
		})
	}
	return metrics
}

func ConvertServiceToSlice(s types2.Service, granularity string) [][]string {
	var r [][]string
	for _, v := range s.Metrics {
		t := []string{s.Keys[0], ReturnIfPresent(s.Keys), v.Name,
			granularity, s.Start, s.End,
			v.Amount, v.Unit}
		r = append(r, t)
	}
	return r
}

func SortFunction(sortBy string) func(r map[int]types2.Service) []types2.Service {
	switch sortBy {
	case "date":
		return SortServicesByStartDate
	case "cost":
		return SortServicesByMetricAmount
	default:
		return SortServicesByMetricAmount
	}
}

func SortServicesByStartDate(r map[int]types2.Service) []types2.Service {
	// Create a slice of key-value pairs
	pairs := make([]struct {
		Key   int
		Value types2.Service
	}, len(r))
	i := 0
	for k, v := range r {
		pairs[i] = struct {
			Key   int
			Value types2.Service
		}{k, v}
		i++
	}

	sort.SliceStable(pairs, func(i, j int) bool {

		t1, _ := time.Parse("2006-01-02", pairs[i].Value.Start)
		t2, _ := time.Parse("2006-01-02", pairs[j].Value.Start)
		return t1.After(t2)
	})

	result := make([]types2.Service, len(pairs))
	for i, pair := range pairs {
		result[i] = pair.Value
	}
	return result
}

func SortServicesByMetricAmount(r map[int]types2.Service) []types2.Service {
	// Create a slice of key-value pairs
	pairs := make([]struct {
		Key   int
		Value types2.Service
	}, len(r))
	i := 0
	for k, v := range r {
		pairs[i] = struct {
			Key   int
			Value types2.Service
		}{k, v}
		i++
	}

	// Sort the slice by the Value.Metrics[0].Amount field
	sort.SliceStable(pairs, func(i, j int) bool {
		return pairs[i].Value.Metrics[0].NumericAmount > pairs[j].Value.
			Metrics[0].NumericAmount
	})

	result := make([]types2.Service, len(pairs))
	for i, pair := range pairs {
		result[i] = pair.Value
	}
	return result
}

func ConvertServiceMapToArray(s map[int]types2.Service,
	granularity string) [][]string {
	var rows [][]string
	for _, v := range s {
		rows = append(rows, ConvertServiceToSlice(v, granularity)...)
	}
	return rows
}

func ConvertToStdoutType(s []types2.Service,
	granularity string) types2.CostAndUsageStdoutType {

	outputType := types2.CostAndUsageStdoutType{
		Granularity: granularity,
	}

	var services []types2.Service
	for _, v := range s {
		var metrics []types2.Metrics
		for _, m := range v.Metrics {
			metrics = append(metrics, types2.Metrics{
				Name:          m.Name,
				Amount:        m.Amount,
				Unit:          m.Unit,
				NumericAmount: m.NumericAmount,
			})
		}
		services = append(services, types2.Service{
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

func ConvertToChartInputType(r types2.CostAndUsageOutputType,
	s []types2.Service) types2.InputType {

	input := types2.InputType{
		Granularity: r.Granularity,
		Start:       r.Start,
		End:         r.End,
		Dimensions:  r.Dimensions,
		Tags:        r.Tags,
	}

	var services []types2.Service
	for _, service := range s {
		var metrics []types2.Metrics
		for _, metric := range service.Metrics {
			metrics = append(metrics, types2.Metrics{
				Name:          metric.Name,
				Amount:        metric.Amount,
				Unit:          metric.Unit,
				UsageQuantity: metric.UsageQuantity,
				NumericAmount: metric.NumericAmount,
			})
		}

		services = append(services, types2.Service{
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

func ConvertToForecastStdoutType(r types2.ForecastPrintData,
	filteredBy string) types2.ForecastStdoutType {
	var forecast []types2.ForecastResults
	for _, v := range r.Forecast.ForecastResultsByTime {
		forecast = append(forecast, types2.ForecastResults{
			TimePeriod: types2.DateInterval{
				Start: *v.TimePeriod.Start,
				End:   *v.TimePeriod.End,
			},
			MeanValue:                    *v.MeanValue,
			PredictionIntervalLowerBound: *v.PredictionIntervalLowerBound,
			PredictionIntervalUpperBound: *v.PredictionIntervalUpperBound,
		})
	}

	return types2.ForecastStdoutType{
		Forecast: forecast,
		Total: types2.Total{
			Amount: *r.Forecast.Total.Amount,
			Unit:   *r.Forecast.Total.Unit,
		},
		FilteredBy: filteredBy,
	}

}

func EncodeString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	hashedString := hex.EncodeToString(hashed)
	return hashedString
}

func ConvertToPineconeStruct(items []*types2.VectorStoreItem) []pinecone.PineconeStruct {
	var pineconeStruct []pinecone.PineconeStruct
	for _, v := range items {
		pineconeStruct = append(pineconeStruct, pinecone.PineconeStruct{
			ID:     v.ID,
			Values: v.EmbeddingVector,
			Metadata: pinecone.Metadata{
				PageContent: v.EmbeddingText,
				Source:      "aws cost explorer",
				Dimensions:  v.Metadata.Dimensions,
				Start:       v.Metadata.StartDate,
				End:         v.Metadata.EndDate,
				Cost:        v.Metadata.Cost,
			},
		})
	}
	return pineconeStruct
}
