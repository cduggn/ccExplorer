package aws

import (
	"context"
	"github.com/cduggn/ccexplorer/pkg/helpers"
	"github.com/cduggn/ccexplorer/pkg/service/aws"
	"time"
)

type CommandError struct {
	msg string
}

func (e CommandError) Error() string {
	return e.msg
}

func GetDimensionValues(c *aws.APIClient, d string) ([]string, error) {
	services, err := c.GetDimensionValues(context.TODO(), c.Client, aws.GetDimensionValuesRequest{
		Dimension: d,
		Time: aws.Time{
			Start: helpers.DefaultStartDate(helpers.DayOfCurrentMonth, helpers.SubtractDays),
			End:   helpers.Format(time.Now()),
		},
	})
	if err != nil {
		return nil, CommandError{
			msg: err.Error(),
		}
	}
	return services, nil
}
