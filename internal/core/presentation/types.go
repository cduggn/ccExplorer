package presentation

type Printer interface {
	Print(interface{}, interface{}) error
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
