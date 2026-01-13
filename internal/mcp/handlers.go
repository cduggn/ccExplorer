package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/cduggn/ccexplorer/internal/utils"
	"github.com/mark3labs/mcp-go/mcp"
)

// handleGetCostAndUsage handles the get_cost_and_usage MCP tool call
func (s *Server) handleGetCostAndUsage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slog.Info("Handling get_cost_and_usage request", "arguments", request.Params.Arguments)

	// Type assert the arguments
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid arguments type")
	}

	// Parse and validate arguments
	mcpParams, err := s.parseGetCostAndUsageParams(args)
	if err != nil {
		return nil, fmt.Errorf("invalid parameters: %w", err)
	}

	// Translate MCP parameters to internal request type
	internalRequest, err := s.translateMCPToInternalRequest(mcpParams)
	if err != nil {
		return nil, fmt.Errorf("failed to translate parameters: %w", err)
	}

	// Call AWS Cost Explorer service
	result, err := s.awsService.GetCostAndUsage(ctx, internalRequest)
	if err != nil {
		return nil, fmt.Errorf("AWS service error: %w", err)
	}

	// Transform the response to a format suitable for MCP
	response := utils.ToCostAndUsageOutputType(result, internalRequest)

	// Convert to JSON for MCP response
	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	// Return MCP tool result
	return mcp.NewToolResultText(string(responseJSON)), nil
}

// parseGetCostAndUsageParams parses the MCP tool arguments into MCPToolParameters
func (s *Server) parseGetCostAndUsageParams(args map[string]interface{}) (types.MCPToolParameters, error) {
	var params types.MCPToolParameters

	// Required parameters
	if startDate, ok := args["start_date"].(string); ok {
		params.StartDate = startDate
	} else {
		return params, fmt.Errorf("start_date is required and must be a string")
	}

	if endDate, ok := args["end_date"].(string); ok {
		params.EndDate = endDate
	} else {
		return params, fmt.Errorf("end_date is required and must be a string")
	}

	// Optional parameters with defaults
	if granularity, ok := args["granularity"].(string); ok {
		params.Granularity = granularity
	} else {
		params.Granularity = "MONTHLY"
	}

	// Parse metrics - can be array or comma-separated string
	if metricsInterface, ok := args["metrics"]; ok {
		if metricsArray, ok := metricsInterface.([]interface{}); ok {
			// Handle array format (typically from HTTP JSON-RPC)
			for _, metric := range metricsArray {
				if metricStr, ok := metric.(string); ok {
					params.Metrics = append(params.Metrics, metricStr)
				}
			}
		} else if metricsStr, ok := metricsInterface.(string); ok && metricsStr != "" {
			// Handle string format (typically from stdio)
			if strings.Contains(metricsStr, ",") {
				// Comma-separated values
				for _, metric := range strings.Split(metricsStr, ",") {
					metric = strings.TrimSpace(metric)
					if metric != "" {
						params.Metrics = append(params.Metrics, metric)
					}
				}
			} else {
				// Single metric
				params.Metrics = append(params.Metrics, metricsStr)
			}
		}
	}
	if len(params.Metrics) == 0 {
		params.Metrics = []string{"UnblendedCost"}
	}

	// Parse group_by - can be array or comma-separated string
	if groupByInterface, ok := args["group_by"]; ok {
		if groupByArray, ok := groupByInterface.([]interface{}); ok {
			// Handle array format (typically from HTTP JSON-RPC)
			for _, groupBy := range groupByArray {
				if groupByStr, ok := groupBy.(string); ok {
					params.GroupBy = append(params.GroupBy, groupByStr)
				}
			}
		} else if groupByStr, ok := groupByInterface.(string); ok && groupByStr != "" {
			// Handle string format (typically from stdio)
			if strings.Contains(groupByStr, ",") {
				// Comma-separated values
				for _, groupBy := range strings.Split(groupByStr, ",") {
					groupBy = strings.TrimSpace(groupBy)
					if groupBy != "" {
						params.GroupBy = append(params.GroupBy, groupBy)
					}
				}
			} else {
				// Single group_by
				params.GroupBy = append(params.GroupBy, groupByStr)
			}
		}
	}

	// Optional filter parameters
	if filterService, ok := args["filter_by_service"].(string); ok && filterService != "" {
		params.FilterByService = filterService
	}

	if excludeDiscounts, ok := args["exclude_discounts"].(bool); ok {
		params.ExcludeDiscounts = excludeDiscounts
	}

	return params, nil
}

