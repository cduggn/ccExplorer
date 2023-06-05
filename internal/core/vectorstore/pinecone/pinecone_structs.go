package pinecone

import "github.com/cduggn/ccexplorer/internal/core/requestbuilder"

type PineconeClient struct {
	RequestBuilder requestbuilder.Builder
	Config         ClientConfig
}

type PineconeStruct struct {
	ID       string    `json:"id"`
	Values   []float32 `json:"values"`
	Metadata Metadata  `json:"metadata"`
}

type Metadata struct {
	PageContent string `json:"page_content"`
	Source      string `json:"source"`
}

type UpsertVectorsRequest struct {
	Message []PineconeStruct `json:"vectors"`
}
