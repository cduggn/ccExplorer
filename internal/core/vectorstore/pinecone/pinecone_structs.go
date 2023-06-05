package pinecone

import (
	"github.com/cduggn/ccexplorer/internal/core/requestbuilder"
	"github.com/cduggn/ccexplorer/internal/core/service/openai"
)

type ClientAPI struct {
	RequestBuilder requestbuilder.Builder
	Config         ClientConfig
	LLMClient      openai.OpenAI
}

type PineconeStruct struct {
	ID       string    `json:"id"`
	Values   []float32 `json:"values"`
	Metadata Metadata  `json:"metadata"`
}

type Metadata struct {
	PageContent string `json:"page_content"`
	Source      string `json:"source"`
	Service     string `json:"service"`
	Year        string `json:"year"`
}

type UpsertVectorsRequest struct {
	Message []PineconeStruct `json:"vectors"`
}
