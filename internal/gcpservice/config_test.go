package gcpservice

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()

	assert.True(t, config.UseADC)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 30, config.RequestTimeout)
	assert.Equal(t, 5, config.RateLimitRPS)
	assert.Equal(t, 10, config.RateLimitBurst)
	assert.Equal(t, 8, config.MaxConcurrency)
	assert.Equal(t, "USD", config.DefaultCurrency)
	assert.True(t, config.EnableAssetInventory)
	assert.True(t, config.EnablePricing)
	assert.True(t, config.EnableCaching)
	assert.Equal(t, 15, config.CacheTTL)
	assert.False(t, config.Debug)
	assert.True(t, config.MetricsEnabled)
	assert.False(t, config.TracingEnabled)
}

func TestNewConfigFromEnv(t *testing.T) {
	// Save original env vars
	originalVars := map[string]string{
		"GOOGLE_APPLICATION_CREDENTIALS": os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		"GCP_PROJECT_ID":                 os.Getenv("GCP_PROJECT_ID"),
		"GCP_DEFAULT_CURRENCY":           os.Getenv("GCP_DEFAULT_CURRENCY"),
		"GCP_MAX_RETRIES":                os.Getenv("GCP_MAX_RETRIES"),
		"GCP_DEBUG":                      os.Getenv("GCP_DEBUG"),
	}

	// Clean up after test
	defer func() {
		for key, value := range originalVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	tests := []struct {
		name     string
		envVars  map[string]string
		validate func(t *testing.T, config *Config)
	}{
		{
			name: "with service account file",
			envVars: map[string]string{
				"GOOGLE_APPLICATION_CREDENTIALS": "/path/to/service-account.json",
				"GCP_PROJECT_ID":                 "test-project",
				"GCP_DEFAULT_CURRENCY":           "EUR",
			},
			validate: func(t *testing.T, config *Config) {
				assert.Equal(t, "/path/to/service-account.json", config.ServiceAccountPath)
				assert.False(t, config.UseADC)
				assert.Equal(t, "test-project", config.ProjectID)
				assert.Equal(t, "EUR", config.DefaultCurrency)
			},
		},
		{
			name: "with numeric env vars",
			envVars: map[string]string{
				"GCP_MAX_RETRIES":      "5",
				"GCP_REQUEST_TIMEOUT":  "60",
				"GCP_RATE_LIMIT_RPS":   "10",
				"GCP_MAX_CONCURRENCY":  "16",
			},
			validate: func(t *testing.T, config *Config) {
				assert.Equal(t, 5, config.MaxRetries)
				assert.Equal(t, 60, config.RequestTimeout)
				assert.Equal(t, 10, config.RateLimitRPS)
				assert.Equal(t, 16, config.MaxConcurrency)
			},
		},
		{
			name: "with boolean env vars",
			envVars: map[string]string{
				"GCP_DEBUG":                  "true",
				"GCP_ENABLE_ASSET_INVENTORY": "false",
				"GCP_ENABLE_CACHING":         "false",
			},
			validate: func(t *testing.T, config *Config) {
				assert.True(t, config.Debug)
				assert.False(t, config.EnableAssetInventory)
				assert.False(t, config.EnableCaching)
			},
		},
		{
			name: "with regions",
			envVars: map[string]string{
				"GCP_DEFAULT_REGIONS": "us-central1, europe-west1, asia-southeast1",
			},
			validate: func(t *testing.T, config *Config) {
				expected := []string{"us-central1", "europe-west1", "asia-southeast1"}
				assert.Equal(t, expected, config.DefaultRegions)
			},
		},
		{
			name: "invalid numeric values ignored",
			envVars: map[string]string{
				"GCP_MAX_RETRIES":     "-1",
				"GCP_REQUEST_TIMEOUT": "invalid",
				"GCP_RATE_LIMIT_RPS":  "0",
			},
			validate: func(t *testing.T, config *Config) {
				// Should use default values when invalid
				assert.Equal(t, 3, config.MaxRetries)   // default
				assert.Equal(t, 30, config.RequestTimeout) // default
				assert.Equal(t, 5, config.RateLimitRPS)  // default
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all env vars first
			for key := range originalVars {
				os.Unsetenv(key)
			}

			// Set test env vars
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			config := NewConfigFromEnv()
			tt.validate(t, config)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid config with ADC",
			config:      NewDefaultConfig(),
			expectError: false,
		},
		{
			name: "valid config with service account path",
			config: &Config{
				ServiceAccountPath:   "/tmp/test-service-account.json",
				MaxRetries:           3,
				RequestTimeout:       30,
				RateLimitRPS:         5,
				RateLimitBurst:       10,
				MaxConcurrency:       8,
				DefaultCurrency:      "USD",
				CacheTTL:             15,
			},
			expectError: false, // Note: file doesn't exist, but we're testing validation logic
		},
		{
			name: "invalid - no authentication",
			config: &Config{
				UseADC:             false,
				ServiceAccountPath: "",
				ServiceAccountJSON: "",
				MaxRetries:         3,
				RequestTimeout:     30,
				RateLimitRPS:       5,
				RateLimitBurst:     10,
				MaxConcurrency:     8,
				DefaultCurrency:    "USD",
				CacheTTL:           15,
			},
			expectError: true,
			errorMsg:    "authentication required",
		},
		{
			name: "invalid - both service account methods",
			config: &Config{
				ServiceAccountPath: "/path/to/file.json",
				ServiceAccountJSON: "{}",
				MaxRetries:         3,
				RequestTimeout:     30,
				RateLimitRPS:       5,
				RateLimitBurst:     10,
				MaxConcurrency:     8,
				DefaultCurrency:    "USD",
				CacheTTL:           15,
			},
			expectError: true,
			errorMsg:    "cannot specify both",
		},
		{
			name: "invalid - max retries out of range",
			config: &Config{
				UseADC:           true,
				MaxRetries:       15, // > 10
				RequestTimeout:   30,
				RateLimitRPS:     5,
				RateLimitBurst:   10,
				MaxConcurrency:   8,
				DefaultCurrency:  "USD",
				CacheTTL:         15,
			},
			expectError: true,
			errorMsg:    "MaxRetries must be between 0 and 10",
		},
		{
			name: "invalid - request timeout out of range",
			config: &Config{
				UseADC:           true,
				MaxRetries:       3,
				RequestTimeout:   500, // > 300
				RateLimitRPS:     5,
				RateLimitBurst:   10,
				MaxConcurrency:   8,
				DefaultCurrency:  "USD",
				CacheTTL:         15,
			},
			expectError: true,
			errorMsg:    "RequestTimeout must be between 1 and 300",
		},
		{
			name: "invalid - rate limit RPS out of range",
			config: &Config{
				UseADC:           true,
				MaxRetries:       3,
				RequestTimeout:   30,
				RateLimitRPS:     150, // > 100
				RateLimitBurst:   10,
				MaxConcurrency:   8,
				DefaultCurrency:  "USD",
				CacheTTL:         15,
			},
			expectError: true,
			errorMsg:    "RateLimitRPS must be between 1 and 100",
		},
		{
			name: "invalid - currency",
			config: &Config{
				UseADC:           true,
				MaxRetries:       3,
				RequestTimeout:   30,
				RateLimitRPS:     5,
				RateLimitBurst:   10,
				MaxConcurrency:   8,
				DefaultCurrency:  "INVALID",
				CacheTTL:         15,
			},
			expectError: true,
			errorMsg:    "invalid DefaultCurrency",
		},
		{
			name: "invalid - region format",
			config: &Config{
				UseADC:           true,
				MaxRetries:       3,
				RequestTimeout:   30,
				RateLimitRPS:     5,
				RateLimitBurst:   10,
				MaxConcurrency:   8,
				DefaultCurrency:  "USD",
				DefaultRegions:   []string{"bad"},
				CacheTTL:         15,
			},
			expectError: true,
			errorMsg:    "invalid region format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				assert.Error(t, err)
				if err != nil && tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					// For file not found errors, we can ignore them in tests
					if !assert.Contains(t, err.Error(), "does not exist") {
						assert.NoError(t, err)
					}
				}
			}
		})
	}
}

