package aws

import (
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	"time"
)

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func ValidateInput(input model.CommandLineInput) error {

	isValidGranularity := IsValidGranularity(input.Interval)
	if !isValidGranularity {
		return ValidationError{
			Message: "Invalid granularity. Valid values are: DAILY, MONTHLY, HOURLY",
		}
	}

	isValidPrintFormat := IsValidPrintFormat(input.PrintFormat)
	if !isValidPrintFormat {
		return ValidationError{
			Message: "Invalid print format. " +
				"Please use one of the following: stdout, csv, chart, pinecone",
		}
	}

	if input.PrintFormat == "pinecone" && input.OpenAIAPIKey == "" {
		return ValidationError{
			Message: "OpenAI API key not set. " +
				"Please set the OPEN_AI_API_KEY in the config file or environment variable",
		}
	}

	if input.PrintFormat == "pinecone" && input.OpenAIAPIKey != "" {
		if HasAccountInformation(input.GroupByDimension) {
			return ValidationError{
				Message: "Cannot use Pinecone with account information. " +
					"Please remove the account dimension",
			}
		}
	}

	IsValid := IsValidMetric(input.Metrics[0])
	if !IsValid {
		return ValidationError{
			Message: "Invalid metric. " +
				"Please use one of the following: AmortizedCost, BlendedCost, NetAmortizedCost, NetUnblendedCost, NormalizedUsageAmount, UnblendedCost, UsageQuantity",
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

func IsValidPrintFormat(f string) bool {
	return f == "stdout" || f == "csv" || f == "chart" || f == "pinecone"
}

func IsValidGranularity(g string) bool {
	return g == "DAILY" || g == "MONTHLY" || g == "HOURLY"
}

func IsValidMetric(m string) bool {
	return m == "AmortizedCost" || m == "BlendedCost" || m == "NetAmortizedCost" ||
		m == "NetUnblendedCost" || m == "NormalizedUsageAmount" || m == "UnblendedCost" ||
		m == "UsageQuantity"
}

func HasAccountInformation(groupBy []string) bool {
	for _, v := range groupBy {
		if v == "LINKED_ACCOUNT" {
			return true
		}
	}
	return false

}
