package model

type APIError struct {
	Msg string
}

func (e APIError) Error() string {
	return e.Msg
}
