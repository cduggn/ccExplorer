package pinecone

import "net/http"

type ClientConfig struct {
	apiKey     string
	BaseURL    string
	HTTPClient *http.Client
}

func DefaultConfig(indexUrl, apiKey string) ClientConfig {
	return ClientConfig{
		BaseURL:    indexUrl,
		HTTPClient: &http.Client{},
		apiKey:     apiKey,
	}
}
