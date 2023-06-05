package writers

import (
	"github.com/cduggn/ccexplorer/internal/core/requestbuilder"
	"github.com/cduggn/ccexplorer/internal/core/vectorstore/pinecone"
)

func NewVectorStoreClient(builder requestbuilder.Builder,
	apiKey string, indexURL string) pinecone.PineconeDB {

	return &pinecone.PineconeClient{
		RequestBuilder: builder,
		Config:         pinecone.DefaultConfig(indexURL, apiKey),
	}
}
