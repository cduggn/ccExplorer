package util

import (
	"fmt"
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	"os"
	"strconv"
	"strings"
	"time"
)

func ConvertToFloat(amount string) float64 {
	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func ReturnIfPresent(s []string) string {
	if len(s) == 1 {
		return ""
	} else {
		return s[1]
	}
}

func ToPrintWriterType(s string) model.PrintWriterType {
	switch s {
	case "csv":
		return model.CSV
	case "stdout":
		return model.Stdout
	case "chart":
		return model.Chart
	case "gpt3":
		return model.OpenAPI
	default:
		return model.Stdout
	}
}

func NewFile(dir string, file string) (*os.File, error) {
	filePath := BuildOutputFilePath(dir, file)
	return os.Create(filePath)
}

func BuildOutputFilePath(dir string, fileName string) string {
	return dir + "/" + fileName
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

func LastDayOfMonth() string {
	return time.Now().AddDate(0, 1, -1).Format("2006-01-02")
}

func SortByFn(sortByDate bool) string {
	if sortByDate {
		return "date"
	}
	return "cost"
}

func SplitCommaSeparatedString(value string) []string {
	var args []string
	if strings.Contains(value, ",") {
		args = strings.Split(value, ",")
	} else {
		args = []string{value}
	}
	return args
}

func SplitNameValuePair(value string) ([]string, error) {
	parts := strings.Split(value, "=")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid group by flag: %s", value)
	}
	return parts, nil
}
