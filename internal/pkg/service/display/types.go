package display

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
	Amount float64
	Unit   string
}
