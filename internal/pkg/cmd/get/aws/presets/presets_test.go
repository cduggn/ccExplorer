package presets

import (
	aws2 "github.com/cduggn/ccexplorer/internal/pkg/cmd/get/aws"
	"github.com/cduggn/ccexplorer/internal/pkg/service/aws"
	"testing"
)

func TestGeneratePresetQuery(t *testing.T) {
	type args struct {
		preset PresetParams
	}
	tests := []struct {
		name    string
		args    args
		want    aws.CostAndUsageRequestType
		wantErr bool
	}{

		{
			name: "Test 1",
			args: args{
				preset: PresetParams{
					Dimension:         []string{"SERVICE", "USAGE_TYPE"},
					Tag:               "Name",
					Filter:            map[string]string{"SERVICE": "Amazon Elastic Compute Cloud - Compute"},
					FilterType:        "SERVICE",
					FilterByDimension: true,
					FilterByTag:       true,
				},
			},
			want: aws.CostAndUsageRequestType{
				GroupBy: []string{
					"SERVICE",
					"USAGE_TYPE",
				},
				DimensionFilter:            map[string]string{"SERVICE": "Amazon Elastic Compute Cloud - Compute"},
				IsFilterByTagEnabled:       true,
				IsFilterByDimensionEnabled: true,
				Time: aws.Time{
					Start: aws2.DefaultStartDate(aws2.DayOfCurrentMonth, aws2.SubtractDays),
					End:   aws2.DefaultEndDate(aws2.Format),
				},
				Granularity: "MONTHLY",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePresetQuery(tt.args.preset)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePresetQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equals(tt.want) {
				t.Errorf("GeneratePresetQuery() got = %v, want %v", got,
					tt.want)
			}
		})
	}

}
