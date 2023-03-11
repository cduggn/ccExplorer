package flags

type FlagError struct {
	Msg string
}

type DimensionAndTagFlagError struct {
	Msg string
}

type DimensionAndTagFlag DimensionAndTagFlagType

type DimensionAndTagFlagType struct {
	Dimensions []string
	Tags       []string
}

type DimensionAndTagFilterFlagError struct {
	msg string
}

type DimensionAndTagFilterFlag DimensionAndTagFilterFlagType

type DimensionAndTagFilterFlagType struct {
	Dimensions map[string]string
	Tags       []string
}

type DimensionFilterByFlagError struct {
	msg string
}

type DimensionFilterByFlagType struct {
	Dimensions map[string]string
}

type DimensionFilterByFlag DimensionFilterByFlagType
