package util

import (
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
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
		dayOfCurrentMonth func(time time.Time) int
		subtractDays      func(time time.Time, days int) string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Today",
			args: args{
				dayOfCurrentMonth: DayOfCurrentMonth,
				subtractDays:      SubtractDays,
			},
			want: DefaultStartDate(DayOfCurrentMonth, SubtractDays),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultStartDate(tt.args.dayOfCurrentMonth, tt.args.subtractDays); got != tt.want {
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

func TestServiceToString(t *testing.T) {
	type args struct {
		service model.Service
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ServiceToString",
			args: args{
				service: model.Service{
					Name:  "test",
					Start: "2021-01-01",
					End:   "2021-01-02",
					Metrics: []model.Metrics{
						{
							Name:          "test",
							Amount:        "0.10",
							Unit:          "USD",
							UsageQuantity: 0.10,
						},
					},
					Keys: []string{
						"SERVICE", "USAGE_QUANTITY",
					},
				},
			},
			want: "SERVICE,USAGE_QUANTITY,2021-01-01,2021-01-02,test,0.10," +
				"USD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			S := ServiceToString(tt.args.service)
			if S != tt.want {
				t.Errorf("ServiceToString() Got: %v, want: %v", S, tt.want)
			}
		})
	}
}
