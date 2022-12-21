package aws

import (
	"github.com/cduggn/cloudcost/internal/pkg/service/aws"
	"github.com/spf13/cobra"
	"time"
)

func CostForecast(cmd *cobra.Command, args []string) error {
	req := NewGetCostForecastRequestType()
	res, _ := aws.GetCostForecast(req)
	aws.PrintGetCostForecastReport(res)
	return nil
}

func NewGetCostForecastRequestType() aws.GetCostForecastRequest {

	services, err := aws.GetDimensionValues(aws.GetDimensionValuesRequest{
		Dimension: "SERVICE",
		Time: aws.Time{
			Start: DefaultStartDate(DayOfCurrentMonth, SubtractDays),
			End:   Format(time.Now()),
		},
	})
	if err != nil {
		panic(err)
	}

	return aws.GetCostForecastRequest{
		Granularity:             "MONTHLY",
		Metric:                  "UNBLENDED_COST",
		PredictionIntervalLevel: 95,
		Time: aws.Time{
			Start: Format(time.Now()),
			End:   "2023-04-30",
		},
		Filter: aws.Filter{
			Dimensions: []aws.Dimension{
				{
					Key:   "SERVICE",
					Value: services,
				},
				//{
				//	Key:   "REGION",
				//	Value: []string{"eu-west-1", "us-east-1", "us-west-1"},
				//},
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
