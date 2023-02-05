package printer

import (
	"encoding/csv"
	"os"
)

var (
	header = []string{"Dimension/Tag", "Dimension/Tag", "Metric", "Granularity",
		"Start",
		"End", "USD Amount", "Unit"}
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

	_, err := CreateOutPutDir()
	if err != nil {
		return PrinterError{
			msg: "Error creating output directory: " + err.Error(),
		}
	}

	f, err := os.Create("./output/ccexplorer.csv")
	if err != nil {
		return PrinterError{
			msg: "Error creating CSV file: " + err.Error()}
	}
	defer f.Close()

	// Write the header row
	w := csv.NewWriter(f)
	err = w.Write(header)
	if err != nil {
		return PrinterError{
			msg: "Error writing header to CSV file: " + err.Error()}
	}

	var rows [][]string
	for _, v := range r.Services {
		rows = append(rows, ConvertServiceToSlice(v, r.Granularity)...)
	}

	if err := w.WriteAll(rows); err != nil {
		return PrinterError{
			msg: "Error writing to CSV file: " + err.Error()}
	}

	return nil
}

func CostAndUsageToChart(sortFn func(r map[int]Service) []Service,
	r CostAndUsageOutputType) error {

	_, err := CreateOutPutDir()
	if err != nil {
		return PrinterError{
			msg: "Error creating output directory: " + err.Error(),
		}
	}

	render := Renderer{}
	err = render.Charts(r)
	if err != nil {
		return err
	}

	return nil
}
