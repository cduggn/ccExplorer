package billing

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"log"
	"strconv"
)

func FetchCloudCost() *Billable {
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
	return ToBillable(result)
}

func ToBillable(result *costexplorer.GetCostAndUsageOutput) *Billable {
	billable := &Billable{
		Services: make(map[string]Service),
	}
	for _, v := range result.ResultsByTime {
		billable.Start = *v.TimePeriod.Start
		billable.End = *v.TimePeriod.End
		for _, g := range v.Groups {
			service := Service{}
			for _, k := range g.Keys {
				service.Name = k
			}
			for key, m := range g.Metrics {
				metrics := Metrics{
					Name:   key,
					Amount: *m.Amount,
					Unit:   *m.Unit,
				}
				service.Metrics = append(service.Metrics, metrics)
			}
			billable.Services[service.Name] = service
		}
	}
	return billable
}

func (b Billable) Print(s Service) {
	fmt.Println("Service: ", s.Name)
	for _, m := range s.Metrics {
		fmt.Println(m.Name, m.Amount, m.Unit)
	}
}

func (b Billable) total() float64 {
	var total float64
	for _, v := range b.Services {
		for _, m := range v.Metrics {
			if m.Name == "UnblendedCost" {
				total += toFloatingPoint(m.Amount)
			}
		}
	}
	return total
}

func toFloatingPoint(amount string) float64 {
	floatNumAmount, _ := strconv.ParseFloat(amount, 64)
	return floatNumAmount
}

func (b Billable) ForService(s string) Service {
	return b.Services[s]
}
