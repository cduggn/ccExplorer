package aws

import (
	"testing"
	"time"
)

// test for DayOfCurrentMonth function
func TestDayOfCurrentMonth(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Today",
			args: args{
				t: time.Now(),
			},
			want: time.Now().Day(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DayOfCurrentMonth(tt.args.t); got != tt.want {
				t.Errorf("DayOfCurrentMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test subtractDays function
func TestSubtractDays(t *testing.T) {
	type args struct {
		today time.Time
		days  int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Today",
			args: args{
				today: time.Now(),
				days:  1,
			},
			want: time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubtractDays(tt.args.today, tt.args.days); got != tt.want {
				t.Errorf("SubtractDays() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test DefaultStartDate function
func TestDefaultStartDate(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Today",
			want: time.Now().AddDate(0, 0, -time.Now().Day()+1).Format("2006-01-02"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultStartDate(DayOfCurrentMonth, SubtractDays); got != tt.want {
				t.Errorf("DefaultStartDate() = %v, want %v", got, tt.want)
			} else {
				t.Logf("DefaultStartDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
