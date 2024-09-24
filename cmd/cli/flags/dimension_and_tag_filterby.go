package flags

import (
	"fmt"
	"github.com/cduggn/ccexplorer/internal/utils"
	"strings"
)

func (e DimensionAndTagFilterFlagError) Error() string {
	return e.msg
}

func NewFilterBy() DimensionAndTagFilterFlag {
	return DimensionAndTagFilterFlag{
		Dimensions: make(map[string]string),
		Tags:       make([]string, 0),
	}
}

func (f *DimensionAndTagFilterFlag) Set(value string) error {

	args := utils.SplitCommaSeparatedString(value)

	for _, arg := range args {

		parts, err := utils.SplitNameValuePair(arg)
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
			return DimensionAndTagFilterFlagError{
				msg: fmt.Sprintf("invalid groupBy type selected: %s", value),
			}
		}
	}
	return nil
}

func (f *DimensionAndTagFilterFlagType) Value() DimensionAndTagFilterFlagType {
	return DimensionAndTagFilterFlagType(*f)
}

func (f *DimensionAndTagFilterFlagType) Equals(other DimensionAndTagFilterFlagType) bool {
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

func (f *DimensionAndTagFilterFlag) Type() string {
	return "FilterBy"
}

func (f *DimensionAndTagFilterFlag) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *DimensionAndTagFilterFlag) Value() DimensionAndTagFilterFlag {
	return *f
}
