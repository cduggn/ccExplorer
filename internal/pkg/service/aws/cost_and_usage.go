package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

var (
	metrics = []string{"UNBLENDED_COST"}
)

type GetCostAndUsageAPI interface {
	GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error)
}

func (*APIClient) GetCostAndUsage(ctx context.Context, api GetCostAndUsageAPI,
	req CostAndUsageRequestType) (
	*CostAndUsageReport,
	error) {

	result, err := api.GetCostAndUsage(context.TODO(),
		&costexplorer.GetCostAndUsageInput{
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
	c.CurateCostAndUsageReport(result)
	return c, nil
}
