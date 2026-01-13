package types

// MCPRequest represents the incoming MCP tool request
type MCPRequest struct {
	ToolName   string                 `json:"tool_name"`
	Parameters MCPToolParameters      `json:"parameters"`
}

// MCPToolParameters represents the parameters for get_cost_and_usage tool
type MCPToolParameters struct {
	StartDate         string   `json:"start_date"`
	EndDate           string   `json:"end_date"`
	Granularity       string   `json:"granularity"`
	Metrics           []string `json:"metrics"`
	GroupBy           []string `json:"group_by"`
	FilterByService   string   `json:"filter_by_service,omitempty"`
	FilterByDimension map[string]string `json:"filter_by_dimension,omitempty"`
	FilterByTag       map[string]string `json:"filter_by_tag,omitempty"`
	ExcludeDiscounts  bool     `json:"exclude_discounts,omitempty"`
}

// MCPResponse represents the MCP tool response
type MCPResponse struct {
	Content []MCPContent `json:"content"`
	IsError bool         `json:"isError,omitempty"`
}

// MCPContent represents the content within an MCP response
type MCPContent struct {
	Type string      `json:"type"`
	Text string      `json:"text,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// MCPError represents an MCP error response
type MCPError struct {
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

