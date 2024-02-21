package awsservice

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/ccexplorer/internal/types"
)

type Service struct {
	*costexplorer.Client
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
