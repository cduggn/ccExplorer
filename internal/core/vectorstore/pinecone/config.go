package pinecone

import "net/http"

type ClientConfig struct {
	apiKey     string
	BaseURL    string
	HTTPClient *http.Client
}

func DefaultConfig(apiKey, baseUrl string) ClientConfig {
	return ClientConfig{
		BaseURL:    baseUrl,
		HTTPClient: &http.Client{},
		apiKey:     apiKey,
	}
}
