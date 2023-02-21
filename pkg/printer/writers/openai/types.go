package openai

type TrainingData struct {
	Dimension   string
	Tag         string
	Metric      string
	Granularity string
	Start       string
	End         string
	USDAmount   string
	Unit        string
}

type Error struct {
	msg string
}

func (e Error) Error() string {
	return e.msg
}
