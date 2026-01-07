package cli

import (
	"fmt"
	"log/slog"

	"github.com/cduggn/ccexplorer/internal/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

// mcpCommand creates the MCP command structure
func mcpCommand() *cobra.Command {
	mcpCmd := &cobra.Command{
		Use:   "mcp",
		Short: "Model Context Protocol server for ccExplorer",
		Long: `Start a Model Context Protocol (MCP) server that exposes ccExplorer functionality
to AI systems through a standardized stdio interface.

The MCP server provides:
- get_cost_and_usage tool for AWS Cost Explorer queries
- Stdio transport for VSCode and other MCP clients`,
	}

	// Add serve subcommand
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the MCP server with stdio transport",
		Long: `Start the Model Context Protocol server with stdio transport.

This is the recommended transport mode for MCP clients like VSCode.

Examples:
  # Start MCP server for stdio transport (VSCode integration)
  ccexplorer mcp serve

The server will run until the MCP client disconnects or the process is terminated.`,
		RunE: runMCPServe,
	}

	mcpCmd.AddCommand(serveCmd)
	return mcpCmd
}

// runMCPServe handles the MCP serve command
func runMCPServe(cmd *cobra.Command, args []string) error {
	slog.Info("Starting ccExplorer MCP server with stdio transport")

	// Configure AWS service (reuse existing service initialization)
	if srv == nil {
		Initialize()
	}

	// Create MCP server
	mcpServer := mcp.NewServer(srv.aws)

	// Register tools
	if err := mcpServer.RegisterTools(); err != nil {
		return fmt.Errorf("failed to register MCP tools: %w", err)
	}

	slog.Info("MCP tools registered successfully, starting stdio server")

	// Start MCP server with stdio transport (recommended for MCP clients)
	if err := server.ServeStdio(mcpServer.MCPServer()); err != nil {
		return fmt.Errorf("MCP stdio server error: %w", err)
	}

	return nil
}