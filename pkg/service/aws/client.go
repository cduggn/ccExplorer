package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/ccexplorer/pkg/logger"
)

var (
	awsAPIClient *APIClient
)

func init() {
	awsAPIClient = &APIClient{}
	err := awsAPIClient.newAWSClient()
	if err != nil {
		logger.Error(err.Error())
	}
}

func NewAPIClient() *APIClient {
	return awsAPIClient
}

func (c *APIClient) newAWSClient() error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return APIError{
			msg: "unable to load SDK config, " + err.Error(),
		}
	}
	c.Client = costexplorer.NewFromConfig(cfg)

	return nil
}
