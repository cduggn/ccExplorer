package presentation

import (
	"github.com/cduggn/ccexplorer/pkg/presentation/writers"
)

type Writer interface {
	Write(interface{}, interface{}) error
}

type Port interface {
	NewPrintWriter(printType writers.PrintWriterType, variant string) Writer
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
