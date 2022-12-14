package aws

type Time struct {
	Start string
	End   string
}

type CostAndUsageRequest struct {
	Granularity     string
	GroupBy         []string
	Tag             string
	Time            Time
	IsFilterEnabled bool
	FilterType      string
	TagFilterValue  string
	Rates           []string
	ExcludeCredits  bool
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
