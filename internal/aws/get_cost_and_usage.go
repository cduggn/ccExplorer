package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	types2 "github.com/cduggn/ccexplorer/internal/types"
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

func ToSlice(d costexplorer.GetDimensionValuesOutput) []string {
	var servicesSlice []string
	for _, service := range d.DimensionValues {
		servicesSlice = append(servicesSlice, *service.Value)
	}
	return servicesSlice
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
