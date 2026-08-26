package awsservice

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	types2 "github.com/cduggn/ccexplorer/internal/types"
)

func (srv *Service) GetAnomalies(ctx context.Context,
	req types2.GetAnomaliesRequestType) (
	*costexplorer.GetAnomaliesOutput,
	error) {

	input := &costexplorer.GetAnomaliesInput{
		DateInterval: &types.AnomalyDateInterval{
			StartDate: aws.String(req.Time.Start),
			EndDate:   aws.String(req.Time.End),
		},
	}
	if req.MonitorArn != "" {
		input.MonitorArn = aws.String(req.MonitorArn)
	}
	if req.Feedback != "" {
		input.Feedback = types.AnomalyFeedbackType(req.Feedback)
	}
	if req.MaxResults != 0 {
		input.MaxResults = aws.Int32(req.MaxResults)
	}

	result, err := paginate(
		func(token *string) (*costexplorer.GetAnomaliesOutput, error) {
			input.NextPageToken = token
			return srv.Client.GetAnomalies(ctx, input)
		},
		func(out *costexplorer.GetAnomaliesOutput) *string {
			return out.NextPageToken
		},
		func(acc, page *costexplorer.GetAnomaliesOutput) *costexplorer.GetAnomaliesOutput {
			acc.Anomalies = append(acc.Anomalies, page.Anomalies...)
			return acc
		},
	)
	if err != nil {
		return nil, err
	}

	// The token describes the last page fetched, not the aggregate we return.
	result.NextPageToken = nil

	return result, nil
}
