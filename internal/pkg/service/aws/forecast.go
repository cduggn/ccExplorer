package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

type GetCostForecastAPI interface {
	GetCostForecast(ctx context.Context, params *costexplorer.GetCostForecastInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostForecastOutput, error)
}

func (*APIClient) GetCostForecast(ctx context.Context,
	api GetCostForecastAPI, req GetCostForecastRequest) (*costexplorer.GetCostForecastOutput, error) {

	result, err := api.GetCostForecast(context.TODO(),
		&costexplorer.GetCostForecastInput{
			Granularity: types.Granularity(req.Granularity),
			Metric:      types.Metric(req.Metric),
			TimePeriod: &types.DateInterval{
				Start: aws.String(req.Time.Start),
				End:   aws.String(req.Time.End),
			},
			PredictionIntervalLevel: aws.Int32(req.PredictionIntervalLevel),
			Filter:                  GenerateFilterExpression(req),
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

func GenerateFilterExpression(req GetCostForecastRequest) *types.Expression {
	var filterExpression types.Expression
	var expList []types.Expression
	var exp types.Expression

	var isMultiFilter bool
	if len(req.Filter.Dimensions) > 1 {
		isMultiFilter = true
	}

	for _, dimension := range req.Filter.Dimensions {
		temp := &types.DimensionValues{
			Key:    types.Dimension(dimension.Key),
			Values: dimension.Value,
		}

		if len(req.Filter.Dimensions) == 1 {
			expList = append(expList, types.Expression{
				Dimensions: temp,
			})
		} else if len(req.Filter.Dimensions) > 1 {
			exp.And = append(exp.And, types.Expression{
				Dimensions: temp,
			})
		}
	}

	if isMultiFilter {
		expList = append(expList, exp)
	}

	//
	//for _, tag := range req.Filter.Tags {
	//	temp := &types.TagValues{
	//		Key:    aws.String(tag.Key),
	//		Values: tag.Value,
	//	}
	//	exp = append(exp, types.Expression{
	//		Tags: temp,
	//	})
	//}

	filterExpression = expList[0]

	return &filterExpression
}
