package custom_flags

import (
	"fmt"
	"strings"
)

var (
	validDimensions = map[string]bool{
		"AZ": true, "SERVICE": true, "USAGE_TYPE": true, "INSTANCE_TYPE": true,
		"LINKED_ACCOUNT": true, "OPERATION": true, "PURCHASE_TYPE": true,
		"PLATFORM": true, "TENANCY": true, "RECORD_TYPE": true,
		"LEGAL_ENTITY_NAME": true, "INVOICING_ENTITY": true, "DEPLOYMENT_OPTION": true,
		"DATABASE_ENGINE": true, "CACHE_ENGINE": true, "INSTANCE_TYPE_FAMILY": true,
		"REGION": true, "BILLING_ENTITY": true, "RESERVATION_ID": true,
		"SAVINGS_PLANS_TYPE": true, "SAVINGS_PLAN_ARN": true, "OPERATING_SYSTEM": true,
	}
	DIMENSIONS = "AZ, SERVICE, " +
		"USAGE_TYPE, INSTANCE_TYPE, LINKED_ACCOUNT, OPERATION, " +
		"PURCHASE_TYPE, PLATFORM, TENANCY, RECORD_TYPE, " +
		"LEGAL_ENTITY_NAME, INVOICING_ENTITY, DEPLOYMENT_OPTION, " +
		"DATABASE_ENGINE, CACHE_ENGINE, INSTANCE_TYPE_FAMILY, " +
		"REGION, BILLING_ENTITY, RESERVATION_ID, " +
		"SAVINGS_PLANS_TYPE, SAVINGS_PLAN_ARN, OPERATING_SYSTEM"
	validTypes = "DIMENSION, TAG"
)

func IsValidDimension(d string) bool {
	if _, ok := validDimensions[strings.ToUpper(d)]; !ok {
		return false
	}
	return true
}

func SplitCommaSeparatedString(value string) []string {
	var args []string
	if strings.Contains(value, ",") {
		args = strings.Split(value, ",")
	} else {
		args = []string{value}
	}
	return args
}

func SplitNameValuePair(value string) ([]string, error) {
	parts := strings.Split(value, "=")
	if len(parts) != 2 {
		return nil, GroupByFlagError{
			msg: fmt.Sprintf("invalid group by flag: %s", value),
		}
	}
	return parts, nil
}
