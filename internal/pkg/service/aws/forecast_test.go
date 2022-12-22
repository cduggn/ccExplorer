package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"strconv"
	"testing"
)

type mockGetCostForecastAPI func(ctx context.Context,
	params *costexplorer.GetCostForecastInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostForecastOutput, error)

func (m mockGetCostForecastAPI) GetCostForecast(ctx context.Context,
	params *costexplorer.GetCostForecastInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostForecastOutput, error) {
	return m(ctx, params, optFns...)
}

var mockGetForecastOutput = &costexplorer.GetCostForecastOutput{
	Total: &types.MetricValue{
		Amount: aws.String("100"),
		Unit:   aws.String("USD"),
	},
}

func TestGetCostForecast(t *testing.T) {
	cases := []struct {
		client func(t *testing.T) GetCostForecastAPI
		bucket string
		key    string
		expect *costexplorer.GetCostForecastOutput
	}{
		{
			client: func(t *testing.T) GetCostForecastAPI {
				return mockGetCostForecastAPI(func(ctx context.Context,
					params *costexplorer.GetCostForecastInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostForecastOutput, error) {
					t.Helper()

					return mockGetForecastOutput, nil
				})
			},
			expect: mockGetForecastOutput,
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			client := tt.client(t)
			output, err := client.GetCostForecast(context.Background(), &costexplorer.GetCostForecastInput{})
			if err != nil {
				t.Fatal(err)
			}
			if output.Total.Amount != tt.expect.Total.Amount {
				t.Errorf("expected %v, got %v", tt.expect.Total.Amount, output.Total.Amount)
			}
		})
	}

}