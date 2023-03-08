package domain

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/ccexplorer/pkg/domain/aws"
	"github.com/cduggn/ccexplorer/pkg/domain/model"
)

type Port interface {
	GetCostAndUsage(ctx context.Context, api aws.GetCostAndUsageAPI,
		req model.CostAndUsageRequestType) (
		*costexplorer.GetCostAndUsageOutput,
		error)
	GetCostForecast(ctx context.Context,
		api aws.GetCostForecastAPI, req model.GetCostForecastRequest) (
		*costexplorer.GetCostForecastOutput, error)
}
