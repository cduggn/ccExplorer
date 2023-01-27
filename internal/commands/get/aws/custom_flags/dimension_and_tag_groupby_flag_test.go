package custom_flags

import (
	"reflect"
	"testing"
)

func TestGroupBySetMethod_ARGValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		wantErr bool
	}{
		{
			name:    "Valid",
			args:    "dimension=SERVICE",
			wantErr: false,
		},
		{
			name:    "Invalid",
			args:    "something=SERVICE",
			wantErr: true,
		},
		{
			name:    "valid",
			args:    "tag=ApplicationName",
			wantErr: false,
		},
		{
			name:    "Valid",
			args:    "DIMENSION=SERVICE",
			wantErr: false,
		},
		{
			name:    "Invalid",
			args:    "SOMETHING=SERVICE",
			wantErr: true,
		},
		{
			name:    "valid",
			args:    "TAG=ApplicationName",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := groupByFlag.Set(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("GroupByFlag.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGroupBySetMethod_ARGParsing(t *testing.T) {
	tests := []struct {
		name string
		args string
		want GroupBy
	}{
		{
			name: "Valid1",
			args: "dimension=REGION",
			want: GroupBy{
				Dimensions: []string{"REGION"},
			},
		},
		{
			name: "Valid2",
			args: "tag=ApplicationName1",
			want: GroupBy{
				Tags: []string{"ApplicationName1"},
			},
		},
		{
			name: "Valid3",
			args: "DIMENSION=REGION",
			want: GroupBy{
				Dimensions: []string{"REGION"},
			},
		},
		{
			name: "Valid4",
			args: "TAG=ApplicationName2",
			want: GroupBy{
				Tags: []string{"ApplicationName2"},
			},
		},
		{
			name: "Valid5",
			args: "TAG=ApplicationName2,Dimension=REGION",
			want: GroupBy{
				Tags:       []string{"ApplicationName2"},
				Dimensions: []string{"REGION"},
			},
		},
		{
			name: "Valid5",
			args: "TAG=ApplicationName2,Dimension=REGION,DIMEnsion=DATABASE_ENGINE",
			want: GroupBy{
				Tags:       []string{"ApplicationName2"},
				Dimensions: []string{"REGION", "DATABASE_ENGINE"},
			},
		},
	}

	for _, tt := range tests {
		var got GroupBy
		t.Run(tt.name, func(t *testing.T) {
			if err := got.Set(tt.args); err != nil {
				t.Errorf("GroupByFlag.Set() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupByFlag.Set() got = %v, want %v", groupByFlag,
					tt.want)
			}
		})
	}
}

func TestSplitByIndividualArgument(t *testing.T) {
	tests := []struct {
		name string
		args string
		want []string
	}{
		{
			name: "Valid1",
			args: "dimension=SERVICE1",
			want: []string{"dimension=SERVICE1"},
		},
		{
			name: "Valid2",
			args: "tag=ApplicationName1",
			want: []string{"tag=ApplicationName1"},
		},
		{
			name: "Valid3",
			args: "SERVICE=AMAZON SIMPLE STORAGE SERVICE",
			want: []string{"SERVICE=AMAZON SIMPLE STORAGE SERVICE"},
		},
		{
			name: "Valid4",
			args: "TAG=ApplicationName2",
			want: []string{"TAG=ApplicationName2"},
		},
		{
			name: "Valid5",
			args: "TAG=ApplicationName2,Dimension=SERVICE2",
			want: []string{"TAG=ApplicationName2", "Dimension=SERVICE2"},
		},
		{
			name: "Valid5",
			args: "TAG=ApplicationName2,Dimension=SERVICE2,DIMEnsion=SERVICE3",
			want: []string{"TAG=ApplicationName2", "Dimension=SERVICE2", "DIMEnsion=SERVICE3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitCommaSeparatedString(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitByIndividualArgument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitIndividualArgument(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				value: "dimension=InstanceId",
			},
			want:    []string{"dimension", "InstanceId"},
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				value: "dimension",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitNameValuePair(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitIndividualArgument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitIndividualArgument() = %v, want %v", got, tt.want)
			}
		})
	}

}
