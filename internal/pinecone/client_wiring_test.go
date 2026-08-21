package pinecone

import (
	"testing"

	http2 "github.com/cduggn/ccexplorer/internal/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NewVectorStoreClient takes (builder, indexURL, pineconeAPIKey, openAIAPIKey)
// while its writer-package namesake takes (builder, openAIAPIKey, indexURL,
// pineconeAPIKey). Both call sites used to get that mapping wrong, so pin the
// wiring down here.
func TestNewVectorStoreClientWiresIndexAndKeyToTheRightPlaces(t *testing.T) {
	c := NewVectorStoreClient(http2.NewRequestBuilder(),
		"https://idx.svc.pinecone.io", "pinecone-key", "openai-key")

	require.NotNil(t, c)
	assert.Equal(t, "https://idx.svc.pinecone.io", c.Config.BaseURL,
		"the index URL must become the request base URL")
	assert.Equal(t, "pinecone-key", c.Config.apiKey,
		"the Pinecone key must become the Api-Key header value")
	assert.NotNil(t, c.LLMClient)
	assert.NotNil(t, c.Encoder)
}
