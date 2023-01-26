package aws

import (
	"time"
)

type ValidationError struct {
	msg string
}

func (e ValidationError) Error() string {
	return e.msg
}

//func ValidateGroupByMap(groupBy map[string]string) ([]string, error) {
//
//	var tag map[string]string
//
//	for key, val := range groupBy {
//		if val == "DIMENSION" {
//			//dimension = map[string]string{key: val}
//			err := ValidateGroupByDimension(val)
//			if err != nil {
//				return nil, err
//			}
//		}
//		if val == "TAG" {
//			tag = map[string]string{key: val}
//			//ValidateGroupByTag(tag)
//		} else {
//			return nil, ValidationError{
//				msg: "GroupBy must be one of the following: DIMENSION, TAG",
//			}
//		}
//	}
//	return nil, nil
//}

func ValidateGroupByDimension(dimension string) error {

	switch dimension {
	case "AZ", "SERVICE", "USAGE_TYPE", "INSTANCE_TYPE", "LINKED_ACCOUNT", "OPERATION",
		"PURCHASE_TYPE", "PLATFORM", "TENANCY", "RECORD_TYPE", "LEGAL_ENTITY_NAME",
		"INVOICING_ENTITY", "DEPLOYMENT_OPTION", "DATABASE_ENGINE",
		"CACHE_ENGINE", "INSTANCE_TYPE_FAMILY", "REGION", "BILLING_ENTITY",
		"RESERVATION_ID", "SAVINGS_PLANS_TYPE", "SAVINGS_PLAN_ARN",
		"OPERATING_SYSTEM":
	default:
		return ValidationError{
			msg: "Dimension must be one of the following: AZ, SERVICE, " +
				"USAGE_TYPE, INSTANCE_TYPE, LINKED_ACCOUNT, OPERATION, " +
				"PURCHASE_TYPE, PLATFORM, TENANCY, RECORD_TYPE, " +
				"LEGAL_ENTITY_NAME, INVOICING_ENTITY, DEPLOYMENT_OPTION, " +
				"DATABASE_ENGINE, CACHE_ENGINE, INSTANCE_TYPE_FAMILY, " +
				"REGION, BILLING_ENTITY, RESERVATION_ID, " +
				"SAVINGS_PLANS_TYPE, SAVINGS_PLAN_ARN, OPERATING_SYSTEM",
		}
	}

	return nil
}

func ValidateTagFilterValue(tagFilterKey string, tag string) error {
	if tagFilterKey == "" {
		return nil
	}
	if tag == "" {
		return ValidationError{
			msg: "Tag must be specified",
		}
	}
	return nil
}

func ValidateGroupByTag(tag map[string]string) (string, error) {
	if len(tag) == 0 {
		return "", nil
	}
	if len(tag) > 1 {
		return "", ValidationError{
			msg: "At most 1 tag can be specified",
		}
	}
	// extract the tag from the map
	var tagValue string
	for _, val := range tag {
		tagValue = val
	}
	return tagValue, nil
}

func ValidateStartDate(startDate string) error {
	if startDate == "" {
		return ValidationError{
			msg: "Start date must be specified",
		}
	}

	start, _ := time.Parse("2006-01-02", startDate)
	today := time.Now()
	if start.After(today) {
		return ValidationError{
			msg: "Start date must be before today's date",
		}
	}

	return nil
}

func ValidateEndDate(endDate, startDate string) error {
	if endDate == "" {
		return ValidationError{
			msg: "End date must be specified",
		}
	}

	end, _ := time.Parse("2006-01-02", endDate)
	today := time.Now()
	if end.After(today) {
		return ValidationError{
			msg: "End date must be before today's date",
		}
	}

	start, _ := time.Parse("2006-01-02", startDate)
	if end.Before(start) {
		return ValidationError{
			msg: "End date must not be before start date",
		}
	}

	return nil
}

func validateForecastDimensionKey(dimensions map[string]string) error {

	if len(dimensions) == 0 {
		return ValidationError{
			msg: "At least 1 dimension must be specified",
		}
	}

	for key := range dimensions {
		switch key {
		case "AZ", "SERVICE", "USAGE_TYPE", "INSTANCE_TYPE", "LINKED_ACCOUNT", "OPERATION",
			"PURCHASE_TYPE", "PLATFORM", "TENANCY", "RECORD_TYPE", "LEGAL_ENTITY_NAME",
			"INVOICING_ENTITY", "DEPLOYMENT_OPTION", "DATABASE_ENGINE",
			"CACHE_ENGINE", "INSTANCE_TYPE_FAMILY", "REGION", "BILLING_ENTITY",
			"RESERVATION_ID", "SAVINGS_PLANS_TYPE", "SAVINGS_PLAN_ARN",
			"OPERATING_SYSTEM":
			continue
		default:
			return ValidationError{
				msg: "Dimension KEY must be one of the following: AZ, " +
					"SERVICE,USAGE_TYPE, INSTANCE_TYPE, LINKED_ACCOUNT, OPERATION, " +
					"PURCHASE_TYPE, PLATFORM, TENANCY, RECORD_TYPE, " +
					"LEGAL_ENTITY_NAME, INVOICING_ENTITY, DEPLOYMENT_OPTION, " +
					"DATABASE_ENGINE, CACHE_ENGINE, INSTANCE_TYPE_FAMILY, " +
					"REGION, BILLING_ENTITY, RESERVATION_ID, " +
					"SAVINGS_PLANS_TYPE, SAVINGS_PLAN_ARN, OPERATING_SYSTEM",
			}
		}
	}
	return nil
}
