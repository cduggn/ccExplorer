package aws

type Time struct {
	Start string
	End   string
}

type Dimension struct {
	Key   string
	Value []string
}

type Tag struct {
	Key   string
	Value []string
}

type Filter struct {
	Dimensions []Dimension
	Tags       []Tag
}

type GetCostForecastRequest struct {
	Time                    Time
	Granularity             string
	Metric                  string
	Filter                  Filter
	PredictionIntervalLevel int32
}

type GetCostForecastReport struct {
}

type CostAndUsageRequestType struct {
	Granularity      string
	GroupBy          []string
	Tag              string
	Time             Time
	IsFilterEnabled  bool
	FilterType       string
	TagFilterValue   string
	Rates            []string
	ExcludeDiscounts bool
}

type CostAndUsageRequestWithResourcesType struct {
	Granularity      string
	GroupBy          []string
	Tag              string
	Time             Time
	IsFilterEnabled  bool
	FilterType       string
	TagFilterValue   string
	Rates            []string
	ExcludeDiscounts bool
}

type CostAndUsageReport struct {
	Services    map[int]Service
	Start       string
	End         string
	Granularity string
}

type Service struct {
	Keys    []string
	Name    string
	Metrics []Metrics
	Start   string
	End     string
}

type Metrics struct {
	Name   string
	Amount string
	Unit   string
}
