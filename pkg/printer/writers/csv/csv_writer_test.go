package csv

import (
	"github.com/cduggn/ccexplorer/pkg/printer/writers/openai"
	"testing"
)

func TestToCSVString(t *testing.T) {
	csvHeader := "Dimension/Tag,Dimension/Tag,Metric,Granularity,Start,End," +
		"USD Amount,Unit;"

	type args struct {
		data [][]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				data: [][]string{
					{"a", "b", "c"},
					{"d", "e", "f"},
				},
			},
			want: csvHeader + "a,b,c;d,e,f",
		},
		{
			name: "test2",
			args: args{
				data: [][]string{
					{"a", "b", "c"},
					{"d", "e", "f"},
					{"g", "h", "i"},
				},
			},
			want: csvHeader + "a,b,c;d,e,f;g,h,i",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := openai.ConvertToCommaDelimitedString(tt.args.data); got != tt.want {
				t.Errorf("Expected %v, got %v", tt.want, got)
			}
		})
	}
}
