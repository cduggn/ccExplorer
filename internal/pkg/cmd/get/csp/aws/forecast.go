package aws

import (
	"github.com/cduggn/cloudcost/internal/pkg/service/aws"
	"github.com/spf13/cobra"
)

func CostForecast(cmd *cobra.Command, args []string) error {
	req := NewGetCostForecastRequestType()
	res, _ := aws.GetCostForecast(req)
	aws.PrintGetCostForecastReport(res)
	return nil
}

func NewGetCostForecastRequestType() aws.GetCostForecastRequest {
	return aws.GetCostForecastRequest{
		Granularity:             "MONTHLY",
		Metric:                  "UNBLENDED_COST",
		PredictionIntervalLevel: 95,
		Time: aws.Time{
			Start: "2022-12-20",
			End:   "2023-04-30",
		},
		Filter: aws.Filter{
			Dimensions: []aws.Dimension{
				{
					Key:   "REGION",
					Value: []string{"eu-west-1", "us-east-1", "us-west-1"},
				},
			},
			//Tags: []aws.Tag{
			//	{
			//		Key:   "Name",
			//		Value: "test",
			//	},
			//},
		},
	}
}
