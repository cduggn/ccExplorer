package aws

import (
	"context"
	"github.com/cduggn/cloudcost/internal/pkg/service/aws"
	"github.com/spf13/cobra"
	"time"
)

func CostForecast(cmd *cobra.Command, args []string) error {

	apiClient := aws.NewAPIClient()
	req := NewGetCostForecastRequestType(GetDimensionValues(apiClient))
	res, _ := apiClient.GetCostForecast(context.TODO(), apiClient.Client, req)
	aws.PrintGetCostForecastReport(res)
	return nil
}

func GetDimensionValues(c *aws.APIClient) []string {
	services, err := c.GetDimensionValues(context.TODO(), c.Client, aws.
		GetDimensionValuesRequest{
		Dimension: "SERVICE",
		Time: aws.Time{
			Start: DefaultStartDate(DayOfCurrentMonth, SubtractDays),
			End:   Format(time.Now()),
		},
	})
	if err != nil {
		panic(err)
	}
	return services
}

func NewGetCostForecastRequestType(dimensions []string) aws.
	GetCostForecastRequest {
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
					Value: dimensions,
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
