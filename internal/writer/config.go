package writer

import (
	"github.com/cduggn/ccexplorer/internal/http"
)

// Config holds configuration for writer components
type Config struct {
	// OutputDir is the directory where files are written
	OutputDir string

	// Vector database configuration
	OpenAIAPIKey   string
	PineconeAPIKey string
	PineconeIndex  string

	// HTTP client for external services
	HTTPClient *http.HttpRequestBuilder

	// Sorting configuration
	SortBy string
}

// Option is a functional option for configuring writers
type Option func(*Config)

// WithOutputDir sets the output directory for file-based writers
func WithOutputDir(dir string) Option {
	return func(c *Config) {
		c.OutputDir = dir
	}
}

// WithVectorDBConfig sets the vector database configuration
func WithVectorDBConfig(openAIKey, pineconeKey, pineconeIndex string) Option {
	return func(c *Config) {
		c.OpenAIAPIKey = openAIKey
		c.PineconeAPIKey = pineconeKey
		c.PineconeIndex = pineconeIndex
	}
}

// WithHTTPClient sets the HTTP client for external services
func WithHTTPClient(client *http.HttpRequestBuilder) Option {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

// WithSortBy sets the sorting preference
func WithSortBy(sortBy string) Option {
	return func(c *Config) {
		c.SortBy = sortBy
	}
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		OutputDir:  "./writer",
		SortBy:     "amount",
		HTTPClient: http.NewRequestBuilder(),
	}
}

// NewConfig creates a Config with the given options
func NewConfig(opts ...Option) *Config {
	config := DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}
	return config
}
