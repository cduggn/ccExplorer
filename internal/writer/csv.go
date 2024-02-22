package writer

import (
	"encoding/csv"
	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/cduggn/ccexplorer/internal/utils"
	"io"
	"os"
)

func WriteToCSV(f *os.File, header []string, rows [][]string) error {
	w, err := NewCSVWriter(f, header)
	if err != nil {
		return types.Error{
			Msg: "Error creating CSV writer: " + err.Error()}
	}
	defer w.Flush()

	if err := w.WriteAll(rows); err != nil {
		return types.Error{
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
	path := utils.BuildOutputFilePath(dir, file)
	f, err := os.Create(path)
	if err != nil {
		return nil, types.Error{
			Msg: "Error creating CSV file: " + err.Error()}
	}
	return f, nil
}
