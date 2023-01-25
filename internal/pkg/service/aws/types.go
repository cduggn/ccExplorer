package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/ccexplorer/internal/pkg/storage"
)

type AWSClient interface {
	GetCostAndUsage(ctx context.Context, api GetCostAndUsageAPI,
		req CostAndUsageRequestType) (
		*costexplorer.GetCostAndUsageOutput,
		error)
	GetDimensionValues(ctx context.Context, api GetDimensionValuesAPI,
		d GetDimensionValuesRequest) ([]string, error)
	GetCostForecast(ctx context.Context,
		api GetCostForecastAPI, req GetCostForecastRequest) (
		*costexplorer.GetCostForecastOutput, error)
}

type APIClient struct {
	*costexplorer.Client
}

type DatabaseManager struct {
	dbClient *storage.CostDataStorage
}

type DBError struct {
	msg string
}

func (e DBError) Error() string {
	return e.msg
}

type APIError struct {
	msg string
}

func (e APIError) Error() string {
	return e.msg
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
	Granularity                string
	GroupBy                    []string
	Tag                        string
	Time                       Time
	IsFilterByTagEnabled       bool
	IsFilterByDimensionEnabled bool
	TagFilterValue             string
	DimensionFilter            map[string]string
	ExcludeDiscounts           bool
	Alias                      string
	Rates                      []string
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

func (c CostAndUsageRequestType) Equals(c2 CostAndUsageRequestType) bool {
	if c.Granularity != c2.Granularity {
		return false
	}
	if c.Tag != c2.Tag {
		return false
	}
	if c.Time.Start != c2.Time.Start {
		return false
	}
	if c.Time.End != c2.Time.End {
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
		if c2.DimensionFilter[k] != v {
			return false
		}
	}
	if c.ExcludeDiscounts != c2.ExcludeDiscounts {
		return false
	}
	if c.Alias != c2.Alias {
		return false
	}
	return true
}
