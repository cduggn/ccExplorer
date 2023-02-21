package printer

import (
	"github.com/cduggn/ccexplorer/pkg/printer/writers/chart"
	"github.com/go-echarts/go-echarts/v2/opts"
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
		CostAndUsage CostAndUsageOutputType
		SortFunc     func(map[int]Service) []Service
	}

	cases := []struct {
		input  args
		expect bool
	}{
		{
			input: args{
				CostAndUsage: CostAndUsageOutputType{
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

func TestGeneratePieItemsC(t *testing.T) {
	type args struct {
		Services []chart.Service
		Key      int
	}

	cases := []struct {
		input  args
		expect []opts.PieData
	}{
		{
			input: args{
				Services: []chart.Service{
					0: {
						Name: "SERVICE",
						Keys: []string{"Amazon Simple Storage Service",
							"PutObject"},
						Metrics: []chart.Metrics{
							{
								Name:          "UNBLENDED",
								Amount:        "0.0000147",
								NumericAmount: 0.0000147,
								Unit:          "USD",
							},
						},
						Start: "2019-01-01",
						End:   "2019-01-31",
					},
					//1: {
					//	Name: "SERVICE",
					//	Keys: []string{"Amazon EC2", "CommittedThroughput"},
					//	Metrics: []Metrics{
					//		{
					//			Name:          "UNBLENDED",
					//			Amount:        "0.0000147",
					//			NumericAmount: 0.0000147,
					//			Unit:          "USD",
					//		},
					//	},
					//},
				},
				Key: 0,
			},
			expect: []opts.PieData{
				{
					Value: 0.0000147,
					Name:  "Amazon Simple Storage Service",
				},
			},
		},
	}
	for _, c := range cases {
		result := chart.PopulatePieDate(c.input.Services, c.input.Key)
		if !reflect.DeepEqual(result, c.expect) {
			t.Errorf("expected %v, got %v", c.expect, result)
		}
	}
}

func TestCreateTitle(t *testing.T) {
	type args struct {
		Dimension string
	}
	cases := []struct {
		input  args
		expect string
	}{
		{
			input: args{
				Dimension: "OPERATION",
			},
			expect: "Pie chart for dimension: [ OPERATION ]",
		},
		{
			input: args{
				Dimension: "SERVICE",
			},
			expect: "Pie chart for dimension: [ SERVICE ]",
		},
	}

	for _, c := range cases {
		result := chart.CreateTitle(c.input.Dimension)
		if !reflect.DeepEqual(result, c.expect) {
			t.Errorf("expected %v, got %v", c.expect, result)
		}
	}

}

func TestCreateSubTitle(t *testing.T) {
	type args struct {
		granularity string
		start       string
		end         string
	}
	cases := []struct {
		input  args
		expect string
	}{
		{
			input: args{
				granularity: "DAILY",
				start:       "2019-01-01",
				end:         "2019-01-31",
			},
			expect: "Response granularity: DAILY. " +
				"Timeframe: 2019-01-01-2019-01-31",
		},
		{
			input: args{
				granularity: "MONTHLY",
				start:       "2019-01-01",
				end:         "2019-01-31",
			},
			expect: "Response granularity: MONTHLY. " +
				"Timeframe: 2019-01-01-2019-01-31",
		},
	}

	for _, c := range cases {
		result := CreateSubTitle(c.input.granularity, c.input.start,
			c.input.end)
		if !reflect.DeepEqual(result, c.expect) {
			t.Errorf("expected %v, got %v", c.expect, result)
		}
	}

}
