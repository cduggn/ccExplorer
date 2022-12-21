package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/cloudcost/internal/pkg/logger"
	"github.com/cduggn/cloudcost/internal/pkg/storage"
)

type AWSClient interface {
	GetCostAndUsage(req CostAndUsageRequestType) (*CostAndUsageReport, error)
	GetDimensionValues(ctx context.Context, api GetDimensionValuesAPI,
		d GetDimensionValuesRequest) ([]string, error)
	GetCostForecast(req GetCostForecastRequest) (*costexplorer.GetCostForecastOutput, error)
}

type APIClient struct {
	*costexplorer.Client
}

var (
	apiClient         *APIClient
	connectionManager DatabaseManager
)

type DatabaseManager struct {
	dbClient *storage.CostDataStorage
}

func init() {
	//db
	connectionManager = DatabaseManager{}
	connectionManager.newDBClient()
	// aws client
	apiClient = &APIClient{}
	err := apiClient.newAWSClient()
	if err != nil {
		logger.Error(err.Error())
	}
}

func NewAPIClient() *APIClient {
	return apiClient
}

func (c *DatabaseManager) newDBClient() {
	c.dbClient = &storage.CostDataStorage{}
	err := c.dbClient.New("./cloudcost.db")
	if err != nil {
		logger.Error(err.Error())
	}
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
