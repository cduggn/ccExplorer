package awsservice

import (
	"context"
	"fmt"
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
					Values: []string{"Refund", "Credit", "DiscountedUsage", "Discount", "BundledDiscount", "SavingsPlanCoveredUsage", "SavingsPlanNegation"},
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

// maxPages bounds pagination. Cost Explorer bills per request, so a runaway
// token loop is expensive as well as slow.
const maxPages = 100

func (srv *Service) GetCostAndUsage(ctx context.Context,
	req types2.CostAndUsageRequestType) (
	*costexplorer.GetCostAndUsageOutput,
	error) {

	input := &costexplorer.GetCostAndUsageInput{
		Granularity: types.Granularity(req.Granularity),
		Metrics:     req.Metrics,
		TimePeriod: &types.DateInterval{
			Start: aws.String(req.Time.Start),
			End:   aws.String(req.Time.End),
		},
		GroupBy: CostAndUsageGroupByGenerator(req),
		Filter:  CostAndUsageFilterGenerator(req),
	}

	var result *costexplorer.GetCostAndUsageOutput

	for page := 1; ; page++ {
		out, err := srv.Client.GetCostAndUsage(ctx, input)
		if err != nil {
			return nil, types2.APIError{
				Msg: err.Error(),
			}
		}

		if result == nil {
			result = out
		} else {
			result.ResultsByTime = append(result.ResultsByTime,
				out.ResultsByTime...)
			result.DimensionValueAttributes = append(
				result.DimensionValueAttributes,
				out.DimensionValueAttributes...)
		}

		if out.NextPageToken == nil || *out.NextPageToken == "" {
			break
		}

		if page >= maxPages {
			return nil, types2.APIError{
				Msg: fmt.Sprintf(
					"result set exceeded %d pages; narrow the date range, "+
						"granularity or grouping", maxPages),
			}
		}

		input.NextPageToken = out.NextPageToken
	}

	// The token describes the last page fetched, not the aggregate we return.
	result.NextPageToken = nil

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
