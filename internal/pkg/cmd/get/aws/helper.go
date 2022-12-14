package aws

import "time"

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

func isFilterEnabled(filterByTag string) bool {
	if filterByTag != "" {
		return true
	} else {
		return false
	}
}

func isFilterDimensionEnabled(dimensionFilterMap map[string]string) bool {
	if len(dimensionFilterMap) > 0 {
		return true
	} else {
		return false
	}
}