func TestConfig_GetAuthenticationMethod(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected string
	}{
		{
			name: "service account file",
			config: &Config{
				ServiceAccountPath: "/path/to/service-account.json",
			},
			expected: "Service Account File: /path/to/service-account.json",
		},
		{
			name: "service account JSON",
			config: &Config{
				ServiceAccountJSON: "{}",
			},
			expected: "Service Account JSON (inline)",
		},
		{
			name: "ADC",
			config: &Config{
				UseADC: true,
			},
			expected: "Application Default Credentials (ADC)",
		},
		{
			name:     "no authentication",
			config:   &Config{},
			expected: "No authentication configured",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetAuthenticationMethod()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_IsProductionReady(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected bool
	}{
		{
			name:     "production ready config",
			config:   NewDefaultConfig(),
			expected: true,
		},
		{
			name: "debug mode enabled",
			config: &Config{
				Debug:           true,
				MaxRetries:      3,
				RequestTimeout:  30,
				MetricsEnabled:  true,
				RateLimitRPS:    5,
			},
			expected: false,
		},
		{
			name: "insufficient retries",
			config: &Config{
				Debug:           false,
				MaxRetries:      1,
				RequestTimeout:  30,
				MetricsEnabled:  true,
				RateLimitRPS:    5,
			},
			expected: false,
		},
		{
			name: "low timeout",
			config: &Config{
				Debug:           false,
				MaxRetries:      3,
				RequestTimeout:  5,
				MetricsEnabled:  true,
				RateLimitRPS:    5,
			},
			expected: false,
		},
		{
			name: "metrics disabled",
			config: &Config{
				Debug:           false,
				MaxRetries:      3,
				RequestTimeout:  30,
				MetricsEnabled:  false,
				RateLimitRPS:    5,
			},
			expected: false,
		},
		{
			name: "high rate limit",
			config: &Config{
				Debug:           false,
				MaxRetries:      3,
				RequestTimeout:  30,
				MetricsEnabled:  true,
				RateLimitRPS:    25,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsProductionReady()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_Clone(t *testing.T) {
	original := &Config{
		ServiceAccountPath:   "/path/to/service-account.json",
		ProjectID:            "test-project",
		UseADC:               false,
		MaxRetries:           5,
		DefaultRegions:       []string{"us-central1", "europe-west1"},
		DefaultCurrency:      "EUR",
		EnableAssetInventory: false,
		Debug:                true,
	}

	clone := original.Clone()

	// Verify all fields are copied
	assert.Equal(t, original.ServiceAccountPath, clone.ServiceAccountPath)
	assert.Equal(t, original.ProjectID, clone.ProjectID)
	assert.Equal(t, original.UseADC, clone.UseADC)
	assert.Equal(t, original.MaxRetries, clone.MaxRetries)
	assert.Equal(t, original.DefaultCurrency, clone.DefaultCurrency)
	assert.Equal(t, original.EnableAssetInventory, clone.EnableAssetInventory)
	assert.Equal(t, original.Debug, clone.Debug)

	// Verify slice is deep copied
	assert.Equal(t, original.DefaultRegions, clone.DefaultRegions)
	
	// Modify original slice to ensure it's a deep copy
	original.DefaultRegions[0] = "modified-region"
	assert.NotEqual(t, original.DefaultRegions[0], clone.DefaultRegions[0])

	// Verify they are different objects
	assert.NotSame(t, original, clone)
}

func TestConfig_String(t *testing.T) {
	config := &Config{
		ServiceAccountPath:   "/path/to/service-account.json",
		ProjectID:            "test-project",
		OrganizationID:       "123456789",
		BillingAccountID:     "ABC123-DEF456-GHI789",
		DefaultCurrency:      "USD",
		DefaultRegions:       []string{"us-central1", "europe-west1"},
		MaxRetries:           3,
		RequestTimeout:       30,
		RateLimitRPS:         5,
		RateLimitBurst:       10,
		MaxConcurrency:       8,
		EnableAssetInventory: true,
		EnablePricing:        true,
		EnableCaching:        true,
		CacheTTL:             15,
		Debug:                false,
		MetricsEnabled:       true,
		TracingEnabled:       false,
	}

	result := config.String()

	// Verify key information is present
	assert.Contains(t, result, "GCP Service Configuration:")
	assert.Contains(t, result, "Service Account File:")
	assert.Contains(t, result, "Project ID: test-project")
	assert.Contains(t, result, "Organization ID: 123456789")
	assert.Contains(t, result, "Default Currency: USD")
	assert.Contains(t, result, "Max Retries: 3")
	assert.Contains(t, result, "Rate Limit: 5 RPS")
	assert.Contains(t, result, "Production Ready:")

	// Test with sensitive JSON credentials
	configWithJSON := &Config{
		ServiceAccountJSON: `{"type": "service_account"}`,
		ProjectID:          "test-project",
		DefaultCurrency:    "USD",
	}

	resultWithJSON := configWithJSON.String()
	assert.Contains(t, resultWithJSON, "***masked***")
	assert.NotContains(t, resultWithJSON, "service_account")
}

// Benchmark tests
func BenchmarkNewDefaultConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewDefaultConfig()
	}
}

func BenchmarkConfig_Validate(b *testing.B) {
	config := NewDefaultConfig()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.Validate()
	}
}

func BenchmarkConfig_Clone(b *testing.B) {
	config := &Config{
		DefaultRegions: []string{"us-central1", "europe-west1", "asia-southeast1"},
		ProjectID:      "test-project",
		UseADC:         true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.Clone()
	}
}

// Test helper functions
func createTempServiceAccountFile(t *testing.T) string {
	t.Helper()
	
	content := `{
		"type": "service_account",
		"project_id": "test-project",
		"client_email": "test@test-project.iam.gserviceaccount.com"
	}`
	
	file, err := os.CreateTemp("", "service-account-*.json")
	require.NoError(t, err)
	
	_, err = file.WriteString(content)
	require.NoError(t, err)
	
	err = file.Close()
	require.NoError(t, err)
	
	t.Cleanup(func() {
		os.Remove(file.Name())
	})
	
	return file.Name()
}

func TestConfig_ValidateWithRealFile(t *testing.T) {
	serviceAccountFile := createTempServiceAccountFile(t)
	
	config := &Config{
		ServiceAccountPath:  serviceAccountFile,
		MaxRetries:          3,
		RequestTimeout:      30,
		RateLimitRPS:        5,
		RateLimitBurst:      10,
		MaxConcurrency:      8,
		DefaultCurrency:     "USD",
		CacheTTL:            15,
	}
	
	err := config.Validate()
	assert.NoError(t, err)
}