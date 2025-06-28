package utils

import (
	"reflect"
	"testing"

	"github.com/cduggn/ccexplorer/internal/types"
)

func TestSortBy(t *testing.T) {
	// Test with Service map - sorting by Start date
	serviceMap := map[int]types.Service{
		0: {Start: "2023-01-01", Metrics: []types.Metrics{{NumericAmount: 100.0}}},
		1: {Start: "2023-01-03", Metrics: []types.Metrics{{NumericAmount: 200.0}}},
		2: {Start: "2023-01-02", Metrics: []types.Metrics{{NumericAmount: 150.0}}},
	}

	// Test sorting by date (ascending)
	sortedByDate := SortBy(serviceMap, func(s types.Service) string {
		return s.Start
	}, false)

	expectedOrder := []string{"2023-01-01", "2023-01-02", "2023-01-03"}
	for i, service := range sortedByDate {
		if service.Start != expectedOrder[i] {
			t.Errorf("SortBy date ascending failed. Expected %s, got %s at index %d", expectedOrder[i], service.Start, i)
		}
	}

	// Test sorting by amount (descending)
	sortedByAmount := SortBy(serviceMap, func(s types.Service) float64 {
		if len(s.Metrics) > 0 {
			return s.Metrics[0].NumericAmount
		}
		return 0.0
	}, true)

	expectedAmounts := []float64{200.0, 150.0, 100.0}
	for i, service := range sortedByAmount {
		if len(service.Metrics) > 0 && service.Metrics[0].NumericAmount != expectedAmounts[i] {
			t.Errorf("SortBy amount descending failed. Expected %f, got %f at index %d", expectedAmounts[i], service.Metrics[0].NumericAmount, i)
		}
	}
}

func TestTransform(t *testing.T) {
	// Test transforming slice of strings to slice of ints
	strings := []string{"1", "2", "3"}
	ints := Transform(strings, func(s string) int {
		switch s {
		case "1":
			return 1
		case "2":
			return 2
		case "3":
			return 3
		default:
			return 0
		}
	})

	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(ints, expected) {
		t.Errorf("Transform failed. Expected %v, got %v", expected, ints)
	}
}

func TestFilterSlice(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6}
	evens := FilterSlice(numbers, func(n int) bool {
		return n%2 == 0
	})

	expected := []int{2, 4, 6}
	if !reflect.DeepEqual(evens, expected) {
		t.Errorf("FilterSlice failed. Expected %v, got %v", expected, evens)
	}
}

func TestMapSliceToMap(t *testing.T) {
	services := []types.Service{
		{Keys: []string{"service1"}, Metrics: []types.Metrics{{NumericAmount: 100.0}}},
		{Keys: []string{"service2"}, Metrics: []types.Metrics{{NumericAmount: 200.0}}},
	}

	serviceMap := MapSliceToMap(services,
		func(s types.Service) string { return s.Keys[0] },
		func(s types.Service) float64 {
			if len(s.Metrics) > 0 {
				return s.Metrics[0].NumericAmount
			}
			return 0.0
		},
	)

	expected := map[string]float64{
		"service1": 100.0,
		"service2": 200.0,
	}

	if !reflect.DeepEqual(serviceMap, expected) {
		t.Errorf("MapSliceToMap failed. Expected %v, got %v", expected, serviceMap)
	}
}

func TestConvertSlice(t *testing.T) {
	// Test converting slice of pointers
	type TestStruct struct {
		Value int
	}

	pointers := []*TestStruct{
		{Value: 1},
		{Value: 2},
		{Value: 3},
	}

	values := ConvertSlice(pointers, func(t *TestStruct) int {
		return t.Value
	})

	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(values, expected) {
		t.Errorf("ConvertSlice failed. Expected %v, got %v", expected, values)
	}
}

func TestMaxLen(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e"}
	
	// Test when slice is longer than max
	max := MaxLen(items, 3)
	if max != 3 {
		t.Errorf("MaxLen failed. Expected 3, got %d", max)
	}
	
	// Test when slice is shorter than max
	max = MaxLen(items, 10)
	if max != 5 {
		t.Errorf("MaxLen failed. Expected 5, got %d", max)
	}
}