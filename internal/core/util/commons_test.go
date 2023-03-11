package util

import (
	"testing"
	"time"
)

func TestDefaultEndDate(t *testing.T) {
	type args struct {
		f func(date time.Time) string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Today",
			args: args{
				f: Format,
			},
			want: Format(time.Now()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultEndDate(tt.args.f); got != tt.want {
				t.Errorf("DefaultEndDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultStartDate(t *testing.T) {
	type args struct {
		d func(time time.Time) int
		s func(time time.Time, days int) string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Today",
			args: args{
				d: DayOfCurrentMonth,
				s: SubtractDays,
			},
			want: SubtractDays(time.Now(), DayOfCurrentMonth(time.Now())-1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultStartDate(tt.args.d, tt.args.s); got != tt.want {
				t.Errorf("DefaultStartDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestFormat(t *testing.T) {
	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Today",
			args: args{
				date: time.Now(),
			},
			want: time.Now().Format("2006-01-02"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.args.date); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestLastDayOfMonth(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "LastDayOfMonth",
			want: time.Now().AddDate(0, 1, -1).Format("2006-01-02"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LastDayOfMonth(); got != tt.want {
				t.Errorf("LastDayOfMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}
