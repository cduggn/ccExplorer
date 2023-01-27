package custom_flags

import (
	"fmt"
	"strings"
)

func (e ForecastFilterError) Error() string {
	return e.msg
}

func NewForecastFilterBy() ForecastFilterBy {
	return ForecastFilterBy{
		Dimensions: make(map[string]string),
	}
}

func (f *ForecastFilterBy) Set(value string) error {

	args := SplitCommaSeparatedString(value)

	for _, arg := range args {

		parts, err := SplitNameValuePair(arg)
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
			return ForecastFilterError{
				msg: fmt.Sprintf("invalid groupBy type selected: %s", value),
			}
		}
	}

	return nil
}

func (f *ForecastFilterByType) Value() ForecastFilterByType {
	return ForecastFilterByType(*f)
}

func (f *ForecastFilterByType) Equals(other ForecastFilterByType) bool {
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

func (f *ForecastFilterBy) Type() string {
	return "FilterBy"
}

func (f *ForecastFilterBy) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *ForecastFilterBy) Value() ForecastFilterBy {
	return *f
}
