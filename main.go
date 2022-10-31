package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"log"
)

func main() {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := costexplorer.NewFromConfig(cfg)

	result, err := client.GetCostAndUsage(context.TODO(), &costexplorer.GetCostAndUsageInput{
		Granularity: types.GranularityMonthly,
		Metrics: []string{
			"UnblendedCost",
			"BlendedCost",
			"UsageQuantity",
		},
		TimePeriod: &types.DateInterval{
			Start: aws.String("2022-10-01"),
			End:   aws.String("2022-10-31"),
		},
		GroupBy: []types.GroupDefinition{
			{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String("SERVICE"),
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
	b := ToBillable(result)

	bill := b.billingSetails("Amazon Route 53")
	b.print(bill)
}
