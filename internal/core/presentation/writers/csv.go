package writers

import (
	"encoding/csv"
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	"io"
	"os"
)

func WriteToCSV(f *os.File, header []string, rows [][]string) error {
	w, err := NewCSVWriter(f, header)
	if err != nil {
		return model.Error{
			Msg: "Error creating CSV writer: " + err.Error()}
	}
	defer w.Flush()

	if err := w.WriteAll(rows); err != nil {
		return model.Error{
			Msg: "Error writing to CSV file: " + err.Error()}
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
	path := buildOutputFilePath(dir, file)
	f, err := os.Create(path)
	if err != nil {
		return nil, model.Error{
			Msg: "Error creating CSV file: " + err.Error()}
	}
	return f, nil
}
