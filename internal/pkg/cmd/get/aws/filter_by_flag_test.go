package aws

import (
	"reflect"
	"testing"
)

func TestFilterBy_Set(t *testing.T) {
	type fields struct {
		Dimensions map[string]string
		Tags       []string
	}
	type args struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   FilterBy
	}{
		{
			name: "Valid1",
			args: args{
				value: "SERVICE=Amazon Simple Storage Service,TAG=ApplicationName",
			},
			want: FilterBy{
				Dimensions: map[string]string{
					"SERVICE": "Amazon Simple Storage Service",
				},
				Tags: []string{"ApplicationName"},
			},
		},
		{
			name: "valid2",
			args: args{
				value: "OPERATION=PUTOBJECT," +
					"SERVICE=Amazon Simple Storage Service," +
					"TAG=ApplicationName1",
			},
			want: FilterBy{
				Dimensions: map[string]string{
					"OPERATION": "PUTOBJECT",
					"SERVICE":   "Amazon Simple Storage Service",
				},
				Tags: []string{"ApplicationName1"},
			},
		},
	}
	for _, tt := range tests {
		f := NewFilterBy()
		t.Run(tt.name, func(t *testing.T) {
			if err := f.Set(tt.args.value); err != nil {
				t.Errorf("FilterBy.Set() error = %v", err)
			}
			if !reflect.DeepEqual(f, tt.want) {
				t.Errorf("FilterBy.Set() got = %v, want %v", f, tt.want)
			}
		})

	}
}
