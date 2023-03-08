package forecast

import (
	"github.com/cduggn/ccexplorer/internal/commands/get/aws"
	"github.com/cduggn/ccexplorer/internal/commands/get/aws/custom_flags"
	"github.com/spf13/cobra"
	"time"
)

var (
	forecastStartDate               string
	forecastEndDate                 string
	forecastGranularity             string
	forecastPredictionIntervalLevel int32
)

func ForecastCommand(c *cobra.Command) *cobra.Command {

	forecastFilterBy := custom_flags.NewForecastFilterBy()
	c.Flags().VarP(&forecastFilterBy, "filterBy", "f",
		"Filter by DIMENSION  (default: none)")

	c.Flags().StringVarP(&forecastStartDate, "start", "s",
		aws.Format(time.Now()), "Start date (defaults to the present day)")

	c.Flags().StringVarP(&forecastEndDate, "end", "e",
		aws.LastDayOfMonth(), "End date (defaults to one month from the start date)")

	// Optional flag to dictate the granularity of the data returned
	c.Flags().StringVarP(&forecastGranularity, "granularity", "g", "MONTHLY",
		"Valid values: DAILY, MONTHLY, HOURLY (default: MONTHLY)")

	c.Flags().Int32VarP(&forecastPredictionIntervalLevel, "predictionIntervalLevel",
		"p", 95, "Prediction interval level (default: 95)")
	return c
}
