package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"log"
)

func GetAWSCostAndUsageWithResources(req CostAndUsageRequest) *CostAndUsageReport {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := costexplorer.NewFromConfig(cfg)

	result, err := client.GetCostAndUsageWithResources(context.TODO(), &costexplorer.GetCostAndUsageWithResourcesInput{
		Granularity: types.Granularity(req.Granularity), //todo: add option to pass HOURLY granularity as well
		Metrics:     req.Rates,
		TimePeriod: &types.DateInterval{
			Start: aws.String(req.Time.Start),
			End:   aws.String(req.Time.End),
		},
		//GroupBy: groupBy(req),
		//Filter:  filter(req),
		//Filter: &types.Expression{
		//	Dimensions: &types.DimensionValues{
		//		Key: "SERVICE",
		//		Values: []string{
		//			"Amazon Elastic Compute Cloud - Compute",
		//			"Amazon Elastic Block Store",
		//		},
		//	},
		//},
		GroupBy: []types.GroupDefinition{
			types.GroupDefinition{
				Type: "DIMENSION",
				Key:  aws.String("SERVICE"),
			},
			//types.GroupDefinition{
			//	Type: "TAG",
			//	Key:  aws.String("ApplicationName"),
			//},
			types.GroupDefinition{
				Type: "DIMENSION",
				Key:  aws.String("RESOURCE_ID"),
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
	c := &CostAndUsageReport{
		Services: make(map[int]Service),
	}
	c.Granularity = req.Granularity
	fmt.Println(result)
	//c.CurateReport(result)
	return c
}
