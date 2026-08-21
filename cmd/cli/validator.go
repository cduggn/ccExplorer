package cli

import (
	"github.com/cduggn/ccexplorer/internal/types"
	"time"
)

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func ValidateInput(input types.CommandLineInput) error {

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

	if input.PrintFormat == "pinecone" {
		for _, missing := range []struct{ value, env string }{
			{input.OpenAIAPIKey, "OPENAI_API_KEY"},
			{input.PineconeAPIKey, "PINECONE_API_KEY"},
			{input.PineconeIndex, "PINECONE_INDEX"},
		} {
			if missing.value == "" {
				return ValidationError{
					Message: missing.env + " is not set. It is required for " +
						"-p pinecone; export it before running the query",
				}
			}
		}

		if HasAccountInformation(input.GroupByDimension) {
			return ValidationError{
				Message: "Cannot use Pinecone with account information. " +
					"Please remove the account dimension",
			}
		}
	}

	if err := validateGranularityAgainstDates(input.Interval, input.Start,
		input.End); err != nil {
		return err
	}

	if len(input.Metrics) == 0 {
		return ValidationError{
			Message: "A cost metric must be specified",
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

const (
	dateLayout      = "2006-01-02"
	timestampLayout = time.RFC3339
)

// ParseDate accepts either a plain calendar date or an ISO 8601 timestamp,
// which is what Cost Explorer requires for HOURLY granularity. The second
// return value reports whether the value carried a time component.
func ParseDate(value string) (time.Time, bool, error) {
	if t, err := time.Parse(timestampLayout, value); err == nil {
		return t, true, nil
	}
	t, err := time.Parse(dateLayout, value)
	if err != nil {
		return time.Time{}, false, ValidationError{
			Message: "Invalid date " + value +
				". Use YYYY-MM-DD, or an ISO 8601 timestamp such as " +
				"2006-01-02T15:04:05Z for HOURLY granularity",
		}
	}
	return t, false, nil
}

func ValidateStartDate(startDate string) error {
	if startDate == "" {
		return ValidationError{
			Message: "Start date must be specified",
		}
	}

	start, _, err := ParseDate(startDate)
	if err != nil {
		return err
	}

	if start.After(time.Now()) {
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

	end, _, err := ParseDate(endDate)
	if err != nil {
		return err
	}

	if end.After(time.Now()) {
		return ValidationError{
			Message: "End date must be before today's date",
		}
	}

	start, _, err := ParseDate(startDate)
	if err != nil {
		return err
	}

	if end.Before(start) {
		return ValidationError{
			Message: "End date must not be before start date",
		}
	}

	return nil
}

// validateGranularityAgainstDates rejects the combination Cost Explorer
// itself rejects: HOURLY granularity requires both bounds to carry a time
// component, and the flag defaults are plain calendar dates, so `-m HOURLY`
// with default dates could only ever fail at the API.
func validateGranularityAgainstDates(interval, start, end string) error {
	if interval != "HOURLY" {
		return nil
	}

	for label, value := range map[string]string{"start": start, "end": end} {
		_, hasTime, err := ParseDate(value)
		if err != nil {
			return err
		}
		if !hasTime {
			return ValidationError{
				Message: "HOURLY granularity requires an ISO 8601 " +
					label + " date with a time component, for example " +
					"2006-01-02T15:04:05Z (got " + value + ")",
			}
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
