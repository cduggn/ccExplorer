package custom_flags

import (
	"fmt"
	"strings"
)

var (
	groupByFlag DimensionAndTagFlag
)

func (e DimensionAndTagFlagError) Error() string {
	return e.msg
}

func (f *DimensionAndTagFlag) Set(value string) error {

	args := SplitCommaSeparatedString(value)

	for _, arg := range args {

		parts, err := SplitNameValuePair(arg)
		if err != nil {
			return err
		}
		switch strings.ToUpper(parts[0]) {
		case "DIMENSION":

			if ok := IsValidDimension(parts[1]); !ok {
				return DimensionAndTagFlagError{
					msg: fmt.Sprintf("Invalid dimension: %s . "+
						"Must be one of %s", parts[1], DIMENSIONS),
				}
			}

			f.Dimensions = append(f.Dimensions, parts[1])
		case "TAG":
			f.Tags = append(f.Tags, parts[1])
		default:
			return DimensionAndTagFlagError{
				msg: fmt.Sprintf("invalid groupBy type selected: %s must be"+
					" one of: %s ",
					value, validTypes),
			}
		}
	}

	return nil
}

func (f *DimensionAndTagFlagType) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *DimensionAndTagFlagType) Value() DimensionAndTagFlagType {
	return DimensionAndTagFlagType(*f)
}

func (f *DimensionAndTagFlagType) Equals(other DimensionAndTagFlagType) bool {
	if len(f.Dimensions) != len(other.Dimensions) {
		return false
	}
	if len(f.Tags) != len(other.Tags) {
		return false
	}
	for i := range f.Dimensions {
		if f.Dimensions[i] != other.Dimensions[i] {
			return false
		}
	}
	for i := range f.Tags {
		if f.Tags[i] != other.Tags[i] {
			return false
		}
	}
	return true
}

func (f *DimensionAndTagFlag) Type() string {
	return "GroupBy"
}

func (f *DimensionAndTagFlag) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *DimensionAndTagFlag) Value() DimensionAndTagFlag {
	return *f
}
