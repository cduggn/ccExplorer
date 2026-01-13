package mcp

import (
	"fmt"
	"strings"

	"github.com/cduggn/ccexplorer/internal/types"
)

// translateMCPToInternalRequest converts MCP parameters to internal CostAndUsageRequestType
func (s *Server) translateMCPToInternalRequest(params types.MCPToolParameters) (types.CostAndUsageRequestType, error) {
	var request types.CostAndUsageRequestType

	// Basic time and granularity settings
	request.Time = types.Time{
		Start: params.StartDate,
		End:   params.EndDate,
	}
	request.Granularity = params.Granularity
	request.Metrics = params.Metrics
	request.ExcludeDiscounts = params.ExcludeDiscounts

	// Parse group_by parameters
	var groupByDimension []string
	var groupByTag []string

	for _, groupBy := range params.GroupBy {
		if strings.HasPrefix(groupBy, "TAG:") {
			// Extract tag name from "TAG:TagName" format
			tagName := strings.TrimPrefix(groupBy, "TAG:")
			groupByTag = append(groupByTag, tagName)
		} else {
			// Regular dimension
			groupByDimension = append(groupByDimension, groupBy)
		}
	}

	request.GroupBy = groupByDimension
	request.GroupByTag = groupByTag

	// Handle filtering
	if params.FilterByService != "" {
		request.IsFilterByDimensionEnabled = true
		request.DimensionFilter = map[string]string{
			"SERVICE": params.FilterByService,
		}
	}

	// Handle additional dimension filters
	if len(params.FilterByDimension) > 0 {
		request.IsFilterByDimensionEnabled = true
		if request.DimensionFilter == nil {
			request.DimensionFilter = make(map[string]string)
		}
		for key, value := range params.FilterByDimension {
			request.DimensionFilter[key] = value
		}
	}

	// Handle tag filters
	if len(params.FilterByTag) > 0 {
		request.IsFilterByTagEnabled = true
		// For now, we'll handle the first tag filter
		// This could be enhanced to support multiple tag filters
		for key, value := range params.FilterByTag {
			request.GroupByTag = append(request.GroupByTag, key)
			request.TagFilterValue = value
			break // Only handle first tag for now
		}
	}

	// Set default print format for internal processing
	request.PrintFormat = "json"

	// Validate the request
	if err := s.validateInternalRequest(request); err != nil {
		return request, fmt.Errorf("invalid request: %w", err)
	}

	return request, nil
}

// validateInternalRequest validates the internal request structure
func (s *Server) validateInternalRequest(request types.CostAndUsageRequestType) error {
	// Validate required fields
	if request.Time.Start == "" {
		return fmt.Errorf("start date is required")
	}
	if request.Time.End == "" {
		return fmt.Errorf("end date is required")
	}

	// Validate granularity
	validGranularities := []string{"DAILY", "MONTHLY", "HOURLY"}
	granularityValid := false
	for _, valid := range validGranularities {
		if request.Granularity == valid {
			granularityValid = true
			break
		}
	}
	if !granularityValid {
		return fmt.Errorf("invalid granularity: %s, must be one of %v", request.Granularity, validGranularities)
	}

	// Validate metrics
	if len(request.Metrics) == 0 {
		return fmt.Errorf("at least one metric is required")
	}
	validMetrics := []string{
		"AmortizedCost", "BlendedCost", "NetAmortizedCost",
		"NetUnblendedCost", "NormalizedUsageAmount", "UnblendedCost", "UsageQuantity",
	}
	for _, metric := range request.Metrics {
		metricValid := false
		for _, valid := range validMetrics {
			if metric == valid {
				metricValid = true
				break
			}
		}
		if !metricValid {
			return fmt.Errorf("invalid metric: %s, must be one of %v", metric, validMetrics)
		}
	}

	return nil
}