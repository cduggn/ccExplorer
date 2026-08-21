package writer

import (
	"encoding/json"
	"testing"

	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func vectorInput() types.CostAndUsageOutputType {
	return types.CostAndUsageOutputType{
		Granularity: "MONTHLY",
		Dimensions:  []string{"SERVICE"},
		Start:       "2024-01-01",
		End:         "2024-02-01",
		Services: map[int]types.Service{
			0: {
				Name:  "Amazon EC2",
				Start: "2024-01-01",
				End:   "2024-02-01",
				Keys:  []string{"Amazon EC2"},
				Metrics: []types.Metrics{{
					Name:          "UnblendedCost",
					Amount:        "12.34",
					NumericAmount: 12.34,
					Unit:          "USD",
				}},
			},
		},
		OpenAIAPIKey:   "openai-key",
		PineconeAPIKey: "pinecone-key",
		PineconeIndex:  "https://idx.svc.pinecone.io",
	}
}

// Regression: Transform used to call NewVectorStoreClient with the three
// credentials rotated, and then hand the renderer only the index name. Each
// value must land in its own field.
func TestVectorTransformerCarriesEachCredentialSeparately(t *testing.T) {
	out, err := NewCostUsageToVectorTransformer().Transform(vectorInput())
	require.NoError(t, err)

	assert.Equal(t, "https://idx.svc.pinecone.io", out.IndexURL)
	assert.Equal(t, "pinecone-key", out.PineconeAPIKey)
	assert.Equal(t, "openai-key", out.OpenAIAPIKey)
	assert.NotEmpty(t, out.Items, "the embedding payload should be built")
}

// The renderer used to receive "" for both keys, so no credential could ever
// reach Pinecone or OpenAI regardless of what the user exported.
func TestVectorOutputNeverHoldsBlankCredentialsForAConfiguredRun(t *testing.T) {
	out, err := NewCostUsageToVectorTransformer().Transform(vectorInput())
	require.NoError(t, err)

	assert.NotEmpty(t, out.IndexURL)
	assert.NotEmpty(t, out.PineconeAPIKey)
	assert.NotEmpty(t, out.OpenAIAPIKey)
}

// CostAndUsageOutputType is marshalled verbatim into MCP tool results, which
// are handed to a model. Credentials must not appear there.
func TestCostAndUsageOutputDoesNotSerialiseCredentials(t *testing.T) {
	b, err := json.Marshal(vectorInput())
	require.NoError(t, err)

	body := string(b)
	assert.NotContains(t, body, "openai-key")
	assert.NotContains(t, body, "pinecone-key")
	assert.NotContains(t, body, "OpenAIAPIKey")
	assert.NotContains(t, body, "PineconeAPIKey")
	assert.Contains(t, body, "Amazon EC2", "real report data must survive")
}
