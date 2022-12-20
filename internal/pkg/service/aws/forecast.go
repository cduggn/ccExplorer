package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

func GetCostForecast(req GetCostForecastRequest) (*costexplorer.GetCostForecastOutput, error) {

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
				Start: aws.String(req.Time.Start),
				End:   aws.String(req.Time.End),
			},
		})

	// copy add services.DimensionValues to a slice of strings
	var servicesSlice []string
	for _, service := range services.DimensionValues {
		servicesSlice = append(servicesSlice, *service.Value)
	}

	fmt.Println(servicesSlice)

	result, err := client.GetCostForecast(context.TODO(), &costexplorer.GetCostForecastInput{
		Granularity: types.Granularity(req.Granularity),
		Metric:      types.Metric(req.Metric),
		TimePeriod: &types.DateInterval{
			Start: aws.String(req.Time.Start),
			End:   aws.String(req.Time.End),
		},
		PredictionIntervalLevel: aws.Int32(req.PredictionIntervalLevel),
		Filter: &types.Expression{
			And: GenerateFilterExpression(req),
		},
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

func GenerateFilterExpression(req GetCostForecastRequest) []types.Expression {

	var exp []types.Expression

	for _, dimension := range req.Filter.Dimensions {
		temp := &types.DimensionValues{
			Key:    types.Dimension(dimension.Key),
			Values: dimension.Value,
		}
		exp = append(exp, types.Expression{
			Dimensions: temp,
		})
	}

	for _, tag := range req.Filter.Tags {
		temp := &types.TagValues{
			Key:    aws.String(tag.Key),
			Values: tag.Value,
		}
		exp = append(exp, types.Expression{
			Tags: temp,
		})
	}

	return exp
}
