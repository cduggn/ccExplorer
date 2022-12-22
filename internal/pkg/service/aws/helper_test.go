package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertToFloat(t *testing.T) {
	cases := []struct {
		input  string
		expect float64
	}{
		{
			input:  "100",
			expect: 100,
		},
		{
			input:  "100.00",
			expect: 100,
		},
		{
			input:  "100.01",
			expect: 100.01,
		},
		{
			input:  "100.1",
			expect: 100.1,
		},
		{
			input:  "100.10",
			expect: 100.1,
		},
		{
			input:  "100.11",
			expect: 100.11,
		},
		{
			input:  "100.111",
			expect: 100.111,
		},
		{
			input:  "100.111111111111",
			expect: 100.111111111111,
		},
	}
	for _, c := range cases {
		result := ConvertToFloat(c.input)
		if result != c.expect {
			t.Errorf("ConvertToFloat(%s) == %f, want %f", c.input, result, c.expect)
		}
	}
}

func TestToSlice(t *testing.T) {
	mockDimensionValuesOutput := costexplorer.GetDimensionValuesOutput{
		DimensionValues: []types.DimensionValuesWithAttributes{
			{
				Attributes: map[string]string{
					"ServiceName": "Amazon Elastic Compute Cloud - Compute",
				},
				Value: aws.String("value"),
			},
		},
	}

	slice := ToSlice(mockDimensionValuesOutput)
	assert.Equal(t, "value", slice[0])
	assert.Equal(t, slice, []string{"value"})
}

func TestCostAndUsageFilterGenerator_FilterByTag(t *testing.T) {
	cases := []struct {
		input  CostAndUsageRequestType
		expect *types.Expression
	}{
		{
			input: CostAndUsageRequestType{
				Granularity: "MONTHLY",
				GroupBy:     []string{"SERVICE"},
				Tag:         "ApplicationName",
				Time: Time{
					Start: "2020-01-01",
					End:   "2020-01-01",
				},
				IsFilterEnabled:  true,
				TagFilterValue:   "MyApp",
				Rates:            []string{"UNBLENDED"},
				ExcludeDiscounts: false,
			},
			expect: &types.Expression{
				Tags: &types.TagValues{
					Key:    aws.String("ApplicationName"),
					Values: []string{"MyApp"},
				},
			},
		},
	}
	for _, c := range cases {
		result := CostAndUsageFilterGenerator(c.input)
		if *result.Tags.Key != *c.expect.Tags.Key {
			t.Errorf("CostAndUsageFilterGenerator(%v) == %v, want %v",
				c.input, result.Tags.Key, c.expect.Tags.Key)
		}
	}
}

func TestCostAndUsageFilterGenerator_FilterByTagAndDiscounts(t *testing.T) {
	cases := []struct {
		input  CostAndUsageRequestType
		expect *types.Expression
	}{
		{
			input: CostAndUsageRequestType{
				Granularity: "MONTHLY",
				GroupBy:     []string{"SERVICE"},
				Tag:         "ApplicationName",
				Time: Time{
					Start: "2020-01-01",
					End:   "2020-01-01",
				},
				IsFilterEnabled:  true,
				TagFilterValue:   "MyApp",
				Rates:            []string{"UNBLENDED"},
				ExcludeDiscounts: true,
			},
			expect: &types.Expression{
				And: []types.Expression{
					{
						Not: &types.Expression{
							Dimensions: &types.DimensionValues{
								Key:    "RECORD_TYPE",
								Values: []string{"Refund", "Credit", "DiscountedUsage"},
							},
						},
					},
					{
						Tags: &types.TagValues{
							Key:    aws.String("ApplicationName"),
							Values: []string{"MyApp"},
						},
					},
				},
			},
		},
	}
	for _, c := range cases {
		result := CostAndUsageFilterGenerator(c.input)
		if result.And[0].Not.Dimensions.Key != c.expect.And[0].Not.Dimensions.Key {
			t.Errorf("CostAndUsageFilterGenerator(%v) == %v, want %v",
				c.input, result.And[0].Not.Dimensions.Key, c.expect.And[0].Not.Dimensions.Key)
		}

		if *result.And[1].Tags.Key != *c.expect.And[1].Tags.Key {
			t.Errorf("CostAndUsageFilterGenerator(%v) == %v, want %v",
				c.input, result.And[1].Tags.Key, c.expect.And[1].Tags.Key)
		}

	}
}

func TestCostAndUsageFilterGenerator_NoFilter(t *testing.T) {
	cases := []struct {
		input  CostAndUsageRequestType
		expect *types.Expression
	}{
		{
			input: CostAndUsageRequestType{
				Granularity: "MONTHLY",
				GroupBy:     []string{"SERVICE"},
				Tag:         "",
				Time: Time{
					Start: "2020-01-01",
					End:   "2020-01-01",
				},
				IsFilterEnabled:  false,
				TagFilterValue:   "",
				Rates:            []string{"UNBLENDED"},
				ExcludeDiscounts: false,
			},
			expect: nil,
		},
	}
	for _, c := range cases {
		result := CostAndUsageFilterGenerator(c.input)
		if result != c.expect {
			t.Errorf("CostAndUsageFilterGenerator(%v) == %v, want %v",
				c.input, result, c.expect)
		}
	}

}

func TestCostAndUsageFilterGenerator_FilterByDiscount(t *testing.T) {
	cases := []struct {
		input  CostAndUsageRequestType
		expect *types.Expression
	}{
		{
			input: CostAndUsageRequestType{
				Granularity: "MONTHLY",
				GroupBy:     []string{"SERVICE"},
				Tag:         "",
				Time: Time{
					Start: "2020-01-01",
					End:   "2020-01-01",
				},
				IsFilterEnabled:  false,
				TagFilterValue:   "",
				Rates:            []string{"UNBLENDED"},
				ExcludeDiscounts: true,
			},
			expect: &types.Expression{
				Not: &types.Expression{
					Dimensions: &types.DimensionValues{
						Key:    "RECORD_TYPE",
						Values: []string{"Refund", "Credit", "DiscountedUsage"},
					},
				},
			},
		},
	}
	for _, c := range cases {
		result := CostAndUsageFilterGenerator(c.input)
		if result.Not.Dimensions.Key != c.expect.Not.Dimensions.
			Key {
			t.Errorf("CostAndUsageFilterGenerator(%v) == %v, want %v",
				c.input, result.Not.Dimensions.Key, c.expect.Not.Dimensions.Key)
		}
		
	}
}
