package aws

import (
	"fmt"
	"github.com/cduggn/cloudcost/internal/pkg/csp/aws"
	"github.com/spf13/cobra"
	"time"
)

var (
	groupBy     []string
	groupByTag  string
	granularity string
	filterBy    string
	rates       []string
	startDate   string
	endDate     string
	report      *aws.CostAndUsageReport
)

func CostSummary(cmd *cobra.Command, args []string) {
	req := NewCostAndUsageRequest(cmd)
	report = aws.GetAWSCostAndUsage(req)
	report.Print()
}

func NewCostAndUsageRequest(cmd *cobra.Command) aws.CostAndUsageRequest {

	dimensions, err := cmd.Flags().GetStringSlice("dimensions")
	if err != nil {
		fmt.Println(err)
	}

	filterBy, _ := cmd.Flags().GetString("filter-by")

	return aws.CostAndUsageRequest{
		Granularity: cmd.Flags().Lookup("granularity").Value.String(),
		GroupBy:     dimensions,
		Tag:         cmd.Flags().Lookup("tags").Value.String(),
		Time: aws.Time{
			Start: cmd.Flags().Lookup("start").Value.String(),
			End:   cmd.Flags().Lookup("end").Value.String(),
		},
		IsFilterEnabled: isFilterEnabled(filterBy),
		TagFilterValue:  filterBy,
		Rates:           rates,
	}

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
