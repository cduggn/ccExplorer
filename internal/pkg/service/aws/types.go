package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/cloudcost/internal/pkg/storage"
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
	DimensionTag               string
	Tag                        string
	Time                       Time
	IsFilterByTagEnabled       bool
	IsFilterByDimensionEnabled bool
	TagFilterValue             string
	DimensionFilterName        string
	DimensionFilterValue       string
	Rates                      []string
	ExcludeDiscounts           bool
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
