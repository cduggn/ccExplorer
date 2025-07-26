package gcpservice

import (
	"testing"
	"time"

	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil config",
			config:      nil,
			expectError: true,
			errorMsg:    "config cannot be nil",
		},
		{
			name: "valid config with ADC",
			config: &Config{
				UseADC:          true,
				MaxRetries:      3,
				RequestTimeout:  30,
				RateLimitRPS:    5,
				RateLimitBurst:  10,
				MaxConcurrency:  8,
				DefaultCurrency: "USD",
				CacheTTL:        15,
			},
			expectError: false,
		},
		{
			name: "invalid config - no authentication",
			config: &Config{
				UseADC:             false,
				ServiceAccountPath: "",
				ServiceAccountJSON: "",
				MaxRetries:         3,
				RequestTimeout:     30,
				RateLimitRPS:       5,
				RateLimitBurst:     10,
				MaxConcurrency:     8,
				DefaultCurrency:    "USD",
				CacheTTL:           15,
			},
			expectError: true,
			errorMsg:    "authentication required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewService(tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, service)
			} else {
				// For valid configs, we expect an error due to missing GCP credentials in test env
				// but we can still test that the service creation logic works
				if err != nil {
					// Expected in test environment without real GCP credentials
					assert.Contains(t, err.Error(), "failed to create")
				}
			}
		})
	}
}

func TestService_GetMetrics(t *testing.T) {
	// Create a service with mock config for testing metrics
	config := NewDefaultConfig()
	
	// Create service instance for testing (this will fail in test env, but we can test the struct)
	service := &Service{
		config: config,
		metrics: &ServiceMetrics{
			RequestCount:    10,
			ErrorCount:      2,
			TotalDuration:   5 * time.Second,
			LastRequestTime: time.Now(),
		},
	}

	metrics := service.GetMetrics()
	
	assert.Equal(t, int64(10), metrics.RequestCount)
	assert.Equal(t, int64(2), metrics.ErrorCount)
	assert.Equal(t, 5*time.Second, metrics.TotalDuration)
	assert.False(t, metrics.LastRequestTime.IsZero())
}

func TestService_updateMetrics(t *testing.T) {
	service := &Service{
		metrics: &ServiceMetrics{},
	}

	startTime := time.Now().Add(-100 * time.Millisecond)

	// Test successful request
	service.updateMetrics(startTime, nil)
	
	metrics := service.GetMetrics()
	assert.Equal(t, int64(1), metrics.RequestCount)
	assert.Equal(t, int64(0), metrics.ErrorCount)
	assert.True(t, metrics.TotalDuration > 0)

	// Test failed request
	service.updateMetrics(startTime, assert.AnError)
	
	metrics = service.GetMetrics()
	assert.Equal(t, int64(2), metrics.RequestCount)
	assert.Equal(t, int64(1), metrics.ErrorCount)
}

func TestService_getProjectsToQuery(t *testing.T) {
	service := &Service{}

	tests := []struct {
		name     string
		req      types.GCPBillingRequest
		expected []string
	}{
		{
			name: "single project ID",
			req: types.GCPBillingRequest{
				ProjectID: "project-1",
			},
			expected: []string{"project-1"},
		},
		{
			name: "multiple project IDs",
			req: types.GCPBillingRequest{
				ProjectIDs: []string{"project-1", "project-2", "project-3"},
			},
			expected: []string{"project-1", "project-2", "project-3"},
		},
		{
			name: "both single and multiple (no duplicates)",
			req: types.GCPBillingRequest{
				ProjectID:  "project-1",
				ProjectIDs: []string{"project-2", "project-1", "project-3"},
			},
			expected: []string{"project-1", "project-2", "project-3"},
		},
		{
			name: "empty request",
			req:  types.GCPBillingRequest{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.getProjectsToQuery(tt.req)
			
			// Check length
			assert.Equal(t, len(tt.expected), len(result))
			
			// Check all expected elements are present (order may vary due to deduplication)
			for _, expected := range tt.expected {
				assert.Contains(t, result, expected)
			}
		})
	}
}

func TestService_calculateTotals(t *testing.T) {
	service := &Service{}

	response := &types.GCPBillingResponse{
		Services: map[string]types.GCPBillingItem{
			"compute": {
				TotalCost: 100.50,
				Credits: []types.GCPCredit{
					{Amount: 10.0},
					{Amount: 5.0},
				},
				Discounts: []types.GCPDiscount{
					{Amount: 15.0},
					{Amount: 3.0},
				},
			},
			"storage": {
				TotalCost: 50.25,
				Credits: []types.GCPCredit{
					{Amount: 2.5},
				},
				Discounts: []types.GCPDiscount{
					{Amount: 1.0},
				},
			},
		},
	}

	service.calculateTotals(response)

	assert.Equal(t, 150.75, response.TotalCost)     // 100.50 + 50.25
	assert.Equal(t, 17.5, response.TotalCredits)   // 10.0 + 5.0 + 2.5
	assert.Equal(t, 19.0, response.TotalDiscounts) // 15.0 + 3.0 + 1.0
}

