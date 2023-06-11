package requestbuilder

import (
	"context"
	"io"
	"net/http"
)

type Builder interface {
	Build(ctx context.Context, method, url string,
		request io.Reader) (*http.Request,
		error)
}

type HttpRequestBuilder struct {
	marshaller marshaller
}

func NewRequestBuilder() *HttpRequestBuilder {
	return &HttpRequestBuilder{
		marshaller: &jsonMarshaller{},
	}
}

func (b *HttpRequestBuilder) Build(ctx context.Context, method, url string,
	payload io.Reader) (*http.Request, error) {
	if payload == nil {
		return http.NewRequestWithContext(ctx, method, url, nil)
	}

	return http.NewRequest(
		method,
		url,
		payload,
	)

}
