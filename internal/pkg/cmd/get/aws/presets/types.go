package presets

type PresetError struct {
	msg string
}

func (e PresetError) Error() string {
	return e.msg
}

type Preset struct {
	Name string
	ID   int
}

type PresetParams struct {
	Alias             string
	Dimension         []string
	Tag               string
	Filter            map[string]string
	FilterType        string
	FilterByDimension bool
	FilterByTag       bool
	ExcludeDiscounts  bool
}
