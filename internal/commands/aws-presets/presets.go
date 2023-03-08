package aws_presets

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	aws2 "github.com/cduggn/ccexplorer/internal/commands/get/aws"
	aws3 "github.com/cduggn/ccexplorer/internal/commands/get/aws/cost_and_usage"
	"github.com/cduggn/ccexplorer/pkg/domain/model"
	"github.com/spf13/cobra"
)

func (e PresetError) Error() string {
	return e.msg
}

func AddAWSPresetCommands() *cobra.Command {
	return &cobra.Command{
		Use:   "run-query",
		Short: "Predefined AWS Cost and Usage queries",
		Run:   runCommand,
	}
}

func runCommand(cmd *cobra.Command, args []string) {
	var presets = PresetList()
	optionsNameList := PromptQueryList(presets)

	var prompt = Prompt(optionsNameList)
	selection := 0
	err := survey.AskOne(prompt, &selection)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	selectedOption := Selected(presets, selection)
	apiRequest, err := SynthesizeQuery(selectedOption)
	if err != nil {
		err := PresetError{
			msg: fmt.Sprintf("Error synthesizing query %v\n",
				err),
		}
		fmt.Print(err)
	}

	DisplaySynthesizedQuery(selectedOption)
	err = execute(apiRequest)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func PromptQueryList(p []PresetParams) []string {
	queries := make([]string, len(p))
	for i, preset := range p {
		queries[i] = preset.Alias
	}
	return queries
}

func Prompt(o []string) *survey.Select {
	return &survey.Select{
		Message: "Choose a query to execute:",
		Options: o,
	}
}

func Selected(p []PresetParams, s int) PresetParams {
	return p[s]
}

func SynthesizeQuery(p PresetParams) (model.CostAndUsageRequestType,
	error) {
	return model.CostAndUsageRequestType{
		GroupBy:                    p.Dimension,
		DimensionFilter:            p.Filter,
		IsFilterByTagEnabled:       p.FilterByTag,
		IsFilterByDimensionEnabled: p.FilterByDimension,
		Time: model.Time{
			Start: aws2.DefaultStartDate(aws2.DayOfCurrentMonth, aws2.SubtractDays),
			End:   aws2.DefaultEndDate(aws2.Format),
		},
		Granularity:      p.Granularity,
		ExcludeDiscounts: p.ExcludeDiscounts,
		PrintFormat:      p.PrintFormat,
		SortByDate:       false,
		Metrics:          p.Metric,
	}, nil
}

func DisplaySynthesizedQuery(p PresetParams) {
	fmt.Println("")
	fmt.Printf("Synthesized Query: %v \n", p.CommandSyntax)
	fmt.Println("")
}

func execute(q model.CostAndUsageRequestType) error {
	err := aws3.ExecutePreset(q)
	if err != nil {
		err := PresetError{
			msg: fmt.Sprintf("Error executing preset query %v\n", err),
		}
		return err
	}
	return nil
}
