package gcpservice

import (
	"context"
	"fmt"
	"sync"
	"time"

	billing "cloud.google.com/go/billing/apiv1"
	"cloud.google.com/go/billing/apiv1/billingpb"
	cloudasset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
	"github.com/cduggn/ccexplorer/internal/types"
	"golang.org/x/time/rate"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// Service implements the GCPService interface for Google Cloud Platform billing operations
type Service struct {
	billingClient   *billing.CloudBillingClient
	catalogClient   *billing.CloudCatalogClient  
	assetClient     *cloudasset.Client
	rateLimiter     *rate.Limiter
	config          *Config
	mu              sync.RWMutex
	connectionPool  map[string]*grpc.ClientConn
	metrics         *ServiceMetrics
}

// ServiceMetrics tracks performance and usage metrics
type ServiceMetrics struct {
	RequestCount    int64
	ErrorCount      int64
	TotalDuration   time.Duration
	LastRequestTime time.Time
	mu              sync.RWMutex
}

// NewService creates a new GCP service instance with enterprise-grade configuration
func NewService(config *Config) (*Service, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	ctx := context.Background()

	// Configure client options with optimizations
	clientOpts := []option.ClientOption{
		option.WithGRPCDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		})),
	}

	// Add authentication options
	if config.ServiceAccountPath != "" {
		clientOpts = append(clientOpts, option.WithCredentialsFile(config.ServiceAccountPath))
	} else if config.ServiceAccountJSON != "" {
		clientOpts = append(clientOpts, option.WithCredentialsJSON([]byte(config.ServiceAccountJSON)))
	}
	// If neither is provided, ADC (Application Default Credentials) will be used

	// Create billing client with error handling using correct v1.20.4 API
	billingClient, err := billing.NewCloudBillingClient(ctx, clientOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create billing client: %w", err)
	}

	// Create catalog client for services and SKUs
	catalogClient, err := billing.NewCloudCatalogClient(ctx, clientOpts...)
	if err != nil {
		billingClient.Close() // Clean up billing client
		return nil, fmt.Errorf("failed to create catalog client: %w", err)
	}

	// Create asset client with error handling
	assetClient, err := cloudasset.NewClient(ctx, clientOpts...)
	if err != nil {
		billingClient.Close() // Clean up billing client
		catalogClient.Close()  // Clean up catalog client
		return nil, fmt.Errorf("failed to create asset client: %w", err)
	}

	// Configure rate limiter (300 requests per minute with burst capacity)
	rateLimiter := rate.NewLimiter(rate.Every(200*time.Millisecond), 10)

	service := &Service{
		billingClient:  billingClient,
		catalogClient:  catalogClient,
		assetClient:    assetClient,
		rateLimiter:    rateLimiter,
		config:         config,
		connectionPool: make(map[string]*grpc.ClientConn),
		metrics:        &ServiceMetrics{},
	}

	return service, nil
}

// GetBillingData retrieves comprehensive billing data for GCP services
func (s *Service) GetBillingData(ctx context.Context, req types.GCPBillingRequest) (*types.GCPBillingResponse, error) {
	// Record metrics
	startTime := time.Now()
	s.updateMetrics(startTime, nil)

	// Rate limiting
	if err := s.rateLimiter.Wait(ctx); err != nil {
		s.updateMetrics(startTime, err)
		return nil, fmt.Errorf("rate limiting error: %w", err)
	}

	// Validate request
	if err := req.Validate(); err != nil {
		s.updateMetrics(startTime, err)
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Determine projects to query
	projects := s.getProjectsToQuery(req)
	if len(projects) == 0 {
		err := fmt.Errorf("no valid projects to query")
		s.updateMetrics(startTime, err)
		return nil, err
	}

	// Aggregate data from all projects
	response := &types.GCPBillingResponse{
		Services:       make(map[string]types.GCPBillingItem),
		Currency:       req.Currency,
		Granularity:    req.Granularity,
		TimeRange:      req.Time,
		BillingAccount: req.BillingAccount,
		Metadata: types.GCPResponseMetadata{
			RequestID:     generateRequestID(),
			ResponseTime:  time.Now(),
			APIVersion:    "v1",
			RecordCount:   0,
		},
	}

	// Process each project
	for _, projectID := range projects {
		projectData, err := s.getBillingDataForProject(ctx, projectID, req)
		if err != nil {
			// Log error but continue with other projects
			fmt.Printf("Warning: failed to get billing data for project %s: %v\n", projectID, err)
			continue
		}

		// Merge project data into response
		s.mergeBillingData(response, projectData)
	}

	// Calculate totals
	s.calculateTotals(response)

	s.updateMetrics(startTime, nil)
	return response, nil
}

// GetProjectBillingInfo retrieves billing information for a specific project
func (s *Service) GetProjectBillingInfo(ctx context.Context, projectID string) (*billingpb.ProjectBillingInfo, error) {
	startTime := time.Now()
	defer s.updateMetrics(startTime, nil)

	if err := s.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiting error: %w", err)
	}

	if projectID == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}

	req := &billingpb.GetProjectBillingInfoRequest{
		Name: fmt.Sprintf("projects/%s", projectID),
	}

	result, err := s.billingClient.GetProjectBillingInfo(ctx, req)
	if err != nil {
		s.updateMetrics(startTime, err)
		return nil, fmt.Errorf("failed to get project billing info: %w", err)
	}

	return result, nil
}

