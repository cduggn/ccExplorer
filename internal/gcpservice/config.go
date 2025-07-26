package gcpservice

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config represents the configuration for GCP service client
type Config struct {
	// Authentication settings
	ServiceAccountPath string // Path to service account JSON file
	ServiceAccountJSON string // Service account JSON content
	ProjectID          string // Default project ID
	UseADC             bool   // Use Application Default Credentials

	// API and performance settings
	MaxRetries       int    // Maximum number of retries for failed requests
	RequestTimeout   int    // Request timeout in seconds (default: 30)
	RateLimitRPS     int    // Rate limit requests per second (default: 5)
	RateLimitBurst   int    // Rate limit burst capacity (default: 10)
	MaxConcurrency   int    // Maximum concurrent requests (default: 8)
	
	// Organization and billing settings
	OrganizationID   string   // Default organization ID
	BillingAccountID string   // Default billing account ID
	DefaultRegions   []string // Default regions for queries
	DefaultCurrency  string   // Default currency (default: USD)
	
	// Feature flags
	EnableAssetInventory bool // Enable Cloud Asset Inventory integration
	EnablePricing       bool // Enable Pricing API integration
	EnableCaching       bool // Enable response caching
	CacheTTL            int  // Cache TTL in minutes (default: 15)
	
	// Debug and monitoring
	Debug              bool   // Enable debug logging
	MetricsEnabled     bool   // Enable metrics collection
	TracingEnabled     bool   // Enable request tracing
}

// NewDefaultConfig creates a new configuration with default values
func NewDefaultConfig() *Config {
	return &Config{
		UseADC:               true,
		MaxRetries:           3,
		RequestTimeout:       30,
		RateLimitRPS:         5,
		RateLimitBurst:       10,
		MaxConcurrency:       8,
		DefaultCurrency:      "USD",
		EnableAssetInventory: true,
		EnablePricing:        true,
		EnableCaching:        true,
		CacheTTL:             15,
		Debug:                false,
		MetricsEnabled:       true,
		TracingEnabled:       false,
	}
}

// NewConfigFromEnv creates a configuration from environment variables
func NewConfigFromEnv() *Config {
	config := NewDefaultConfig()

	// Authentication
	if path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); path != "" {
		config.ServiceAccountPath = path
		config.UseADC = false
	}
	
	if json := os.Getenv("GCP_SERVICE_ACCOUNT_JSON"); json != "" {
		config.ServiceAccountJSON = json
		config.UseADC = false
	}

	if projectID := os.Getenv("GOOGLE_CLOUD_PROJECT"); projectID != "" {
		config.ProjectID = projectID
	}
	if projectID := os.Getenv("GCP_PROJECT_ID"); projectID != "" {
		config.ProjectID = projectID
	}

	// Organization and billing
	if orgID := os.Getenv("GCP_ORGANIZATION_ID"); orgID != "" {
		config.OrganizationID = orgID
	}
	
	if billingID := os.Getenv("GCP_BILLING_ACCOUNT_ID"); billingID != "" {
		config.BillingAccountID = billingID
	}

	if regions := os.Getenv("GCP_DEFAULT_REGIONS"); regions != "" {
		config.DefaultRegions = strings.Split(regions, ",")
		// Trim whitespace
		for i, region := range config.DefaultRegions {
			config.DefaultRegions[i] = strings.TrimSpace(region)
		}
	}

	if currency := os.Getenv("GCP_DEFAULT_CURRENCY"); currency != "" {
		config.DefaultCurrency = currency
	}

	// Performance settings
	if maxRetries := os.Getenv("GCP_MAX_RETRIES"); maxRetries != "" {
		if val, err := strconv.Atoi(maxRetries); err == nil && val > 0 {
			config.MaxRetries = val
		}
	}

	if timeout := os.Getenv("GCP_REQUEST_TIMEOUT"); timeout != "" {
		if val, err := strconv.Atoi(timeout); err == nil && val > 0 {
			config.RequestTimeout = val
		}
	}

	if rps := os.Getenv("GCP_RATE_LIMIT_RPS"); rps != "" {
		if val, err := strconv.Atoi(rps); err == nil && val > 0 {
			config.RateLimitRPS = val
		}
	}

	if burst := os.Getenv("GCP_RATE_LIMIT_BURST"); burst != "" {
		if val, err := strconv.Atoi(burst); err == nil && val > 0 {
			config.RateLimitBurst = val
		}
	}

	if concurrency := os.Getenv("GCP_MAX_CONCURRENCY"); concurrency != "" {
		if val, err := strconv.Atoi(concurrency); err == nil && val > 0 {
			config.MaxConcurrency = val
		}
	}

	// Feature flags
	if enableAssets := os.Getenv("GCP_ENABLE_ASSET_INVENTORY"); enableAssets != "" {
		config.EnableAssetInventory = strings.ToLower(enableAssets) == "true"
	}

	if enablePricing := os.Getenv("GCP_ENABLE_PRICING"); enablePricing != "" {
		config.EnablePricing = strings.ToLower(enablePricing) == "true"
	}

	if enableCaching := os.Getenv("GCP_ENABLE_CACHING"); enableCaching != "" {
		config.EnableCaching = strings.ToLower(enableCaching) == "true"
	}

	if cacheTTL := os.Getenv("GCP_CACHE_TTL"); cacheTTL != "" {
		if val, err := strconv.Atoi(cacheTTL); err == nil && val > 0 {
			config.CacheTTL = val
		}
	}

	// Debug and monitoring
	if debug := os.Getenv("GCP_DEBUG"); debug != "" {
		config.Debug = strings.ToLower(debug) == "true"
	}

	if metrics := os.Getenv("GCP_ENABLE_METRICS"); metrics != "" {
		config.MetricsEnabled = strings.ToLower(metrics) == "true"
	}

	if tracing := os.Getenv("GCP_ENABLE_TRACING"); tracing != "" {
		config.TracingEnabled = strings.ToLower(tracing) == "true"
	}

	return config
}

