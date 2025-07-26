package types

import (
	"context"
	"fmt"
	"time"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/spf13/cobra"
)

type Command interface {
	Run(cmd *cobra.Command, args []string) error
	SynthesizeRequest(input interface{}) (interface{}, error)
	InputHandler(cmd *cobra.Command) interface{}
	Execute(req interface{}) error
	DefineFlags()
}

// Generic command interface for type-safe command implementations
type GenericCommand[TInput, TOutput any] interface {
	Run(cmd *cobra.Command, args []string) error
	SynthesizeRequest(input TInput) (TOutput, error)
	InputHandler(cmd *cobra.Command) TInput
	Execute(req TOutput) error
	DefineFlags()
}

type CostDataReader interface {
	ExtractGroupBySelections() ([]string, []string)
	ExtractFilterBySelection() (FilterBySelections, error)
	ExtractStartAndEndDates() (string, string, error)
	ExtractPrintPreferences() PrintOptions
}

type CommandLineInput struct {
	GroupByDimension    []string
	GroupByTag          []string
	FilterByValues      map[string]string
	IsFilterByTag       bool
	TagFilterValue      string
	IsFilterByDimension bool
	Start               string
	End                 string
	ExcludeDiscounts    bool
	Interval            string
	PrintFormat         string
	Metrics             []string
	SortByDate          bool
	OpenAIAPIKey        string
	PineconeIndex       string
	PineconeAPIKey      string
}

type FilterBySelections struct {
	Tags                string
	Dimensions          map[string]string
	IsFilterByTag       bool
	IsFilterByDimension bool
}

type PrintOptions struct {
	IsSortByDate     bool
	ExcludeDiscounts bool
	Format           string
	OpenAIKey        string
	Granularity      string
	Metric           string
	PineconeIndex    string
	PineconeAPIKey   string
}

type ForecastCommandLineInput struct {
	FilterByValues          Filter
	Granularity             string
	PredictionIntervalLevel int32
	Start                   string
	End                     string
}

type PresetParams struct {
	Alias             string
	Dimension         []string
	Tag               string
	Filter            map[string]string
	FilterType        string
	FilterByDimension bool
	FilterByTag       bool
	ExcludeDiscounts  bool
	CommandSyntax     string
	Description       []string
	Granularity       string
	PrintFormat       string
	Metric            []string
}

type GetCostForecastAPI interface {
	GetCostForecast(ctx context.Context, params *costexplorer.GetCostForecastInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostForecastOutput, error)
}

type GetCostForecastRequest struct {
	Time                    Time
	Granularity             string
	Metric                  string
	Filter                  Filter
	PredictionIntervalLevel int32
}

type GetCostForecastReport struct{}

type GetCostAndUsageAPI interface {
	GetCostAndUsage(ctx context.Context,
		optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error)
}

type Time struct {
	Start string
	End   string
}

type Dimension struct {
	Key   string
	Value []string
}

type Tag struct {
	Key   string
	Value []string
}

type Filter struct {
	Dimensions []Dimension
	Tags       []Tag
}

type CostAndUsageRequestType struct {
	Granularity                string
	GroupBy                    []string
	GroupByTag                 []string
	Time                       Time
	IsFilterByTagEnabled       bool
	IsFilterByDimensionEnabled bool
	TagFilterValue             string
	DimensionFilter            map[string]string
	ExcludeDiscounts           bool
	Alias                      string
	Rates                      []string
	PrintFormat                string
	Metrics                    []string
	SortByDate                 bool
	OpenAIAPIKey               string
	PineconeIndex              string
	PineconeAPIKey             string
}

type CostAndUsageRequestWithResourcesType struct {
	Granularity      string
	GroupBy          []string
	Tag              string
	Time             Time
	IsFilterEnabled  bool
	FilterType       string
	TagFilterValue   string
	Rates            []string
	ExcludeDiscounts bool
}

type GetDimensionValuesRequest struct {
	Dimension string
	Time      Time
}

func (t Time) Equals(other Time) bool {
	return t.Start == other.Start && t.End == other.End
}

