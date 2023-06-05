package writers

import (
	"github.com/cduggn/ccexplorer/internal/core/service/openai"
)

func NewLLMClient(apiKey string) openai.OpenAI {

	return openai.NewClient(apiKey)
}
