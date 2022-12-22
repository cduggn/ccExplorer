package aws

import (
	"context"
	"github.com/cduggn/cloudcost/internal/pkg/service/aws"
	"time"
)

func GetDimensionValues(c *aws.APIClient) []string {
	services, err := c.GetDimensionValues(context.TODO(), c.Client, aws.
		GetDimensionValuesRequest{
		Dimension: "SERVICE",
		Time: aws.Time{
			Start: DefaultStartDate(DayOfCurrentMonth, SubtractDays),
			End:   Format(time.Now()),
		},
	})
	if err != nil {
		panic(err)
	}
	return services
}
