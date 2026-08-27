package cli

import (
	"testing"
	"time"

	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateStartDate(t *testing.T) {
	future := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{"valid calendar date", "2024-01-01", ""},
		{"valid ISO 8601 timestamp", "2024-01-01T00:00:00Z", ""},
		{"empty", "", "must be specified"},
		{"future", future, "before today"},
		// Regression: time.Parse's error used to be discarded, so garbage
		// became the zero time, which is before today, and validation passed.
		// The query then failed at the AWS API with an opaque message.
		{"not a date at all", "garbage", "Invalid date"},
		{"impossible calendar date", "2024-13-45", "Invalid date"},
		{"US ordering", "01/02/2024", "Invalid date"},
		{"partial", "2024-01", "Invalid date"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateStartDate(tc.input)
			if tc.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestValidateEndDate(t *testing.T) {
	tests := []struct {
		name       string
		end, start string
		wantErr    string
	}{
		{"valid range", "2024-02-01", "2024-01-01", ""},
		{"same day", "2024-01-01", "2024-01-01", ""},
		{"end before start", "2024-01-01", "2024-02-01", "not be before start"},
		{"empty end", "", "2024-01-01", "must be specified"},
		{"malformed end", "nonsense", "2024-01-01", "Invalid date"},
		// A malformed start reaching this function used to silently become
		// the zero time, making every end date look valid.
		{"malformed start", "2024-02-01", "nonsense", "Invalid date"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEndDate(tc.end, tc.start)
			if tc.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

// -m HOURLY with the default flag values could only ever fail at the API:
// the defaults are plain calendar dates and Cost Explorer requires a time
// component for hourly data.
func TestValidateGranularityAgainstDates(t *testing.T) {
	tests := []struct {
		name                 string
		interval, start, end string
		wantErr              string
	}{
		{"monthly with dates", "MONTHLY", "2024-01-01", "2024-02-01", ""},
		{"daily with dates", "DAILY", "2024-01-01", "2024-02-01", ""},
		{"hourly with timestamps", "HOURLY",
			"2024-01-01T00:00:00Z", "2024-01-02T00:00:00Z", ""},
		{"hourly with plain dates", "HOURLY", "2024-01-01", "2024-01-02",
			"requires an ISO 8601"},
		{"hourly with only start timestamped", "HOURLY",
			"2024-01-01T00:00:00Z", "2024-01-02", "requires an ISO 8601"},
		{"hourly with only end timestamped", "HOURLY",
			"2024-01-01", "2024-01-02T00:00:00Z", "requires an ISO 8601"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateGranularityAgainstDates(tc.interval, tc.start, tc.end)
			if tc.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestValidateAnomaliesDateRange(t *testing.T) {
	tests := []struct {
		name       string
		start, end string
		wantErr    string
	}{
		{"valid 90 day range", "2024-01-01", "2024-03-30", ""},
		{"valid short range", "2024-01-01", "2024-01-10", ""},
		{"exceeds 90 days", "2024-01-01", "2024-04-15", "must not exceed 90 days"},
		{"end before start", "2024-02-01", "2024-01-01", "not be before start"},
		{"malformed start", "nonsense", "2024-01-10", "Invalid date"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateAnomaliesDateRange(tc.start, tc.end)
			if tc.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestIsValidFeedbackType(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty is valid (unset filter)", "", true},
		{"YES", "YES", true},
		{"NO", "NO", true},
		{"PLANNED_ACTIVITY", "PLANNED_ACTIVITY", true},
		{"unknown value", "MAYBE", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, IsValidFeedbackType(tc.input))
		})
	}
}

func TestValidateAnomaliesInput(t *testing.T) {
	base := types.AnomaliesCommandLineInput{
		Start: "2024-01-01",
		End:   "2024-01-10",
	}

	t.Run("accepts a well formed request", func(t *testing.T) {
		assert.NoError(t, ValidateAnomaliesInput(base))
	})

	t.Run("rejects an invalid feedback value", func(t *testing.T) {
		in := base
		in.Feedback = "MAYBE"
		assert.ErrorContains(t, ValidateAnomaliesInput(in), "Invalid feedback value")
	})

	t.Run("rejects a range over 90 days", func(t *testing.T) {
		in := base
		in.End = "2024-06-01"
		assert.ErrorContains(t, ValidateAnomaliesInput(in), "must not exceed 90 days")
	})
}

func baseInput() types.CommandLineInput {
	return types.CommandLineInput{
		Interval:    "MONTHLY",
		PrintFormat: "stdout",
		Metrics:     []string{"UnblendedCost"},
		Start:       "2024-01-01",
		End:         "2024-02-01",
	}
}

func TestValidateInput(t *testing.T) {
	t.Run("accepts a well formed request", func(t *testing.T) {
		assert.NoError(t, ValidateInput(baseInput()))
	})

	t.Run("rejects bad granularity", func(t *testing.T) {
		in := baseInput()
		in.Interval = "WEEKLY"
		assert.ErrorContains(t, ValidateInput(in), "Invalid granularity")
	})

	t.Run("rejects bad print format", func(t *testing.T) {
		in := baseInput()
		in.PrintFormat = "yaml"
		assert.ErrorContains(t, ValidateInput(in), "Invalid print format")
	})

	t.Run("rejects hourly with plain dates", func(t *testing.T) {
		in := baseInput()
		in.Interval = "HOURLY"
		assert.ErrorContains(t, ValidateInput(in), "requires an ISO 8601")
	})

	// Metrics[0] was indexed unconditionally.
	t.Run("rejects empty metrics without panicking", func(t *testing.T) {
		in := baseInput()
		in.Metrics = nil
		assert.NotPanics(t, func() {
			assert.ErrorContains(t, ValidateInput(in), "metric must be specified")
		})
	})

	t.Run("rejects pinecone as a print format", func(t *testing.T) {
		in := baseInput()
		in.PrintFormat = "pinecone"
		assert.ErrorContains(t, ValidateInput(in), "Invalid print format")
	})
}
