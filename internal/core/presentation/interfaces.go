package presentation

type Printer interface {
	Print(interface{}, interface{}) error
}
