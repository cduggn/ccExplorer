package aws

import (
	"fmt"
	"strings"
)

type GroupByFlagError struct {
	msg string
}

func (e GroupByFlagError) Error() string {
	return e.msg
}

var grouppByFlag GroupBy

type GroupBy GroupByType

type GroupByType struct {
	Dimensions []string
	Tags       []string
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

func (f *GroupBy) Set(value string) error {

	args := splitByIndividualArgument(value)

	for _, arg := range args {

		parts, err := splitIndividualArgument(arg)
		if err != nil {
			return err
		}
		switch strings.ToUpper(parts[0]) {
		case "DIMENSION":
			f.Dimensions = append(f.Dimensions, parts[1])
		case "TAG":
			f.Tags = append(f.Tags, parts[1])
		default:
			return GroupByFlagError{
				msg: fmt.Sprintf("invalid groupBy type selected: %s", value),
			}
		}
	}

	return nil
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

func splitByIndividualArgument(value string) []string {
	var args []string
	if strings.Contains(value, ",") {
		args = strings.Split(value, ",")
	} else {
		args = strings.Split(value, " ")
	}
	return args
}

func splitIndividualArgument(value string) ([]string, error) {
	parts := strings.Split(value, "=")
	if len(parts) != 2 {
		return nil, GroupByFlagError{
			msg: fmt.Sprintf("invalid group by flag: %s", value),
		}
	}
	return parts, nil
}
