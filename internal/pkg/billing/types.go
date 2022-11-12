package billing

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
}

type CostAndUsageReport struct {
	Services map[int]Service
	Start    string
	End      string
}

type Service struct {
	Keys    []string
	Name    string
	Metrics []Metrics
}

type Metrics struct {
	Name   string
	Amount string
	Unit   string
}
