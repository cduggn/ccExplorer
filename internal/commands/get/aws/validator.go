package aws

import (
	"time"
)

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

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
			Message: "Dimension must be one of the following: AZ, SERVICE, " +
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

func ValidateStartDate(startDate string) error {
	if startDate == "" {
		return ValidationError{
			Message: "Start date must be specified",
		}
	}

	start, _ := time.Parse("2006-01-02", startDate)
	today := time.Now()
	if start.After(today) {
		return ValidationError{
			Message: "Start date must be before today's date",
		}
	}

	return nil
}

func ValidateEndDate(endDate, startDate string) error {
	if endDate == "" {
		return ValidationError{
			Message: "End date must be specified",
		}
	}

	end, _ := time.Parse("2006-01-02", endDate)
	today := time.Now()
	if end.After(today) {
		return ValidationError{
			Message: "End date must be before today's date",
		}
	}

	start, _ := time.Parse("2006-01-02", startDate)
	if end.Before(start) {
		return ValidationError{
			Message: "End date must not be before start date",
		}
	}

	return nil
}
