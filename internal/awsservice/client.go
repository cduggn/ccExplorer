package awsservice

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/ccexplorer/internal/types"
)

// CostExplorerAPI is the subset of the Cost Explorer client used by this
// package. It exists so the pagination and filter logic can be exercised
// without a live AWS account.
type CostExplorerAPI interface {
	GetCostAndUsage(ctx context.Context, params *costexplorer.GetCostAndUsageInput,
		optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error)
	GetCostForecast(ctx context.Context, params *costexplorer.GetCostForecastInput,
		optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostForecastOutput, error)
	GetAnomalies(ctx context.Context, params *costexplorer.GetAnomaliesInput,
		optFns ...func(*costexplorer.Options)) (*costexplorer.GetAnomaliesOutput, error)
}

type Service struct {
	Client CostExplorerAPI
}

func New() (*Service, error) {

	var err error
	var cfg aws.Config

	if profile := Profile(); profile == "not-provided" {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithSharedConfigProfile(profile))
	}

	if err != nil {
		return nil, types.APIError{
			Msg: "unable to load SDK config, " + err.Error(),
		}
	}
	return &Service{
		Client: costexplorer.NewFromConfig(cfg),
	}, nil
}
