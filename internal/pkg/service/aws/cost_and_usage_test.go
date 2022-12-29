package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

type mockGetCostAndUsageAPI func(ctx context.Context,
	params *costexplorer.GetCostAndUsageInput,
	optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error)

func (m mockGetCostAndUsageAPI) GetCostAndUsage(ctx context.Context,
	params *costexplorer.GetCostAndUsageInput, optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error) {
	return m(ctx, params, optFns...)
}

var mockOutput = &costexplorer.GetCostAndUsageOutput{
	ResultsByTime: []types.ResultByTime{
		{
			Estimated: true,
			TimePeriod: &types.DateInterval{
				End:   aws.String("2020-01-01"),
				Start: aws.String("2020-12-31"),
			},
			Total: map[string]types.MetricValue{
				"UnblendedCost": {
					Amount: aws.String("100"),
					Unit:   aws.String("USD"),
				},
				"BlendedCost": {
					Amount: aws.String("100"),
					Unit:   aws.String("USD"),
				},
			},
		},
	},
}

func TestAPIClient_GetCostAndUsage(t *testing.T) {

	cases := []struct {
		client func(t *testing.T) GetCostAndUsageAPI
		bucket string
		key    string
		expect *costexplorer.GetCostAndUsageOutput
	}{
		{
			client: func(t *testing.T) GetCostAndUsageAPI {
				return mockGetCostAndUsageAPI(func(ctx context.Context,
					params *costexplorer.GetCostAndUsageInput,
					optFns ...func(*costexplorer.Options)) (*costexplorer.GetCostAndUsageOutput, error) {
					t.Helper()

					return mockOutput, nil
				})
			},
			expect: mockOutput,
		},
	}

	api := &APIClient{}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			ctx := context.TODO()

			results, err := api.GetCostAndUsage(ctx, tt.client(t),
				CostAndUsageRequestType{
					Granularity: "MONTHLY",
					Time: Time{
						Start: "2020-01-01",
						End:   "2020-12-31",
					},
					GroupBy: []string{"SERVICE"},
				})
			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}
			//if e, a := tt.expect, results; results.Services != 0 {
			//	t.Errorf("expect %v, got %v", e, a)
			//}
			assert.Equal(t, len(tt.expect.ResultsByTime), len(results.ResultsByTime))
		})
	}

}
