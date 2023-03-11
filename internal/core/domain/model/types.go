package model

type APIError struct {
	Msg string
}

func (e APIError) Error() string {
	return e.Msg
}

type PresetError struct {
	Msg string
}

func (e PresetError) Error() string {
	return e.Msg
}
