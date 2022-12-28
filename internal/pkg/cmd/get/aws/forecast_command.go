package aws

import (
	"github.com/spf13/cobra"
	"time"
)

var (
	forecastFilterDimension         map[string]string
	forecastFilterTag               map[string]string
	forecastStartDate               string
	forecastEndDate                 string
	forecastGranularity             string
	forecastPredictionIntervalLevel int32
)

func ForecastCommand(c *cobra.Command) *cobra.Command {

	// create flag with keyvalue type
	c.Flags().StringToStringVarP(&forecastFilterDimension,
		"forecast-filter-dimensions",
		"d",
		nil, "Filter by dimension. "+
			"Example: -d SERVICE='Amazon EC2' Dimension values can be found in"+
			" the AWS Cost Explorer UI")

	c.Flags().StringToStringVarP(&forecastFilterTag, "filter-tags",
		"t",
		nil, "Filter by tag key and value")

	c.Flags().StringVarP(&forecastStartDate, "start", "s",
		Format(time.Now()), "Must start from today's date")

	c.Flags().StringVarP(&forecastEndDate, "end", "e", DefaultEndDate(Format),
		"Defaults to the present day")

	// Optional flag to dictate the granularity of the data returned
	c.Flags().StringVarP(&forecastGranularity, "granularity", "g", "MONTHLY",
		"Granularity of billing information to fetch. Monthly, Daily or Hourly")

	c.Flags().Int32VarP(&forecastPredictionIntervalLevel, "predictionIntervalLevel",
		"p", 95, "Prediction interval level")
	return c
}
