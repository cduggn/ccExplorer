package ports

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
)

type AWSService interface {
	GetCostAndUsage(ctx context.Context,
		req model.CostAndUsageRequestType) (
		*costexplorer.GetCostAndUsageOutput,
		error)
	GetCostForecast(ctx context.Context,
		req model.GetCostForecastRequest) (
		*costexplorer.
		GetCostForecastOutput, error)
}
