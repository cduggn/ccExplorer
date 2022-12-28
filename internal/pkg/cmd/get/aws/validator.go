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

func ValidateDimension(dimension []string) error {

	if len(dimension) > 2 || len(dimension) == 0 {
		return ValidationError{
			msg: "At most 2 dimensions are allowed. " +
				"When grouping by tag, only 1 dimension is allowed",
		}
	}

	for _, d := range dimension {
		switch d {
		case "AZ", "SERVICE", "USAGE_TYPE", "INSTANCE_TYPE", "LINKED_ACCOUNT", "OPERATION",
			"PURCHASE_TYPE", "PLATFORM", "TENANCY", "RECORD_TYPE", "LEGAL_ENTITY_NAME",
			"INVOICING_ENTITY", "DEPLOYMENT_OPTION", "DATABASE_ENGINE",
			"CACHE_ENGINE", "INSTANCE_TYPE_FAMILY", "REGION", "BILLING_ENTITY",
			"RESERVATION_ID", "SAVINGS_PLANS_TYPE", "SAVINGS_PLAN_ARN",
			"OPERATING_SYSTEM":
			continue
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
	}
	return nil
}

func ValidateTag(tag string, dimension []string) error {
	if tag != "" && len(dimension) != 1 {
		return ValidationError{
			msg: "When grouping by tag, 1 dimension is allowed",
		}
	}
	return nil
}

func ValidateFilterBy(filterBy string, tag string) error {
	if filterBy != "" && tag == "" {
		return ValidationError{
			msg: "When filtering by tag value, a tag must be specified",
		}
	}
	return nil
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
