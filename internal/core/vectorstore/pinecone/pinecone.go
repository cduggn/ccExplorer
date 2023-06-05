package pinecone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	"github.com/cduggn/ccexplorer/internal/core/requestbuilder"
	"github.com/cduggn/ccexplorer/internal/core/service/openai"
	"io"
	"net/http"
)

func NewVectorStoreClient(builder requestbuilder.Builder,
	apiKey string, indexURL string, openAIAPIKey string) PineconeDB {

	return &ClientAPI{
		RequestBuilder: builder,
		Config:         DefaultConfig(indexURL, apiKey),
		LLMClient:      openai.NewClient(openAIAPIKey),
	}
}

func (p *ClientAPI) Upsert(ctx context.Context,
	data []PineconeStruct) error {

	batches := splitIntoBatches(data)

	for _, batch := range batches {
		message := UpsertVectorsRequest{
			Message: batch,
		}
		err := p.sendBatchRequest(ctx, message)
		if err != nil {
			return err
		}
	}
	return nil
}

func splitIntoBatches(data []PineconeStruct) [][]PineconeStruct {
	var batches [][]PineconeStruct
	batchSize := 25
	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}
		batches = append(batches, data[i:end])
	}
	return batches
}

func (p *ClientAPI) sendBatchRequest(ctx context.Context,
	message UpsertVectorsRequest) error {

	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := p.RequestBuilder.Build(ctx, http.MethodPost,
		p.Config.BaseURL+"/vectors/upsert", bytes.NewReader(payload))
	if err != nil {
		return err
	}

	err = p.sendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}

func (p *ClientAPI) ConvertToPineconeStruct(
	data []*model.CostAndUsage) []PineconeStruct {

	var pineconeSlice []PineconeStruct

	//for index, d := range data {
	//	pinecone := PineconeStruct{
	//		ID:       strconv.Itoa(index),
	//		Values:   d.Embeddings,
	//		Metadata: Metadata{PageContent: d.Combined, Source: "AWS"},
	//	}
	//	pineconeSlice = append(pineconeSlice, pinecone)
	//}
	return pineconeSlice
}

func (p *ClientAPI) sendRequest(req *http.Request, v any) error {
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Api-Key", p.Config.apiKey)

	res, err := p.Config.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	return decodeResponse(res.Body, v)
}

func decodeResponse(body io.Reader, v any) error {
	if v == nil {
		return nil
	}
	if result, ok := v.(*string); ok {
		return decodeString(body, result)
	}
	return json.NewDecoder(body).Decode(v)
}

func decodeString(body io.Reader, output *string) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	*output = string(b)
	return nil
}
