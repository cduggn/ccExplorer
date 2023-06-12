package usecases

import (
	"context"
	"fmt"
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	"github.com/cduggn/ccexplorer/internal/core/logger"
	"github.com/cduggn/ccexplorer/internal/core/requestbuilder"
	"github.com/cduggn/ccexplorer/internal/core/util"
	"os"
	"strings"
)

var (
	csvFileName = "ccexplorer.csv"
	csvHeader   = []string{"Dimension/Tag", "Dimension/Tag", "Metric",
		"Granularity",
		"Start",
		"End", "USD Amount", "Unit"}
	OutputDir = "./output"
)

type Builder struct {
}

func CostAndUsageToVectorMapper(r model.CostAndUsageOutputType) error {

	client := NewVectorStoreClient(requestbuilder.NewRequestBuilder(),
		r.OpenAIAPIKey, r.PineconeIndex, r.PineconeAPIKey)
	items, err := client.CreateVectorStoreInput(r)
	if err != nil {
		return model.Error{
			Msg: "Error writing to vector store: " + err.Error()}
	}

	vectors, err := client.CreateEmbeddings(items)
	if err != nil {
		return model.Error{
			Msg: "Error writing to vector store: " + err.Error()}
	}

	for index, m := range vectors {
		items[index].EmbeddingVector = m.Embedding
		items[index].ID = util.EncodeString(items[index].EmbeddingText)
	}

	input := util.ConvertToPineconeStruct(items)

	resp, err := client.Upsert(context.Background(), input)
	if err != nil {
		return model.Error{
			Msg: "Error writing to vector store: " + err.Error()}
	}

	logger.Info(fmt.Sprintf("Upserted %d items to vector store", resp.UpsertedCount))
	return nil
}

func CostAndUsageToStdoutMapper(sortFn func(r map[int]model.Service) []model.
	Service,
	r model.CostAndUsageOutputType) error {

	sortedServices := sortFn(r.Services)
	output := util.ConvertToStdoutType(sortedServices, r.Granularity)

	w, err := NewStdoutWriter("costAndUsage")
	if err != nil {
		return model.Error{
			Msg: "Error writing to stdout : " + err.Error()}
	}
	w.Writer(output)
	return nil
}

func CostAndUsageToCSVMapper(sortFn func(r map[int]model.Service) []model.
	Service,
	r model.CostAndUsageOutputType) error {

	f, err := NewCSVFile(OutputDir, csvFileName)
	if err != nil {
		return model.Error{
			Msg: "Error creating CSV file: " + err.Error()}
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	rows := util.ConvertServiceMapToArray(r.Services, r.Granularity)
	err = WriteToCSV(f, csvHeader, rows)
	if err != nil {
		return model.Error{
			Msg: "Error writing to CSV file: " + err.Error()}
	}
	return nil
}

func CostAndUsageToChartMapper(sortFn func(r map[int]model.Service) []model.
	Service,
	r model.CostAndUsageOutputType) error {

	builder := Builder{}
	s := sortFn(r.Services)
	input := util.ConvertToChartInputType(r, s)

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

func ForecastToStdoutMapper(r model.ForecastPrintData,
	dimensions []string) {

	filteredBy := strings.Join(dimensions, " | ")
	output := util.ConvertToForecastStdoutType(r, filteredBy)
	w, err := NewStdoutWriter("forecast")
	if err != nil {
		return
	}
	w.Writer(output)
}
