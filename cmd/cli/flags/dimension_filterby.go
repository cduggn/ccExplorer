package flags

import (
	"fmt"
	"github.com/cduggn/ccexplorer/internal/utils"
	"strings"
)

func (e DimensionFilterByFlagError) Error() string {
	return e.msg
}

func NewForecastFilterBy() DimensionFilterByFlag {
	return DimensionFilterByFlag{
		Dimensions: make(map[string]string),
	}
}

func (f *DimensionFilterByFlag) Set(value string) error {

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
		default:
			return DimensionFilterByFlagError{
				msg: fmt.Sprintf("invalid groupBy type selected: %s", value),
			}
		}
	}

	return nil
}

func (f *DimensionFilterByFlagType) Value() DimensionFilterByFlagType {
	return DimensionFilterByFlagType(*f)
}

func (f *DimensionFilterByFlagType) Equals(other DimensionFilterByFlagType) bool {
	if len(f.Dimensions) != len(other.Dimensions) {
		return false
	}
	for i := range f.Dimensions {
		if f.Dimensions[i] != other.Dimensions[i] {
			return false
		}
	}
	return true

}

func (f *DimensionFilterByFlag) Type() string {
	return "FilterBy"
}

func (f *DimensionFilterByFlag) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *DimensionFilterByFlag) Value() DimensionFilterByFlag {
	return *f
}
