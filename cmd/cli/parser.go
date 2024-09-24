package cli

import (
	"github.com/cduggn/ccexplorer/cmd/cli/flags"
	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/spf13/viper"
	"strings"
)

func (c *CostCommandType) ExtractGroupBySelections() ([]string, []string) {
	// groupBY dimensions and tags
	groupByValues := c.Cmd.Flags().Lookup("groupBy").Value
	groupBy, _ := groupByValues.(*flags.DimensionAndTagFlag)

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

func (c *CostCommandType) ExtractFilterBySelection() (types.FilterBySelections, error) {

	var filterSelections types.FilterBySelections

	filterByValues := c.Cmd.Flags().Lookup("filterBy").Value
	filterBy, _ := filterByValues.(*flags.DimensionAndTagFilterFlag)

	if len(filterBy.Tags) > 1 {
		return types.FilterBySelections{}, ValidationError{
			Message: "Results can be filtered by a single TAG filter.",
		}
	} else if len(filterBy.Tags) == 1 {
		filterSelections.IsFilterByTag = true
		filterSelections.Tags = filterBy.Tags[0]
	}

	if len(filterBy.Dimensions) > 2 {
		return types.FilterBySelections{}, ValidationError{
			Message: "Results can be filtered by at most two DIMENSION" +
				" filters.",
		}
	} else if len(filterBy.Dimensions) > 0 {
		filterSelections.IsFilterByDimension = true
		filterSelections.Dimensions = filterBy.Dimensions
	}

	return filterSelections, nil
}

func (c *CostCommandType) ExtractStartAndEndDates() (
	string, string, error) {
	start := c.Cmd.Flags().Lookup("startDate").Value.String()
	err := ValidateStartDate(start)
	if err != nil {
		return "", "", err
	}

	end := c.Cmd.Flags().Lookup("endDate").Value.String()
	err = ValidateEndDate(end, start)
	if err != nil {
		return "", "", err
	}

	return start, end, nil
}

func (c *CostCommandType) ExtractPrintPreferences() types.PrintOptions {

	var printOptions types.PrintOptions

	printFormat := c.Cmd.Flags().Lookup("printFormat").Value.String()
	printOptions.Format = strings.ToLower(printFormat)

	sortByDate, _ := c.Cmd.Flags().GetBool("sortByDate")
	printOptions.IsSortByDate = sortByDate

	openAIKey := viper.GetString("openai_api_key")
	printOptions.OpenAIKey = openAIKey

	pineconeAPIKey := viper.GetString("PINECONE_API_KEY")
	printOptions.PineconeAPIKey = pineconeAPIKey

	pineconeIndex := viper.GetString("PINECONE_INDEX")
	printOptions.PineconeIndex = pineconeIndex

	excludeDiscounts, _ := c.Cmd.Flags().GetBool("excludeDiscounts")
	printOptions.ExcludeDiscounts = excludeDiscounts

	granularity := c.Cmd.Flags().Lookup("granularity").Value.String()
	granularity = strings.ToUpper(granularity)
	printOptions.Granularity = granularity

	metric := c.Cmd.Flags().Lookup("metric").Value.String()
	printOptions.Metric = metric

	return printOptions
}
