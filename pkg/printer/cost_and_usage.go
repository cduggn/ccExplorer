package printer

import (
	"github.com/cduggn/ccexplorer/pkg/printer/writers/chart"
	"github.com/cduggn/ccexplorer/pkg/printer/writers/csv"
	"github.com/cduggn/ccexplorer/pkg/printer/writers/openai"
	"github.com/cduggn/ccexplorer/pkg/printer/writers/stdout"
)

var (
	maxDisplayRows = 10
	csvFileName    = "ccexplorer.csv"
	csvHeader      = []string{"Dimension/Tag", "Dimension/Tag", "Metric",
		"Granularity",
		"Start",
		"End", "USD Amount", "Unit"}
)

func CostAndUsageToStdout(sortFn func(r map[int]Service) []Service,
	r CostAndUsageOutputType) {

	sortedServices := sortFn(r.Services)
	output := ConvertToStdoutType(sortedServices, r.Granularity)

	w, err := stdout.NewStdoutWriter("costAndUsage")
	if err != nil {
		return
	}

	w.Writer(output)
}

func CostAndUsageToCSV(sortFn func(r map[int]Service) []Service,
	r CostAndUsageOutputType) error {

	f, err := csv.NewCSVFile(OutputDir, csvFileName)
	if err != nil {
		return Error{
			msg: "Error creating CSV file: " + err.Error()}
	}
	defer f.Close()

	rows := ConvertServiceMapToArray(r.Services, r.Granularity)

	err = csv.Writer(f, csvHeader, rows)
	if err != nil {
		return Error{
			msg: "Error writing to CSV file: " + err.Error()}
	}

	return nil
}

func CostAndUsageToChart(sortFn func(r map[int]Service) []Service,
	r CostAndUsageOutputType) error {

	builder := chart.Builder{}

	s := sortFn(r.Services)

	input := ConvertToChartInputType(r, s)

	charts, err := builder.NewCharts(input)
	if err != nil {
		return err
	}

	err = chart.Writer(charts)
	if err != nil {
		return err
	}
	return nil
}

func CostAndUsageToOpenAI(sortFn func(r map[int]Service) []Service,
	r CostAndUsageOutputType) error {

	sortedData := sortFn(r.Services)
	rows := ConvertServiceSliceToArray(sortedData, r.Granularity)

	maxRows := MaxRows(rows, maxDisplayRows)

	data := openai.BuildPromptText(rows[:maxRows])
	resp, err := openai.Summarize(r.OpenAIAPIKey, data)
	if err != nil {
		return err
	}
	err = openai.Writer(resp.Choices[0].Message.Content)

	if err != nil {
		return err
	}
	return nil
}

func MaxRows(rows [][]string, maxRows int) int {
	if len(rows) > maxRows {
		return maxRows
	}
	return len(rows)
}
