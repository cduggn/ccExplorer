package writer

import (
	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/go-echarts/go-echarts/v2/components"
)

// TableOutput represents data formatted for table display
type TableOutput struct {
	Headers      []string
	Rows         [][]string
	Footer       []string
	Title        string
	Total        string
	Granularity  string
	DateRange    string
	RowCount     int
	Metric       string
	Unit         string
	SortedBy     string
	HasTagColumn bool // Whether the Tag/Dimension column has any values
}

// CSVOutput represents data formatted for CSV export
type CSVOutput struct {
	Headers []string
	Rows    [][]string
	Filename string
}

// ChartOutput represents data formatted for chart visualization
type ChartOutput struct {
	Page      *components.Page
	Title     string
	Filename  string
}

// VectorOutput represents data formatted for vector database storage
type VectorOutput struct {
	Items       []*types.VectorStoreItem
	IndexName   string
	BatchSize   int
}

// ForecastTableOutput represents forecast data for table display
type ForecastTableOutput struct {
	Headers    []string
	Rows       [][]string
	Footer     []string
	FilterInfo string
	Total      types.Total
}

// NewTableOutput creates a new TableOutput instance
func NewTableOutput(headers []string, rows [][]string, total string, granularity string, dateRange string, metric string, unit string, sortedBy string, hasTagColumn bool) *TableOutput {
	return &TableOutput{
		Headers:      headers,
		Rows:         rows,
		Total:        total,
		Granularity:  granularity,
		DateRange:    dateRange,
		RowCount:     len(rows),
		Metric:       metric,
		Unit:         unit,
		SortedBy:     sortedBy,
		HasTagColumn: hasTagColumn,
	}
}

// NewCSVOutput creates a new CSVOutput instance
func NewCSVOutput(headers []string, rows [][]string, filename string) *CSVOutput {
	return &CSVOutput{
		Headers:  headers,
		Rows:     rows,
		Filename: filename,
	}
}

// NewChartOutput creates a new ChartOutput instance
func NewChartOutput(page *components.Page, title, filename string) *ChartOutput {
	return &ChartOutput{
		Page:     page,
		Title:    title,
		Filename: filename,
	}
}

// NewVectorOutput creates a new VectorOutput instance
func NewVectorOutput(items []*types.VectorStoreItem, indexName string) *VectorOutput {
	return &VectorOutput{
		Items:     items,
		IndexName: indexName,
		BatchSize: 25, // default batch size
	}
}

// NewForecastTableOutput creates a new ForecastTableOutput instance
func NewForecastTableOutput(headers []string, rows [][]string, filterInfo string, total types.Total) *ForecastTableOutput {
	return &ForecastTableOutput{
		Headers:    headers,
		Rows:       rows,
		FilterInfo: filterInfo,
		Total:      total,
	}
}