package chart

type Error struct {
	msg string
}

func (e Error) Error() string {
	return e.msg
}

type Builder struct {
}

type InputType struct {
	Services     []Service
	Granularity  string
	Start        string
	End          string
	Dimensions   []string
	Tags         []string
	SortBy       string
	OpenAIAPIKey string
}

// todo remove this duplication
type Service struct {
	Keys    []string
	Name    string
	Metrics []Metrics
	Start   string
	End     string
}

type Metrics struct {
	Name          string
	Amount        string
	NumericAmount float64
	Unit          string
	UsageQuantity float64
}
