package usecases

import (
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	"github.com/cduggn/ccexplorer/internal/core/requestbuilder"
	"github.com/cduggn/ccexplorer/internal/core/vectorstore/pinecone"
	gogpt "github.com/sashabaranov/go-openai"
)

type VectorStore interface {
	CreateVectorStoreInput(r model.CostAndUsageOutputType) ([]*model.
		VectorStoreItem, error)
	CreateEmbeddings(items []*model.VectorStoreItem) ([]gogpt.Embedding, error)
}

type VectorStoreClient struct {
	apikey         string
	indexUrl       string
	openAIAPIKey   string
	requestbuilder requestbuilder.Builder

	client *pinecone.ClientAPI
}

func NewVectorStoreClient(builder requestbuilder.Builder, apikey, indexUrl,
	openAIAPIKey string) VectorStore {
	return &VectorStoreClient{
		apikey:         apikey,
		indexUrl:       indexUrl,
		openAIAPIKey:   openAIAPIKey,
		requestbuilder: builder,
		client:         pinecone.NewVectorStoreClient(builder, indexUrl, apikey, openAIAPIKey),
	}
}

func (v *VectorStoreClient) CreateVectorStoreInput(r model.
	CostAndUsageOutputType) ([]*model.VectorStoreItem, error) {
	items := v.client.ConvertToVectorStoreItem(r)

	return items, nil
}

func (v *VectorStoreClient) CreateEmbeddings(items []*model.VectorStoreItem) (
	[]gogpt.Embedding,
	error) {

	var batch []string
	for _, item := range items {
		batch = append(batch, item.EmbeddingText)
	}

	vectors, err := v.client.LLMClient.GenerateEmbeddings(batch)
	if err != nil {
		return nil, err
	}

	return vectors, nil
}
