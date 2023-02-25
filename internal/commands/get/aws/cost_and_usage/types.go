package cost_and_usage

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
