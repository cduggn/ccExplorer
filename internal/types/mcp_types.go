package types

// GetCostAndUsageInput is the typed input for the get_cost_and_usage MCP tool.
// The SDK auto-unmarshals JSON into this struct and infers the JSON schema from tags.
type GetCostAndUsageInput struct {
	StartDate        string `json:"start_date" jsonschema:"required,description=Start date in YYYY-MM-DD format"`
	EndDate          string `json:"end_date" jsonschema:"required,description=End date in YYYY-MM-DD format"`
	Granularity      string `json:"granularity" jsonschema:"enum=DAILY,enum=MONTHLY,enum=HOURLY,description=Query granularity"`
	Metrics          string `json:"metrics" jsonschema:"description=Cost metrics (comma-separated)"`
	GroupBy          string `json:"group_by" jsonschema:"description=Group by dimensions (comma-separated)"`
	FilterByService  string `json:"filter_by_service" jsonschema:"description=Filter by AWS service name"`
	ExcludeDiscounts bool   `json:"exclude_discounts" jsonschema:"description=Exclude discount line items"`
}

// MCPToolParameters represents the parsed parameters for get_cost_and_usage tool
type MCPToolParameters struct {
	StartDate         string            `json:"start_date"`
	EndDate           string            `json:"end_date"`
	Granularity       string            `json:"granularity"`
	Metrics           []string          `json:"metrics"`
	GroupBy           []string          `json:"group_by"`
	FilterByService   string            `json:"filter_by_service,omitempty"`
	FilterByDimension map[string]string `json:"filter_by_dimension,omitempty"`
	FilterByTag       map[string]string `json:"filter_by_tag,omitempty"`
	ExcludeDiscounts  bool              `json:"exclude_discounts,omitempty"`
}

