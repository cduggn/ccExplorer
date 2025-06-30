# VSCode + GitHub Copilot Integration with ccExplorer MCP Server

This guide explains how to integrate ccExplorer's built-in MCP (Model Context Protocol) server with VSCode and GitHub Copilot Chat for AI-powered AWS cost analysis.

## Prerequisites

- VSCode 1.99 or later
- GitHub Copilot extension enabled
- ccExplorer binary built (`make build`)
- AWS credentials configured (AWS CLI, environment variables, or IAM roles)

## Quick Setup

### 1. Enable MCP in VSCode

Add this setting to your VSCode configuration:

```json
{
  "chat.mcp.enabled": true
}
```

### 2. Workspace Configuration (Recommended)

The repository includes a pre-configured `.vscode/mcp.json` file that will automatically connect ccExplorer's MCP server to VSCode Copilot Chat.

When you first use the MCP server, VSCode will prompt you for:
- **AWS Access Key ID** (optional - leave empty to use default credential chain)
- **AWS Secret Access Key** (optional - leave empty to use default credential chain) 
- **AWS Region** (e.g., `us-east-1`, `us-west-2`)

### 3. Build ccExplorer

Ensure the binary is available:

```bash
make build
```

## Alternative Setup Methods

### User Settings (Global Configuration)

To enable ccExplorer MCP across all workspaces, add this to your VSCode user settings:

```json
{
  "chat.mcp.enabled": true,
  "chat.mcp.servers": {
    "ccExplorer": {
      "type": "stdio",
      "command": "/path/to/your/ccexplorer",
      "args": ["mcp", "serve"],
      "env": {
        "AWS_REGION": "us-east-1"
      }
    }
  }
}
```

## Usage with GitHub Copilot Chat

### 1. Open Copilot Chat in Agent Mode

1. Open the Chat view (`Ctrl+Shift+P` → "Chat: Open")
2. Select **Agent Mode** from the dropdown
3. Click the **Tools** button and ensure ccExplorer tools are enabled

### 2. Example Queries

Here are practical examples of how to use ccExplorer through Copilot Chat:

#### Basic Cost Analysis
```
@agent What were my AWS costs for the last 30 days grouped by service?
```

#### Monthly Cost Trends
```
@agent Show my AWS costs for the last 6 months with monthly granularity, grouped by service
```

#### Service-Specific Analysis
```
@agent What are my EC2 costs for the last quarter, excluding discounts?
```

#### Tag-Based Cost Analysis
```
@agent Show costs grouped by ProjectName tag for the last month
```

#### Cost Comparison
```
@agent Compare my AWS costs between last month and this month
```

### 3. Available MCP Tools

The ccExplorer MCP server provides the following tool:

#### `get_cost_and_usage`

**Parameters:**
- `start_date` (required): Start date in YYYY-MM-DD format
- `end_date` (required): End date in YYYY-MM-DD format  
- `granularity` (optional): DAILY, MONTHLY, or HOURLY
- `metrics` (optional): Cost metrics to retrieve (UnblendedCost, AmortizedCost, etc.)
- `group_by` (optional): Group results by SERVICE, AZ, INSTANCE_TYPE, or TAG:TagName
- `filter_by_service` (optional): Filter to specific AWS services
- `exclude_discounts` (optional): Exclude discount information

## Troubleshooting

### Check MCP Server Status

Use VSCode's Command Palette:
```
MCP: List Servers
```

This shows configured servers and their status, with options to start/stop/restart and view logs.

### Common Issues

1. **Server Not Starting**
   - Ensure `./ccexplorer` binary exists and is executable
   - Check AWS credentials are configured correctly
   - Verify AWS region is set

2. **Permission Denied**
   - Make sure ccExplorer binary has execute permissions: `chmod +x ./ccexplorer`

3. **AWS Authentication Errors**
   - Verify AWS credentials with: `aws sts get-caller-identity`
   - Check that your AWS credentials have Cost Explorer permissions

4. **Tool Not Available in Chat**
   - Ensure you're in Agent Mode
   - Check that ccExplorer tools are enabled in the Tools panel
   - Restart the MCP server if needed

### Debug Mode

For detailed logging, you can run the MCP server manually with environment variables:

```bash
# Enable debug logging
CCEXPLORER_LOG_LEVEL=debug ./ccexplorer mcp serve
```

### Test MCP Server

You can test the MCP server directly using stdio:

```bash
# Test tools list
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | ./ccexplorer mcp serve

# Test tool call (requires AWS credentials)
echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_cost_and_usage","arguments":{"start_date":"2024-01-01","end_date":"2024-01-31"}}}' | ./ccexplorer mcp serve
```

## Security Considerations

- The MCP server runs locally and connects directly to AWS APIs
- AWS credentials are handled through standard AWS credential chain
- No cost data is sent to external services beyond AWS
- VSCode may prompt for tool execution confirmation for security

## Advanced Configuration

### Environment Variables

You can configure the MCP server using environment variables instead of VSCode inputs:

```bash
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=your-key-id
export AWS_SECRET_ACCESS_KEY=your-secret-key
```

### Auto-Discovery

VSCode can automatically discover MCP servers configured in other tools like Claude Desktop. Enable with:

```json
{
  "chat.mcp.discovery.enabled": true
}
```

## Support

For issues with the MCP integration:
1. Check the [ccExplorer GitHub issues](https://github.com/cduggn/ccExplorer/issues)
2. Review VSCode MCP documentation
3. Verify AWS permissions and configuration

The MCP server provides comprehensive AWS cost analysis capabilities directly within your development environment, making it easy to understand and optimize cloud spending without leaving your IDE.