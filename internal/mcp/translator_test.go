package mcp

import (
	"testing"

	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseGetCostAndUsageParams(t *testing.T) {
	server := &Server{}

	tests := []struct {
		name     string
		args     map[string]interface{}
		expected types.MCPToolParameters
		wantErr  bool
	}{
		{
			name: "valid basic parameters",
			args: map[string]interface{}{
				"start_date": "2024-01-01",
				"end_date":   "2024-01-31",
			},
			expected: types.MCPToolParameters{
				StartDate:   "2024-01-01",
				EndDate:     "2024-01-31",
				Granularity: "MONTHLY",
				Metrics:     []string{"UnblendedCost"},
			},
			wantErr: false,
		},
		{
			name: "valid parameters with all options",
			args: map[string]interface{}{
				"start_date":         "2024-01-01",
				"end_date":           "2024-01-31",
				"granularity":        "DAILY",
				"metrics":            []interface{}{"AmortizedCost", "BlendedCost"},
				"group_by":           []interface{}{"SERVICE", "TAG:Project"},
				"filter_by_service":  "Amazon Simple Storage Service",
				"exclude_discounts":  true,
			},
			expected: types.MCPToolParameters{
				StartDate:         "2024-01-01",
				EndDate:           "2024-01-31",
				Granularity:       "DAILY",
				Metrics:           []string{"AmortizedCost", "BlendedCost"},
				GroupBy:           []string{"SERVICE", "TAG:Project"},
				FilterByService:   "Amazon Simple Storage Service",
				ExcludeDiscounts:  true,
			},
			wantErr: false,
		},
		{
			name: "missing start_date",
			args: map[string]interface{}{
				"end_date": "2024-01-31",
			},
			wantErr: true,
		},
		{
			name: "missing end_date",
			args: map[string]interface{}{
				"start_date": "2024-01-01",
			},
			wantErr: true,
		},
		{
			name: "invalid start_date type",
			args: map[string]interface{}{
				"start_date": 123,
				"end_date":   "2024-01-31",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := server.parseGetCostAndUsageParams(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected.StartDate, result.StartDate)
			assert.Equal(t, tt.expected.EndDate, result.EndDate)
			assert.Equal(t, tt.expected.Granularity, result.Granularity)
			assert.Equal(t, tt.expected.Metrics, result.Metrics)
			assert.Equal(t, tt.expected.GroupBy, result.GroupBy)
			assert.Equal(t, tt.expected.FilterByService, result.FilterByService)
			assert.Equal(t, tt.expected.ExcludeDiscounts, result.ExcludeDiscounts)
		})
	}
}

func TestTranslateMCPToInternalRequest(t *testing.T) {
	server := &Server{}

	tests := []struct {
		name     string
		params   types.MCPToolParameters
		wantErr  bool
		validate func(t *testing.T, req types.CostAndUsageRequestType)
	}{
		{
			name: "basic parameters",
			params: types.MCPToolParameters{
				StartDate:   "2024-01-01",
				EndDate:     "2024-01-31",
				Granularity: "MONTHLY",
				Metrics:     []string{"UnblendedCost"},
				GroupBy:     []string{"SERVICE"},
			},
			wantErr: false,
			validate: func(t *testing.T, req types.CostAndUsageRequestType) {
				assert.Equal(t, "2024-01-01", req.Time.Start)
				assert.Equal(t, "2024-01-31", req.Time.End)
				assert.Equal(t, "MONTHLY", req.Granularity)
				assert.Equal(t, []string{"UnblendedCost"}, req.Metrics)
				assert.Equal(t, []string{"SERVICE"}, req.GroupBy)
				assert.False(t, req.IsFilterByDimensionEnabled)
				assert.False(t, req.IsFilterByTagEnabled)
			},
		},
		{
			name: "with dimension and tag grouping",
			params: types.MCPToolParameters{
				StartDate:   "2024-01-01",
				EndDate:     "2024-01-31",
				Granularity: "DAILY",
				Metrics:     []string{"AmortizedCost"},
				GroupBy:     []string{"SERVICE", "TAG:Project"},
			},
			wantErr: false,
			validate: func(t *testing.T, req types.CostAndUsageRequestType) {
				assert.Equal(t, []string{"SERVICE"}, req.GroupBy)
				assert.Equal(t, []string{"Project"}, req.GroupByTag)
			},
		},
		{
			name: "with service filter",
			params: types.MCPToolParameters{
				StartDate:       "2024-01-01",
				EndDate:         "2024-01-31",
				Granularity:     "MONTHLY",
				Metrics:         []string{"UnblendedCost"},
				FilterByService: "Amazon S3",
			},
			wantErr: false,
			validate: func(t *testing.T, req types.CostAndUsageRequestType) {
				assert.True(t, req.IsFilterByDimensionEnabled)
				assert.Equal(t, "Amazon S3", req.DimensionFilter["SERVICE"])
			},
		},
		{
			name: "invalid granularity",
			params: types.MCPToolParameters{
				StartDate:   "2024-01-01",
				EndDate:     "2024-01-31",
				Granularity: "INVALID",
				Metrics:     []string{"UnblendedCost"},
			},
			wantErr: true,
		},
		{
			name: "invalid metric",
			params: types.MCPToolParameters{
				StartDate:   "2024-01-01",
				EndDate:     "2024-01-31",
				Granularity: "MONTHLY",
				Metrics:     []string{"InvalidMetric"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := server.translateMCPToInternalRequest(tt.params)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

func TestValidateInternalRequest(t *testing.T) {
	server := &Server{}

	tests := []struct {
		name    string
		request types.CostAndUsageRequestType
		wantErr bool
	}{
		{
			name: "valid request",
			request: types.CostAndUsageRequestType{
				Time:        types.Time{Start: "2024-01-01", End: "2024-01-31"},
				Granularity: "MONTHLY",
				Metrics:     []string{"UnblendedCost"},
			},
			wantErr: false,
		},
		{
			name: "missing start date",
			request: types.CostAndUsageRequestType{
				Time:        types.Time{End: "2024-01-31"},
				Granularity: "MONTHLY",
				Metrics:     []string{"UnblendedCost"},
			},
			wantErr: true,
		},
		{
			name: "invalid granularity",
			request: types.CostAndUsageRequestType{
				Time:        types.Time{Start: "2024-01-01", End: "2024-01-31"},
				Granularity: "YEARLY",
				Metrics:     []string{"UnblendedCost"},
			},
			wantErr: true,
		},
		{
			name: "no metrics",
			request: types.CostAndUsageRequestType{
				Time:        types.Time{Start: "2024-01-01", End: "2024-01-31"},
				Granularity: "MONTHLY",
				Metrics:     []string{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := server.validateInternalRequest(tt.request)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}