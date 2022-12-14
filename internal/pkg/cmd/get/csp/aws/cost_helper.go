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

	dimensions, err := cmd.Flags().GetStringSlice("group-by-dimension")
	if err != nil {
		fmt.Println(err)
	}

	filterBy, _ := cmd.Flags().GetString("filter-by")

	return aws.CostAndUsageRequest{
		Granularity: cmd.Flags().Lookup("granularity").Value.String(),
		GroupBy:     dimensions,
		Tag:         cmd.Flags().Lookup("group-by-tag").Value.String(),
		Time: aws.Time{
			Start: cmd.Flags().Lookup("start-date").Value.String(),
			End:   cmd.Flags().Lookup("end-date").Value.String(),
		},
		IsFilterEnabled: isFilterEnabled(filterBy),
		TagFilterValue:  filterBy,
		Rates:           rates,
	}

}

func Time() time.Time {
	return time.Now()
}

func Today() string {
	return Format(Time())
}

func Format(date time.Time) string {
	return date.Format("2006-01-02")
}

func PastMonth() string {
	today := Time()
	monthAgo := SubtractDays(today, 30)
	return Format(monthAgo)
}

func SubtractDays(today time.Time, days int) time.Time {
	return today.AddDate(0, 0, -days)
}

func isFilterEnabled(filterBy string) bool {
	if filterBy != "" {
		return true
	} else {
		return false
	}
}
