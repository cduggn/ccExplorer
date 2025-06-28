package types

import (
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Error struct {
	Msg string
}

func (e Error) Error() string {
	return e.Msg
}

type PrintWriterType int

type SortBy int

const (
	Stdout PrintWriterType = iota
	CSV
	Chart
	OpenAPI
	Pinecone
)

type InputType struct {
	Services     []Service
	Granularity  string
	Start        string
	End          string
	Dimensions   []string
	Tags         []string
	SortBy       string
	OpenAIAPIKey string
}

type Service struct {
	Keys    []string
	Name    string
	Metrics []Metrics
	Start   string
	End     string
}

type Metrics struct {
	Name          string
	Amount        string
	NumericAmount float64
	Unit          string
	UsageQuantity float64
}

type Table interface {
	Writer(interface{})
	Header()
	Footer(row table.Row)
	AddRows(rows []table.Row)
	Style()
}

type CostAndUsage struct {
	Rows  []table.Row
	Total string
}

type CostAndUsageStdoutType struct {
	Granularity string
	Services    []Service
}

type ForecastStdoutType struct {
	Forecast   []ForecastResults
	FilteredBy string
	Total      Total
}

type Total struct {
	Amount string
	Unit   string
}

type ForecastResults struct {
	MeanValue                    string
	PredictionIntervalLowerBound string
	PredictionIntervalUpperBound string
	TimePeriod                   DateInterval
}

type DateInterval struct {
	End   string
	Start string
}
type CostAndUsageOutputType struct {
	Services       map[int]Service
	Granularity    string
	Start          string
	End            string
	Dimensions     []string
	Tags           []string
	SortBy         string
	OpenAIAPIKey   string
	PineconeAPIKey string
	PineconeIndex  string
}

type ForecastPrintData struct {
	Forecast *costexplorer.GetCostForecastOutput
	Filters  []string
}

// Generic output types for improved type safety and reduced interface{} usage

// GenericOutputType provides a type-safe output container
type GenericOutputType[T any] struct {
	Data        T
	Granularity string
	Start       string
	End         string
	Metadata    map[string]interface{}
}

func NewGenericOutputType[T any](data T, granularity, start, end string) *GenericOutputType[T] {
	return &GenericOutputType[T]{
		Data:        data,
		Granularity: granularity,
		Start:       start,
		End:         end,
		Metadata:    make(map[string]interface{}),
	}
}

func (g *GenericOutputType[T]) SetMetadata(key string, value interface{}) {
	g.Metadata[key] = value
}

func (g *GenericOutputType[T]) GetMetadata(key string) (interface{}, bool) {
	value, exists := g.Metadata[key]
	return value, exists
}

// GenericServiceOutput provides type-safe service output
type GenericServiceOutput[T any] struct {
	Services map[int]T
	Total    float64
	Count    int
}

func NewGenericServiceOutput[T any]() *GenericServiceOutput[T] {
	return &GenericServiceOutput[T]{
		Services: make(map[int]T),
	}
}

func (g *GenericServiceOutput[T]) AddService(index int, service T) {
	g.Services[index] = service
	g.Count++
}

func (g *GenericServiceOutput[T]) GetServices() []T {
	result := make([]T, 0, len(g.Services))
	for i := 0; i < len(g.Services); i++ {
		if service, exists := g.Services[i]; exists {
			result = append(result, service)
		}
	}
	return result
}

// GenericTableRow provides type-safe table row operations
type GenericTableRow[T any] struct {
	Data   T
	Format func(T) []string
}

func NewGenericTableRow[T any](data T, formatter func(T) []string) *GenericTableRow[T] {
	return &GenericTableRow[T]{
		Data:   data,
		Format: formatter,
	}
}

func (r *GenericTableRow[T]) ToStringSlice() []string {
	return r.Format(r.Data)
}