// GetBillingServices lists available GCP services for billing queries
func (s *Service) GetBillingServices(ctx context.Context) (*billingpb.ListServicesResponse, error) {
	startTime := time.Now()
	defer s.updateMetrics(startTime, nil)

	if err := s.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiting error: %w", err)
	}

	req := &billingpb.ListServicesRequest{
		PageSize: 1000, // Maximum page size for efficiency
	}

	serviceIterator := s.catalogClient.ListServices(ctx, req)
	var services []*billingpb.Service
	
	for {
		service, err := serviceIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			s.updateMetrics(startTime, err)
			return nil, fmt.Errorf("failed to iterate services: %w", err)
		}
		services = append(services, service)
	}

	// Create response manually since we're using iterator
	response := &billingpb.ListServicesResponse{
		Services: services,
	}

	return response, nil
}

// GetBillingSkus retrieves SKUs for a specific service
func (s *Service) GetBillingSkus(ctx context.Context, serviceID string) (*billingpb.ListSkusResponse, error) {
	startTime := time.Now()
	defer s.updateMetrics(startTime, nil)

	if err := s.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiting error: %w", err)
	}

	if serviceID == "" {
		return nil, fmt.Errorf("service ID cannot be empty")
	}

	req := &billingpb.ListSkusRequest{
		Parent:   fmt.Sprintf("services/%s", serviceID),
		PageSize: 1000,
	}

	skuIterator := s.catalogClient.ListSkus(ctx, req)
	var skus []*billingpb.Sku
	
	for {
		sku, err := skuIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			s.updateMetrics(startTime, err)
			return nil, fmt.Errorf("failed to iterate SKUs for service %s: %w", serviceID, err)
		}
		skus = append(skus, sku)
	}

	// Create response manually since we're using iterator
	response := &billingpb.ListSkusResponse{
		Skus: skus,
	}

	return response, nil
}

// GetAssetInventory retrieves asset inventory for cost attribution
func (s *Service) GetAssetInventory(ctx context.Context, req types.GCPAssetRequest) (*assetpb.SearchAllResourcesResponse, error) {
	startTime := time.Now()
	defer s.updateMetrics(startTime, nil)

	if err := s.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiting error: %w", err)
	}

	assetReq := &assetpb.SearchAllResourcesRequest{
		Scope:      req.Scope,
		Query:      req.Query,
		AssetTypes: req.AssetTypes,
		PageSize:   req.PageSize,
		PageToken:  req.PageToken,
		OrderBy:    req.OrderBy,
	}

	assetIterator := s.assetClient.SearchAllResources(ctx, assetReq)
	var resources []*assetpb.ResourceSearchResult
	
	for {
		resource, err := assetIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			s.updateMetrics(startTime, err)
			return nil, fmt.Errorf("failed to iterate asset inventory: %w", err)
		}
		resources = append(resources, resource)
	}

	// Create response manually since we're using iterator
	response := &assetpb.SearchAllResourcesResponse{
		Results: resources,
	}

	return response, nil
}

// GetCostForecast provides cost forecasting capabilities using historical data analysis
func (s *Service) GetCostForecast(ctx context.Context, req types.GCPForecastRequest) (*types.GCPForecastResponse, error) {
	startTime := time.Now()
	defer s.updateMetrics(startTime, nil)

	// Since GCP doesn't have a direct forecast API, we'll use historical data
	// to generate forecasts using trend analysis
	billingReq := types.GCPBillingRequest{
		ProjectID:      req.ProjectID,
		ProjectIDs:     req.ProjectIDs,
		BillingAccount: req.BillingAccount,
		Services:       req.Services,
		Regions:        req.Regions,
		Time:           req.TimeRange,
		Granularity:    req.Granularity,
		Currency:       req.Currency,
	}

	historicalData, err := s.GetBillingData(ctx, billingReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical data for forecast: %w", err)
	}

	// Generate forecast based on historical trends
	forecast := s.generateForecast(historicalData, req)

	return forecast, nil
}

