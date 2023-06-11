package requestbuilder

import "encoding/json"

type marshaller interface {
	marshal(value any) ([]byte, error)
}

type jsonMarshaller struct{}

func (j *jsonMarshaller) marshal(value any) ([]byte, error) {
	return json.Marshal(value)
}
