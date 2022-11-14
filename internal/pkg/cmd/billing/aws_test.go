package billing

import (
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	today := time.Now()
	updatedTime := today.AddDate(0, 0, -4)
	updatedTime.Format("2006-01-02")

	if updatedTime.Before(today) {
		t.Logf("updatedTime %s is before today %s", updatedTime.Format("2006-01-02"), today.Format("2006-01-02"))
	} else {
		t.Errorf("updatedTime %s is not before today %s", updatedTime.Format("2006-01-02"), today.Format("2006-01-02"))
	}
}
