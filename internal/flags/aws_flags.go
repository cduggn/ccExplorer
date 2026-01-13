package flags

// Type aliases for cleaner usage
type GroupByFlag = Flag[GroupByType, DimensionValidator]
type FilterByFlag = Flag[FilterByType, FilterValidator] 
type DimensionFilterFlag = Flag[map[string]string, DimensionOnlyValidator]

// Factory functions for creating AWS-specific flags
func NewGroupByFlag() *GroupByFlag {
	return NewFlag[GroupByType](DimensionValidator{})
}

func NewFilterByFlag() *FilterByFlag {
	return NewFlag[FilterByType](FilterValidator{})
}

func NewDimensionFilterFlag() *DimensionFilterFlag {
	return NewFlag[map[string]string](DimensionOnlyValidator{})
}

// Legacy compatibility types and functions for gradual migration
// These will be removed after all usage is updated

// DimensionAndTagFlag provides backward compatibility
type DimensionAndTagFlag struct {
	*GroupByFlag
}

func (f *DimensionAndTagFlag) Value() GroupByType {
	return f.GroupByFlag.Value()
}

// DimensionAndTagFilterFlag provides backward compatibility  
type DimensionAndTagFilterFlag struct {
	*FilterByFlag
}

func (f *DimensionAndTagFilterFlag) Value() FilterByType {
	return f.FilterByFlag.Value()
}

// DimensionFilterByFlag provides backward compatibility
type DimensionFilterByFlag struct {
	*DimensionFilterFlag
}

func (f *DimensionFilterByFlag) Value() map[string]string {
	return f.DimensionFilterFlag.Value()
}

// Factory functions for backward compatibility
func NewDimensionAndTagFlag() *DimensionAndTagFlag {
	return &DimensionAndTagFlag{
		GroupByFlag: NewGroupByFlag(),
	}
}

func NewFilterBy() *DimensionAndTagFilterFlag {
	return &DimensionAndTagFilterFlag{
		FilterByFlag: NewFilterByFlag(),
	}
}

func NewDimensionFilterBy() *DimensionFilterByFlag {
	return &DimensionFilterByFlag{
		DimensionFilterFlag: NewDimensionFilterFlag(),
	}
}

// NewForecastFilterBy provides backward compatibility for forecast commands
func NewForecastFilterBy() *DimensionFilterByFlag {
	return NewDimensionFilterBy()
}