package aws

import (
	"context"
	"github.com/cduggn/cloudcost/internal/pkg/service/aws"
	"time"
)

type CommandError struct {
	msg string
}

func (e CommandError) Error() string {
	return e.msg
}

func GetDimensionValues(c *aws.APIClient, d string) ([]string, error) {
	services, err := c.GetDimensionValues(context.TODO(), c.Client, aws.
		GetDimensionValuesRequest{
		Dimension: d,
		Time: aws.Time{
			Start: DefaultStartDate(DayOfCurrentMonth, SubtractDays),
			End:   Format(time.Now()),
		},
	})
	if err != nil {
		return nil, CommandError{
			msg: err.Error(),
		}
	}
	return services, nil
}
