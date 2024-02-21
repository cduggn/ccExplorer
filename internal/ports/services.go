package ports

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/ccexplorer/internal/types"
)

type AWSService interface {
	GetCostAndUsage(ctx context.Context,
		req types.CostAndUsageRequestType) (
		*costexplorer.GetCostAndUsageOutput,
		error)
	GetCostForecast(ctx context.Context,
		req types.GetCostForecastRequest) (
		*costexplorer.
			GetCostForecastOutput, error)
}
