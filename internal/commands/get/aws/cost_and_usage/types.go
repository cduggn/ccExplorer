package cost_and_usage

import "github.com/cduggn/ccexplorer/internal/commands/get/aws/custom_flags"

type CommandLineInput struct {
	GroupByValues       *custom_flags.DimensionAndTagFlag
	GroupByTag          []string
	FilterByValues      *custom_flags.DimensionAndTagFilterFlag
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
