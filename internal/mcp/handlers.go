package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/cduggn/ccexplorer/internal/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// handleGetCostAndUsage handles the get_cost_and_usage MCP tool call
func (s *Server) handleGetCostAndUsage(ctx context.Context, req *mcp.CallToolRequest, input types.GetCostAndUsageInput) (*mcp.CallToolResult, any, error) {
	slog.Info("Handling get_cost_and_usage request", "input", input)

	// Convert typed input to internal MCPToolParameters
	mcpParams := inputToParams(input)

	// Translate MCP parameters to internal request type
	internalRequest, err := s.translateMCPToInternalRequest(mcpParams)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to translate parameters: %w", err)
	}

	// Call AWS Cost Explorer service
	result, err := s.awsService.GetCostAndUsage(ctx, internalRequest)
	if err != nil {
		return nil, nil, fmt.Errorf("AWS service error: %w", err)
	}

	// Transform the response to a format suitable for MCP
	response := utils.ToCostAndUsageOutputType(result, internalRequest)

	// Convert to JSON for MCP response
	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	// Return MCP tool result
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(responseJSON)}},
	}, nil, nil
}

// inputToParams converts the typed SDK input to internal MCPToolParameters
func inputToParams(input types.GetCostAndUsageInput) types.MCPToolParameters {
	params := types.MCPToolParameters{
		StartDate:        input.StartDate,
		EndDate:          input.EndDate,
		Granularity:      input.Granularity,
		FilterByService:  input.FilterByService,
		ExcludeDiscounts: input.ExcludeDiscounts,
	}

	if params.Granularity == "" {
		params.Granularity = "MONTHLY"
	}

	// Parse comma-separated metrics
	if input.Metrics != "" {
		for _, m := range strings.Split(input.Metrics, ",") {
			m = strings.TrimSpace(m)
			if m != "" {
				params.Metrics = append(params.Metrics, m)
			}
		}
	}
	if len(params.Metrics) == 0 {
		params.Metrics = []string{"UnblendedCost"}
	}

	// Parse comma-separated group_by
	if input.GroupBy != "" {
		for _, g := range strings.Split(input.GroupBy, ",") {
			g = strings.TrimSpace(g)
			if g != "" {
				params.GroupBy = append(params.GroupBy, g)
			}
		}
	}

	return params
}