// Validate validates the configuration and returns an error if invalid
func (c *Config) Validate() error {
	// Validate authentication configuration
	if !c.UseADC && c.ServiceAccountPath == "" && c.ServiceAccountJSON == "" {
		return fmt.Errorf("authentication required: either set UseADC=true, provide ServiceAccountPath, or ServiceAccountJSON")
	}

	if c.ServiceAccountPath != "" && c.ServiceAccountJSON != "" {
		return fmt.Errorf("cannot specify both ServiceAccountPath and ServiceAccountJSON")
	}

	if c.ServiceAccountPath != "" {
		if _, err := os.Stat(c.ServiceAccountPath); os.IsNotExist(err) {
			return fmt.Errorf("service account file does not exist: %s", c.ServiceAccountPath)
		}
	}

	// Validate numeric parameters
	if c.MaxRetries < 0 || c.MaxRetries > 10 {
		return fmt.Errorf("MaxRetries must be between 0 and 10, got %d", c.MaxRetries)
	}

	if c.RequestTimeout < 1 || c.RequestTimeout > 300 {
		return fmt.Errorf("RequestTimeout must be between 1 and 300 seconds, got %d", c.RequestTimeout)
	}

	if c.RateLimitRPS < 1 || c.RateLimitRPS > 100 {
		return fmt.Errorf("RateLimitRPS must be between 1 and 100, got %d", c.RateLimitRPS)
	}

	if c.RateLimitBurst < 1 || c.RateLimitBurst > 50 {
		return fmt.Errorf("RateLimitBurst must be between 1 and 50, got %d", c.RateLimitBurst)
	}

	if c.MaxConcurrency < 1 || c.MaxConcurrency > 50 {
		return fmt.Errorf("MaxConcurrency must be between 1 and 50, got %d", c.MaxConcurrency)
	}

	if c.CacheTTL < 1 || c.CacheTTL > 1440 { // Max 24 hours
		return fmt.Errorf("CacheTTL must be between 1 and 1440 minutes, got %d", c.CacheTTL)
	}

	// Validate currency format
	validCurrencies := []string{"USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "SEK", "NOK", "DKK"}
	validCurrency := false
	for _, currency := range validCurrencies {
		if c.DefaultCurrency == currency {
			validCurrency = true
			break
		}
	}
	if !validCurrency {
		return fmt.Errorf("invalid DefaultCurrency: %s, must be one of %v", c.DefaultCurrency, validCurrencies)
	}

	// Validate region format (basic validation)
	for _, region := range c.DefaultRegions {
		if len(region) < 5 || !strings.Contains(region, "-") {
			return fmt.Errorf("invalid region format: %s (expected format like 'us-central1')", region)
		}
	}

	return nil
}

// GetAuthenticationMethod returns a description of the authentication method being used
func (c *Config) GetAuthenticationMethod() string {
	if c.ServiceAccountPath != "" {
		return fmt.Sprintf("Service Account File: %s", c.ServiceAccountPath)
	}
	if c.ServiceAccountJSON != "" {
		return "Service Account JSON (inline)"
	}
	if c.UseADC {
		return "Application Default Credentials (ADC)"
	}
	return "No authentication configured"
}

// IsProductionReady returns true if the configuration is suitable for production use
func (c *Config) IsProductionReady() bool {
	// Production readiness checks
	if c.Debug {
		return false // Debug mode should not be enabled in production
	}
	
	if c.MaxRetries < 2 {
		return false // Should have retries enabled
	}
	
	if c.RequestTimeout < 10 {
		return false // Timeout should be reasonable
	}
	
	if !c.MetricsEnabled {
		return false // Metrics should be enabled for monitoring
	}
	
	if c.RateLimitRPS > 20 {
		return false // Rate limiting should be conservative
	}

	return true
}

// Clone creates a deep copy of the configuration
func (c *Config) Clone() *Config {
	clone := &Config{
		ServiceAccountPath:   c.ServiceAccountPath,
		ServiceAccountJSON:   c.ServiceAccountJSON,
		ProjectID:            c.ProjectID,
		UseADC:               c.UseADC,
		MaxRetries:           c.MaxRetries,
		RequestTimeout:       c.RequestTimeout,
		RateLimitRPS:         c.RateLimitRPS,
		RateLimitBurst:       c.RateLimitBurst,
		MaxConcurrency:       c.MaxConcurrency,
		OrganizationID:       c.OrganizationID,
		BillingAccountID:     c.BillingAccountID,
		DefaultCurrency:      c.DefaultCurrency,
		EnableAssetInventory: c.EnableAssetInventory,
		EnablePricing:        c.EnablePricing,
		EnableCaching:        c.EnableCaching,
		CacheTTL:             c.CacheTTL,
		Debug:                c.Debug,
		MetricsEnabled:       c.MetricsEnabled,
		TracingEnabled:       c.TracingEnabled,
	}

	// Deep copy slice
	if c.DefaultRegions != nil {
		clone.DefaultRegions = make([]string, len(c.DefaultRegions))
		copy(clone.DefaultRegions, c.DefaultRegions)
	}

	return clone
}

// String returns a string representation of the configuration (sensitive data masked)
func (c *Config) String() string {
	authMethod := c.GetAuthenticationMethod()
	if strings.Contains(authMethod, "JSON") {
		authMethod = "Service Account JSON (***masked***)"
	}

	return fmt.Sprintf(`GCP Service Configuration:
  Authentication: %s
  Project ID: %s
  Organization ID: %s
  Billing Account: %s
  Default Currency: %s
  Default Regions: %v
  Max Retries: %d
  Request Timeout: %ds
  Rate Limit: %d RPS (burst: %d)
  Max Concurrency: %d
  Asset Inventory: %t
  Pricing API: %t
  Caching: %t (TTL: %dm)
  Debug: %t
  Metrics: %t
  Tracing: %t
  Production Ready: %t`,
		authMethod,
		c.ProjectID,
		c.OrganizationID,
		c.BillingAccountID,
		c.DefaultCurrency,
		c.DefaultRegions,
		c.MaxRetries,
		c.RequestTimeout,
		c.RateLimitRPS,
		c.RateLimitBurst,
		c.MaxConcurrency,
		c.EnableAssetInventory,
		c.EnablePricing,
		c.EnableCaching,
		c.CacheTTL,
		c.Debug,
		c.MetricsEnabled,
		c.TracingEnabled,
		c.IsProductionReady(),
	)
}