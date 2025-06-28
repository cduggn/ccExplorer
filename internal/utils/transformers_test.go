package utils

import (
	"reflect"
	"testing"

	"github.com/cduggn/ccexplorer/internal/types"
)

func TestTransformServicesToRows(t *testing.T) {
	services := []types.Service{
		{
			Keys:  []string{"service1", "tag1"},
			Start: "2023-01-01",
			End:   "2023-01-02",
			Metrics: []types.Metrics{
				{Name: "BlendedCost", Amount: "100.00", Unit: "USD"},
			},
		},
		{
			Keys:  []string{"service2"},
			Start: "2023-01-01",
			End:   "2023-01-02",
			Metrics: []types.Metrics{
				{Name: "BlendedCost", Amount: "50.00", Unit: "USD"},
			},
		},
	}

	rows := TransformServicesToRows(services, "DAILY")
	
	// Should have 2 rows (one per service)
	if len(rows) != 2 {
		t.Errorf("TransformServicesToRows failed. Expected 2 rows, got %d", len(rows))
	}

	// Check first row structure
	expectedFirstRow := []string{"service1", "tag1", "BlendedCost", "DAILY", "2023-01-01", "2023-01-02", "100.00", "USD"}
	if !reflect.DeepEqual(rows[0], expectedFirstRow) {
		t.Errorf("TransformServicesToRows failed. Expected %v, got %v", expectedFirstRow, rows[0])
	}
}

func TestTransformServiceMapToRows(t *testing.T) {
	serviceMap := map[int]types.Service{
		0: {
			Keys:  []string{"service1"},
			Start: "2023-01-01",
			End:   "2023-01-02",
			Metrics: []types.Metrics{
				{Name: "BlendedCost", Amount: "100.00", Unit: "USD"},
			},
		},
		1: {
			Keys:  []string{"service2"},
			Start: "2023-01-01",
			End:   "2023-01-02",
			Metrics: []types.Metrics{
				{Name: "BlendedCost", Amount: "50.00", Unit: "USD"},
			},
		},
	}

	rows := TransformServiceMapToRows(serviceMap, "DAILY")
	
	// Should have 2 rows
	if len(rows) != 2 {
		t.Errorf("TransformServiceMapToRows failed. Expected 2 rows, got %d", len(rows))
	}
}

func TestFilterServicesByMetricType(t *testing.T) {
	services := []types.Service{
		{
			Metrics: []types.Metrics{
				{Name: "BlendedCost", Amount: "100.00", Unit: "USD"},
				{Name: "UsageQuantity", Amount: "10", Unit: "Hours"},
			},
		},
		{
			Metrics: []types.Metrics{
				{Name: "UnblendedCost", Amount: "50.00", Unit: "USD"},
			},
		},
		{
			Metrics: []types.Metrics{
				{Name: "BlendedCost", Amount: "75.00", Unit: "USD"},
			},
		},
	}

	filtered := FilterServicesByMetricType(services, "BlendedCost")
	
	// Should have 2 services with BlendedCost metric
	if len(filtered) != 2 {
		t.Errorf("FilterServicesByMetricType failed. Expected 2 services, got %d", len(filtered))
	}

	// Check that all returned services have the BlendedCost metric
	for _, service := range filtered {
		hasMetric := false
		for _, metric := range service.Metrics {
			if metric.Name == "BlendedCost" {
				hasMetric = true
				break
			}
		}
		if !hasMetric {
			t.Errorf("FilterServicesByMetricType failed. Service doesn't have BlendedCost metric")
		}
	}
}

func TestFilterServicesByDateRange(t *testing.T) {
	services := []types.Service{
		{Start: "2023-01-01", End: "2023-01-02"},
		{Start: "2023-01-03", End: "2023-01-04"},
		{Start: "2023-01-05", End: "2023-01-06"},
	}

	filtered := FilterServicesByDateRange(services, "2023-01-01", "2023-01-04")
	
	// Should have 2 services within the date range
	if len(filtered) != 2 {
		t.Errorf("FilterServicesByDateRange failed. Expected 2 services, got %d", len(filtered))
	}
}

func TestGroupServicesByKey(t *testing.T) {
	services := []types.Service{
		{Keys: []string{"service1", "tag1"}},
		{Keys: []string{"service2", "tag1"}},
		{Keys: []string{"service1", "tag2"}},
	}

	grouped := GroupServicesByKey(services, 0) // Group by first key
	
	// Should have 2 groups: service1 and service2
	if len(grouped) != 2 {
		t.Errorf("GroupServicesByKey failed. Expected 2 groups, got %d", len(grouped))
	}

	// Check service1 group has 2 services
	if len(grouped["service1"]) != 2 {
		t.Errorf("GroupServicesByKey failed. Expected 2 services in service1 group, got %d", len(grouped["service1"]))
	}

	// Check service2 group has 1 service
	if len(grouped["service2"]) != 1 {
		t.Errorf("GroupServicesByKey failed. Expected 1 service in service2 group, got %d", len(grouped["service2"]))
	}
}

func TestCalculateTotalCost(t *testing.T) {
	services := []types.Service{
		{
			Metrics: []types.Metrics{
				{Name: "BlendedCost", NumericAmount: 100.0, Unit: "USD"},
				{Name: "UsageQuantity", NumericAmount: 10.0, Unit: "Hours"},
			},
		},
		{
			Metrics: []types.Metrics{
				{Name: "BlendedCost", NumericAmount: 50.0, Unit: "USD"},
			},
		},
	}

	total := CalculateTotalCost(services)
	
	// Should be 150.0 (100.0 + 50.0, UsageQuantity ignored)
	expected := 150.0
	if total != expected {
		t.Errorf("CalculateTotalCost failed. Expected %f, got %f", expected, total)
	}
}

func TestExtractUniqueKeys(t *testing.T) {
	services := []types.Service{
		{Keys: []string{"service1", "tag1"}},
		{Keys: []string{"service2", "tag1"}},
		{Keys: []string{"service1", "tag2"}},
	}

	keys := ExtractUniqueKeys(services)
	
	// Should have 4 unique keys: service1, service2, tag1, tag2
	if len(keys) != 4 {
		t.Errorf("ExtractUniqueKeys failed. Expected 4 unique keys, got %d", len(keys))
	}

	// Convert to map for easier checking
	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}

	expectedKeys := []string{"service1", "service2", "tag1", "tag2"}
	for _, expectedKey := range expectedKeys {
		if !keyMap[expectedKey] {
			t.Errorf("ExtractUniqueKeys failed. Missing expected key: %s", expectedKey)
		}
	}
}