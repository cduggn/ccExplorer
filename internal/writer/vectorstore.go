package writer

import (
	"context"
	"github.com/cduggn/ccexplorer/internal/http"
	pinecone2 "github.com/cduggn/ccexplorer/internal/pinecone"
	"github.com/cduggn/ccexplorer/internal/types"
	gogpt "github.com/sashabaranov/go-openai"
)

type VectorStore interface {
	CreateVectorStoreInput(r types.CostAndUsageOutputType) ([]*types.VectorStoreItem, error)
	CreateEmbeddings(items []*types.VectorStoreItem) ([]gogpt.Embedding, error)
	Upsert(context context.Context, data []pinecone2.PineconeStruct) (resp types.UpsertResponse, err error)
}

type VectorStoreClient struct {
	apikey         string
	indexUrl       string
	openAIAPIKey   string
	requestbuilder http.Builder

	client *pinecone2.ClientAPI
}

func NewVectorStoreClient(builder http.Builder, openAIAPIKey,
	indexUrl,
	pineconeAPIKey string) VectorStore {
	return &VectorStoreClient{
		apikey:         openAIAPIKey,
		indexUrl:       indexUrl,
		openAIAPIKey:   pineconeAPIKey,
		requestbuilder: builder,
		client: pinecone2.NewVectorStoreClient(builder, indexUrl,
			pineconeAPIKey, openAIAPIKey),
	}
}

func (v *VectorStoreClient) CreateVectorStoreInput(r types.CostAndUsageOutputType) ([]*types.VectorStoreItem, error) {
	items := v.client.ConvertToVectorStoreItem(r)

	return items, nil
}

func (v *VectorStoreClient) CreateEmbeddings(items []*types.VectorStoreItem) (
	[]gogpt.Embedding,
	error) {

	batch := make([]string, len(items))
	for index, item := range items {
		batch[index] = item.EmbeddingText
	}

	vectors, err := v.client.LLMClient.GenerateEmbeddings(batch)
	if err != nil {
		return nil, err
	}

	return vectors, nil
}

func (v *VectorStoreClient) Upsert(context context.Context,
	items []pinecone2.PineconeStruct) (resp types.UpsertResponse, err error) {

	return v.client.Upsert(context, items)
}
