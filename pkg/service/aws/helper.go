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
	groupByTag = func(tag []string) []types.GroupDefinition {

		var groups []types.GroupDefinition
		for _, t := range tag {
			groups = append(groups, types.GroupDefinition{
				Type: types.GroupDefinitionTypeTag,
				Key:  aws.String(t),
			})
		}
		return groups
	}
	groupByTagAndDimension = func(tag []string, dimensions []string) []types.
				GroupDefinition {
		var groups []types.GroupDefinition
		for _, d := range dimensions {
			groups = append(groups, types.GroupDefinition{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String(d),
			})
		}
		for _, t := range tag {
			groups = append(groups, types.GroupDefinition{
				Type: types.GroupDefinitionTypeTag,
				Key:  aws.String(t),
			})
		}

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
	filterByTag = func(tag []string, value string) *types.Expression {
		if len(tag) == 0 || tag[0] == "" {
			return &types.Expression{}
		}
		return &types.Expression{
			Tags: &types.TagValues{
				Key:    aws.String(tag[0]),
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
		filters = append(filters, *filterByTag(req.GroupByTag, req.TagFilterValue))
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
	if len(req.GroupByTag) == 1 && len(req.GroupBy) == 1 {
		return groupByTagAndDimension(req.GroupByTag, req.GroupBy)
	} else if len(req.GroupByTag) >= 1 {
		return groupByTag(req.GroupByTag)
	} else {
		return groupByDimension(req.GroupBy)
	}

}
