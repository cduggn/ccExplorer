package printer

import (
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"os"
)

var (
	aiPrompt   = "Create a styled HTML page summary of the following data. Use a table with alternate color for each row. Dispaly the headings which is the first row for csv data. Make three cost reducing recommendations"
	aiFileName = "ccexplorer_ai.html"
)

func AIWriter(f *os.File, completions string) error {
	_, err := f.WriteString(completions)
	if err != nil {
		return PrinterError{
			msg: "Failed writing to AI HTML: " + err.Error(),
		}
	}
	return nil
}

func SummarizeWIthAI(apiKey string, data string) (gogpt.CompletionResponse,
	error) {

	fmt.Println("Generating costAndUsage report with gpt3...")

	c := gogpt.NewClient(apiKey)
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: 400,
		Prompt:    BuildPromptText(data),
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return gogpt.CompletionResponse{}, err
	}

	return resp, nil
}

func BuildPromptText(data string) string {
	return aiPrompt + data
}
