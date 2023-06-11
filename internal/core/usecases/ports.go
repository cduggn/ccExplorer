package usecases

type Printer interface {
	Write(interface{}, interface{}) error
}
