package flags

import (
	"fmt"
	"strings"

	"github.com/cduggn/ccexplorer/internal/utils"
)

// ValidDimensions contains all valid AWS Cost Explorer dimensions
var ValidDimensions = map[string]bool{
	"AZ": true, "SERVICE": true, "USAGE_TYPE": true, "INSTANCE_TYPE": true,
	"LINKED_ACCOUNT": true, "OPERATION": true, "PURCHASE_TYPE": true,
	"PLATFORM": true, "TENANCY": true, "RECORD_TYPE": true,
	"LEGAL_ENTITY_NAME": true, "INVOICING_ENTITY": true, "DEPLOYMENT_OPTION": true,
	"DATABASE_ENGINE": true, "CACHE_ENGINE": true, "INSTANCE_TYPE_FAMILY": true,
	"REGION": true, "BILLING_ENTITY": true, "RESERVATION_ID": true,
	"SAVINGS_PLANS_TYPE": true, "SAVINGS_PLAN_ARN": true, "OPERATING_SYSTEM": true,
}

// DimensionNames contains all valid dimension names as a slice
var DimensionNames = []string{
	"AZ", "SERVICE", "USAGE_TYPE", "INSTANCE_TYPE", "LINKED_ACCOUNT", "OPERATION",
	"PURCHASE_TYPE", "PLATFORM", "TENANCY", "RECORD_TYPE", "LEGAL_ENTITY_NAME",
	"INVOICING_ENTITY", "DEPLOYMENT_OPTION", "DATABASE_ENGINE", "CACHE_ENGINE",
	"INSTANCE_TYPE_FAMILY", "REGION", "BILLING_ENTITY", "RESERVATION_ID",
	"SAVINGS_PLANS_TYPE", "SAVINGS_PLAN_ARN", "OPERATING_SYSTEM",
}

// GroupByType represents the structure for dimension and tag grouping
type GroupByType struct {
	Dimensions []string
	Tags       []string
}

// FilterByType represents the structure for dimension and tag filtering
type FilterByType struct {
	Dimensions map[string]string
	Tags       []string
}

// DimensionValidator validates AWS dimensions for groupBy operations
type DimensionValidator struct{}

func (v DimensionValidator) Validate(value string) (GroupByType, error) {
	result := GroupByType{
		Dimensions: make([]string, 0),
		Tags:       make([]string, 0),
	}

	args := utils.SplitCommaSeparatedString(value)
	for _, arg := range args {
		parts, err := utils.SplitNameValuePair(arg)
		if err != nil {
			return result, err
		}

		switch strings.ToUpper(parts[0]) {
		case "DIMENSION":
			dimension := strings.ToUpper(parts[1])
			if !ValidDimensions[dimension] {
				return result, ValidationError{
					Field:   "dimension",
					Value:   parts[1],
					Allowed: DimensionNames,
				}
			}
			result.Dimensions = append(result.Dimensions, parts[1])
		case "TAG":
			result.Tags = append(result.Tags, parts[1])
		default:
			return result, ValidationError{
				Field:   "groupBy type",
				Value:   parts[0],
				Allowed: []string{"DIMENSION", "TAG"},
			}
		}
	}

	return result, nil
}

func (v DimensionValidator) AllowedValues() []string {
	return []string{"DIMENSION=<dimension_name>", "TAG=<tag_name>"}
}

func (v DimensionValidator) Type() string {
	return "GroupBy"
}

// FilterValidator validates AWS dimensions and tags for filtering operations
type FilterValidator struct{}

func (v FilterValidator) Validate(value string) (FilterByType, error) {
	result := FilterByType{
		Dimensions: make(map[string]string),
		Tags:       make([]string, 0),
	}

	args := utils.SplitCommaSeparatedString(value)
	for _, arg := range args {
		parts, err := utils.SplitNameValuePair(arg)
		if err != nil {
			return result, err
		}

		switch strings.ToUpper(parts[0]) {
		case "AZ", "SERVICE", "USAGE_TYPE", "INSTANCE_TYPE", "LINKED_ACCOUNT", "OPERATION",
			"PURCHASE_TYPE", "PLATFORM", "TENANCY", "RECORD_TYPE", "LEGAL_ENTITY_NAME",
			"INVOICING_ENTITY", "DEPLOYMENT_OPTION", "DATABASE_ENGINE",
			"CACHE_ENGINE", "INSTANCE_TYPE_FAMILY", "REGION", "BILLING_ENTITY",
			"RESERVATION_ID", "SAVINGS_PLANS_TYPE", "SAVINGS_PLAN_ARN",
			"OPERATING_SYSTEM":
			result.Dimensions[parts[0]] = parts[1]
		case "TAG":
			result.Tags = append(result.Tags, parts[1])
		default:
			return result, ValidationError{
				Field:   "filterBy type",
				Value:   parts[0],
				Allowed: append(DimensionNames, "TAG"),
			}
		}
	}

	return result, nil
}

func (v FilterValidator) AllowedValues() []string {
	allowed := make([]string, len(DimensionNames)+1)
	for i, dim := range DimensionNames {
		allowed[i] = fmt.Sprintf("%s=<value>", dim)
	}
	allowed[len(DimensionNames)] = "TAG=<tag_value>"
	return allowed
}

func (v FilterValidator) Type() string {
	return "FilterBy"
}

// DimensionOnlyValidator validates only AWS dimensions (no tags)
type DimensionOnlyValidator struct{}

func (v DimensionOnlyValidator) Validate(value string) (map[string]string, error) {
	result := make(map[string]string)

	args := utils.SplitCommaSeparatedString(value)
	for _, arg := range args {
		parts, err := utils.SplitNameValuePair(arg)
		if err != nil {
			return result, err
		}

		dimension := strings.ToUpper(parts[0])
		if !ValidDimensions[dimension] {
			return result, ValidationError{
				Field:   "dimension",
				Value:   parts[0],
				Allowed: DimensionNames,
			}
		}
		result[parts[0]] = parts[1]
	}

	return result, nil
}

func (v DimensionOnlyValidator) AllowedValues() []string {
	allowed := make([]string, len(DimensionNames))
	for i, dim := range DimensionNames {
		allowed[i] = fmt.Sprintf("%s=<value>", dim)
	}
	return allowed
}

func (v DimensionOnlyValidator) Type() string {
	return "DimensionFilter"
}