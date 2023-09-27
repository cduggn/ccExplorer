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
	data []PineconeStruct) (model.UpsertResponse, error) {

	batches := splitIntoBatches(data)

	var resp model.UpsertResponse

	for _, batch := range batches {
		message := UpsertVectorsRequest{
			Message: batch,
		}
		res, err := p.sendBatchRequest(ctx, message)
		resp.UpsertedCount += res.UpsertedCount
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

		dimensions := strings.Join(r.Dimensions, ",")
		tags := strings.Join(r.Tags, ",")

		item := model.VectorStoreItem{
			EmbeddingText: p.AddSemanticMeaning(d, dimensions, tags),
			Metadata: model.VectorStoreItemMetadata{
				StartDate:   d.Start,
				EndDate:     d.End,
				Granularity: r.Granularity,
				Dimensions:  dimensions,
				Tags:        tags,
			},
		}
		s = append(s, &item)
	}
	return s
}

func (p *ClientAPI) AddSemanticMeaning(s model.Service, dimensions, tags string) string {
	var r strings.Builder

	// append keys
	fmt.Fprintf(&r, "AWS Cost explorer cost and usage results grouped by dimensions and tags named %s %s ", dimensions, tags)
	fmt.Fprintf(&r, "and with the following key values %s,", strings.Join(s.Keys, ","))
	fmt.Fprintf(&r, " over the time period which starts and ends at %s,%s,%s,", s.Start, s.End, s.Name)

	// append metrics
	metrics := make([]string, len(s.Metrics))
	for i, v := range s.Metrics {
		encodedAmount := p.Encoder.CategorizeCostsWithBinning(v.NumericAmount)
		metrics[i] = fmt.Sprintf(
			"the metrics values include the cost category dataset name: %s, the cost associated with this grouped dimension and/or tag for this time period: %s, the currency unit used to represent the cost: %s, and an encoded value to normalize the cost into a binning category: %s",
			v.Name, v.Amount, v.Unit, encodedAmount)
	}

	r.WriteString(strings.Join(metrics, ","))
	return r.String()
}

//func (p *ClientAPI) AddSemanticMeaning(s model.Service, dimensions, tags string) string {
//	var r strings.Builder
//
//	// append keys
//	r.WriteString("AWS Cost explorer cost and usage results grouped by dimensions and tags named ")
//	r.WriteString(dimensions + " " + tags + " ")
//	r.WriteString(" and with the following key values ")
//	keys := strings.Join(s.Keys, ",")
//	r.WriteString(keys)
//	r.WriteString(",")
//	r.WriteString(" over the time period which starts and ends at  ")
//	// append start, end, and name
//	r.WriteString(s.Start)
//	r.WriteString(",")
//	r.WriteString(s.End)
//	r.WriteString(",")
//	r.WriteString(s.Name)
//	r.WriteString(",")
//
//	// append metrics
//	metrics := make([]string, len(s.Metrics))
//	for i, v := range s.Metrics {
//
//		encodedAmount := p.Encoder.CategorizeCostsWithBinning(v.NumericAmount)
//		metrics[i] = fmt.Sprintf(" the metrics values include the cost category dataset name: %s, the cost associated with this grouped dimension and/or tag for this time period: %s, the currency unit used to represent the cost: %s, and an encoded value to normalize the cost into a binning category: %s", v.Name, v.Amount, v.Unit,
//			encodedAmount)
//	}
//	r.WriteString(strings.Join(metrics, ","))
//	return r.String()
//}

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
