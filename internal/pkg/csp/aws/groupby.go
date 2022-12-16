package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
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
)

func groupBy(req CostAndUsageRequestType) []types.GroupDefinition {
	if req.Tag != "" && len(req.GroupBy) == 1 {
		return groupByTagAndDimension(req.Tag, req.GroupBy)
	} else if req.Tag != "" {
		return groupByTag(req.Tag)
	} else {
		return groupByDimension(req.GroupBy)
	}
}

func groupByArn(req CostAndUsageRequestType) []types.GroupDefinition {
	return groupByDimension(req.GroupBy)
}
