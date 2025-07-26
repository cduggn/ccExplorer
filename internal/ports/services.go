package ports

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/cduggn/ccexplorer/internal/types"
	billing "cloud.google.com/go/billing/apiv1/billingpb"
	cloudasset "cloud.google.com/go/asset/apiv1/assetpb"
)

// AWSService interface for AWS Cost Explorer API integration
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

// GCPService interface for GCP Cloud Billing API integration
// Follows the same patterns as AWSService but adapted for GCP's billing ecosystem
type GCPService interface {
	// GetBillingData retrieves billing data for GCP services
	// Equivalent to AWS GetCostAndUsage but for GCP billing accounts and projects
	GetBillingData(ctx context.Context, req types.GCPBillingRequest) (*types.GCPBillingResponse, error)
	
	// GetProjectBillingInfo retrieves billing information for specific projects
	// Useful for multi-project enterprise scenarios
	GetProjectBillingInfo(ctx context.Context, projectID string) (*billing.ProjectBillingInfo, error)
	
	// GetBillingServices lists available GCP services for billing queries
	// Provides discovery capabilities for service-level filtering
	GetBillingServices(ctx context.Context) (*billing.ListServicesResponse, error)
	
	// GetBillingSkus retrieves SKUs for specific services
	// Enables detailed cost analysis at SKU level
	GetBillingSkus(ctx context.Context, serviceID string) (*billing.ListSkusResponse, error)
	
	// GetAssetInventory retrieves asset inventory for cost attribution
	// Integrates with Cloud Asset Inventory API for resource-level analysis
	GetAssetInventory(ctx context.Context, req types.GCPAssetRequest) (*cloudasset.SearchAllResourcesResponse, error)
	
	// GetCostForecast provides cost forecasting capabilities for GCP
	// Note: GCP doesn't have direct forecast API, this will use historical data analysis
	GetCostForecast(ctx context.Context, req types.GCPForecastRequest) (*types.GCPForecastResponse, error)
}