func (c CostAndUsageRequestType) Equals(c2 CostAndUsageRequestType) bool {
	if c.Granularity != c2.Granularity {
		return false
	}
	if !c.Time.Equals(c2.Time) {
		return false
	}
	if c.IsFilterByTagEnabled != c2.IsFilterByTagEnabled {
		return false
	}
	if c.IsFilterByDimensionEnabled != c2.IsFilterByDimensionEnabled {
		return false
	}
	if c.TagFilterValue != c2.TagFilterValue {
		return false
	}
	if len(c.DimensionFilter) != len(c2.DimensionFilter) {
		return false
	}
	for k, v := range c.DimensionFilter {
		if v2, ok := c2.DimensionFilter[k]; !ok || v != v2 {
			return false
		}
	}
	if c.ExcludeDiscounts != c2.ExcludeDiscounts {
		return false
	}
	if c.Alias != c2.Alias {
		return false
	}
	if len(c.Rates) != len(c2.Rates) {
		return false
	}
	for i, v := range c.Rates {
		if v2 := c2.Rates[i]; v != v2 {
			return false
		}
	}
	return true
}

type APIError struct {
	Msg string
}

func (e APIError) Error() string {
	return e.Msg
}

type PresetError struct {
	Msg string
}

func (e PresetError) Error() string {
	return e.Msg
}

// Generic collection types for improved type safety

// GenericFilter provides type-safe filtering capabilities
type GenericFilter[T any] struct {
	Predicate func(T) bool
}

