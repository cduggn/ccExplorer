package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
)

type Service struct {
	*costexplorer.Client
}

func New() (*Service, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, model.APIError{
			Msg: "unable to load SDK config, " + err.Error(),
		}
	}
	return &Service{
		Client: costexplorer.NewFromConfig(cfg),
	}, nil
}

func (srv *Service) GetCostAndUsage(ctx context.Context,
	req model.CostAndUsageRequestType) (
	*costexplorer.GetCostAndUsageOutput,
	error) {

	result, err := srv.Client.GetCostAndUsage(context.TODO(),
		&costexplorer.GetCostAndUsageInput{
			Granularity: types.Granularity(req.Granularity), //todo: add option to pass HOURLY granularity as well
			Metrics:     req.Metrics,
			TimePeriod: &types.DateInterval{
				Start: aws.String(req.Time.Start),
				End:   aws.String(req.Time.End),
			},
			GroupBy: CostAndUsageGroupByGenerator(req),
			Filter:  CostAndUsageFilterGenerator(req),
		})

	if err != nil {
		return nil, model.APIError{
			Msg: err.Error(),
		}
	}
	return result, nil
}

func (srv *Service) GetCostForecast(ctx context.Context,
	req model.GetCostForecastRequest) (
	*costexplorer.
		GetCostForecastOutput, error) {

	result, err := srv.Client.GetCostForecast(context.TODO(),
		&costexplorer.GetCostForecastInput{
			Granularity: types.Granularity(req.Granularity),
			Metric:      types.Metric(req.Metric),
			TimePeriod: &types.DateInterval{
				Start: aws.String(req.Time.Start),
				End:   aws.String(req.Time.End),
			},
			PredictionIntervalLevel: aws.Int32(req.PredictionIntervalLevel),
			Filter:                  CostForecastFilterGenerator(req),
		})

	if err != nil {
		return nil, model.APIError{
			Msg: err.Error(),
		}
	}

	return result, nil
}
