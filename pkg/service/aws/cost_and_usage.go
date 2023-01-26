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

//"BLENDED_COST", "AMORTIZED_COST", "NET_UNBLENDED_COST",
//"NET_AMORTIZED_COST", "USAGE_QUANTITY", "NORMALIZED_USAGE_AMOUNT","USAGE_QUANTITY"

type GetCostAndUsageAPI interface {
	GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput,
		optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error)
}

func (*APIClient) GetCostAndUsage(ctx context.Context, api GetCostAndUsageAPI, req CostAndUsageRequestType) (
	*costexplorer.GetCostAndUsageOutput,
	error) {

	result, err := api.GetCostAndUsage(context.TODO(),
		&costexplorer.GetCostAndUsageInput{
			Granularity: types.Granularity(req.Granularity), //todo: add option to pass HOURLY granularity as well
			Metrics:     metrics,
			TimePeriod: &types.DateInterval{
				Start: aws.String(req.Time.Start),
				End:   aws.String(req.Time.End),
			},
			GroupBy: CostAndUsageGroupByGenerator(req),
			Filter:  CostAndUsageFilterGenerator(req),
		})

	if err != nil {
		return nil, APIError{
			msg: err.Error(),
		}
	}
	return result, nil
}
