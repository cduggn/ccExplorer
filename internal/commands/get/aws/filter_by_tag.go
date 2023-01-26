package aws

import (
	"fmt"
	"strings"
)

type FilterByFlagError struct {
	msg string
}

func (e FilterByFlagError) Error() string {
	return e.msg
}

type FilterBy FilterByType

type FilterByType struct {
	Dimensions map[string]string
	Tags       []string
}

func NewFilterBy() FilterBy {
	return FilterBy{
		Dimensions: make(map[string]string),
		Tags:       make([]string, 0),
	}
}

func (f *FilterByType) Value() FilterByType {
	return FilterByType(*f)
}

func (f *FilterByType) Equals(other FilterByType) bool {
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

func (f *FilterBy) Set(value string) error {

	args := splitByIndividualArgument(value)

	for _, arg := range args {

		parts, err := splitIndividualArgument(arg)
		if err != nil {
			return err
		}
		switch strings.ToUpper(parts[0]) {
		case "AZ", "SERVICE", "USAGE_TYPE", "INSTANCE_TYPE", "LINKED_ACCOUNT", "OPERATION",
			"PURCHASE_TYPE", "PLATFORM", "TENANCY", "RECORD_TYPE", "LEGAL_ENTITY_NAME",
			"INVOICING_ENTITY", "DEPLOYMENT_OPTION", "DATABASE_ENGINE",
			"CACHE_ENGINE", "INSTANCE_TYPE_FAMILY", "REGION", "BILLING_ENTITY",
			"RESERVATION_ID", "SAVINGS_PLANS_TYPE", "SAVINGS_PLAN_ARN",
			"OPERATING_SYSTEM":
			f.Dimensions[parts[0]] = parts[1]
		case "TAG":
			f.Tags = append(f.Tags, parts[1])
		default:
			return FilterByFlagError{
				msg: fmt.Sprintf("invalid groupBy type selected: %s", value),
			}
		}
	}

	return nil
}

func (f *FilterBy) Type() string {
	return "FilterBy"
}

func (f *FilterBy) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *FilterBy) Value() FilterBy {
	return *f
}
