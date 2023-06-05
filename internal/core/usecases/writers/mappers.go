package writers

import (
	"fmt"
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	"github.com/cduggn/ccexplorer/internal/core/service/openai"
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

func CostAndUsageToPineconeMapper(r model.CostAndUsageOutputType) error {

	embeddingClient := openai.NewClient(r.OpenAIAPIKey)

	e, err := embeddingClient.GenerateEmbeddings("test")
	if err != nil {
		return model.Error{
			Msg: "Error generating embeddings for pinecone : " + err.Error()}
	}

	fmt.Print(e)
	//dbclient := NewVectorStoreClient(requestbuilder.NewRequestBuilder(),
	//	"ccexplorer","")
	//
	//r.OpenAIAPIKey = ""
	//
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
