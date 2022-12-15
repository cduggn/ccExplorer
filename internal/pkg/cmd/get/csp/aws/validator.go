package aws

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
