package openai

import (
	"bytes"
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"html/template"
	"os"
	"strings"
)

var (
	maxModelTokens = 4097
	rowHeader      = "Dimension/Tag,Dimension/Tag,Metric," +
		"Granularity,Start,End,USD Amount,Unit;"
	OutputDir  = "./output"
	aiFileName = "ccexplorer_ai.html"
)

func Writer(completions string) error {
	f, err := newFile(OutputDir, aiFileName)
	if err != nil {
		return Error{
			msg: "Failed creating AI HTML: " + err.Error(),
		}
	}
	defer f.Close()

	_, err = f.WriteString(completions)
	if err != nil {
		return Error{
			msg: "Failed writing to AI HTML: " + err.Error(),
		}
	}
	return nil
}

func SummarizeCompletionsAPI(apiKey string, promptData string) (gogpt.
	CompletionResponse,
	error) {

	fmt.Println("Generating costAndUsage report with gpt3...")

	c := gogpt.NewClient(apiKey)
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: maxModelTokens - 840, //todo - make this value dynamic
		Prompt:    promptData,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return gogpt.CompletionResponse{}, Error{
			msg: "GPT-3 failed to generate report: " + err.Error(),
		}
	}

	return resp, nil
}

func Summarize(apiKey string, userMessage string) (gogpt.ChatCompletionResponse,
	error) {

	fmt.Println("Generating costAndUsage report with gpt3...")

	c := gogpt.NewClient(apiKey)
	ctx := context.Background()

	message := CreateChatCompletionMessage(userMessage)

	req := gogpt.ChatCompletionRequest{
		Model:     gogpt.GPT3Dot5Turbo,
		Messages:  message,
		MaxTokens: maxModelTokens - 840,
		Stream:    false,
	}

	resp, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		return gogpt.ChatCompletionResponse{}, Error{
			msg: "GPT-3 failed to generate report: " + err.Error(),
		}
	}

	return resp, nil
}

func CreateChatCompletionMessage(userMessage string) []gogpt.
	ChatCompletionMessage {
	return []gogpt.ChatCompletionMessage{

		{
			Role: "system",
			Content: "You are an AWS cost optimization expert that recommends" +
				" specific cost saving measure in HTML format for each row" +
				" of CSV" +
				" data.",
		},
		{
			Role:    "user",
			Content: userMessage,
		},
	}
}

func ConvertToCommaDelimitedString(rows [][]string) string {
	var buf bytes.Buffer

	buf.WriteString(rowHeader)

	for i, row := range rows {
		for j, col := range row {
			buf.WriteString(col)
			if j < len(row)-1 {
				buf.WriteByte(',')
			}
		}
		if i < len(rows)-1 {
			buf.WriteByte(';')
		}
	}
	cvsString := buf.String()
	return cvsString
}

func BuildPromptText(rows [][]string) string {
	var builder strings.Builder
	builder.WriteString("<!DOCTYPE html> Generate table that looks like this: ")

	builder.WriteString("# Table headers [Dimension/Tag, " +
		"Dimension/Tag, " +
		"Metric, Granularity, Start, End, USD Amount, Unit, " +
		"Percentage of Total\t] ")

	builder.WriteString(" Use the following csv data to display the top 10 rows: ")

	costAndUsageData := ConvertToCommaDelimitedString(rows)
	builder.WriteString(costAndUsageData)

	builder.WriteString(" Display a title named Cost and Usage Report above" +
		" the table centered. " +
		" Include a subtitle with the date range in smaller font.")

	builder.WriteString(" Use HTML, CSS and modern libraries to create a simple, " +
		"minimalistic design with alternating row colors, " +
		"and hover effect. Left align table row text. " +
		" Use a simple grey theme for the table. " +
		"Text font should be no more than size 18. ")

	builder.WriteString(" Add a column to number each row.")

	builder.WriteString(" Add a column to display the percentage of the total " +
		"cost for each row. ")

	builder.WriteString(" Detail a cost optimization" +
		" recommendation sentence for each table row based on the costs" +
		" shown in a new column. ")
	builder.WriteString(" Add the AWS well architected framework principle" +
		" which applies to the cost recommendation. ")

	return builder.String()
}

func CreateTrainingData(t *template.Template, data []TrainingData) (string,
	error) {
	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func BuildTrainingDataRow(rows [][]string) []TrainingData {
	return []TrainingData{
		{
			Dimension:   rows[0][0],
			Tag:         rows[0][1],
			Metric:      rows[0][2],
			Granularity: rows[0][3],
			Start:       rows[0][4],
			End:         rows[0][5],
			USDAmount:   rows[0][6],
			Unit:        rows[0][7],
		},
	}
}

// todo remove duplication
func newFile(dir string, file string) (*os.File, error) {
	filePath := buildOutputFilePath(dir, file)
	return os.Create(filePath)
}

// todo remove duplication
func buildOutputFilePath(dir string, fileName string) string {
	return dir + "/" + fileName
}
