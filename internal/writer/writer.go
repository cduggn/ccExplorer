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

type StdoutPrinter struct {
	Variant string
}

type CsvPrinter struct {
	Variant string
}

type OpenAIPrinter struct {
	Variant string
}

type ChartPrinter struct {
	Variant string
}

type PineconePrinter struct {
	Variant string
}

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
		return &StdoutPrinter{variant}
	case types.CSV:
		return &CsvPrinter{variant}
	case types.Chart:
		return &ChartPrinter{variant}
	case types.Pinecone:
		return &PineconePrinter{variant}
	default:
		panic("Invalid print type")
	}
}

func (p *PineconePrinter) Write(f interface{}, c interface{}) error {
	switch p.Variant {
	case "costAndUsage":
		err := CostAndUsageToVectorMapper(c.(types.CostAndUsageOutputType))
		if err != nil {
			return err
		}
		/// working with CostAndUsageOutputType
	}
	return nil
}

func (p *StdoutPrinter) Write(f interface{}, c interface{}) error {
	switch p.Variant {
	case "forecast":
		ForecastToStdoutMapper(f.(types.ForecastPrintData),
			c.([]string))
	case "costAndUsage":
		fn := utils.SortFunction(f.(string))
		err := CostAndUsageToStdoutMapper(fn,
			c.(types.CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *CsvPrinter) Write(f interface{}, c interface{}) error {
	switch p.Variant {
	case "costAndUsage":
		fn := utils.SortFunction(f.(string))
		err := CostAndUsageToCSVMapper(fn,
			c.(types.CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ChartPrinter) Write(f interface{}, c interface{}) error {
	switch p.Variant {
	case "costAndUsage":
		fn := utils.SortFunction(f.(string))
		err := CostAndUsageToChartMapper(fn,
			c.(types.CostAndUsageOutputType))
		if err != nil {
			return err
		}
	}
	return nil
}

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
