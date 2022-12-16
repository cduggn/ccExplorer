package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

var (
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
)

func filter(req CostAndUsageRequestType) *types.Expression {
	expression := &types.Expression{}

	if req.IncludeDiscounts && req.IsFilterEnabled {
		expression.And = []types.Expression{*filterCredits(),
			*filterByTag(req.Tag, req.TagFilterValue)}
	} else if req.IncludeDiscounts {
		expression = filterCredits()
	} else if req.IsFilterEnabled {
		expression = filterByTag(req.Tag, req.TagFilterValue)
	} else {
		return nil
	}
	return expression
}
