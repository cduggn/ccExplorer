package cli

import (
	"github.com/cduggn/ccexplorer/internal/flags"
	"github.com/cduggn/ccexplorer/internal/types"
	"github.com/spf13/viper"
	"strings"
)

func (c *CostCommandType) ExtractGroupBySelections() ([]string, []string) {
	// groupBY dimensions and tags
	groupByValues := c.Cmd.Flags().Lookup("groupBy").Value
	groupBy, _ := groupByValues.(*flags.GroupByFlag)
	
	groupByData := groupBy.Value()

	// groupBy TAGs
	var groupByTag []string
	if len(groupByData.Tags) > 0 {
		groupByTag = groupByData.Tags
	}

	// groupBy DIMENSIONS
	var groupByDimension []string
	if len(groupByData.Dimensions) > 0 {
		groupByDimension = groupByData.Dimensions
	}

	return groupByTag, groupByDimension
}

func (c *CostCommandType) ExtractFilterBySelection() (types.FilterBySelections, error) {

	var filterSelections types.FilterBySelections

	filterByValues := c.Cmd.Flags().Lookup("filterBy").Value
	filterBy, _ := filterByValues.(*flags.FilterByFlag)
	
	filterByData := filterBy.Value()

	if len(filterByData.Tags) > 1 {
		return types.FilterBySelections{}, ValidationError{
			Message: "Results can be filtered by a single TAG filter.",
		}
	} else if len(filterByData.Tags) == 1 {
		filterSelections.IsFilterByTag = true
		filterSelections.Tags = filterByData.Tags[0]
	}

	if len(filterByData.Dimensions) > 2 {
		return types.FilterBySelections{}, ValidationError{
			Message: "Results can be filtered by at most two DIMENSION" +
				" filters.",
		}
	} else if len(filterByData.Dimensions) > 0 {
		filterSelections.IsFilterByDimension = true
		filterSelections.Dimensions = filterByData.Dimensions
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
