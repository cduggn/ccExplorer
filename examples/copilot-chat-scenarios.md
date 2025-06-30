# GitHub Copilot Chat Scenarios with ccExplorer MCP

This document provides practical examples of how to use ccExplorer's MCP server with GitHub Copilot Chat in VSCode for AI-powered AWS cost analysis.

## Prerequisites

- VSCode with MCP enabled (`"chat.mcp.enabled": true`)
- ccExplorer MCP server configured (see [VSCode MCP Integration Guide](../docs/vscode-mcp-integration.md))
- VSCode in Agent Mode with ccExplorer tools enabled

## Basic Cost Analysis Scenarios

### Monthly Cost Overview

```
@agent What were my total AWS costs for the last month grouped by service?
```

**Expected Result**: Copilot will use the `get_cost_and_usage` tool to retrieve costs from the previous month, automatically calculating the date range and grouping by AWS service.

### Service-Specific Analysis

```
@agent Show me my EC2 costs for the last 3 months with monthly breakdown
```

**Expected Result**: Retrieves EC2-specific costs with monthly granularity, filtered to only include EC2 services.

```
@agent What are my S3 storage costs for the current quarter, excluding any discounts?
```

**Expected Result**: Filters for Amazon S3 services over the current quarter with discounts excluded.

## Advanced Cost Analysis

### Cost Optimization Scenarios

```
@agent Help me identify my top 5 most expensive AWS services this month and suggest optimization opportunities
```

**Expected Result**: Copilot retrieves current month costs grouped by service, analyzes the top 5, and provides optimization recommendations based on the cost patterns.

### Cost Trend Analysis

```
@agent Compare my AWS costs between last month and this month. Which services had the biggest increase?
```

**Expected Result**: Copilot will make two separate API calls for different date ranges and compare the results to identify cost increases.

```
@agent Show me daily AWS costs for the last 7 days. Are there any unusual spikes?
```

**Expected Result**: Retrieves daily granularity data and analyzes for cost anomalies or spikes.

## Tag-Based Cost Analysis

### Project Cost Tracking

```
@agent What are the costs for my ProjectName tag "web-app" for the last month?
```

**Expected Result**: Groups costs by the ProjectName tag and filters for "web-app" value.

```
@agent Compare costs between my different projects (ProjectName tags) for the current month
```

**Expected Result**: Groups all costs by ProjectName tag to show project-by-project breakdown.

### Environment Cost Analysis

```
@agent Show me production vs development environment costs for the last quarter
```

**Expected Result**: Assumes Environment tags exist and compares costs between production and development environments.

## Operational Cost Scenarios

### Budget Monitoring

```
@agent My AWS budget is $1000/month. How much have I spent so far this month and am I on track?
```

**Expected Result**: Copilot calculates month-to-date spending and projects end-of-month costs based on current usage patterns.

### Cost Forecasting

```
@agent Based on my current AWS usage, what will my costs likely be next month?
```

**Expected Result**: While ccExplorer currently focuses on historical data, Copilot can analyze trends and provide projections.

## Troubleshooting and Analysis

### Unexpected Cost Investigation

```
@agent I notice my AWS bill increased significantly this month. Can you help me identify what's causing the increase?
```

**Expected Result**: Copilot compares current month to previous months, identifies services with the largest increases, and suggests investigation areas.

```
@agent Show me all costs over $100 for individual services in the last month
```

**Expected Result**: Retrieves and filters cost data to show only high-cost services.

### Resource Utilization Analysis

```
@agent What are my data transfer costs for the last month? Break it down by operation type.
```

**Expected Result**: Filters for data transfer operations and groups by operation type to show bandwidth costs.

```
@agent Show me my storage costs across all AWS services for the last quarter
```

**Expected Result**: Identifies storage-related operations across multiple AWS services.

## Complex Multi-Dimensional Analysis

### Cross-Service Analysis

```
@agent Compare my compute costs (EC2, Lambda, ECS) vs storage costs (S3, EBS) for the last 3 months
```

**Expected Result**: Copilot categorizes services into compute and storage buckets and compares their costs over time.

### Geographic Cost Analysis

```
@agent Show me costs by AWS region for the last month. Which regions are most expensive?
```

**Expected Result**: Groups costs by availability zone or region dimension to show geographic cost distribution.

## Interactive Follow-Up Scenarios

### Drill-Down Analysis

**Initial Query:**
```
@agent What were my AWS costs last month grouped by service?
```

**Follow-up:**
```
@agent Can you break down the EC2 costs by operation type?
```

**Further Drill-down:**
```
@agent For the EC2 RunInstances operation, show me the daily breakdown
```

### Comparative Analysis

**Initial Query:**
```
@agent Show me my current month AWS costs
```

**Follow-up:**
```
@agent How does this compare to the same month last year?
```

## Tips for Effective Queries

### Specific Date Ranges
- Use specific months: "January 2024", "last quarter", "past 6 months"
- Be explicit about time periods: "from January 1 to March 31, 2024"

### Clear Grouping Preferences
- Specify how you want data grouped: "by service", "by tag", "by operation"
- Ask for specific dimensions: "group by ProjectName tag and service"

### Output Preferences
- Request specific formats: "show as a table", "summarize the key findings"
- Ask for explanations: "explain why these costs occurred"

### Filter Specifications
- Be specific about services: "only EC2 costs", "exclude data transfer charges"
- Mention discount preferences: "including discounts", "excluding credits and refunds"

## Error Handling and Troubleshooting

If Copilot indicates it cannot access cost data:

1. **Check MCP Server Status**:
```
@agent Can you check if the ccExplorer MCP server is running?
```

2. **Verify AWS Credentials**:
```
@agent Try to get costs for yesterday to test the connection
```

3. **Simplify the Query**:
```
@agent Just show me total AWS costs for last month
```

## Best Practices

1. **Start Simple**: Begin with basic queries and add complexity
2. **Be Specific**: Provide clear date ranges and grouping preferences
3. **Use Follow-ups**: Build on previous queries for deeper analysis
4. **Verify Results**: Cross-check important findings with AWS Console
5. **Consider Context**: Remember that Copilot maintains conversation context for follow-up questions

These scenarios demonstrate the power of combining ccExplorer's cost analysis capabilities with Copilot's natural language interface and analytical capabilities.