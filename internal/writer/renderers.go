package writer

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/cduggn/ccexplorer/internal/http"
	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/cduggn/ccexplorer/internal/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"log/slog"
)

// StdoutTableRenderer renders table data to stdout
type StdoutTableRenderer struct {
	variant string
}

// NewStdoutTableRenderer creates a new stdout table renderer
func NewStdoutTableRenderer(variant string) *StdoutTableRenderer {
	return &StdoutTableRenderer{variant: variant}
}

// Render implements the Renderer interface for stdout tables
func (r *StdoutTableRenderer) Render(data *TableOutput) error {
	switch r.variant {
	case "costAndUsage":
		return r.renderCostUsageTable(data)
	default:
		return fmt.Errorf("unknown table variant: %s", r.variant)
	}
}

func (r *StdoutTableRenderer) renderCostUsageTable(data *TableOutput) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 6, WidthMax: 8},
	})
	t.SetStyle(table.StyleColoredGreenWhiteOnBlack)
	t.SuppressEmptyColumns()

	// Convert headers to table.Row
	headerRow := make(table.Row, len(data.Headers))
	for i, h := range data.Headers {
		headerRow[i] = h
	}
	t.AppendHeader(headerRow)

	// Convert data rows to table.Row
	for _, row := range data.Rows {
		tableRow := make(table.Row, len(row))
		for i, cell := range row {
			tableRow[i] = cell
		}
		t.AppendRow(tableRow)
	}

	// Add divider and footer
	divider := make(table.Row, len(data.Headers))
	t.AppendRow(divider)
	
	footer := table.Row{"", "", "", "", "Cost", data.Total, "", "", "", ""}
	t.AppendFooter(footer)
	
	t.Render()
	return nil
}

// ForecastTableRenderer renders forecast table data to stdout  
type ForecastTableRenderer struct{}

// NewForecastTableRenderer creates a new forecast table renderer
func NewForecastTableRenderer() *ForecastTableRenderer {
	return &ForecastTableRenderer{}
}

// Render implements the Renderer interface for forecast tables
func (r *ForecastTableRenderer) Render(data *ForecastTableOutput) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredGreenWhiteOnBlack)

	// Convert headers to table.Row
	headerRow := make(table.Row, len(data.Headers))
	for i, h := range data.Headers {
		headerRow[i] = h
	}
	t.AppendHeader(headerRow)

	// Convert data rows to table.Row
	for _, row := range data.Rows {
		tableRow := make(table.Row, len(row))
		for i, cell := range row {
			tableRow[i] = cell
		}
		t.AppendRow(tableRow)
	}

	// Add footer with filter info and total
	footer := table.Row{
		"FilteredBy", data.FilterInfo, "", "", "",
		data.Total.Unit, data.Total.Amount,
	}
	t.AppendFooter(footer)
	
	t.Render()
	return nil
}

// CSVRenderer renders CSV data to files
type CSVRenderer struct{}

// NewCSVRenderer creates a new CSV renderer
func NewCSVRenderer() *CSVRenderer {
	return &CSVRenderer{}
}

// Render implements the Renderer interface for CSV files
func (r *CSVRenderer) Render(data *CSVOutput) error {
	filePath := utils.BuildOutputFilePath(OutputDir, data.Filename)
	file, err := os.Create(filePath)
	if err != nil {
		return types.Error{Msg: "Error creating CSV file: " + err.Error()}
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write(data.Headers); err != nil {
		return types.Error{Msg: "Error writing CSV header: " + err.Error()}
	}

	// Write data rows
	if err := writer.WriteAll(data.Rows); err != nil {
		return types.Error{Msg: "Error writing CSV data: " + err.Error()}
	}

	return nil
}

// ChartRenderer renders chart data to HTML files
type ChartRenderer struct{}

// NewChartRenderer creates a new chart renderer
func NewChartRenderer() *ChartRenderer {
	return &ChartRenderer{}
}

// Render implements the Renderer interface for charts
func (r *ChartRenderer) Render(data *ChartOutput) error {
	filePath := utils.BuildOutputFilePath(OutputDir, data.Filename)
	file, err := os.Create(filePath)
	if err != nil {
		return types.Error{Msg: "Failed creating chart HTML file: " + err.Error()}
	}
	defer file.Close()

	return data.Page.Render(io.MultiWriter(file))
}

// VectorRenderer renders vector data to vector databases
type VectorRenderer struct{}

// NewVectorRenderer creates a new vector renderer
func NewVectorRenderer() *VectorRenderer {
	return &VectorRenderer{}
}

// Render implements the Renderer interface for vector databases
func (r *VectorRenderer) Render(data *VectorOutput) error {
	// Create the vector store client - this would need to be passed in or configured
	client := NewVectorStoreClient(
		http.NewRequestBuilder(),
		data.IndexName,
		"", // API keys would need to be provided
		"",
	)

	// Create embeddings for the vector items
	vectors, err := client.CreateEmbeddings(data.Items)
	if err != nil {
		return types.Error{Msg: "Error creating embeddings: " + err.Error()}
	}

	// Update items with embedding vectors
	for index, vector := range vectors {
		if index < len(data.Items) {
			data.Items[index].EmbeddingVector = vector.Embedding
			data.Items[index].ID = utils.EncodeString(data.Items[index].EmbeddingText)
		}
	}

	// Convert to Pinecone format and upsert
	pineconeData := utils.ConvertToPineconeStruct(data.Items)
	resp, err := client.Upsert(context.Background(), pineconeData)
	if err != nil {
		return types.Error{Msg: "Error upserting to vector store: " + err.Error()}
	}

	slog.Info(fmt.Sprintf("Upserted %d items to vector store", resp.UpsertedCount))
	return nil
}