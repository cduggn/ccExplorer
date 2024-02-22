package writer

type Printer interface {
	Write(interface{}, interface{}) error
}
