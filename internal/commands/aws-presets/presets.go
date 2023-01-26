package aws_presets

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	aws3 "github.com/cduggn/ccexplorer/internal/commands/get/aws"
	aws2 "github.com/cduggn/ccexplorer/pkg/helpers"
	"github.com/cduggn/ccexplorer/pkg/service/aws"
	"github.com/spf13/cobra"
)

func AddAWSPresetCommands() *cobra.Command {
	return &cobra.Command{
		Use:   "aws-presets",
		Short: "Preset AWS queries",
		Run: func(cmd *cobra.Command, args []string) {

			var presets = PresetList()

			queries := make([]string, len(presets))
			for i, preset := range presets {
				queries[i] = preset.Alias
			}

			selected := 0
			var prompt = &survey.Select{
				Message: "Choose a preset:",
				Options: queries,
				Description: func(value string, index int) string {
					return presets[index].CommandSyntax
				},
			}
			err := survey.AskOne(prompt, &selected)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			query := presets[selected]
			apiRequest, err := GeneratePresetQuery(query)
			if err != nil {
				err := PresetError{
					msg: fmt.Sprintf("Error generating preset query %v\n", err),
				}
				fmt.Print(err)
			}

			fmt.Printf("Executing %v\n", query.CommandSyntax)

			err = aws3.ExecuteCostCommand(apiRequest)
			if err != nil {
				err := PresetError{
					msg: fmt.Sprintf("Error executing preset query %v\n", err),
				}
				fmt.Print(err)
			}

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
		Granularity:      "MONTHLY",
		ExcludeDiscounts: p.ExcludeDiscounts,
	}, nil
}
