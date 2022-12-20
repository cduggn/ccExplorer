package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

func GetCostForecast(req GetCostForecastRequestType) (*costexplorer.GetCostForecastOutput, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, APIError{
			msg: "unable to load SDK config, " + err.Error(),
		}
	}
	client := costexplorer.NewFromConfig(cfg)

	services, _ := client.GetDimensionValues(context.TODO(),
		&costexplorer.GetDimensionValuesInput{
			Dimension: "SERVICE",
			TimePeriod: &types.DateInterval{
				Start: aws.String("2022-12-01"),
				End:   aws.String("2022-12-30"),
			},
		})

	// copy add services.DimensionValues to a slice of strings
	var servicesSlice []string
	for _, service := range services.DimensionValues {
		servicesSlice = append(servicesSlice, *service.Value)
	}

	result, err := client.GetCostForecast(context.TODO(), &costexplorer.GetCostForecastInput{
		Granularity: types.Granularity("MONTHLY"),
		Metric:      "UNBLENDED_COST",
		TimePeriod: &types.DateInterval{
			Start: aws.String("2022-12-20"),
			End:   aws.String("2023-04-30"),
		},
		Filter: &types.Expression{
			Dimensions: &types.DimensionValues{
				Key:    "SERVICE",
				Values: servicesSlice,

				//And: []types.Expression{
				//	{
				//		Dimensions: &types.DimensionValues{
				//			Key:    "SERVICE",
				//			Values: servicesSlice,
				//		},
				//	},
				//{
				//	Dimensions: &types.DimensionValues{
				//		Key:    "REGION",
				//		Values: []string{"eu-west-1"},
				//	},
				//},
			},
		},
		//
		//Filter: &types.Expression{
		//	Tags: &types.TagValues{
		//		Key:    aws.String("ApplicationName"),
		//
		//	},
		//},

	})

	// convert result to GetCostForecastResult struct

	if err != nil {
		return nil, APIError{
			msg: "Error while fetching cost and usage data from AWS",
		}
	}

	//c := &CostAndUsageReport{
	//	Services: make(map[int]Service),
	//}
	//c.Granularity = req.Granularity

	return result, nil
}
