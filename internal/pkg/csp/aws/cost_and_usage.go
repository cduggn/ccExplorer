package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

var (
	metrics = []string{"UNBLENDED_COST"}
)

type APIError struct {
	msg string
}

func (e APIError) Error() string {
	return e.msg
}

func GetCostAndUsage(req CostAndUsageRequestType) (*CostAndUsageReport, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, APIError{
			msg: "unable to load SDK config, " + err.Error(),
		}
	}
	client := costexplorer.NewFromConfig(cfg)

	result, err := client.GetCostAndUsage(context.TODO(), &costexplorer.GetCostAndUsageInput{
		Granularity: types.Granularity(req.Granularity), //todo: add option to pass HOURLY granularity as well
		Metrics:     metrics,
		TimePeriod: &types.DateInterval{
			Start: aws.String(req.Time.Start),
			End:   aws.String(req.Time.End),
		},
		GroupBy: groupBy(req),
		Filter:  filter(req),
	})

	if err != nil {
		return nil, APIError{
			msg: "Error while fetching cost and usage data from AWS",
		}
	}
	c := &CostAndUsageReport{
		Services: make(map[int]Service),
	}
	c.Granularity = req.Granularity
	c.CurateReport(result)
	return c, nil
}