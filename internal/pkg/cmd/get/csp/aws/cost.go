package aws

import (
	"github.com/cduggn/cloudcost/internal/pkg/service/aws"
	"github.com/spf13/cobra"
	"time"
)

var (
	groupBy          []string
	groupByTag       string
	granularity      string
	filterBy         string
	rates            []string
	startDate        string
	endDate          string
	report           *aws.CostAndUsageReport
	withoutDiscounts bool
)

func CostAndUsageSummary(cmd *cobra.Command, args []string) error {

	req, err := NewCostAndUsageRequest(cmd)
	if err != nil {
		return err
	}

	report, err = aws.GetCostAndUsage(req)
	if err != nil {
		return err
	}
	report.PrintCostAndUsageReport()

	return nil
}

func NewCostAndUsageRequest(cmd *cobra.Command) (aws.CostAndUsageRequestType, error) {

	dimensions, err := cmd.Flags().GetStringSlice("dimensions")
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}
	err = ValidateDimension(dimensions)
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}

	tag := cmd.Flags().Lookup("tags").Value.String()
	err = ValidateTag(tag, dimensions)
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}

	filter, _ := cmd.Flags().GetString("filter-by")
	err = ValidateFilterBy(filter, tag)
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}

	start := cmd.Flags().Lookup("start").Value.String()
	err = ValidateStartDate(start)
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}

	end := cmd.Flags().Lookup("end").Value.String()
	err = ValidateEndDate(end, start)
	if err != nil {
		return aws.CostAndUsageRequestType{}, err
	}

	excludeDiscounts, _ := cmd.Flags().GetBool("exclude-discounts")
	interval := cmd.Flags().Lookup("granularity").Value.String()

	return aws.CostAndUsageRequestType{
		Granularity: interval,
		GroupBy:     dimensions,
		Tag:         tag,
		Time: aws.Time{
			Start: start,
			End:   end,
		},
		IsFilterEnabled:  isFilterEnabled(filterBy),
		TagFilterValue:   filter,
		Rates:            rates,
		ExcludeDiscounts: excludeDiscounts,
	}, nil

}

func DefaultEndDate(f func(date time.Time) string) string {
	return f(time.Now())
}

func Format(date time.Time) string {
	return date.Format("2006-01-02")
}

func DefaultStartDate(d func(time time.Time) int, s func(time time.Time, days int) string) string {
	today := time.Now()
	dayOfMonth := d(today)
	return s(today, dayOfMonth-1) // subtract 1 to get the first day of the month
}

func DayOfCurrentMonth(time time.Time) int {
	return time.Day()
}

func SubtractDays(today time.Time, days int) string {
	return today.AddDate(0, 0, -days).Format("2006-01-02")
}

func isFilterEnabled(filterBy string) bool {
	if filterBy != "" {
		return true
	} else {
		return false
	}
}
