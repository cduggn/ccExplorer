package usecases

import (
	"context"
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	"github.com/cduggn/ccexplorer/internal/core/requestbuilder"
	pinecone2 "github.com/cduggn/ccexplorer/internal/pinecone"
	gogpt "github.com/sashabaranov/go-openai"
)

type VectorStore interface {
	CreateVectorStoreInput(r model.CostAndUsageOutputType) ([]*model.
		VectorStoreItem, error)
	CreateEmbeddings(items []*model.VectorStoreItem) ([]gogpt.Embedding, error)
	Upsert(context context.Context, data []pinecone2.PineconeStruct) (resp model.UpsertResponse, err error)
}

type VectorStoreClient struct {
	apikey         string
	indexUrl       string
	openAIAPIKey   string
	requestbuilder requestbuilder.Builder

	client *pinecone2.ClientAPI
}

func NewVectorStoreClient(builder requestbuilder.Builder, openAIAPIKey,
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

func (v *VectorStoreClient) CreateVectorStoreInput(r model.
	CostAndUsageOutputType) ([]*model.VectorStoreItem, error) {
	items := v.client.ConvertToVectorStoreItem(r)

	return items, nil
}

func (v *VectorStoreClient) CreateEmbeddings(items []*model.VectorStoreItem) (
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
	items []pinecone2.PineconeStruct) (resp model.UpsertResponse, err error) {

	return v.client.Upsert(context, items)
}
