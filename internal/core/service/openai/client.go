package openai

import (
	"fmt"
	gogpt "github.com/sashabaranov/go-openai"
)

var (
	OutputDir = "./output"
	//maxModelTokens = 4097
)

type OpenAI interface {
	GenerateEmbeddings(text string) (string, error)
}

type OpenAIClient struct {
	client *gogpt.Client
}

func NewClient(apiKey string) OpenAI {
	return &OpenAIClient{
		client: gogpt.NewClient(apiKey),
	}
}

func (o *OpenAIClient) GenerateEmbeddings(text string) (string, error) {
	fmt.Print("Generate Embeddings....")
	return "nil", nil
}
