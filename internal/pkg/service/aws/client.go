package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/cloudcost/internal/pkg/logger"
	"github.com/cduggn/cloudcost/internal/pkg/storage"
)

var (
	awsAPIClient      *APIClient
	connectionManager DatabaseManager
)

func init() {

	connectionManager = DatabaseManager{}
	err := connectionManager.newDBClient()
	if err != nil {
		logger.Error(err.Error())
	}

	awsAPIClient = &APIClient{}
	err = awsAPIClient.newAWSClient()
	if err != nil {
		logger.Error(err.Error())
	}
}

func NewAPIClient() *APIClient {
	return awsAPIClient
}

func (c *DatabaseManager) newDBClient() error {
	c.dbClient = &storage.CostDataStorage{}
	err := c.dbClient.New("./cloudcost.db")
	if err != nil {
		return DBError{
			msg: "unable to create database client, " + err.Error(),
		}
	}
	logger.Info("database connection established")
	return nil
}

func (c *APIClient) newAWSClient() error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return APIError{
			msg: "unable to load SDK config, " + err.Error(),
		}
	}
	c.Client = costexplorer.NewFromConfig(cfg)
	logger.Info("aws cost explorer client created")
	return nil
}
