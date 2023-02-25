package cost_and_usage

import (
	"github.com/cduggn/ccexplorer/internal/commands/get/aws/custom_flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

func ExtractGroupBySelections(cmd *cobra.Command) ([]string, []string) {
	// groupBY dimensions and tags
	groupByValues := cmd.Flags().Lookup("groupBy").Value
	groupBy, _ := groupByValues.(*custom_flags.DimensionAndTagFlag)

	// groupBy TAGs
	var groupByTag []string
	if len(groupBy.Tags) > 0 {
		groupByTag = groupBy.Tags
	}

	// groupBy DIMENSIONS
	var groupByDimension []string
	if len(groupBy.Dimensions) > 0 {
		groupByDimension = groupBy.Dimensions
	}

	return groupByTag, groupByDimension
}

func ExtractFilterBySelection(cmd *cobra.Command) (FilterBySelections, error) {

	var filterSelections FilterBySelections

	filterByValues := cmd.Flags().Lookup("filterBy").Value
	filterBy, _ := filterByValues.(*custom_flags.DimensionAndTagFilterFlag)

	if len(filterBy.Tags) > 1 {
		return FilterBySelections{}, ValidationError{
			Message: "Results can be filtered by a single TAG filter.",
		}
	} else if len(filterBy.Tags) == 1 {
		filterSelections.IsFilterByTag = true
		filterSelections.Tags = filterBy.Tags[0]
	}

	if len(filterBy.Dimensions) > 2 {
		return FilterBySelections{}, ValidationError{
			Message: "Results can be filtered by at most two DIMENSION" +
				" filters.",
		}
	} else if len(filterBy.Dimensions) > 0 {
		filterSelections.IsFilterByDimension = true
		filterSelections.Dimensions = filterBy.Dimensions
	}

	return filterSelections, nil
}

func ExtractStartAndEndDates(cmd *cobra.Command) (string, string, error) {
	start := cmd.Flags().Lookup("startDate").Value.String()
	err := ValidateStartDate(start)
	if err != nil {
		return "", "", err
	}

	end := cmd.Flags().Lookup("endDate").Value.String()
	err = ValidateEndDate(end, start)
	if err != nil {
		return "", "", err
	}

	return start, end, nil
}

func ExtractPrintPreferences(cmd *cobra.Command) PrintOptions {

	var printOptions PrintOptions

	printFormat := cmd.Flags().Lookup("printFormat").Value.String()
	printOptions.Format = strings.ToLower(printFormat)

	sortByDate, _ := cmd.Flags().GetBool("sortByDate")
	printOptions.IsSortByDate = sortByDate

	openAIKey := viper.GetString("open_ai_api_key")
	printOptions.OpenAIKey = openAIKey

	excludeDiscounts, _ := cmd.Flags().GetBool("excludeDiscounts")
	printOptions.ExcludeDiscounts = excludeDiscounts

	granularity := cmd.Flags().Lookup("granularity").Value.String()
	granularity = strings.ToUpper(granularity)
	printOptions.Granularity = granularity

	metric := cmd.Flags().Lookup("metric").Value.String()
	printOptions.Metric = metric

	return printOptions
}