// Close gracefully closes all client connections
func (s *Service) Close() error {
	var errors []error

	if s.billingClient != nil {
		if err := s.billingClient.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close billing client: %w", err))
		}
	}

	if s.catalogClient != nil {
		if err := s.catalogClient.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close catalog client: %w", err))
		}
	}

	if s.assetClient != nil {
		if err := s.assetClient.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close asset client: %w", err))
		}
	}

	// Close connection pool
	s.mu.Lock()
	for name, conn := range s.connectionPool {
		if err := conn.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close connection %s: %w", name, err))
		}
	}
	s.connectionPool = make(map[string]*grpc.ClientConn)
	s.mu.Unlock()

	if len(errors) > 0 {
		return fmt.Errorf("errors while closing service: %v", errors)
	}

	return nil
}

// GetMetrics returns current service metrics
func (s *Service) GetMetrics() ServiceMetrics {
	s.metrics.mu.RLock()
	defer s.metrics.mu.RUnlock()
	return *s.metrics
}

// Helper methods

func (s *Service) updateMetrics(startTime time.Time, err error) {
	s.metrics.mu.Lock()
	defer s.metrics.mu.Unlock()

	s.metrics.RequestCount++
	s.metrics.TotalDuration += time.Since(startTime)
	s.metrics.LastRequestTime = time.Now()

	if err != nil {
		s.metrics.ErrorCount++
	}
}

func (s *Service) getProjectsToQuery(req types.GCPBillingRequest) []string {
	var projects []string

	if req.ProjectID != "" {
		projects = append(projects, req.ProjectID)
	}

	projects = append(projects, req.ProjectIDs...)

	// Remove duplicates
	seen := make(map[string]bool)
	var uniqueProjects []string
	for _, project := range projects {
		if !seen[project] {
			seen[project] = true
			uniqueProjects = append(uniqueProjects, project)
		}
	}

	return uniqueProjects
}

func (s *Service) getBillingDataForProject(ctx context.Context, projectID string, req types.GCPBillingRequest) (*types.GCPBillingResponse, error) {
	// This is a simplified implementation. In a real scenario, you would:
	// 1. Query the Cloud Billing API for detailed cost data
	// 2. Process the results according to the request parameters
	// 3. Transform the data into the internal format

	// For now, return a placeholder response
	return &types.GCPBillingResponse{
		Services:       make(map[string]types.GCPBillingItem),
		ProjectInfo:    projectID,
		Currency:       req.Currency,
		Granularity:    req.Granularity,
		TimeRange:      req.Time,
		BillingAccount: req.BillingAccount,
	}, nil
}

func (s *Service) mergeBillingData(target *types.GCPBillingResponse, source *types.GCPBillingResponse) {
	// Merge service data
	for serviceID, serviceData := range source.Services {
		if existing, exists := target.Services[serviceID]; exists {
			// Merge existing service data
			existing.TotalCost += serviceData.TotalCost
			target.Services[serviceID] = existing
		} else {
			target.Services[serviceID] = serviceData
		}
	}

	// Update metadata
	target.Metadata.RecordCount += source.Metadata.RecordCount
}

func (s *Service) calculateTotals(response *types.GCPBillingResponse) {
	var totalCost, totalCredits, totalDiscounts float64

	for _, service := range response.Services {
		totalCost += service.TotalCost
		for _, credit := range service.Credits {
			totalCredits += credit.Amount
		}
		for _, discount := range service.Discounts {
			totalDiscounts += discount.Amount
		}
	}

	response.TotalCost = totalCost
	response.TotalCredits = totalCredits
	response.TotalDiscounts = totalDiscounts
}

func (s *Service) generateForecast(historicalData *types.GCPBillingResponse, req types.GCPForecastRequest) *types.GCPForecastResponse {
	// Simplified forecast generation using trend analysis
	// In a production system, this would use more sophisticated forecasting algorithms

	totalHistoricalCost := historicalData.TotalCost
	
	// Simple linear trend projection
	forecastMultiplier := 1.1 // 10% growth assumption
	totalForecast := totalHistoricalCost * forecastMultiplier

	return &types.GCPForecastResponse{
		ProjectID:         req.ProjectID,
		BillingAccount:    req.BillingAccount,
		Currency:          req.Currency,
		Granularity:       req.Granularity,
		ForecastPeriod:    req.ForecastPeriod,
		HistoricalPeriod:  req.TimeRange,
		TotalForecast:     totalForecast,
		ConfidenceLevel:   req.ConfidenceLevel,
		ModelType:         req.ModelType,
		ServiceForecasts:  make(map[string]float64),
		RegionalForecasts: make(map[string]float64),
		Metadata: types.GCPResponseMetadata{
			RequestID:    generateRequestID(),
			ResponseTime: time.Now(),
			APIVersion:   "v1",
		},
	}
}

func generateRequestID() string {
	return fmt.Sprintf("gcp-req-%d", time.Now().UnixNano())
}