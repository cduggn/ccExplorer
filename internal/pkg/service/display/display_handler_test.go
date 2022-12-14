package display

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"reflect"
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

// Test for CurateCostAndUsageReport
func TestCurateCostAndUsageReport(t *testing.T) {
	cases := []struct {
		input  *costexplorer.GetCostAndUsageOutput
		expect CostAndUsageReport
	}{
		{
			input: &costexplorer.GetCostAndUsageOutput{
				ResultsByTime: []types.ResultByTime{
					{
						TimePeriod: &types.DateInterval{
							Start: aws.String("2019-01-01"),
							End:   aws.String("2019-01-02"),
						},
						Groups: []types.Group{
							{
								Keys: []string{"key1", "key2"},
								Metrics: map[string]types.MetricValue{
									"metric1": {
										Amount: aws.String("100"),
										Unit:   aws.String("USD"),
									},
									"metric2": {
										Amount: aws.String("200"),
										Unit:   aws.String("USD"),
									},
								},
							},
						},
					},
				},
			},
			expect: CostAndUsageReport{
				Start:       "2019-01-01",
				End:         "2019-01-02",
				Granularity: "DAILY",
				Services: map[int]Service{
					0: {
						Start: "2019-01-01",
						End:   "2019-01-02",
						Keys:  []string{"key1", "key2"},
						Metrics: []Metrics{
							{
								Name:          "metric1",
								Amount:        "100",
								NumericAmount: 100,
								Unit:          "USD",
							},
							{
								Name:          "metric2",
								Amount:        "200",
								NumericAmount: 200,
								Unit:          "USD",
							},
						},
					},
				},
			},
		},
	}
	for _, c := range cases {
		result := CurateCostAndUsageReport(c.input, "DAILY")
		if !result.Equals(c.expect) {
			t.Errorf("CurateCostAndUsageReport(%v) == %v, want %v", c.input, result, c.expect)
		}
	}
}

// Test for SortServicesByMetricAmount function with input type map[int
//]Service and output type []Service
func TestSortServicesByMetricAmount(t *testing.T) {
	cases := []struct {
		input  map[int]Service
		expect []Service
	}{
		{
			input: map[int]Service{
				0: {
					Metrics: []Metrics{
						{
							Name:          "metric1",
							Amount:        "0.00000147",
							NumericAmount: 0.00000147,
							Unit:          "USD",
						},
					},
				},
				1: {
					Metrics: []Metrics{
						{
							Name:          "metric1",
							Amount:        "0.0000147",
							NumericAmount: 0.0000147,
							Unit:          "USD",
						},
					},
				},
				2: {
					Metrics: []Metrics{
						{
							Name:          "metric1",
							Amount:        "1.5",
							NumericAmount: 1.5,
							Unit:          "USD",
						},
					},
				},
			},
			expect: []Service{
				{
					Metrics: []Metrics{
						{
							Name:          "metric1",
							Amount:        "1.5",
							NumericAmount: 1.5,
							Unit:          "USD",
						},
					},
				},
				{
					Metrics: []Metrics{
						{
							Name:          "metric1",
							Amount:        "0.0000147",
							NumericAmount: 0.0000147,
							Unit:          "USD",
						},
					},
				},
				{
					Metrics: []Metrics{
						{
							Name:          "metric1",
							Amount:        "0.00000147",
							NumericAmount: 0.00000147,
							Unit:          "USD",
						},
					},
				},
			},
		},
	}
	for _, c := range cases {
		result := SortServicesByMetricAmount(c.input)
		if !reflect.DeepEqual(result, c.expect) {
			t.Errorf("expected %v, got %v", c.expect, result)
		}
	}
}
