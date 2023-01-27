package custom_flags

type GroupByFlagError struct {
	msg string
}

type GroupBy GroupByType

type GroupByType struct {
	Dimensions []string
	Tags       []string
}

type FilterByFlagError struct {
	msg string
}

type FilterBy FilterByType

type FilterByType struct {
	Dimensions map[string]string
	Tags       []string
}

type ForecastFilterError struct {
	msg string
}

type ForecastFilterByType struct {
	Dimensions map[string]string
}

type ForecastFilterBy ForecastFilterByType
