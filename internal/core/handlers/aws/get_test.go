package aws

import (
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	"github.com/cduggn/ccexplorer/internal/core/util"
	"testing"
)

func TestGeneratePresetQuery(t *testing.T) {

	presetCommand := &PresetCommandType{}

	type args struct {
		preset model.PresetParams
	}
	tests := []struct {
		name    string
		args    args
		want    model.CostAndUsageRequestType
		wantErr bool
	}{

		{
			name: "Test 1",
			args: args{
				preset: model.PresetParams{
					Dimension:         []string{"SERVICE", "USAGE_TYPE"},
					Tag:               "Name",
					Filter:            map[string]string{"SERVICE": "Amazon Elastic Compute Cloud - Compute"},
					FilterType:        "SERVICE",
					FilterByDimension: true,
					FilterByTag:       true,
					ExcludeDiscounts:  true,
					Granularity:       "MONTHLY",
				},
			},
			want: model.CostAndUsageRequestType{
				GroupBy: []string{
					"SERVICE",
					"USAGE_TYPE",
				},
				DimensionFilter:            map[string]string{"SERVICE": "Amazon Elastic Compute Cloud - Compute"},
				IsFilterByTagEnabled:       true,
				IsFilterByDimensionEnabled: true,
				Time: model.Time{
					Start: util.DefaultStartDate(util.DayOfCurrentMonth,
						util.SubtractDays),
					End: util.DefaultEndDate(util.Format),
				},
				Granularity:      "MONTHLY",
				ExcludeDiscounts: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := presetCommand.SynthesizeRequest(tt.args.preset)
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
