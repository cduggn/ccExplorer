package awsservice

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	types2 "github.com/cduggn/ccexplorer/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestCostAndUsageFilterGenerator_FilterByTagEmptyArray(t *testing.T) {
	cases := []struct {
		input  types2.CostAndUsageRequestType
		expect *types.Expression
	}{
		{
			input: types2.CostAndUsageRequestType{
				Granularity: "MONTHLY",
				GroupBy:     []string{"SERVICE"},
				GroupByTag:  []string{""},
				Time: types2.Time{
					Start: "2020-01-01",
					End:   "2020-01-01",
				},
				IsFilterByTagEnabled: true,
				TagFilterValue:       "",
				Rates:                []string{"UNBLENDED"},
				ExcludeDiscounts:     false,
			},
			expect: &types.Expression{
				Tags: nil,
			},
		},
	}
	for _, c := range cases {
		result := CostAndUsageFilterGenerator(c.input)
		if result.Tags != c.expect.Tags {
			t.Errorf("CostAndUsageFilterGenerator(%v) == %v, want %v",
				c.input, result.Tags.Key, c.expect.Tags.Key)
		}
	}
}

func TestCostAndUsageFilterGenerator_FilterByTag(t *testing.T) {
	cases := []struct {
		input  types2.CostAndUsageRequestType
		expect *types.Expression
	}{
		{
			input: types2.CostAndUsageRequestType{
				Granularity: "MONTHLY",
				GroupBy:     []string{"SERVICE"},
				GroupByTag:  []string{"ApplicationName"},
				Time: types2.Time{
					Start: "2020-01-01",
					End:   "2020-01-01",
				},
				IsFilterByTagEnabled: true,
				TagFilterValue:       "MyApp",
				Rates:                []string{"UNBLENDED"},
				ExcludeDiscounts:     false,
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
		input  types2.CostAndUsageRequestType
		expect *types.Expression
	}{
		{
			input: types2.CostAndUsageRequestType{
				Granularity: "MONTHLY",
				GroupBy:     []string{"SERVICE"},
				GroupByTag:  []string{"ApplicationName"},
				Time: types2.Time{
					Start: "2020-01-01",
					End:   "2020-01-01",
				},
				IsFilterByTagEnabled: true,
				TagFilterValue:       "MyApp",
				Rates:                []string{"UNBLENDED"},
				ExcludeDiscounts:     true,
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

func TestCostAndUsageFilterGenerator_FilterByTagDiscountsAndExcludes(t *testing.T) {
	cases := []struct {
		input  types2.CostAndUsageRequestType
		expect *types.Expression
	}{
		{
			input: types2.CostAndUsageRequestType{
				Granularity: "MONTHLY",
				GroupBy:     []string{"SERVICE"},
				GroupByTag:  []string{"ApplicationName"},
				Time: types2.Time{
					Start: "2020-01-01",
					End:   "2020-01-01",
				},
				IsFilterByDimensionEnabled: true,
				//DimensionFilterName:        "SERVICE",
				//DimensionFilterValue:       "Amazon S3",
				DimensionFilter:      map[string]string{"SERVICE": "Amazon S3"},
				IsFilterByTagEnabled: true,
				TagFilterValue:       "MyApp",
				Rates:                []string{"UNBLENDED"},
				ExcludeDiscounts:     true,
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
					{
						Dimensions: &types.DimensionValues{
							Key:    "SERVICE",
							Values: []string{"Amazon S3"},
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
		if result.And[2].Dimensions.Key != c.expect.And[2].Dimensions.Key {
			t.Errorf("CostAndUsageFilterGenerator(%v) == %v, want %v", c.input, result.And[2].Dimensions.Key, c.expect.And[2].Dimensions.Key)
		}

	}
}

func TestCostAndUsageFilterGenerator_NoFilter(t *testing.T) {
	cases := []struct {
		input  types2.CostAndUsageRequestType
		expect *types.Expression
	}{
		{
			input: types2.CostAndUsageRequestType{
				Granularity: "MONTHLY",
				GroupBy:     []string{"SERVICE"},
				GroupByTag:  []string{},
				Time: types2.Time{
					Start: "2020-01-01",
					End:   "2020-01-01",
				},
				IsFilterByTagEnabled: false,
				TagFilterValue:       "",
				Rates:                []string{"UNBLENDED"},
				ExcludeDiscounts:     false,
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
		input  types2.CostAndUsageRequestType
		expect *types.Expression
	}{
		{
			input: types2.CostAndUsageRequestType{
				Granularity: "MONTHLY",
				GroupBy:     []string{"SERVICE"},
				GroupByTag:  []string{},
				Time: types2.Time{
					Start: "2020-01-01",
					End:   "2020-01-01",
				},
				IsFilterByTagEnabled: false,
				TagFilterValue:       "",
				Rates:                []string{"UNBLENDED"},
				ExcludeDiscounts:     true,
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

func TestCostForecastFilterGenerator_MultiDimension(t *testing.T) {
	cases := []struct {
		input  types2.GetCostForecastRequest
		expect *types.Expression
	}{
		{
			input: types2.GetCostForecastRequest{
				Time: types2.Time{
					Start: "2020-01-01",
					End:   "2020-01-01",
				},
				Granularity: "MONTHLY",
				Metric:      "UNBLENDED_COST",
				Filter: types2.Filter{
					Dimensions: []types2.Dimension{
						{
							Key: "SERVICE",
							Value: []string{"Amazon Elastic Compute Cloud" +
								" - Compute"},
						},
						{
							Key:   "RECORD_TYPE",
							Value: []string{"Usage"},
						},
					},
					Tags: nil,
				},
				PredictionIntervalLevel: 0,
			},
			expect: &types.Expression{
				And: []types.Expression{
					{
						Dimensions: &types.DimensionValues{
							Key:    "SERVICE",
							Values: []string{"Amazon Elastic Compute Cloud - Compute"},
						},
					},
					{
						Dimensions: &types.DimensionValues{
							Key:    "RECORD_TYPE",
							Values: []string{"Usage"},
						},
					},
				},
			},
		},
	}
	for _, c := range cases {
		result := CostForecastFilterGenerator(c.input)
		if result.And[0].Dimensions.Key != c.expect.And[0].Dimensions.Key {
			t.Errorf("CostForecastFilterGenerator(%v) == %v, want %v",
				c.input, result.Tags.Key, c.expect.Tags.Key)
		}
		if result.And[1].Dimensions.Key != c.expect.And[1].Dimensions.Key {
			t.Errorf("CostForecastFilterGenerator(%v) == %v, want %v",
				c.input, result.Tags.Key, c.expect.Tags.Key)
		}
	}
}

func TestCostForecastFilterGenerator_SingleDimension(t *testing.T) {
	cases := []struct {
		input  types2.GetCostForecastRequest
		expect *types.Expression
	}{
		{
			input: types2.GetCostForecastRequest{
				Time: types2.Time{
					Start: "2020-01-01",
					End:   "2020-01-01",
				},
				Granularity: "MONTHLY",
				Metric:      "UNBLENDED_COST",
				Filter: types2.Filter{
					Dimensions: []types2.Dimension{
						{
							Key:   "SERVICE",
							Value: []string{"Amazon Elastic Compute Cloud - Compute"},
						},
					},
					Tags: nil,
				},
				PredictionIntervalLevel: 0,
			},
			expect: &types.Expression{
				Dimensions: &types.DimensionValues{
					Key:    "SERVICE",
					Values: []string{"Amazon Elastic Compute Cloud - Compute"},
				},
			},
		},
	}
	for _, c := range cases {
		result := CostForecastFilterGenerator(c.input)
		if result.Dimensions.Key != c.expect.Dimensions.Key {
			t.Errorf("CostForecastFilterGenerator(%v) == %v, want %v",
				c.input, result.Tags.Key, c.expect.Tags.Key)
		}
	}
}

func TestCostAndUsageGroupByGenerator_SingleDimension(t *testing.T) {
	cases := []struct {
		input  types2.CostAndUsageRequestType
		expect []types.GroupDefinition
	}{
		{
			input: types2.CostAndUsageRequestType{
				Granularity: "MONTHLY",
				GroupBy:     []string{"SERVICE"},
				GroupByTag:  []string{},
				Time: types2.Time{
					Start: "2020-01-01",
					End:   "2020-01-01",
				},
				IsFilterByTagEnabled: false,
				TagFilterValue:       "",
				Rates:                []string{"UNBLENDED"},
				ExcludeDiscounts:     true,
			},
			expect: []types.GroupDefinition{
				{
					Type: "DIMENSION",
					Key:  aws.String("SERVICE"),
				},
			},
		},
	}
	for _, c := range cases {
		result := CostAndUsageGroupByGenerator(c.input)
		if *result[0].Key != *c.expect[0].Key {
			t.Errorf("CostAndUsageGroupByGenerator(%v) == %v, want %v",
				c.input, result[0].Key, c.expect[0].Key)
		}
	}
}

func TestCostAndUsageGroupByGenerator_MultiDimension(t *testing.T) {
	cases := []struct {
		input  types2.CostAndUsageRequestType
		expect []types.GroupDefinition
	}{
		{
			input: types2.CostAndUsageRequestType{
				Granularity: "MONTHLY",
				GroupBy:     []string{"SERVICE", "RECORD_TYPE"},
				GroupByTag:  []string{},
				Time: types2.Time{
					Start: "2020-01-01",
					End:   "2020-01-01",
				},
				IsFilterByTagEnabled: false,
				TagFilterValue:       "",
				Rates:                []string{"UNBLENDED"},
				ExcludeDiscounts:     true,
			},
			expect: []types.GroupDefinition{
				{
					Type: "DIMENSION",
					Key:  aws.String("SERVICE"),
				},
				{
					Type: "DIMENSION",
					Key:  aws.String("RECORD_TYPE"),
				},
			},
		},
	}
	for _, c := range cases {
		result := CostAndUsageGroupByGenerator(c.input)
		if *result[0].Key != *c.expect[0].Key {
			t.Errorf("CostAndUsageGroupByGenerator(%v) == %v, want %v",
				c.input, result[0].Key, c.expect[0].Key)
		}
		if *result[1].Key != *c.expect[1].Key {
			t.Errorf("CostAndUsageGroupByGenerator(%v) == %v, want %v",
				c.input, result[1].Key, c.expect[1].Key)
		}
	}
}

func TestCostAndUsageGroupByGenerator_ByTag(t *testing.T) {
	cases := []struct {
		input  types2.CostAndUsageRequestType
		expect []types.GroupDefinition
	}{
		{
			input: types2.CostAndUsageRequestType{
				Granularity: "MONTHLY",
				GroupBy:     []string{"OPERATION"},
				GroupByTag:  []string{"ApplicationName"},
				Time: types2.Time{
					Start: "2020-01-01",
					End:   "2020-01-01",
				},
				IsFilterByTagEnabled: false,
				TagFilterValue:       "MyApp",
				Rates:                []string{"UNBLENDED"},
				ExcludeDiscounts:     true,
			},
			expect: []types.GroupDefinition{
				{
					Type: "DIMENSION",
					Key:  aws.String("OPERATION"),
				},
				{
					Type: "TAG",
					Key:  aws.String("ApplicationName"),
				},
			},
		},
	}
	for _, c := range cases {
		result := CostAndUsageGroupByGenerator(c.input)
		if *result[0].Key != *c.expect[0].Key {
			t.Errorf("CostAndUsageGroupByGenerator(%v) == %v, want %v",
				c.input, result[0].Key, c.expect[0].Key)
		}
	}
}

func TestCostAndUsageGroupByGenerator_ByTagAndDimesion(t *testing.T) {
	cases := []struct {
		input  types2.CostAndUsageRequestType
		expect []types.GroupDefinition
	}{
		{
			input: types2.CostAndUsageRequestType{
				Granularity: "MONTHLY",
				GroupBy:     []string{"SERVICE"},
				GroupByTag:  []string{"ApplicationName"},
				Time: types2.Time{
					Start: "2020-01-01",
					End:   "2020-01-01",
				},
				IsFilterByTagEnabled: false,
				TagFilterValue:       "MyApp",
				Rates:                []string{"UNBLENDED"},
				ExcludeDiscounts:     true,
			},
			expect: []types.GroupDefinition{
				{
					Type: "DIMENSION",
					Key:  aws.String("SERVICE"),
				},
				{
					Type: "TAG",
					Key:  aws.String("ApplicationName"),
				},
			},
		},
	}
	for _, c := range cases {
		result := CostAndUsageGroupByGenerator(c.input)
		if *result[0].Key != *c.expect[0].Key {
			t.Errorf("CostAndUsageGroupByGenerator(%v) == %v, want %v",
				c.input, result[0].Key, c.expect[0].Key)
		}
		if *result[1].Key != *c.expect[1].Key {
			t.Errorf("CostAndUsageGroupByGenerator(%v) == %v, want %v",
				c.input, result[1].Key, c.expect[1].Key)
		}
	}
}
