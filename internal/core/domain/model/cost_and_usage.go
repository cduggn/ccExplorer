package model

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
)

type GetCostAndUsageAPI interface {
	GetCostAndUsage(ctx context.Context,
		optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error)
}

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

type CostAndUsageRequestType struct {
	Granularity                string
	GroupBy                    []string
	GroupByTag                 []string
	Time                       Time
	IsFilterByTagEnabled       bool
	IsFilterByDimensionEnabled bool
	TagFilterValue             string
	DimensionFilter            map[string]string
	ExcludeDiscounts           bool
	Alias                      string
	Rates                      []string
	PrintFormat                string
	Metrics                    []string
	SortByDate                 bool
	OpenAIAPIKey               string
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

type GetDimensionValuesRequest struct {
	Dimension string
	Time      Time
}

func (t Time) Equals(other Time) bool {
	return t.Start == other.Start && t.End == other.End
}

func (c CostAndUsageRequestType) Equals(c2 CostAndUsageRequestType) bool {
	if c.Granularity != c2.Granularity {
		return false
	}
	if !c.Time.Equals(c2.Time) {
		return false
	}
	if c.IsFilterByTagEnabled != c2.IsFilterByTagEnabled {
		return false
	}
	if c.IsFilterByDimensionEnabled != c2.IsFilterByDimensionEnabled {
		return false
	}
	if c.TagFilterValue != c2.TagFilterValue {
		return false
	}
	if len(c.DimensionFilter) != len(c2.DimensionFilter) {
		return false
	}
	for k, v := range c.DimensionFilter {
		if v2, ok := c2.DimensionFilter[k]; !ok || v != v2 {
			return false
		}
	}
	if c.ExcludeDiscounts != c2.ExcludeDiscounts {
		return false
	}
	if c.Alias != c2.Alias {
		return false
	}
	if len(c.Rates) != len(c2.Rates) {
		return false
	}
	for i, v := range c.Rates {
		if v2 := c2.Rates[i]; v != v2 {
			return false
		}
	}
	return true
}