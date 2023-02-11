package printer

import (
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/jedib0t/go-pretty/v6/table"
)

type PrintWriterType int

type SortBy int

const (
	Amount SortBy = iota
	Date
)

const (
	Stdout PrintWriterType = iota
	CSV
	Chart
)

type PrinterError struct {
	msg string
}

func (e PrinterError) Error() string {
	return e.msg
}

type Printer interface {
	Print(interface{}, interface{}) error
}
type StdoutPrinter struct {
	Variant string
}

type CsvPrinter struct {
	Variant string
}

type ChartPrinter struct {
	Variant string
}

type CostAndUsage struct {
	Rows  []table.Row
	Total float64
}

type CostAndUsageReport struct {
	Services    map[int]Service
	Start       string
	End         string
	Granularity string
	Dimensions  []string
	Tags        []string
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

type Renderer struct {
}

type ForecastPrintData struct {
	Forecast *costexplorer.GetCostForecastOutput
	Filters  []string
}

type CostAndUsageOutputType struct {
	Services    map[int]Service
	Granularity string
	Start       string
	End         string
	Dimensions  []string
	Tags        []string
	SortBy      string
}

type ChartData struct {
	StartDate      string
	EndDate        string
	Granularity    string
	DimensionOrTag string
	Title          string
	SubTitle       string
	NumericValues  float64
}

func (c CostAndUsageReport) Len() int {
	return len(c.Services)
}

func (c CostAndUsageReport) Less(i, j int) bool {
	return c.Services[i].Metrics[0].NumericAmount > c.Services[j].Metrics[0].NumericAmount
}

func (c CostAndUsageReport) Swap(i, j int) {
	c.Services[i], c.Services[j] = c.Services[j], c.Services[i]
}

func (c CostAndUsageReport) Equals(c2 CostAndUsageReport) bool {
	if c.Start != c2.Start || c.End != c2.End {
		return false
	}
	for k, v := range c.Services {
		v2, ok := c2.Services[k]
		if !ok {
			return false
		}
		if !v.Equals(v2) {
			return false
		}
	}
	return true
}

func (s Service) Equals(s2 Service) bool {
	if s.Start != s2.Start || s.End != s2.End {
		return false
	}
	if len(s.Keys) != len(s2.Keys) {
		return false
	}
	for i, v := range s.Keys {
		if v != s2.Keys[i] {
			return false
		}
	}
	if len(s.Metrics) != len(s2.Metrics) {
		return false
	}
	for i, v := range s.Metrics {
		if !v.Equals(s2.Metrics[i]) {
			return false
		}
	}
	return true
}

func (s Service) Less(i, j int) bool {
	return s.Metrics[i].NumericAmount > s.Metrics[j].NumericAmount
}

func (m Metrics) Equals(m2 Metrics) bool {
	if m.Name != m2.Name || m.Amount != m2.Amount || m.Unit != m2.Unit {
		return false
	}
	return true
}
