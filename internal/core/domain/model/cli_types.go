package model

import "github.com/spf13/cobra"

type Command interface {
	Run(cmd *cobra.Command, args []string) error
	SynthesizeRequest(input interface{}) (interface{}, error)
	InputHandler(cmd *cobra.Command) interface{}
	Execute(req interface{}) error
	DefineFlags()
}

type CostDataReader interface {
	ExtractGroupBySelections() ([]string, []string)
	ExtractFilterBySelection() (FilterBySelections, error)
	ExtractStartAndEndDates() (string, string, error)
	ExtractPrintPreferences() PrintOptions
}

type CommandLineInput struct {
	GroupByDimension    []string
	GroupByTag          []string
	FilterByValues      map[string]string
	IsFilterByTag       bool
	TagFilterValue      string
	IsFilterByDimension bool
	Start               string
	End                 string
	ExcludeDiscounts    bool
	Interval            string
	PrintFormat         string
	Metrics             []string
	SortByDate          bool
	OpenAIAPIKey        string
}

type FilterBySelections struct {
	Tags                string
	Dimensions          map[string]string
	IsFilterByTag       bool
	IsFilterByDimension bool
}

type PrintOptions struct {
	IsSortByDate     bool
	ExcludeDiscounts bool
	Format           string
	OpenAIKey        string
	Granularity      string
	Metric           string
}

type ForecastCommandLineInput struct {
	FilterByValues          Filter
	Granularity             string
	PredictionIntervalLevel int32
	Start                   string
	End                     string
}
