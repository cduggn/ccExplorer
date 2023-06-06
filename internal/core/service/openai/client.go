package openai

import (
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-openai"
)

type OpenAI interface {
	GenerateEmbeddings(items []string) ([]gogpt.Embedding,
		error)
}

type Client struct {
	client *gogpt.Client
}

func NewClient(apiKey string) OpenAI {
	return &Client{
		client: gogpt.NewClient(apiKey),
	}
}

func (o *Client) GenerateEmbeddings(items []string) ([]gogpt.
	Embedding,
	error) {

	req := gogpt.EmbeddingRequest{
		Input: items,
		Model: gogpt.AdaEmbeddingV2,
	}

	resp, err := o.client.CreateEmbeddings(context.Background(),
		req)

	if err != nil {
		fmt.Printf("Embedding error: %v\n", err)
		return nil, err
	}

	return resp.Data, nil
}
