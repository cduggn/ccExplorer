package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

var (
	groupByDimension = func(dimensions []string) []types.GroupDefinition {
		var groups []types.GroupDefinition
		for _, d := range dimensions {
			groups = append(groups, types.GroupDefinition{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String(d),
			})
		}
		return groups
	}
	groupByTag = func(tag string) []types.GroupDefinition {
		return []types.GroupDefinition{
			{
				Type: types.GroupDefinitionTypeTag,
				Key:  aws.String(tag),
			},
		}
	}
	groupByTagAndDimension = func(tag string, dimensions []string) []types.GroupDefinition {
		var groups []types.GroupDefinition
		for _, d := range dimensions {
			groups = append(groups, types.GroupDefinition{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String(d),
			})
		}
		groups = append(groups, types.GroupDefinition{
			Type: types.GroupDefinitionTypeTag,
			Key:  aws.String(tag),
		})
		return groups
	}
	filterCredits = func() *types.Expression {
		return &types.Expression{
			Not: &types.Expression{
				Dimensions: &types.DimensionValues{
					Key:    "RECORD_TYPE",
					Values: []string{"Refund", "Credit", "DiscountedUsage"},
				},
			},
		}
	}
	filterByTag = func(tag string, value string) *types.Expression {
		return &types.Expression{
			Tags: &types.TagValues{
				Key:    aws.String(tag),
				Values: []string{value},
			},
		}
	}
	filterByDimension = func(dimension string, value string) *types.Expression {
		return &types.Expression{
			Dimensions: &types.DimensionValues{
				Key:    types.Dimension(dimension),
				Values: []string{value},
			},
		}
	}
)

func ToSlice(d costexplorer.GetDimensionValuesOutput) []string {
	var servicesSlice []string
	for _, service := range d.DimensionValues {
		servicesSlice = append(servicesSlice, *service.Value)
	}
	return servicesSlice
}

func CostAndUsageFilterGenerator(req CostAndUsageRequestType) *types.
	Expression {
	expression := &types.Expression{}
	var filters []types.Expression

	if req.ExcludeDiscounts {
		filters = append(filters, *filterCredits())
	}
	if req.IsFilterByTagEnabled {
		filters = append(filters, *filterByTag(req.Tag, req.TagFilterValue))
	}
	if req.IsFilterByDimensionEnabled {
		for key, value := range req.DimensionFilter {
			filters = append(filters, *filterByDimension(key, value))
		}
		//filters = append(filters, *filterByDimension(req.DimensionFilterName, req.DimensionFilterValue))
	}

	if len(filters) == 0 {
		return nil
	} else if len(filters) == 1 {
		expression = &filters[0]
	} else {
		expression.And = filters
	}
	return expression
}

func CostForecastFilterGenerator(req GetCostForecastRequest) *types.
	Expression {
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

	filterExpression = expList[0]

	return &filterExpression
}

func CostAndUsageGroupByGenerator(req CostAndUsageRequestType) []types.GroupDefinition {
	if req.Tag != "" && len(req.GroupBy) == 1 {
		return groupByTagAndDimension(req.Tag, req.GroupBy)
	} else if req.Tag != "" {
		return groupByTag(req.Tag)
	} else {
		return groupByDimension(req.GroupBy)
	}

	//if len(req.DimensionSubFilterName) > 0 {
	//	// extract key value from map index 0
	//	for k, v := range req.DimensionSubFilterName {
	//		filters = append(filters, *filterByDimension(k, v))
	//	}
	//}
}
