package aws

import (
	"github.com/cduggn/cloudcost/internal/pkg/service/aws"
	"github.com/spf13/cobra"
)

type ForecastRequestType struct {
	Granularity     string
	GroupBy         []string
	IsFilterEnabled bool
	Tag             string
	FilterType      string
	TagFilterValue  string
}

func CostForecast(cmd *cobra.Command, args []string) error {

	res, _ := aws.GetCostForecast(aws.GetCostForecastRequestType{})

	aws.PrintGetCostForecastReport(res)
	return nil
}