func (f GenericFilter[T]) Apply(items []T) []T {
	var result []T
	for _, item := range items {
		if f.Predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// GenericCollection provides common collection operations
type GenericCollection[T any] struct {
	Items []T
}

func NewGenericCollection[T any](items []T) *GenericCollection[T] {
	return &GenericCollection[T]{Items: items}
}

func (c *GenericCollection[T]) Filter(predicate func(T) bool) *GenericCollection[T] {
	var filtered []T
	for _, item := range c.Items {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}
	return &GenericCollection[T]{Items: filtered}
}

func (c *GenericCollection[T]) Map(mapper func(T) T) *GenericCollection[T] {
	mapped := make([]T, len(c.Items))
	for i, item := range c.Items {
		mapped[i] = mapper(item)
	}
	return &GenericCollection[T]{Items: mapped}
}

func (c *GenericCollection[T]) ToSlice() []T {
	return c.Items
}

// GenericServiceMap provides type-safe service operations
type GenericServiceMap[K comparable, V any] map[K]V

func (m GenericServiceMap[K, V]) Filter(predicate func(K, V) bool) GenericServiceMap[K, V] {
	result := make(GenericServiceMap[K, V])
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

func (m GenericServiceMap[K, V]) Transform(transformer func(V) V) GenericServiceMap[K, V] {
	result := make(GenericServiceMap[K, V])
	for k, v := range m {
		result[k] = transformer(v)
	}
	return result
}

// ===== GCP-Specific Types =====
// These types support comprehensive GCP Cloud Billing API integration

// GCPCommandLineInput represents user input for GCP cost analysis commands
type GCPCommandLineInput struct {
	ProjectID          string            // Primary GCP Project ID
	ProjectIDs         []string          // Multiple project IDs for batch operations
	BillingAccount     string            // GCP Billing Account ID
	OrganizationID     string            // GCP Organization ID for enterprise queries
	FolderIDs          []string          // GCP Folder IDs for hierarchical analysis
	Services           []string          // Filter by GCP service names
	Regions            []string          // Filter by GCP regions
	SKUs               []string          // Filter by specific SKUs
	Labels             map[string]string // Filter by resource labels
	Start              string            // Start date for cost analysis
	End                string            // End date for cost analysis
	Granularity        string            // Time granularity (DAILY, MONTHLY, HOURLY)
	Currency           string            // Currency for cost reporting
	GroupBy            []string          // Group results by dimensions
	IncludeCredits     bool              // Include promotional credits in analysis
	IncludeDiscounts   bool              // Include discounts in cost calculation
	PrintFormat        string            // Output format (stdout, csv, chart, pinecone)
	SortByDate         bool              // Sort results by date vs cost
	CostThreshold      float64           // Minimum cost threshold for filtering
	OpenAIAPIKey       string            // API key for AI/vector operations
	PineconeIndex      string            // Pinecone index for vector storage
	PineconeAPIKey     string            // Pinecone API key
}

// GCPBillingRequest represents a request for GCP billing data
type GCPBillingRequest struct {
	ProjectID          string            // Primary project for billing query
	ProjectIDs         []string          // Additional projects for multi-project analysis
	BillingAccount     string            // Billing account filter
	OrganizationID     string            // Organization-level query
	FolderIDs          []string          // Folder-level queries
	Services           []string          // Service name filters
	Regions            []string          // Regional filters
	SKUs               []string          // SKU-level filters
	Labels             map[string]string // Label-based filters
	Time               GCPTimeRange      // Time range for analysis
	Granularity        string            // Time granularity
	Currency           string            // Preferred currency
	GroupBy            []string          // Grouping dimensions
	IncludeCredits     bool              // Include credits in results
	IncludeDiscounts   bool              // Include discounts in results
	CostThreshold      float64           // Minimum cost filter
	PageSize           int32             // Pagination size for large results
	PageToken          string            // Pagination token
	PrintFormat        string            // Output format preference
	SortByDate         bool              // Sorting preference
	OpenAIAPIKey       string            // AI integration key
	PineconeIndex      string            // Vector database index
	PineconeAPIKey     string            // Vector database key
}

// GCPTimeRange represents a time range for GCP billing queries
type GCPTimeRange struct {
	Start string // Start date in YYYY-MM-DD format
	End   string // End date in YYYY-MM-DD format
}

// GCPBillingResponse represents the response from GCP billing queries
type GCPBillingResponse struct {
	Services        map[string]GCPBillingItem `json:"services"`        // Service-level billing data
	ProjectInfo     string                    `json:"project_info"`    // Primary project information
	BillingAccount  string                    `json:"billing_account"` // Associated billing account
	Currency        string                    `json:"currency"`        // Response currency
	Granularity     string                    `json:"granularity"`     // Time granularity used
	TimeRange       GCPTimeRange              `json:"time_range"`      // Query time range
	TotalCost       float64                   `json:"total_cost"`      // Total cost across all services
	TotalCredits    float64                   `json:"total_credits"`   // Total credits applied
	TotalDiscounts  float64                   `json:"total_discounts"` // Total discounts applied
	Metadata        GCPResponseMetadata       `json:"metadata"`        // Response metadata
	NextPageToken   string                    `json:"next_page_token"` // Pagination token for next page
}

// GCPBillingItem represents billing information for a single GCP service
type GCPBillingItem struct {
	ServiceID          string                 `json:"service_id"`          // GCP service identifier
	ServiceDisplayName string                 `json:"service_display_name"` // Human-readable service name
	ProjectID          string                 `json:"project_id"`          // Associated project
	BillingAccount     string                 `json:"billing_account"`     // Associated billing account
	Currency           string                 `json:"currency"`            // Cost currency
	TotalCost          float64                `json:"total_cost"`          // Total service cost
	Granularity        string                 `json:"granularity"`         // Time granularity
	TimeRange          GCPTimeRange           `json:"time_range"`          // Service query time range
	SKUs               []GCPSKU               `json:"skus"`                // Service SKU breakdown
	RegionalCosts      map[string]float64     `json:"regional_costs"`      // Cost by region
	UsageMetrics       map[string]float64     `json:"usage_metrics"`       // Usage statistics
	Credits            []GCPCredit            `json:"credits"`             // Applied credits
	Discounts          []GCPDiscount          `json:"discounts"`           // Applied discounts
	Labels             map[string]string      `json:"labels"`              // Resource labels
	Tags               map[string]string      `json:"tags"`                // Resource tags
}

// GCPSKU represents a single Stock Keeping Unit (pricing unit) in GCP
type GCPSKU struct {
	SKUID          string               `json:"sku_id"`          // Unique SKU identifier
	DisplayName    string               `json:"display_name"`    // Human-readable SKU name
	Description    string               `json:"description"`     // SKU description
	Category       string               `json:"category"`        // Service category
	Family         string               `json:"family"`          // Resource family
	Group          string               `json:"group"`           // Resource group
	Usage          string               `json:"usage"`           // Usage type
	Unit           string               `json:"unit"`            // Pricing unit
	Regions        []string             `json:"regions"`         // Available regions
	PricingInfo    []GCPPricingTier     `json:"pricing_info"`    // Pricing tier information
	EstimatedCost  float64              `json:"estimated_cost"`  // Estimated cost for time period
	ActualUsage    float64              `json:"actual_usage"`    // Actual usage amount
}

// GCPPricingTier represents tiered pricing information for a SKU
type GCPPricingTier struct {
	Currency         string         `json:"currency"`          // Pricing currency
	Summary          string         `json:"summary"`           // Pricing summary
	Tiers            []GCPTierRate  `json:"tiers"`             // Tiered rate structure
	AggregationLevel string         `json:"aggregation_level"` // Aggregation method
	AggregationCount int32          `json:"aggregation_count"` // Aggregation period
}

// GCPTierRate represents a single pricing tier rate
type GCPTierRate struct {
	StartUsageAmount float64  `json:"start_usage_amount"` // Usage amount where tier starts
	UnitPrice        GCPPrice `json:"unit_price"`         // Price per unit in this tier
}

// GCPPrice represents a monetary amount in GCP
type GCPPrice struct {
	CurrencyCode string `json:"currency_code"` // ISO currency code
	Units        int64  `json:"units"`         // Whole units of currency
	Nanos        int32  `json:"nanos"`         // Fractional units (nano-units)
}

// ToFloat64 converts GCPPrice to float64 for calculations
func (p GCPPrice) ToFloat64() float64 {
	return float64(p.Units) + float64(p.Nanos)/1e9
}

// GCPCredit represents promotional credits or discounts applied
type GCPCredit struct {
	CreditID     string  `json:"credit_id"`     // Credit identifier
	Name         string  `json:"name"`          // Credit name/description
	Type         string  `json:"type"`          // Credit type (promotional, sustained_use, etc.)
	Amount       float64 `json:"amount"`        // Credit amount
	Currency     string  `json:"currency"`      // Credit currency
	AppliedDate  string  `json:"applied_date"`  // Date credit was applied
	ExpiryDate   string  `json:"expiry_date"`   // Credit expiration date
}

// GCPDiscount represents discounts applied to billing
type GCPDiscount struct {
	DiscountID   string  `json:"discount_id"`   // Discount identifier
	Name         string  `json:"name"`          // Discount name
	Type         string  `json:"type"`          // Discount type
	Amount       float64 `json:"amount"`        // Discount amount
	Currency     string  `json:"currency"`      // Discount currency
	AppliedDate  string  `json:"applied_date"`  // Date discount was applied
}

// GCPResponseMetadata contains metadata about the API response
type GCPResponseMetadata struct {
	RequestID       string        `json:"request_id"`       // Unique request identifier
	ResponseTime    time.Time     `json:"response_time"`    // Response timestamp
	QueryDuration   time.Duration `json:"query_duration"`   // Time taken to process query
	APIVersion      string        `json:"api_version"`      // API version used
	RecordCount     int           `json:"record_count"`     // Number of records returned
	TotalRecords    int           `json:"total_records"`    // Total records available
	HasMoreResults  bool          `json:"has_more_results"` // Whether more results are available
}

// GCPAssetRequest represents a request for GCP asset inventory data
type GCPAssetRequest struct {
	Scope         string            `json:"scope"`          // Search scope (projects/*, folders/*, organizations/*)
	Query         string            `json:"query"`          // Asset search query
	AssetTypes    []string          `json:"asset_types"`    // Filter by asset types
	ProjectIDs    []string          `json:"project_ids"`    // Filter by project IDs
	Regions       []string          `json:"regions"`        // Filter by regions
	Labels        map[string]string `json:"labels"`         // Filter by labels
	PageSize      int32             `json:"page_size"`      // Results per page
	PageToken     string            `json:"page_token"`     // Pagination token
	OrderBy       string            `json:"order_by"`       // Result ordering
	ReadTime      time.Time         `json:"read_time"`      // Point-in-time read
}

// GCPForecastRequest represents a request for GCP cost forecasting
type GCPForecastRequest struct {
	ProjectID        string       `json:"project_id"`         // Project for forecasting
	ProjectIDs       []string     `json:"project_ids"`        // Multiple projects
	BillingAccount   string       `json:"billing_account"`    // Billing account
	Services         []string     `json:"services"`           // Service filters
	Regions          []string     `json:"regions"`            // Regional filters
	TimeRange        GCPTimeRange `json:"time_range"`         // Historical data range
	ForecastPeriod   GCPTimeRange `json:"forecast_period"`    // Forecast time range
	Granularity      string       `json:"granularity"`        // Forecast granularity
	Currency         string       `json:"currency"`           // Forecast currency
	ConfidenceLevel  float64      `json:"confidence_level"`   // Confidence level (0.0-1.0)
	IncludeCredits   bool         `json:"include_credits"`    // Include credits in forecast
	IncludeDiscounts bool         `json:"include_discounts"`  // Include discounts in forecast
	ModelType        string       `json:"model_type"`         // Forecasting model type
}

// GCPForecastResponse represents the response from cost forecasting
type GCPForecastResponse struct {
	ProjectID          string                  `json:"project_id"`           // Primary project
	BillingAccount     string                  `json:"billing_account"`      // Billing account
	Currency           string                  `json:"currency"`             // Forecast currency
	Granularity        string                  `json:"granularity"`          // Time granularity
	ForecastPeriod     GCPTimeRange            `json:"forecast_period"`      // Forecast time range
	HistoricalPeriod   GCPTimeRange            `json:"historical_period"`    // Historical data period
	TotalForecast      float64                 `json:"total_forecast"`       // Total forecasted cost
	ConfidenceLevel    float64                 `json:"confidence_level"`     // Forecast confidence
	ModelType          string                  `json:"model_type"`           // Model used for forecast
	ForecastBreakdown  []GCPForecastDataPoint  `json:"forecast_breakdown"`   // Time-series forecast data
	ServiceForecasts   map[string]float64      `json:"service_forecasts"`    // Per-service forecasts
	RegionalForecasts  map[string]float64      `json:"regional_forecasts"`   // Per-region forecasts
	VarianceAnalysis   GCPVarianceAnalysis     `json:"variance_analysis"`    // Forecast variance analysis
	Metadata           GCPResponseMetadata     `json:"metadata"`             // Response metadata
}

// GCPForecastDataPoint represents a single point in the forecast time series
type GCPForecastDataPoint struct {
	Date           string  `json:"date"`             // Date for this data point
	ForecastCost   float64 `json:"forecast_cost"`    // Forecasted cost
	LowerBound     float64 `json:"lower_bound"`      // Lower confidence bound
	UpperBound     float64 `json:"upper_bound"`      // Upper confidence bound
	HistoricalCost float64 `json:"historical_cost"`  // Historical cost (if available)
}

// GCPVarianceAnalysis provides analysis of forecast variance and trends
type GCPVarianceAnalysis struct {
	TrendDirection     string  `json:"trend_direction"`     // INCREASING, DECREASING, STABLE
	MonthOverMonth     float64 `json:"month_over_month"`    // Month-over-month change percentage
	YearOverYear       float64 `json:"year_over_year"`      // Year-over-year change percentage
	SeasonalityFactor  float64 `json:"seasonality_factor"`  // Detected seasonality impact
	VolatilityScore    float64 `json:"volatility_score"`    // Cost volatility score (0-1)
	ConfidenceScore    float64 `json:"confidence_score"`    // Overall forecast confidence (0-1)
	TopCostDrivers     []string `json:"top_cost_drivers"`   // Primary cost drivers identified
	RiskFactors        []string `json:"risk_factors"`       // Potential risk factors
}

// GCPFilterSelections represents user filter selections for GCP queries
type GCPFilterSelections struct {
	Services         []string          `json:"services"`          // Service filters
	Regions          []string          `json:"regions"`           // Region filters
	Projects         []string          `json:"projects"`          // Project filters
	SKUs             []string          `json:"skus"`              // SKU filters
	Labels           map[string]string `json:"labels"`            // Label filters
	BillingAccounts  []string          `json:"billing_accounts"`  // Billing account filters
	ResourceTypes    []string          `json:"resource_types"`    // Resource type filters
	CostThreshold    float64           `json:"cost_threshold"`    // Minimum cost threshold
	UsageThreshold   float64           `json:"usage_threshold"`   // Minimum usage threshold
	IsFilterByLabel  bool              `json:"is_filter_by_label"` // Whether label filters are active
	IsFilterByRegion bool              `json:"is_filter_by_region"` // Whether region filters are active
}

// GCPPrintOptions represents output formatting options for GCP data
type GCPPrintOptions struct {
	Format             string  `json:"format"`               // Output format
	SortByDate         bool    `json:"sort_by_date"`         // Sort by date vs cost
	IncludeCredits     bool    `json:"include_credits"`      // Show credits in output
	IncludeDiscounts   bool    `json:"include_discounts"`    // Show discounts in output
	ShowRegionalBreakdown bool `json:"show_regional_breakdown"` // Include regional cost breakdown
	ShowSKUDetails     bool    `json:"show_sku_details"`     // Include SKU-level details
	ShowUsageMetrics   bool    `json:"show_usage_metrics"`   // Include usage statistics
	Granularity        string  `json:"granularity"`          // Time granularity for output
	Currency           string  `json:"currency"`             // Display currency
	CostThreshold      float64 `json:"cost_threshold"`       // Minimum cost to display
	TopN               int     `json:"top_n"`                // Show top N results only
	OpenAIKey          string  `json:"openai_key"`           // OpenAI API key
	PineconeIndex      string  `json:"pinecone_index"`       // Pinecone index name
	PineconeAPIKey     string  `json:"pinecone_api_key"`     // Pinecone API key
}

// Equals method for GCPBillingRequest to support testing and comparison
func (req GCPBillingRequest) Equals(other GCPBillingRequest) bool {
	if req.ProjectID != other.ProjectID ||
		req.BillingAccount != other.BillingAccount ||
		req.Granularity != other.Granularity ||
		req.Currency != other.Currency ||
		req.Time.Start != other.Time.Start ||
		req.Time.End != other.Time.End {
		return false
	}
	
	// Compare slices
	if len(req.Services) != len(other.Services) {
		return false
	}
	for i, service := range req.Services {
		if service != other.Services[i] {
			return false
		}
	}
	
	if len(req.Regions) != len(other.Regions) {
		return false
	}
	for i, region := range req.Regions {
		if region != other.Regions[i] {
			return false
		}
	}
	
	// Compare maps
	if len(req.Labels) != len(other.Labels) {
		return false
	}
	for k, v := range req.Labels {
		if other.Labels[k] != v {
			return false
		}
	}
	
	return true
}

// Validate method for GCPBillingRequest
func (req GCPBillingRequest) Validate() error {
	if req.ProjectID == "" && len(req.ProjectIDs) == 0 {
		return fmt.Errorf("at least one project ID must be specified")
	}
	
	if req.Time.Start == "" || req.Time.End == "" {
		return fmt.Errorf("start and end dates are required")
	}
	
	validGranularities := []string{"DAILY", "MONTHLY", "HOURLY"}
	validGranularity := false
	for _, valid := range validGranularities {
		if req.Granularity == valid {
			validGranularity = true
			break
		}
	}
	if !validGranularity {
		return fmt.Errorf("invalid granularity: %s", req.Granularity)
	}
	
	return nil
}
