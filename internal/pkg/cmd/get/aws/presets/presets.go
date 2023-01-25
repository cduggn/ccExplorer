package presets

import (
	"fmt"
	aws2 "github.com/cduggn/ccexplorer/internal/pkg/cmd/get/aws"
	"github.com/cduggn/ccexplorer/internal/pkg/service/aws"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var templates = &promptui.SelectTemplates{
	Label:    "{{ . }}?",
	Active:   "\U0001F336 {{ .Alias | cyan }} ",
	Inactive: "  {{ .Alias | cyan }} ",
	Selected: "\U0001F336 {{ .Alias | red | cyan }}",
	Details: `
--------- AWS Preset Queries ----------
{{ "Alias:" | faint }}	{{ .Alias }}
{{ "GroupBy:" | faint }}	{{ .Dimension }}
{{ "FilterBy:" | faint }}	{{ .Filter }}
{{ "FilterByDimension:" | faint }}	{{
.FilterByDimension }}
{{ "FilterByTag:" | faint }}	{{ .FilterByTag }}
`,
}

func AddAWSPresetCommands() *cobra.Command {
	return &cobra.Command{
		Use:   "aws-presets",
		Short: "Load preset queries",
		Run: func(cmd *cobra.Command, args []string) {
			prompt := promptui.Select{
				Label:     "Select a preset:",
				Items:     AWSPresets(),
				Templates: templates,
			}

			val, _, err := prompt.Run()
			if err != nil {
				err := PresetError{
					msg: fmt.Sprintf("Prompt failed %v\n", err),
				}
				panic(err)
			}
			presets := prompt.Items.([]PresetParams)
			selected := presets[val]

			fmt.Printf("You chose here %+v", selected)
		},
	}
}

func AWSPresets() []PresetParams {

	p := []PresetParams{
		{
			Alias:             "S3 Operation Costs",
			Dimension:         []string{"SERVICE", "OPERATION"},
			Tag:               "",
			Filter:            map[string]string{"SERVICE": "Amazon Simple Storage Service"},
			FilterByTag:       false,
			FilterByDimension: false,
		},
		{
			Alias:             "S3 Operation Costs by Tag",
			Dimension:         []string{"SERVICE", "USAGE_TYPE"},
			Tag:               "Name",
			Filter:            map[string]string{"SERVICE": "Amazon Simple Storage Service"},
			FilterByTag:       true,
			FilterByDimension: false,
		},
	}
	return p
}

//return []Preset{
//{Name: "S3 Operation Costs", ID: 1},
//{Name: "DynamoDB Operation Costs", ID: 2},
//{Name: "Linked Accounts Costs", ID: 3},
////{Name: "AWS Cost Explorer", ID: 4},
////{Name: "AWS Cost Explorer", ID: 5},
////{Name: "AWS Cost Explorer", ID: 6},
////{Name: "AWS Cost Explorer", ID: 7},
////{Name: "AWS Cost Explorer", ID: 8},
////{Name: "AWS Cost Explorer", ID: 9},
////{Name: "AWS Cost Explorer", ID: 10},
//}
//}

func GeneratePresetQuery(p PresetParams) (aws.CostAndUsageRequestType, error) {
	return aws.CostAndUsageRequestType{
		GroupBy:                    p.Dimension,
		DimensionFilter:            p.Filter,
		IsFilterByTagEnabled:       p.FilterByTag,
		IsFilterByDimensionEnabled: p.FilterByDimension,
		Time: aws.Time{
			Start: aws2.DefaultStartDate(aws2.DayOfCurrentMonth, aws2.SubtractDays),
			End:   aws2.DefaultEndDate(aws2.Format),
		},
		Granularity: "MONTHLY",
	}, nil
}
