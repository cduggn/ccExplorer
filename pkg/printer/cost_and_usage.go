package printer

import (
	"github.com/cduggn/ccexplorer/pkg/printer/writers/chart"
	"github.com/cduggn/ccexplorer/pkg/printer/writers/csv"
	"github.com/cduggn/ccexplorer/pkg/printer/writers/openai"
	"github.com/jedib0t/go-pretty/v6/table"
)

var (
	csvFileName = "ccexplorer.csv"
	csvHeader   = []string{"Dimension/Tag", "Dimension/Tag", "Metric",
		"Granularity",
		"Start",
		"End", "USD Amount", "Unit"}

	costAndUsageHeader = table.Row{"Rank", "Dimension/Tag", "Dimension/Tag",
		"Metric Name", "Truncated USD Amount", "Amount",
		"Unit",
		"Granularity",
		"Start",
		"End"}
	tableDivider = table.Row{"-", "-", "-",
		"-", "-", "-", "-",
		"-",
		"-", ""}
	costAndUsageTableFooter = func(t string) table.Row {
		return table.
			Row{"", "",
			"",
			"",
			"TOTAL COST",
			t, "", "", "", ""}
	}
)

func CostAndUsageToStdout(sortFn func(r map[int]Service) []Service,
	r CostAndUsageOutputType) {
	sortedServices := sortFn(r.Services)

	t := CreateTable(costAndUsageHeader)

	granularity := r.Granularity

	rows := CostUsageToRows(sortedServices, granularity)

	t.AppendRows(rows.Rows)
	t.AppendRow(tableDivider)
	t.AppendRow(costAndUsageTableFooter(rows.Total))

	t.Render()
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
		return nil
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

	sorted := sortFn(r.Services)

	rows := ConvertServiceSliceToArray(sorted, r.Granularity)

	data := openai.BuildPromptText(rows)

	resp, err := openai.Summarize(r.OpenAIAPIKey, data)
	if err != nil {
		return err
	}

	err = openai.Writer(resp.Choices[0].Text)
	if err != nil {
		return err
	}

	return nil
}
