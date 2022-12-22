package aws

import (
	"github.com/spf13/cobra"
	"time"
)

var (
	forecastFilter                  string
	forecastFilterValue             string
	forecastFilterType              string
	forecastStartDate               string
	forecastEndDate                 string
	forecastGranularity             string
	forecastPredictionIntervalLevel int32
)

func ForecastCommand(c *cobra.Command) *cobra.Command {

	c.Flags().StringVarP(&forecastStartDate, "start", "s",
		Format(time.Now()), "Must start from today's date")

	c.Flags().StringVarP(&forecastEndDate, "end", "e", DefaultEndDate(Format),
		"Defaults to the present day")

	c.Flags().StringVarP(&forecastFilterType, "forecast-filter-type", "t",
		DefaultEndDate(Format),
		"Forecast filter type")

	c.Flags().StringVarP(&forecastFilter, "forecast-filter", "f",
		DefaultEndDate(Format),
		"Forecast filter name")

	c.Flags().StringVarP(&forecastFilterValue, "forecast-filter-value", "v",
		DefaultEndDate(Format),
		"Forecast filter value")

	// Optional flag to dictate the granularity of the data returned
	c.Flags().StringVarP(&forecastGranularity, "granularity", "g", "MONTHLY",
		"Granularity of billing information to fetch. Monthly, Daily or Hourly")

	c.Flags().Int32VarP(&forecastPredictionIntervalLevel, "predictionIntervalLevel",
		"p", 95, "Prediction interval level")
	return c
}
