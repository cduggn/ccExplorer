package model

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
