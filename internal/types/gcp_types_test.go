package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGCPPrice_ToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		price    GCPPrice
		expected float64
	}{
		{
			name: "whole units only",
			price: GCPPrice{
				CurrencyCode: "USD",
				Units:        100,
				Nanos:        0,
			},
			expected: 100.0,
		},
		{
			name: "with fractional units",
			price: GCPPrice{
				CurrencyCode: "USD",
				Units:        50,
				Nanos:        500000000, // 0.5
			},
			expected: 50.5,
		},
		{
			name: "fractional only",
			price: GCPPrice{
				CurrencyCode: "USD",
				Units:        0,
				Nanos:        250000000, // 0.25
			},
			expected: 0.25,
		},
		{
			name: "small fractional",
			price: GCPPrice{
				CurrencyCode: "USD",
				Units:        0,
				Nanos:        1, // 0.000000001
			},
			expected: 0.000000001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.price.ToFloat64()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGCPBillingRequest_Validate(t *testing.T) {
	tests := []struct {
		name        string
		req         GCPBillingRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid request with project ID",
			req: GCPBillingRequest{
				ProjectID:   "test-project-123",
				Granularity: "MONTHLY",
				Time: GCPTimeRange{
					Start: "2024-01-01",
					End:   "2024-01-31",
				},
			},
			expectError: false,
		},
		{
			name: "valid request with project IDs",
			req: GCPBillingRequest{
				ProjectIDs:  []string{"test-project-123", "another-project-456"},
				Granularity: "DAILY",
				Time: GCPTimeRange{
					Start: "2024-01-01",
					End:   "2024-01-31",
				},
			},
			expectError: false,
		},
		{
			name: "invalid - no projects",
			req: GCPBillingRequest{
				Granularity: "MONTHLY",
				Time: GCPTimeRange{
					Start: "2024-01-01",
					End:   "2024-01-31",
				},
			},
			expectError: true,
			errorMsg:    "at least one project ID must be specified",
		},
		{
			name: "invalid - no start date",
			req: GCPBillingRequest{
				ProjectID:   "test-project-123",
				Granularity: "MONTHLY",
				Time: GCPTimeRange{
					End: "2024-01-31",
				},
			},
			expectError: true,
			errorMsg:    "start and end dates are required",
		},
		{
			name: "invalid - no end date",
			req: GCPBillingRequest{
				ProjectID:   "test-project-123",
				Granularity: "MONTHLY",
				Time: GCPTimeRange{
					Start: "2024-01-01",
				},
			},
			expectError: true,
			errorMsg:    "start and end dates are required",
		},
		{
			name: "invalid - bad granularity",
			req: GCPBillingRequest{
				ProjectID:   "test-project-123",
				Granularity: "INVALID",
				Time: GCPTimeRange{
					Start: "2024-01-01",
					End:   "2024-01-31",
				},
			},
			expectError: true,
			errorMsg:    "invalid granularity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGCPBillingRequest_Equals(t *testing.T) {
	baseReq := GCPBillingRequest{
		ProjectID:      "test-project",
		BillingAccount: "123456-789012-345678",
		Granularity:    "MONTHLY",
		Currency:       "USD",
		Time: GCPTimeRange{
			Start: "2024-01-01",
			End:   "2024-01-31",
		},
		Services: []string{"compute", "storage"},
		Regions:  []string{"us-central1", "europe-west1"},
		Labels: map[string]string{
			"env":  "production",
			"team": "backend",
		},
	}

	tests := []struct {
		name     string
		req1     GCPBillingRequest
		req2     GCPBillingRequest
		expected bool
	}{
		{
			name:     "identical requests",
			req1:     baseReq,
			req2:     baseReq,
			expected: true,
		},
		{
			name: "different project ID",
			req1: baseReq,
			req2: func() GCPBillingRequest {
				r := baseReq
				r.ProjectID = "different-project"
				return r
			}(),
			expected: false,
		},
		{
			name: "different billing account",
			req1: baseReq,
			req2: func() GCPBillingRequest {
				r := baseReq
				r.BillingAccount = "different-account"
				return r
			}(),
			expected: false,
		},
		{
			name: "different granularity",
			req1: baseReq,
			req2: func() GCPBillingRequest {
				r := baseReq
				r.Granularity = "DAILY"
				return r
			}(),
			expected: false,
		},
		{
			name: "different start time",
			req1: baseReq,
			req2: func() GCPBillingRequest {
				r := baseReq
				r.Time.Start = "2024-02-01"
				return r
			}(),
			expected: false,
		},
		{
			name: "different services length",
			req1: baseReq,
			req2: func() GCPBillingRequest {
				r := baseReq
				r.Services = []string{"compute"}
				return r
			}(),
			expected: false,
		},
		{
			name: "different services content",
			req1: baseReq,
			req2: func() GCPBillingRequest {
				r := baseReq
				r.Services = []string{"compute", "network"}
				return r
			}(),
			expected: false,
		},
		{
			name: "different labels length",
			req1: baseReq,
			req2: func() GCPBillingRequest {
				r := baseReq
				r.Labels = map[string]string{"env": "production"}
				return r
			}(),
			expected: false,
		},
		{
			name: "different labels content",
			req1: baseReq,
			req2: func() GCPBillingRequest {
				r := baseReq
				r.Labels = map[string]string{
					"env":  "staging",
					"team": "backend",
				}
				return r
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.req1.Equals(tt.req2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGCPCommandLineInput_Validation(t *testing.T) {
	tests := []struct {
		name  string
		input GCPCommandLineInput
		valid bool
	}{
		{
			name: "valid input",
			input: GCPCommandLineInput{
				ProjectID:    "test-project",
				Start:        "2024-01-01",
				End:          "2024-01-31",
				Granularity:  "MONTHLY",
				Currency:     "USD",
				PrintFormat:  "stdout",
			},
			valid: true,
		},
		{
			name: "multiple projects",
			input: GCPCommandLineInput{
				ProjectIDs:  []string{"project-1", "project-2"},
				Start:       "2024-01-01",
				End:         "2024-01-31",
				Granularity: "DAILY",
				Currency:    "EUR",
				PrintFormat: "csv",
			},
			valid: true,
		},
		{
			name: "with filters",
			input: GCPCommandLineInput{
				ProjectID:   "test-project",
				Services:    []string{"compute", "storage"},
				Regions:     []string{"us-central1"},
				SKUs:        []string{"CP-COMPUTEENGINE-VMIMAGE-N1-STANDARD-1"},
				Start:       "2024-01-01",
				End:         "2024-01-31",
				Granularity: "HOURLY",
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Since GCPCommandLineInput doesn't have a Validate method,
			// we test basic field presence and format expectations
			if tt.valid {
				assert.True(t, tt.input.ProjectID != "" || len(tt.input.ProjectIDs) > 0)
				if tt.input.Granularity != "" {
					validGranularities := []string{"DAILY", "MONTHLY", "HOURLY"}
					assert.Contains(t, validGranularities, tt.input.Granularity)
				}
			}
		})
	}
}

func TestGCPResponseMetadata(t *testing.T) {
	metadata := GCPResponseMetadata{
		RequestID:       "test-request-123",
		ResponseTime:    time.Now(),
		QueryDuration:   500 * time.Millisecond,
		APIVersion:      "v1",
		RecordCount:     150,
		TotalRecords:    500,
		HasMoreResults:  true,
	}

	assert.Equal(t, "test-request-123", metadata.RequestID)
	assert.False(t, metadata.ResponseTime.IsZero())
	assert.Equal(t, 500*time.Millisecond, metadata.QueryDuration)
	assert.Equal(t, "v1", metadata.APIVersion)
	assert.Equal(t, 150, metadata.RecordCount)
	assert.Equal(t, 500, metadata.TotalRecords)
	assert.True(t, metadata.HasMoreResults)
}

func TestGCPBillingItem_Structure(t *testing.T) {
	item := GCPBillingItem{
		ServiceID:          "compute.googleapis.com",
		ServiceDisplayName: "Compute Engine",
		ProjectID:          "test-project",
		BillingAccount:     "123456-789012-345678",
		Currency:           "USD",
		TotalCost:          150.75,
		Granularity:        "MONTHLY",
		TimeRange: GCPTimeRange{
			Start: "2024-01-01",
			End:   "2024-01-31",
		},
		SKUs: []GCPSKU{
			{
				SKUID:         "CP-COMPUTEENGINE-VMIMAGE-N1-STANDARD-1",
				DisplayName:   "N1 Standard Instance Core running in Americas",
				Description:   "N1 predefined instance with 1 vCPU running in Americas",
				Category:      "Compute",
				Family:        "Compute",
				Group:         "IaaS",
				Usage:         "OnDemand",
				Unit:          "hour",
				Regions:       []string{"us-central1", "us-east1"},
				EstimatedCost: 45.50,
				ActualUsage:   100.0,
			},
		},
		RegionalCosts: map[string]float64{
			"us-central1":   75.25,
			"europe-west1":  50.25,
			"asia-southeast1": 25.25,
		},
		UsageMetrics: map[string]float64{
			"instance_hours": 720.0,
			"storage_gb":     100.0,
		},
		Credits: []GCPCredit{
			{
				CreditID:    "credit-123",
				Name:        "Sustained Use Discount",
				Type:        "sustained_use",
				Amount:      15.50,
				Currency:    "USD",
				AppliedDate: "2024-01-15",
				ExpiryDate:  "2024-12-31",
			},
		},
		Discounts: []GCPDiscount{
			{
				DiscountID:  "discount-456",
				Name:        "Committed Use Discount",
				Type:        "committed_use",
				Amount:      10.25,
				Currency:    "USD",
				AppliedDate: "2024-01-01",
			},
		},
		Labels: map[string]string{
			"environment": "production",
			"team":        "backend",
		},
		Tags: map[string]string{
			"project": "main-app",
			"owner":   "team-lead",
		},
	}

	// Test structure integrity
	assert.Equal(t, "compute.googleapis.com", item.ServiceID)
	assert.Equal(t, "Compute Engine", item.ServiceDisplayName)
	assert.Equal(t, 150.75, item.TotalCost)
	assert.Len(t, item.SKUs, 1)
	assert.Len(t, item.RegionalCosts, 3)
	assert.Len(t, item.UsageMetrics, 2)
	assert.Len(t, item.Credits, 1)
	assert.Len(t, item.Discounts, 1)
	assert.Len(t, item.Labels, 2)
	assert.Len(t, item.Tags, 2)

	// Test SKU structure
	sku := item.SKUs[0]
	assert.Equal(t, "CP-COMPUTEENGINE-VMIMAGE-N1-STANDARD-1", sku.SKUID)
	assert.Equal(t, "hour", sku.Unit)
	assert.Equal(t, 45.50, sku.EstimatedCost)
	assert.Equal(t, 100.0, sku.ActualUsage)
	assert.Len(t, sku.Regions, 2)

	// Test credits and discounts
	assert.Equal(t, "sustained_use", item.Credits[0].Type)
	assert.Equal(t, 15.50, item.Credits[0].Amount)
	assert.Equal(t, "committed_use", item.Discounts[0].Type)
	assert.Equal(t, 10.25, item.Discounts[0].Amount)
}

func TestGCPForecastResponse_Structure(t *testing.T) {
	response := GCPForecastResponse{
		ProjectID:         "test-project",
		BillingAccount:    "123456-789012-345678",
		Currency:          "USD",
		Granularity:       "MONTHLY",
		ForecastPeriod: GCPTimeRange{
			Start: "2024-02-01",
			End:   "2024-02-29",
		},
		HistoricalPeriod: GCPTimeRange{
			Start: "2024-01-01",
			End:   "2024-01-31",
		},
		TotalForecast:   1250.75,
		ConfidenceLevel: 0.95,
		ModelType:       "linear_regression",
		ForecastBreakdown: []GCPForecastDataPoint{
			{
				Date:           "2024-02-01",
				ForecastCost:   42.50,
				LowerBound:     35.25,
				UpperBound:     49.75,
				HistoricalCost: 40.00,
			},
			{
				Date:           "2024-02-02",
				ForecastCost:   43.25,
				LowerBound:     36.00,
				UpperBound:     50.50,
				HistoricalCost: 41.00,
			},
		},
		ServiceForecasts: map[string]float64{
			"compute": 800.50,
			"storage": 250.25,
			"network": 200.00,
		},
		RegionalForecasts: map[string]float64{
			"us-central1":     625.375,
			"europe-west1":    375.225,
			"asia-southeast1": 250.15,
		},
		VarianceAnalysis: GCPVarianceAnalysis{
			TrendDirection:    "INCREASING",
			MonthOverMonth:    5.5,
			YearOverYear:      12.3,
			SeasonalityFactor: 0.85,
			VolatilityScore:   0.25,
			ConfidenceScore:   0.92,
			TopCostDrivers:    []string{"instance-hours", "storage-requests"},
			RiskFactors:       []string{"seasonal-spike", "new-features"},
		},
		Metadata: GCPResponseMetadata{
			RequestID:     "forecast-req-789",
			ResponseTime:  time.Now(),
			QueryDuration: 2 * time.Second,
			APIVersion:    "v1",
		},
	}

	// Test main forecast structure
	assert.Equal(t, "test-project", response.ProjectID)
	assert.Equal(t, 1250.75, response.TotalForecast)
	assert.Equal(t, 0.95, response.ConfidenceLevel)
	assert.Equal(t, "linear_regression", response.ModelType)

	// Test forecast breakdown
	assert.Len(t, response.ForecastBreakdown, 2)
	assert.Equal(t, "2024-02-01", response.ForecastBreakdown[0].Date)
	assert.Equal(t, 42.50, response.ForecastBreakdown[0].ForecastCost)
	assert.Equal(t, 35.25, response.ForecastBreakdown[0].LowerBound)
	assert.Equal(t, 49.75, response.ForecastBreakdown[0].UpperBound)

	// Test service and regional forecasts
	assert.Len(t, response.ServiceForecasts, 3)
	assert.Equal(t, 800.50, response.ServiceForecasts["compute"])
	
	assert.Len(t, response.RegionalForecasts, 3)
	assert.Equal(t, 625.375, response.RegionalForecasts["us-central1"])

	// Test variance analysis
	assert.Equal(t, "INCREASING", response.VarianceAnalysis.TrendDirection)
	assert.Equal(t, 5.5, response.VarianceAnalysis.MonthOverMonth)
	assert.Equal(t, 12.3, response.VarianceAnalysis.YearOverYear)
	assert.Len(t, response.VarianceAnalysis.TopCostDrivers, 2)
	assert.Len(t, response.VarianceAnalysis.RiskFactors, 2)
	assert.Contains(t, response.VarianceAnalysis.TopCostDrivers, "instance-hours")
	assert.Contains(t, response.VarianceAnalysis.RiskFactors, "seasonal-spike")
}

// Benchmark tests for performance validation
func BenchmarkGCPPrice_ToFloat64(b *testing.B) {
	price := GCPPrice{
		CurrencyCode: "USD",
		Units:        100,
		Nanos:        500000000,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		price.ToFloat64()
	}
}

func BenchmarkGCPBillingRequest_Validate(b *testing.B) {
	req := GCPBillingRequest{
		ProjectID:   "test-project-123",
		Granularity: "MONTHLY",
		Time: GCPTimeRange{
			Start: "2024-01-01",
			End:   "2024-01-31",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.Validate()
	}
}

func BenchmarkGCPBillingRequest_Equals(b *testing.B) {
	req1 := GCPBillingRequest{
		ProjectID:      "test-project",
		BillingAccount: "123456-789012-345678",
		Services:       []string{"compute", "storage", "network"},
		Regions:        []string{"us-central1", "europe-west1"},
		Labels: map[string]string{
			"env":  "production",
			"team": "backend",
		},
	}
	req2 := req1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req1.Equals(req2)
	}
}