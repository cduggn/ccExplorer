package usecases

import (
	"fmt"
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	"github.com/cduggn/ccexplorer/internal/core/requestbuilder"
	"github.com/cduggn/ccexplorer/internal/core/vectorstore/pinecone"
)

func WriteToVectorStore(r model.CostAndUsageOutputType, apikey, indexUrl,
	openAIAPIKey string) error {

	client := pinecone.NewVectorStoreClient(requestbuilder.NewRequestBuilder(), indexUrl, apikey, openAIAPIKey)
	items := client.ConvertToVectorStoreItem(r)

	fmt.Print(items)
	return nil
}
