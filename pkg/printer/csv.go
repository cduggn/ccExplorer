package printer

import (
	"encoding/csv"
	"io"
	"os"
)

var (
	csvheader = []string{"Dimension/Tag", "Dimension/Tag", "Metric",
		"Granularity",
		"Start",
		"End", "USD Amount", "Unit"}
	csvFileName           = "ccexplorer.csv"
	csvHeaderPromptFormat = "Dimension/Tag,Dimension/Tag,Metric," +
		"Granularity,Start,End,USD Amount,Unit;"
)

func CSVWriter(f *os.File, header []string, rows [][]string) error {
	w, err := NewCSVWriter(f, header)
	if err != nil {
		return PrinterError{
			msg: "Error creating CSV writer: " + err.Error()}
	}
	defer w.Flush()

	if err := w.WriteAll(rows); err != nil {
		return PrinterError{
			msg: "Error writing to CSV file: " + err.Error()}
	}
	return nil
}

func NewCSVWriter(f io.Writer, header []string) (*csv.Writer, error) {
	w := csv.NewWriter(f)
	err := w.Write(header)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func NewCSVFile(dir string, file string) (*os.File, error) {
	path := BuildOutputFilePath(dir, file)
	f, err := os.Create(path)
	if err != nil {
		return nil, PrinterError{
			msg: "Error creating CSV file: " + err.Error()}
	}
	return f, nil
}
