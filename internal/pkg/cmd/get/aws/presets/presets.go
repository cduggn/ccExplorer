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

			apiRequest, err := GeneratePresetQuery(selected)
			if err != nil {
				err := PresetError{
					msg: fmt.Sprintf("Error generating preset query %v\n", err),
				}
				panic(err)
			}

			err = aws2.ExecuteCostCommand(apiRequest)
			if err != nil {
				err := PresetError{
					msg: fmt.Sprintf("Error executing preset query %v\n", err),
				}
				panic(err)
			}

			//fmt.Printf("You chose here %+v", selected)
		},
	}
}

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

func AWSPresets() []PresetParams {

	p := []PresetParams{
		{
			Alias:             "S3 costs grouped by OPERATION",
			Dimension:         []string{"SERVICE", "OPERATION"},
			Tag:               "",
			Filter:            map[string]string{"SERVICE": "Amazon Simple Storage Service"},
			FilterByTag:       false,
			FilterByDimension: true,
		},
		{
			Alias:             "S3 costs grouped by USAGE_TYPE",
			Dimension:         []string{"SERVICE", "USAGE_TYPE"},
			Tag:               "Name",
			Filter:            map[string]string{"SERVICE": "Amazon Simple Storage Service"},
			FilterByTag:       false,
			FilterByDimension: true,
		},
		{
			Alias:             "S3 costs grouped by LINKED_ACCOUNT",
			Dimension:         []string{"SERVICE", "LINKED_ACCOUNT"},
			Tag:               "Name",
			Filter:            map[string]string{"SERVICE": "Amazon Simple Storage Service"},
			FilterByTag:       false,
			FilterByDimension: true,
		},
	}
	return p
}
