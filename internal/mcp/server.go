package mcp

import (
	"context"
	"log/slog"

	"github.com/cduggn/ccexplorer/internal/ports"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server wraps the MCP server with ccExplorer-specific functionality
type Server struct {
	mcpServer  *mcp.Server
	awsService ports.AWSService
}

// NewServer creates a new MCP server instance for stdio transport
func NewServer(awsService ports.AWSService) *Server {
	slog.Info("Creating new ccExplorer MCP server")

	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "ccExplorer MCP Server",
		Version: "1.0.0",
	}, nil)

	return &Server{
		mcpServer:  mcpServer,
		awsService: awsService,
	}
}

// RegisterTools registers all available MCP tools with the server
func (s *Server) RegisterTools() error {
	slog.Info("Registering MCP tools")

	// Register the get_cost_and_usage tool with typed handler
	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_cost_and_usage",
		Description: "Query AWS Cost Explorer for cost and usage data",
	}, s.handleGetCostAndUsage)

	slog.Info("Successfully registered get_cost_and_usage tool")

	return nil
}

// Run starts the MCP server with the given transport
func (s *Server) Run(ctx context.Context) error {
	return s.mcpServer.Run(ctx, &mcp.StdioTransport{})
}
