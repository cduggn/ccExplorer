package model

type PresetParams struct {
	Alias             string
	Dimension         []string
	Tag               string
	Filter            map[string]string
	FilterType        string
	FilterByDimension bool
	FilterByTag       bool
	ExcludeDiscounts  bool
	CommandSyntax     string
	Description       []string
	Granularity       string
	PrintFormat       string
	Metric            []string
}
