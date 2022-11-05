package billing

import "C"
import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"log"
)

var (
	metrics = []string{
		"UnblendedCost",
		"BlendedCost",
		"UsageQuantity",
	}
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
)

func GetAWSCostAndUsage(req CostAndUsageRequest) *CostAndUsageReport {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := costexplorer.NewFromConfig(cfg)

	result, err := client.GetCostAndUsage(context.TODO(), &costexplorer.GetCostAndUsageInput{
		Granularity: types.Granularity(req.Granularity),
		Metrics:     metrics,
		TimePeriod: &types.DateInterval{
			Start: aws.String("2022-11-01"),
			End:   aws.String("2022-11-06"),
		},
		GroupBy: groupBy(req),
	})

	if err != nil {
		log.Fatal(err)
	}
	c := &CostAndUsageReport{
		Services: make(map[int]Service),
	}
	c.CurateReport(result)
	return c
}

func groupBy(req CostAndUsageRequest) []types.GroupDefinition {
	if req.Tag != "" {
		return groupByTag(req.Tag)
	} else {
		return groupByDimension(req.GroupBy)
	}
}

func (c *CostAndUsageReport) CurateReport(output *costexplorer.GetCostAndUsageOutput) {

	for _, v := range output.ResultsByTime {
		c.Start = *v.TimePeriod.Start
		c.End = *v.TimePeriod.End
		for index, g := range v.Groups {
			keys := make([]string, 0)
			service := Service{}

			keys = append(keys, g.Keys...)

			for key, m := range g.Metrics {
				metrics := Metrics{
					Name:   key,
					Amount: *m.Amount,
					Unit:   *m.Unit,
				}
				service.Metrics = append(service.Metrics, metrics)
			}
			service.Keys = keys
			c.Services[index] = service
		}
	}
}

func (c *CostAndUsageReport) Print() {
	count := 0
	for _, m := range c.Services {
		count++
		fmt.Printf(" \n| %d | GroupedBy", count)
		for _, k := range m.Keys {
			fmt.Print(" | ", k, " ")
		}
		fmt.Print("\nCost And Usage Report\n")
		for _, v := range m.Metrics {
			fmt.Printf("Name: %s: \n", v.Name)
			fmt.Printf("Amount: %s Unit: %s\n", v.Amount, v.Unit)
		}
	}
}
