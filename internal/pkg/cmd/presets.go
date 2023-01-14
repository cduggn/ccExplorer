package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var templates = &promptui.SelectTemplates{
	Label:    "{{ . }}?",
	Active:   "\U0001F336 {{ .Name | cyan }} ({{ .HeatUnit | red }})",
	Inactive: "  {{ .Name | cyan }} ({{ .HeatUnit | red }})",
	Selected: "\U0001F336 {{ .Name | red | cyan }}",
	Details: `
--------- Pepper ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Heat Unit:" | faint }}	{{ .HeatUnit }}
{{ "Peppers:" | faint }}	{{ .Peppers }}`,
}

type UsageRequest struct {
	Name     string
	HeatUnit int
	Peppers  int
}

func AddPresetCommands() *cobra.Command {
	return &cobra.Command{
		Use:   "presets",
		Short: "Load preset queries",
		Run: func(cmd *cobra.Command, args []string) {
			prompt := promptui.Select{
				Label: "Select an preset to use:",
				Items: []UsageRequest{
					{Name: "Bell Pepper", HeatUnit: 0, Peppers: 0},
					{Name: "Banana Pepper", HeatUnit: 100, Peppers: 1},
					{Name: "Poblano", HeatUnit: 1000, Peppers: 2},
					{Name: "Jalapeño", HeatUnit: 3500, Peppers: 3},
					{Name: "Aleppo", HeatUnit: 10000, Peppers: 4},
					{Name: "Tabasco", HeatUnit: 30000, Peppers: 5},
					{Name: "Malagueta", HeatUnit: 50000, Peppers: 6},
					{Name: "Habanero", HeatUnit: 100000, Peppers: 7},
					{Name: "Red Savina Habanero", HeatUnit: 350000, Peppers: 8},
					{Name: "Dragon’s Breath", HeatUnit: 855000, Peppers: 9},
				},
				Templates: templates,
			}

			_, result, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}
			//awsClient := aws.NewAPIClient()
			//aws.RightSizingRecommendationS3(context.Background(), awsClient.Client)
			fmt.Printf("You chose %q\n", result)
		},
	}
}
