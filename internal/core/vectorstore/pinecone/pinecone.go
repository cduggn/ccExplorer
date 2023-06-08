package pinecone

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	"github.com/cduggn/ccexplorer/internal/core/encoder"
	"github.com/cduggn/ccexplorer/internal/core/requestbuilder"
	"github.com/cduggn/ccexplorer/internal/core/service/openai"
	"io"
	"net/http"
	"strings"
)

func NewVectorStoreClient(builder requestbuilder.Builder,
	indexURL, pineconeAPIKey, openAIAPIKey string) *ClientAPI {

	return &ClientAPI{
		RequestBuilder: builder,
		Config:         DefaultConfig(indexURL, pineconeAPIKey),
		LLMClient:      openai.NewClient(openAIAPIKey),
		Encoder:        encoder.NewEncoder(),
	}
}

func (p *ClientAPI) Upsert(ctx context.Context,
	data []PineconeStruct) (resp model.UpsertResponse, err error) {

	batches := splitIntoBatches(data)

	for _, batch := range batches {
		message := UpsertVectorsRequest{
			Message: batch,
		}
		resp, err = p.sendBatchRequest(ctx, message)
		if err != nil {
			return model.UpsertResponse{}, err
		}
	}
	return resp, nil
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
	message UpsertVectorsRequest) (resp model.UpsertResponse, err error) {

	payload, err := json.Marshal(message)
	if err != nil {
		return model.UpsertResponse{}, err
	}

	req, err := p.RequestBuilder.Build(ctx, http.MethodPost,
		p.Config.BaseURL+"/vectors/upsert", bytes.NewReader(payload))
	if err != nil {
		return model.UpsertResponse{}, err
	}

	err = p.sendRequest(req, &resp)
	if err != nil {
		return model.UpsertResponse{}, err
	}
	return
}

func (p *ClientAPI) ConvertToVectorStoreItem(r model.
	CostAndUsageOutputType) []*model.
	VectorStoreItem {
	var s []*model.VectorStoreItem
	for _, d := range r.Services {

		item := model.VectorStoreItem{
			EmbeddingText: p.serviceToString(d),
			Metadata: model.VectorStoreItemMetadata{
				StartDate:   r.Start,
				Granularity: r.Granularity,
				Dimensions:  strings.Join(r.Dimensions, ","),
				Tags:        strings.Join(r.Tags, ","),
			},
		}
		s = append(s, &item)
	}
	return s
}

func (p *ClientAPI) serviceToString(s model.Service) string {
	var r strings.Builder

	// append keys
	keys := strings.Join(s.Keys, ",")
	r.WriteString(keys)
	r.WriteString(",")

	// append start, end, and name
	r.WriteString(s.Start)
	r.WriteString(",")
	r.WriteString(s.End)
	r.WriteString(",")
	r.WriteString(s.Name)
	r.WriteString(",")

	// append metrics
	metrics := make([]string, len(s.Metrics))
	for i, v := range s.Metrics {

		encodedAmount := p.Encoder.CategorizeCostsWithBinning(v.NumericAmount)
		metrics[i] = fmt.Sprintf("%s,%s,%s", v.Amount, v.Unit, encodedAmount)
	}
	r.WriteString(strings.Join(metrics, ","))
	return r.String()
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
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.
		StatusBadRequest {
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
