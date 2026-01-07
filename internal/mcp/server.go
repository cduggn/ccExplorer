package mcp

import (
	"log/slog"

	"github.com/cduggn/ccexplorer/internal/ports"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Server wraps the MCP server with ccExplorer-specific functionality
type Server struct {
	mcpServer  *server.MCPServer
	awsService ports.AWSService
}

// NewServer creates a new MCP server instance for stdio transport
func NewServer(awsService ports.AWSService) *Server {
	slog.Info("Creating new ccExplorer MCP server")
	
	mcpServer := server.NewMCPServer(
		"ccExplorer MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	return &Server{
		mcpServer:  mcpServer,
		awsService: awsService,
	}
}

// RegisterTools registers all available MCP tools with the server
func (s *Server) RegisterTools() error {
	slog.Info("Registering MCP tools")
	
	// Register the get_cost_and_usage tool
	getCostTool := mcp.NewTool("get_cost_and_usage",
		mcp.WithDescription("Query AWS Cost Explorer for cost and usage data"),
		mcp.WithString("start_date", mcp.Required()),
		mcp.WithString("end_date", mcp.Required()),
		mcp.WithString("granularity", mcp.Enum("DAILY", "MONTHLY", "HOURLY")),
		mcp.WithString("metrics"),
		mcp.WithString("group_by"),
		mcp.WithString("filter_by_service"),
		mcp.WithBoolean("exclude_discounts"),
	)

	// Add tool to the MCP server with our handler
	s.mcpServer.AddTool(getCostTool, s.handleGetCostAndUsage)
	slog.Info("Successfully registered get_cost_and_usage tool")
	
	return nil
}

// MCPServer returns the underlying MCP server for stdio transport
func (s *Server) MCPServer() *server.MCPServer {
	return s.mcpServer
}