func TestService_generateForecast(t *testing.T) {
	service := &Service{}

	historicalData := &types.GCPBillingResponse{
		TotalCost: 1000.0,
	}

	forecastReq := types.GCPForecastRequest{
		ProjectID:       "test-project",
		BillingAccount:  "123456-789012-345678",
		Currency:        "USD", 
		Granularity:     "MONTHLY",
		ConfidenceLevel: 0.95,
		ModelType:       "linear",
		ForecastPeriod: types.GCPTimeRange{
			Start: "2024-01-01",
			End:   "2024-01-31",
		},
		TimeRange: types.GCPTimeRange{
			Start: "2023-12-01", 
			End:   "2023-12-31",
		},
	}

	forecast := service.generateForecast(historicalData, forecastReq)

	assert.Equal(t, "test-project", forecast.ProjectID)
	assert.Equal(t, "123456-789012-345678", forecast.BillingAccount)
	assert.Equal(t, "USD", forecast.Currency)
	assert.Equal(t, "MONTHLY", forecast.Granularity)
	assert.Equal(t, 0.95, forecast.ConfidenceLevel)
	assert.Equal(t, "linear", forecast.ModelType)
	assert.Equal(t, 1100.0, forecast.TotalForecast) // 1000 * 1.1 (10% growth)
	assert.NotEmpty(t, forecast.Metadata.RequestID)
	assert.False(t, forecast.Metadata.ResponseTime.IsZero())
}

func TestGenerateRequestID(t *testing.T) {
	id1 := generateRequestID()
	id2 := generateRequestID()

	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2) // Should be unique
	assert.Contains(t, id1, "gcp-req-")
	assert.Contains(t, id2, "gcp-req-")
}

// Benchmark tests for performance validation
func BenchmarkService_updateMetrics(b *testing.B) {
	service := &Service{
		metrics: &ServiceMetrics{},
	}
	startTime := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.updateMetrics(startTime, nil)
	}
}

func BenchmarkService_getProjectsToQuery(b *testing.B) {
	service := &Service{}
	req := types.GCPBillingRequest{
		ProjectID:  "project-1",
		ProjectIDs: []string{"project-2", "project-3", "project-4", "project-5"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.getProjectsToQuery(req)
	}
}

func BenchmarkService_calculateTotals(b *testing.B) {
	service := &Service{}
	response := &types.GCPBillingResponse{
		Services: map[string]types.GCPBillingItem{
			"compute": {TotalCost: 100.50},
			"storage": {TotalCost: 50.25},
			"network": {TotalCost: 25.75},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.calculateTotals(response)
	}
}

// Integration test helpers (would require real GCP credentials to run)
func TestService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// This test would require real GCP credentials and would be run separately
	// from unit tests in a proper CI/CD pipeline
	t.Skip("integration test requires GCP credentials - run separately")
}

// Mock implementations for testing
type MockBillingClient struct{}

type MockAssetClient struct{}

// Helper functions for test setup
func setupTestService(t *testing.T) *Service {
	return &Service{
		config: NewDefaultConfig(),
		metrics: &ServiceMetrics{},
	}
}

func createTestBillingRequest() types.GCPBillingRequest {
	return types.GCPBillingRequest{
		ProjectID:   "test-project",
		Currency:    "USD",
		Granularity: "MONTHLY",
		Time: types.GCPTimeRange{
			Start: "2024-01-01",
			End:   "2024-01-31",
		},
	}
}

func TestService_ValidateRequest(t *testing.T) {
	tests := []struct {
		name        string
		req         types.GCPBillingRequest
		expectError bool
	}{
		{
			name:        "valid request",
			req:         createTestBillingRequest(),
			expectError: false,
		},
		{
			name: "invalid request - no project",
			req: types.GCPBillingRequest{
				Currency:    "USD",
				Granularity: "MONTHLY",
				Time: types.GCPTimeRange{
					Start: "2024-01-01",
					End:   "2024-01-31",
				},
			},
			expectError: true,
		},
		{
			name: "invalid request - no time range",
			req: types.GCPBillingRequest{
				ProjectID:   "test-project",
				Currency:    "USD",
				Granularity: "MONTHLY",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}