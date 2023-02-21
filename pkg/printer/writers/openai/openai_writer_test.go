package openai

import (
	"html/template"
	"strings"
	"testing"
)

var (
	tableTemplate = `
<table>
	<thead>
		<tr>
			<th>Dimension/Tag</th>
			<th>Dimension/Tag</th>
			<th>Start</th>
			<th>End</th>
			<th>USD Amount</th>
		</tr>
	</thead>
	<tbody>
		{{range .}}
			<tr>
			<td>{{.Dimension}}</td>
			<td>{{.Tag}}</td>
			<td>{{.Start}}</td>
			<td>{{.End}}</td>
			<td>{{.USDAmount}}</td>
			</tr>
		{{end}}
	</tbody>
</table>
`
	trainingTemplateOutput = `
<table>
	<thead>
		<tr>
			<th>Dimension/Tag</th>
			<th>Dimension/Tag</th>
			<th>Start</th>
			<th>End</th>
			<th>USD Amount</th>
		</tr>
	</thead>
	<tbody>
		
			<tr>
			<td>a</td>
			<td>b</td>
			<td>e</td>
			<td>f</td>
			<td>0.0</td>
			</tr>
		
	</tbody>
</table>
`
)

func TestNewTrainingExample(t *testing.T) {
	type args struct {
		t *template.Template
		s []TrainingData
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				t: template.Must(template.New("table1").Parse(tableTemplate)),
				s: []TrainingData{
					{
						Dimension: "a",
						Tag:       "b",
						Start:     "e",
						End:       "f",
						USDAmount: "0.0",
					},
				},
			},
			want: trainingTemplateOutput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, _ := CreateTrainingData(tt.args.t, tt.args.s)
			want := removeSpacingAndTabs(tt.want)
			if removeSpacingAndTabs(got) != want {
				//t.Errorf("NewTrainingExample() = %v, want %v", got, tt.want)
				t.Errorf("Expected %v, got %v", want, removeSpacingAndTabs(got))
			}
		})
	}
}

func removeSpacingAndTabs(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", "")
}

func TestBuildPrompt(t *testing.T) {
	type args struct {
		rows [][]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				rows: [][]string{
					{"a", "b", "c", "d", "e", "f", "0.0", "h"},
				},
			},
			want: trainingTemplateOutput,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := removeSpacingAndTabs(tt.want)
			got := removeSpacingAndTabs(BuildCostAndUsagePromptText(tt.args.rows))
			if got != want {
				t.Errorf("Expected %v, got %v", want, got)
			}
		})
	}
}
