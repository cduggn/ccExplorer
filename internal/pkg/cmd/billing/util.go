package billing

import "time"

func Time() time.Time {
	return time.Now()
}

func Format(date time.Time) string {
	return date.Format("2006-01-02")
}

func SubtractDays(today time.Time, days int) time.Time {
	return today.AddDate(0, 0, -days)
}
