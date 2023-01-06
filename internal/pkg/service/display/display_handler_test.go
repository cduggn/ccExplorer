package display

import "testing"

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
