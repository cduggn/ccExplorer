package writer

import (
	"context"
	"fmt"
	"github.com/cduggn/ccexplorer/internal/http"
	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/cduggn/ccexplorer/internal/utils"
	"log/slog"
	"os"
	"strings"
)

var (
	csvFileName = "ccexplorer.csv"
	csvHeader   = []string{"Dimension/Tag", "Dimension/Tag", "Metric",
		"Granularity",
		"Start",
		"End", "USD Amount", "Unit"}
	OutputDir = "./writer"
)

type Builder struct {
}

// Legacy printer types - kept for interface compatibility
// Actual implementations are now in writers.go using generics

func init() {
	if _, err := os.Stat(OutputDir); os.IsNotExist(err) {
		err := os.Mkdir(OutputDir, 0755)
		if err != nil {
			panic("Unable writer directory")
		}
	}
}

func NewPrintWriter(printType types.PrintWriterType, variant string) Printer {
	switch printType {
	case types.Stdout:
		return NewGenericStdoutPrinter(variant)
	case types.CSV:
		return NewGenericCsvPrinter(variant)
	case types.Chart:
		return NewGenericChartPrinter(variant)
	case types.Pinecone:
		return NewGenericPineconePrinter(variant)
	default:
		panic("Invalid print type")
	}
}

// Legacy mapper functions - kept for backward compatibility but will be removed
// These are now handled by the generic transformers and renderers

func CostAndUsageToVectorMapper(r types.CostAndUsageOutputType) error {

	client := NewVectorStoreClient(http.NewRequestBuilder(),
		r.OpenAIAPIKey, r.PineconeIndex, r.PineconeAPIKey)
	items, err := client.CreateVectorStoreInput(r)
	if err != nil {
		return types.Error{
			Msg: "Error writing to vector store: " + err.Error()}
	}

	vectors, err := client.CreateEmbeddings(items)
	if err != nil {
		return types.Error{
			Msg: "Error writing to vector store: " + err.Error()}
	}

	for index, m := range vectors {
		items[index].EmbeddingVector = m.Embedding
		items[index].ID = utils.EncodeString(items[index].EmbeddingText)
	}

	input := utils.ConvertToPineconeStruct(items)

	resp, err := client.Upsert(context.Background(), input)
	if err != nil {
		return types.Error{
			Msg: "Error writing to vector store: " + err.Error()}
	}

	slog.Info(fmt.Sprintf("Upserted %d items to vector store", resp.UpsertedCount))
	return nil
}

func CostAndUsageToStdoutMapper(sortFn func(r map[int]types.Service) []types.Service,
	r types.CostAndUsageOutputType) error {

	sortedServices := sortFn(r.Services)
	output := utils.ConvertToStdoutType(sortedServices, r.Granularity)

	w, err := NewStdoutWriter("costAndUsage")
	if err != nil {
		return types.Error{
			Msg: "Error writing to stdout : " + err.Error()}
	}
	w.Writer(output)
	return nil
}

func CostAndUsageToCSVMapper(sortFn func(r map[int]types.Service) []types.Service,
	r types.CostAndUsageOutputType) error {

	f, err := NewCSVFile(OutputDir, csvFileName)
	if err != nil {
		return types.Error{
			Msg: "Error creating CSV file: " + err.Error()}
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	rows := utils.ConvertServiceMapToArray(r.Services, r.Granularity)
	err = WriteToCSV(f, csvHeader, rows)
	if err != nil {
		return types.Error{
			Msg: "Error writing to CSV file: " + err.Error()}
	}
	return nil
}

func CostAndUsageToChartMapper(sortFn func(r map[int]types.Service) []types.Service,
	r types.CostAndUsageOutputType) error {

	builder := Builder{}
	s := sortFn(r.Services)
	input := utils.ConvertToChartInputType(r, s)

	charts, err := builder.NewCharts(input)
	if err != nil {
		return err
	}

	err = WriteToChart(charts)
	if err != nil {
		return err
	}
	return nil
}

func ForecastToStdoutMapper(r types.ForecastPrintData,
	dimensions []string) {

	filteredBy := strings.Join(dimensions, " | ")
	output := utils.ConvertToForecastStdoutType(r, filteredBy)
	w, err := NewStdoutWriter("forecast")
	if err != nil {
		return
	}
	w.Writer(output)
}
