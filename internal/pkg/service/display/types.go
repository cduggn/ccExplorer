package display

import "sort"

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
	Name          string
	Amount        string
	NumericAmount float64
	Unit          string
}

func (c CostAndUsageReport) Len() int {
	return len(c.Services)
}

func (c CostAndUsageReport) Less(i, j int) bool {
	return c.Services[i].Metrics[0].NumericAmount > c.Services[j].Metrics[0].NumericAmount
}

func (c CostAndUsageReport) Swap(i, j int) {
	c.Services[i], c.Services[j] = c.Services[j], c.Services[i]
}

func (c CostAndUsageReport) Equals(c2 CostAndUsageReport) bool {
	if c.Start != c2.Start || c.End != c2.End {
		return false
	}
	for k, v := range c.Services {
		v2, ok := c2.Services[k]
		if !ok {
			return false
		}
		if !v.Equals(v2) {
			return false
		}
	}
	return true
}

func (s Service) Equals(s2 Service) bool {
	if s.Start != s2.Start || s.End != s2.End {
		return false
	}
	if len(s.Keys) != len(s2.Keys) {
		return false
	}
	for i, v := range s.Keys {
		if v != s2.Keys[i] {
			return false
		}
	}
	if len(s.Metrics) != len(s2.Metrics) {
		return false
	}
	for i, v := range s.Metrics {
		if !v.Equals(s2.Metrics[i]) {
			return false
		}
	}
	return true
}

func (m Metrics) Equals(m2 Metrics) bool {
	if m.Name != m2.Name || m.Amount != m2.Amount || m.Unit != m2.Unit {
		return false
	}
	return true
}

func SortServicesByMetricAmount(r *CostAndUsageReport) {
	// Create a slice of key-value pairs
	pairs := make([]struct {
		Key   int
		Value Service
	}, len(r.Services))
	i := 0
	for k, v := range r.Services {
		pairs[i] = struct {
			Key   int
			Value Service
		}{k, v}
		i++
	}

	// Sort the slice by the Value.Metrics[0].Amount field
	sort.SliceStable(pairs, func(i, j int) bool {
		return pairs[i].Value.Metrics[0].Amount > pairs[j].Value.Metrics[0].
			Amount
	})

	for i, pair := range pairs {
		r.Services[i] = pair.Value
	}

}
