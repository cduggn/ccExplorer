package printer

import (
	"reflect"
	"testing"
)

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

func TestSortServicesByStartDate(t *testing.T) {
	cases := []struct {
		input  map[int]Service
		expect []Service
	}{
		{
			input: map[int]Service{
				0: {
					Start: "2020-01-01",
				},
				1: {
					Start: "2020-01-02",
				},
				2: {
					Start: "2020-01-03",
				},
				3: {
					Start: "2019-01-03",
				},
				4: {
					Start: "2022-12-13",
				},
				5: {
					Start: "2008-01-28",
				},
			},
			expect: []Service{
				{
					Start: "2022-12-13",
				},
				{
					Start: "2020-01-03",
				},
				{
					Start: "2020-01-02",
				},
				{
					Start: "2020-01-01",
				},
				{
					Start: "2019-01-03",
				},
				{
					Start: "2008-01-28",
				},
			},
		},
	}
	for _, c := range cases {
		result := SortServicesByStartDate(c.input)
		if !reflect.DeepEqual(result, c.expect) {
			t.Errorf("expected %v, got %v", c.expect, result)
		}
	}
}
