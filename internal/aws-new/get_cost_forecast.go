package aws_new

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
)

func (srv *Service) GetCostForecast(ctx context.Context,
	req model.GetCostForecastRequest) (
	*costexplorer.
		GetCostForecastOutput, error) {

	result, err := srv.Client.GetCostForecast(context.TODO(),
		&costexplorer.GetCostForecastInput{
			Granularity: types.Granularity(req.Granularity),
			Metric:      types.Metric(req.Metric),
			TimePeriod: &types.DateInterval{
				Start: aws.String(req.Time.Start),
				End:   aws.String(req.Time.End),
			},
			PredictionIntervalLevel: aws.Int32(req.PredictionIntervalLevel),
			Filter:                  CostForecastFilterGenerator(req),
		})

	if err != nil {
		return nil, model.APIError{
			Msg: err.Error(),
		}
	}

	return result, nil
}

func CostForecastFilterGenerator(req model.GetCostForecastRequest) *types.
	Expression {
	var filterExpression types.Expression
	var expList []types.Expression
	var exp types.Expression

	if req.Filter.Dimensions == nil && req.Filter.Tags == nil {
		return nil
	}

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

	filterExpression = expList[0]

	return &filterExpression
}

func ExtractForecastFilters(d map[string]string) model.Filter {

	if len(d) == 0 {
		return model.Filter{}
	}

	dimensions := CreateForecastDimensionFilter(d)

	return model.Filter{
		Dimensions: dimensions,
	}
}

func CreateForecastDimensionFilter(m map[string]string) []model.Dimension {

	if len(m) == 0 {
		return nil
	}
	var dimensions []model.Dimension
	for k, v := range m {
		dimensions = append(dimensions, model.Dimension{
			Key:   k,
			Value: []string{v},
		})
	}
	return dimensions
}
