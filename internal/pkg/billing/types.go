package billing

type CostAndUsageRequest struct {
	Granularity string
	GroupBy     []string
	Tag         string
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
