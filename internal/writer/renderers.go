package writer

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"

	"github.com/cduggn/ccexplorer/internal/http"
	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/cduggn/ccexplorer/internal/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// StdoutTableRenderer renders table data to stdout
type StdoutTableRenderer struct{}

// NewStdoutTableRenderer creates a new stdout table renderer
func NewStdoutTableRenderer() *StdoutTableRenderer {
	return &StdoutTableRenderer{}
}

// Render implements the Renderer interface for stdout tables
func (r *StdoutTableRenderer) Render(data *TableOutput) error {
	// Print summary header with consolidated metadata
	fmt.Println()
	fmt.Printf("  %s\n", text.FgHiCyan.Sprint("AWS Cost Explorer Report"))
	fmt.Printf("  %s %s  |  %s %s  |  %s %d services\n",
		text.FgHiWhite.Sprint("Period:"), data.DateRange,
		text.FgHiWhite.Sprint("Granularity:"), data.Granularity,
		text.FgHiWhite.Sprint("Results:"), data.RowCount)
	fmt.Printf("  %s %s  |  %s %s  |  %s %s %s\n",
		text.FgHiWhite.Sprint("Metric:"), data.Metric,
		text.FgHiWhite.Sprint("Unit:"), data.Unit,
		text.FgHiWhite.Sprint("Sorted by:"), data.SortedBy, text.FgHiWhite.Sprint("▼"))
	fmt.Println()
	// Color legend
	fmt.Printf("  %s  %s >= $5  %s >= $1  %s >= $0.01  %s < $0.01\n",
		text.FgHiWhite.Sprint("Cost:"),
		text.FgHiRed.Sprint("■"),
		text.FgHiYellow.Sprint("■"),
		text.FgHiGreen.Sprint("■"),
		text.FgHiBlack.Sprint("■"))
	fmt.Println()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Column configurations based on whether tag column exists
	if data.HasTagColumn {
		// 7 columns: #, Service, Tag/Dimension, Cost, %, Start, End
		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, Align: text.AlignRight, AlignHeader: text.AlignCenter},                               // #
			{Number: 2, WidthMax: 38, AlignHeader: text.AlignCenter},                                         // Service
			{Number: 3, WidthMax: 20, AlignHeader: text.AlignCenter},                                         // Tag/Dimension
			{Number: 4, Align: text.AlignRight, AlignHeader: text.AlignCenter, AlignFooter: text.AlignRight}, // Cost
			{Number: 5, Align: text.AlignRight, AlignHeader: text.AlignCenter},                               // %
			{Number: 6, Align: text.AlignCenter, AlignHeader: text.AlignCenter},                              // Start
			{Number: 7, Align: text.AlignCenter, AlignHeader: text.AlignCenter},                              // End
		})
	} else {
		// 6 columns: #, Service, Cost, %, Start, End
		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, Align: text.AlignRight, AlignHeader: text.AlignCenter},                               // #
			{Number: 2, WidthMax: 42, AlignHeader: text.AlignCenter},                                         // Service
			{Number: 3, Align: text.AlignRight, AlignHeader: text.AlignCenter, AlignFooter: text.AlignRight}, // Cost
			{Number: 4, Align: text.AlignRight, AlignHeader: text.AlignCenter},                               // %
			{Number: 5, Align: text.AlignCenter, AlignHeader: text.AlignCenter},                              // Start
			{Number: 6, Align: text.AlignCenter, AlignHeader: text.AlignCenter},                              // End
		})
	}

	// Use a cleaner style with rounded borders
	t.SetStyle(table.StyleRounded)
	t.Style().Color.Header = text.Colors{text.FgHiWhite, text.Bold}
	t.Style().Color.Row = text.Colors{text.FgWhite}
	t.Style().Color.RowAlternate = text.Colors{text.FgHiBlack}
	t.Style().Color.Footer = text.Colors{text.FgHiGreen, text.Bold}
	t.Style().Options.SeparateRows = false

	// Convert headers to table.Row
	headerRow := make(table.Row, len(data.Headers))
	for i, h := range data.Headers {
		headerRow[i] = h
	}
	t.AppendHeader(headerRow)

	// Determine cost column index based on layout
	costColIdx := 2 // Without tag column: #(0), Service(1), Cost(2)
	if data.HasTagColumn {
		costColIdx = 3 // With tag column: #(0), Service(1), Tag(2), Cost(3)
	}

	// Convert data rows to table.Row with conditional coloring and separators
	for rowIdx, row := range data.Rows {
		// Add separator every 10 rows for readability
		if rowIdx > 0 && rowIdx%10 == 0 {
			t.AppendSeparator()
		}

		tableRow := make(table.Row, len(row))
		for i, cell := range row {
			tableRow[i] = r.formatCell(cell, i, row, costColIdx)
		}
		t.AppendRow(tableRow)
	}

	// Footer with total
	var footer table.Row
	if data.HasTagColumn {
		footer = table.Row{"", "", "", data.Total, "100%", "", ""}
	} else {
		footer = table.Row{"", "", data.Total, "100%", "", ""}
	}
	t.AppendFooter(footer)

	t.Render()
	fmt.Println()
	return nil
}

// formatCell applies conditional formatting to cells based on value and position
func (r *StdoutTableRenderer) formatCell(cell string, colIndex int, row []string, costColIdx int) interface{} {
	// Apply color coding based on cost amount
	if colIndex == costColIdx && len(row) > costColIdx {
		amount, err := strconv.ParseFloat(row[costColIdx], 64)
		if err == nil {
			switch {
			case amount >= 5.0:
				return text.FgHiRed.Sprint(cell) // High cost - red
			case amount >= 1.0:
				return text.FgHiYellow.Sprint(cell) // Medium cost - yellow
			case amount >= 0.01:
				return text.FgHiGreen.Sprint(cell) // Low cost - green
			default:
				return text.FgHiBlack.Sprint(cell) // Negligible - dim
			}
		}
	}
	return cell
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
	if err := ensureOutputDir(); err != nil {
		return types.Error{Msg: "Error creating output directory: " + err.Error()}
	}

	filePath := outputPath(data.Filename)
	file, err := os.Create(filePath)
	if err != nil {
		return types.Error{Msg: "Error creating CSV file: " + err.Error()}
	}
	defer func() { _ = file.Close() }()

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

	fmt.Printf("CSV written to: %s\n", filePath)
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
	if err := ensureOutputDir(); err != nil {
		return types.Error{Msg: "Error creating output directory: " + err.Error()}
	}

	filePath := outputPath(data.Filename)
	file, err := os.Create(filePath)
	if err != nil {
		return types.Error{Msg: "Failed creating chart HTML file: " + err.Error()}
	}
	defer func() { _ = file.Close() }()

	if err := data.Page.Render(io.MultiWriter(file)); err != nil {
		return err
	}

	fmt.Printf("Chart written to: %s\n", filePath)
	return nil
}

// VectorRenderer renders vector data to vector databases
type VectorRenderer struct{}

// NewVectorRenderer creates a new vector renderer
func NewVectorRenderer() *VectorRenderer {
	return &VectorRenderer{}
}

// Render implements the Renderer interface for vector databases
func (r *VectorRenderer) Render(data *VectorOutput) error {
	// NewVectorStoreClient(builder, openAIAPIKey, indexURL, pineconeAPIKey)
	client := NewVectorStoreClient(
		http.NewRequestBuilder(),
		data.OpenAIAPIKey,
		data.IndexURL,
		data.PineconeAPIKey,
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
