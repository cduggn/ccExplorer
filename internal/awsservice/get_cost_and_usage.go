package awsservice

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	types2 "github.com/cduggn/ccexplorer/internal/types"
	"github.com/cduggn/ccexplorer/internal/utils"
)

var (
	groupByDimension = func(dimensions []string) []types.GroupDefinition {
		return utils.Transform(dimensions, func(d string) types.GroupDefinition {
			return types.GroupDefinition{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String(d),
			}
		})
	}
	groupByTag = func(tag []string) []types.GroupDefinition {
		return utils.Transform(tag, func(t string) types.GroupDefinition {
			return types.GroupDefinition{
				Type: types.GroupDefinitionTypeTag,
				Key:  aws.String(t),
			}
		})
	}
	groupByTagAndDimension = func(tag []string, dimensions []string) []types.GroupDefinition {
		dimensionGroups := groupByDimension(dimensions)
		tagGroups := groupByTag(tag)
		
		// Combine the two slices using generic utilities
		combined := make([]types.GroupDefinition, 0, len(dimensionGroups)+len(tagGroups))
		combined = append(combined, dimensionGroups...)
		combined = append(combined, tagGroups...)
		return combined
	}
	filterCredits = func() *types.Expression {
		return &types.Expression{
			Not: &types.Expression{
				Dimensions: &types.DimensionValues{
					Key:    "RECORD_TYPE", // Note: RECORD_TYPE is the equivalent of CHARGE_TYPE - https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/manage-cost-categories.html#cost-categories-terms
					Values: []string{"Refund", "Credit", "DiscountedUsage", "Discount", "BundledDiscount ", "SavingsPlanCoveredUsage", "SavingsPlanNegation"},
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

func (srv *Service) GetCostAndUsage(ctx context.Context,
	req types2.CostAndUsageRequestType) (
	*costexplorer.GetCostAndUsageOutput,
	error) {

	result, err := srv.Client.GetCostAndUsage(context.TODO(),
		&costexplorer.GetCostAndUsageInput{
			Granularity: types.Granularity(req.Granularity), //todo: add option to pass HOURLY granularity as well
			Metrics:     req.Metrics,
			TimePeriod: &types.DateInterval{
				Start: aws.String(req.Time.Start),
				End:   aws.String(req.Time.End),
			},
			GroupBy: CostAndUsageGroupByGenerator(req),
			Filter:  CostAndUsageFilterGenerator(req),
		})

	if err != nil {
		return nil, types2.APIError{
			Msg: err.Error(),
		}
	}
	return result, nil
}

// ToSlice converts dimension values to string slice using generic transformation
func ToSlice(d costexplorer.GetDimensionValuesOutput) []string {
	return utils.Transform(d.DimensionValues, func(dimension types.DimensionValuesWithAttributes) string {
		return *dimension.Value
	})
}

func CostAndUsageFilterGenerator(req types2.CostAndUsageRequestType) *types.
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

func CostAndUsageGroupByGenerator(req types2.CostAndUsageRequestType) []types.GroupDefinition {
	if len(req.GroupByTag) == 1 && len(req.GroupBy) == 1 {
		return groupByTagAndDimension(req.GroupByTag, req.GroupBy)
	} else if len(req.GroupByTag) >= 1 {
		return groupByTag(req.GroupByTag)
	} else {
		return groupByDimension(req.GroupBy)
	}

}
