package custom_flags

import (
	"fmt"
	"strings"
)

var (
	groupByFlag GroupBy
)

func (e GroupByFlagError) Error() string {
	return e.msg
}

func (f *GroupBy) Set(value string) error {

	args := SplitCommaSeparatedString(value)

	for _, arg := range args {

		parts, err := SplitNameValuePair(arg)
		if err != nil {
			return err
		}
		switch strings.ToUpper(parts[0]) {
		case "DIMENSION":

			if ok := IsValidDimension(parts[1]); !ok {
				return GroupByFlagError{
					msg: fmt.Sprintf("Invalid dimension: %s . "+
						"Must be one of %s", parts[1], DIMENSIONS),
				}
			}

			f.Dimensions = append(f.Dimensions, parts[1])
		case "TAG":
			f.Tags = append(f.Tags, parts[1])
		default:
			return GroupByFlagError{
				msg: fmt.Sprintf("invalid groupBy type selected: %s must be"+
					" one of: %s ",
					value, validTypes),
			}
		}
	}

	return nil
}

func (f *GroupByType) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *GroupByType) Value() GroupByType {
	return GroupByType(*f)
}

func (f *GroupByType) Equals(other GroupByType) bool {
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

func (f *GroupBy) Type() string {
	return "GroupBy"
}

func (f *GroupBy) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *GroupBy) Value() GroupBy {
	return *f
}
