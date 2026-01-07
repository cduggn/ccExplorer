package writer

import (
	"fmt"

	"github.com/cduggn/ccexplorer/internal/types"
)

// WriterFactory provides type-safe writer creation
type WriterFactory struct {
	config *Config
}

// NewWriterFactory creates a new factory with the given configuration
func NewWriterFactory(opts ...Option) *WriterFactory {
	return &WriterFactory{
		config: NewConfig(opts...),
	}
}

// CreateCostUsageWriter creates a writer for cost and usage data
func (f *WriterFactory) CreateCostUsageWriter(format FormatType) (interface{}, error) {
	switch format {
	case FormatTable:
		return NewCostUsageTableWriter(f.configToOptions()...), nil
	case FormatCSV:
		return NewCostUsageCSVWriter(f.configToOptions()...), nil
	case FormatChart:
		return NewCostUsageChartWriter(f.configToOptions()...), nil
	case FormatVector:
		return NewCostUsageVectorWriter(f.configToOptions()...), nil
	default:
		return nil, fmt.Errorf("unsupported format: %v", format)
	}
}

// CreateForecastWriter creates a writer for forecast data
func (f *WriterFactory) CreateForecastWriter() *ForecastTableWriter {
	return NewForecastTableWriter(f.configToOptions()...)
}

// WriteCostUsage is a type-safe convenience method for writing cost usage data
func (f *WriterFactory) WriteCostUsage(format FormatType, data types.CostAndUsageOutputType) error {
	switch format {
	case FormatTable:
		writer := NewCostUsageTableWriter(f.configToOptions()...)
		return writer.Write(data)
	case FormatCSV:
		writer := NewCostUsageCSVWriter(f.configToOptions()...)
		return writer.Write(data)
	case FormatChart:
		writer := NewCostUsageChartWriter(f.configToOptions()...)
		return writer.Write(data)
	case FormatVector:
		writer := NewCostUsageVectorWriter(f.configToOptions()...)
		return writer.Write(data)
	default:
		return fmt.Errorf("unsupported format: %v", format)
	}
}

// WriteForecast is a type-safe convenience method for writing forecast data
func (f *WriterFactory) WriteForecast(data types.ForecastPrintData) error {
	writer := NewForecastTableWriter(f.configToOptions()...)
	return writer.Write(data)
}

// configToOptions converts the factory's config to a slice of options
func (f *WriterFactory) configToOptions() []Option {
	return []Option{
		WithOutputDir(f.config.OutputDir),
		WithSortBy(f.config.SortBy),
		WithVectorDBConfig(f.config.OpenAIAPIKey, f.config.PineconeAPIKey, f.config.PineconeIndex),
		WithHTTPClient(f.config.HTTPClient),
	}
}
