package stdout

import (
	"github.com/jedib0t/go-pretty/v6/table"
)

type Table interface {
	Writer(interface{})
	Header()
	Footer(row table.Row)
	AddRows(rows []table.Row)
}

type CostAndUsageTable struct {
	Table table.Writer
}

type ForecastTable struct {
	Table table.Writer
}

type CostAndUsage struct {
	Rows  []table.Row
	Total string
}

type CostAndUsageStdoutType struct {
	Granularity string
	Services    []Service
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
