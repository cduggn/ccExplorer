package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"log"
)

func GetAWSCostAndUsageWithDiscounts(req CostAndUsageRequest) *CostAndUsageReport {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := costexplorer.NewFromConfig(cfg)

	result, err := client.GetCostAndUsage(context.TODO(), &costexplorer.GetCostAndUsageInput{
		Granularity: types.Granularity(req.Granularity), //todo: add option to pass HOURLY granularity as well
		Metrics:     req.Rates,
		TimePeriod: &types.DateInterval{
			Start: aws.String(req.Time.Start),
			End:   aws.String(req.Time.End),
		},
		GroupBy: groupBy(req),
		Filter:  filter(req),
	})

	if err != nil {
		log.Fatal(err)
	}
	c := &CostAndUsageReport{
		Services: make(map[int]Service),
	}
	c.Granularity = req.Granularity
	c.CurateReport(result)
	return c
}
