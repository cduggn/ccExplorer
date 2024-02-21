package pinecone

import (
	"github.com/cduggn/ccexplorer/internal/codec"
	"github.com/cduggn/ccexplorer/internal/http"
	"github.com/cduggn/ccexplorer/internal/openai"
)

type ClientAPI struct {
	RequestBuilder http.Builder
	Config         ClientConfig
	LLMClient      openai.OpenAI
	Encoder        codec.Encode
}

type PineconeStruct struct {
	ID       string    `json:"id"`
	Values   []float32 `json:"values"`
	Metadata Metadata  `json:"metadata"`
}

type Metadata struct {
	PageContent string `json:"page_content"`
	Source      string `json:"source"`
	Dimensions  string `json:"dimensions"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Cost        string `json:"cost"`
}

type UpsertVectorsRequest struct {
	Message []PineconeStruct `json:"vectors"`
}
