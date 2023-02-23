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
	trainingTemplate = `
        <table>
            <thead>
                <tr>
                    <th>Dimension/Tag</th>
                    <th>Dimension/Tag</th>
                    <th>Start</th>
                    <th>End</th>
                    <th>USD Amount</th>
                </tr>
            </thead>
            <tbody>
                {{range .}}
                <tr>
                    <td>{{.Dimension}}</td>
                    <td>{{.Tag}}</td>
                    <td>{{.Start}}</td>
                    <td>{{.End}}</td>
                    <td>{{.USDAmount}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    `
	rowHeader = "Dimension/Tag,Dimension/Tag,Metric," +
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
	builder.WriteString("Generate a stylish html table that look like this: ")

	trainingData := BuildCostAndUsagePromptText(rows)
	builder.WriteString(trainingData)

	builder.WriteString(" Use the following csv data to display the top 10 rows: ")
	costAndUsageData := ConvertToCommaDelimitedString(rows[:15])
	builder.WriteString(costAndUsageData)

	builder.WriteString(" Use HTML, CSS and modern libraries to create a simple, " +
		"clean , minimalistic design. " +
		"Add a column to number each row. " +
		"Center the table on the page. " +
		"Use simple grey theme for table and font size 16 for columns and" +
		" rows. " +
		"Set a table width no more than 70% of page width. " +
		"Table style should include alternating row colors, a hover effect. " +
		"Left align table text. ")

	builder.WriteString(" Set table header title as Cost and Usage Report. " +
		" Display the date range in a smaller font below the" +
		" title. ")

	builder.WriteString(" Generate cost optimization" +
		" recommendations for each table row entry which are specific to that" +
		" row. This should be in text form now a URL. Display recommendations" +
		" in a new row. " +
		"Recommendations should be specific for each row " + ": ")

	return builder.String()
}

func BuildCostAndUsagePromptText(rows [][]string) string {
	t := CreateTrainingTemplate()
	s, err := CreateTrainingData(t, BuildTrainingDataRow(rows))
	if err != nil {
		fmt.Println("Error populating template: ", err)
	}
	return s
}

func CreateTrainingTemplate() *template.Template {
	t := template.Must(template.New("table").Parse(trainingTemplate))
	return t
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

func Summarize(apiKey string, promptData string) (gogpt.
	CompletionResponse,
	error) {

	fmt.Println("Generating costAndUsage report with gpt3...")

	c := gogpt.NewClient(apiKey)
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: 2400,
		Prompt:    promptData,
		//Temperature:
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return gogpt.CompletionResponse{}, err
	}

	return resp, nil
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
