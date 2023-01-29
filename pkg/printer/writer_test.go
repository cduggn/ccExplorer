package printer

import (
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

// Test for SortServicesByMetricAmount function with input type map[int
// ]Service and output type []Service
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

func TestConvertServiceToSlice(t *testing.T) {
	type args struct {
		services    Service
		granularity string
	}

	cases := []struct {
		input  args
		expect [][]string
	}{
		{
			input: args{
				granularity: "DAILY",
				services: Service{
					Name: "Amazon EC2",
					Keys: []string{"Amazon EC2", "CommittedThroughput"},
					Metrics: []Metrics{
						{
							Name:          "UNBLENDED",
							Amount:        "0.0000147",
							NumericAmount: 0.0000147,
							Unit:          "USD",
						},
						{
							Name:          "BLENDED",
							Amount:        "0.0000147",
							NumericAmount: 0.0000147,
							Unit:          "USD",
						},
					},
					Start: "2019-01-01",
					End:   "2019-01-31",
				},
			},
			expect: [][]string{

				{
					"Amazon EC2", "CommittedThroughput", "UNBLENDED", "DAILY", "2019-01-01",
					"2019-01-31", "0.0000147", "USD",
				},
				{
					"Amazon EC2", "CommittedThroughput", "BLENDED", "DAILY",
					"2019-01-01",
					"2019-01-31", "0.0000147", "USD",
				},
			},
		},
	}
	for _, c := range cases {
		result := ConvertServiceToSlice(c.input.services, "DAILY")
		if !reflect.DeepEqual(result, c.expect) {
			t.Errorf("expected %v, got %v", c.expect, result)
		}
	}

}

func TestCostAndUsageToCSV(t *testing.T) {

	type args struct {
		CostAndUsage CostAndUsageReport
		SortFunc     func(map[int]Service) []Service
	}

	cases := []struct {
		input  args
		expect bool
	}{
		{
			input: args{
				CostAndUsage: CostAndUsageReport{
					Granularity: "DAILY",
					Services: map[int]Service{
						0: {
							Keys: []string{"Amazon Simple Storage Service",
								"PutObject"},

							Metrics: []Metrics{
								{
									Name:          "UNBLENDED",
									Amount:        "0.0000147",
									NumericAmount: 0.0000147,
									Unit:          "USD",
								},
								{
									Name:          "BLENDED",
									Amount:        "0.0000147",
									NumericAmount: 0.0000147,
									Unit:          "USD",
								},
							},
							Start: "2019-01-01",
							End:   "2019-01-31",
						},
					},
				},
				SortFunc: SortServicesByMetricAmount,
			},
			expect: true,
		},
	}
	for _, c := range cases {
		err := CostAndUsageToCSV(c.input.SortFunc, c.input.CostAndUsage)
		if err != nil {
			t.Errorf("Failed writing to CSV %v, got error %v", c.expect, err)
		}
	}

}
