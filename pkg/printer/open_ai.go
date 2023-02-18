package printer

import (
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"os"
)

var (
	aiPrompt = "Create HTML report using the CSV provided which represents" +
		" AWS" +
		" Cost" +
		" Explorer data. Display date-range, " +
		"top 20 costs in table descending with borders," +
		"display a 5 bullet point summary of costs:  "
	aiFileName = "ccexplorer_ai.html"
	//standByText
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

	fmt.Println()

	c := gogpt.NewClient(apiKey)
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: 1500,
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
