package utils

import (
	"github.com/cduggn/ccexplorer/internal/types"
)

// Generic transformation functions using the new generic utilities

// TransformServicesToRows converts services to table rows using generic transformation
func TransformServicesToRows(services []types.Service, granularity string) [][]string {
	var allRows [][]string
	for _, service := range services {
		rows := ConvertServiceToSlice(service, granularity)
		allRows = append(allRows, rows...)
	}
	return allRows
}

// TransformServiceMapToRows converts a service map to rows using generic utilities
func TransformServiceMapToRows(serviceMap map[int]types.Service, granularity string) [][]string {
	// First convert map to slice to maintain order
	services := make([]types.Service, 0, len(serviceMap))
	for _, service := range serviceMap {
		services = append(services, service)
	}
	
	// Then transform each service to its row representation
	var allRows [][]string
	for _, service := range services {
		rows := ConvertServiceToSlice(service, granularity)
		allRows = append(allRows, rows...)
	}
	
	return allRows
}

// TransformServicesWithMetrics transforms services while preserving metrics structure
func TransformServicesWithMetrics(services []types.Service) []types.Service {
	return Transform(services, func(service types.Service) types.Service {
		// Deep copy metrics to avoid mutation
		metrics := Transform(service.Metrics, func(metric types.Metrics) types.Metrics {
			return types.Metrics{
				Name:          metric.Name,
				Amount:        metric.Amount,
				Unit:          metric.Unit,
				NumericAmount: metric.NumericAmount,
				UsageQuantity: metric.UsageQuantity,
			}
		})
		
		return types.Service{
			Name:    service.Name,
			Keys:    service.Keys,
			Start:   service.Start,
			End:     service.End,
			Metrics: metrics,
		}
	})
}

// FilterServicesByMetricType filters services that have a specific metric type
func FilterServicesByMetricType(services []types.Service, metricName string) []types.Service {
	return FilterSlice(services, func(service types.Service) bool {
		for _, metric := range service.Metrics {
			if metric.Name == metricName {
				return true
			}
		}
		return false
	})
}

// FilterServicesByDateRange filters services within a specific date range
func FilterServicesByDateRange(services []types.Service, startDate, endDate string) []types.Service {
	return FilterSlice(services, func(service types.Service) bool {
		return service.Start >= startDate && service.End <= endDate
	})
}

// GroupServicesByKey groups services by a specific key index
func GroupServicesByKey(services []types.Service, keyIndex int) map[string][]types.Service {
	result := make(map[string][]types.Service)
	for _, service := range services {
		key := "unknown"
		if keyIndex < len(service.Keys) {
			key = service.Keys[keyIndex]
		}
		result[key] = append(result[key], service)
	}
	return result
}

// CalculateTotalCost calculates total cost from services using generic aggregation
func CalculateTotalCost(services []types.Service) float64 {
	var total float64
	for _, service := range services {
		for _, metric := range service.Metrics {
			if metric.Unit == "USD" {
				total += metric.NumericAmount
			}
		}
	}
	return total
}

// ExtractUniqueKeys extracts all unique keys from services
func ExtractUniqueKeys(services []types.Service) []string {
	keySet := make(map[string]bool)
	
	for _, service := range services {
		for _, key := range service.Keys {
			keySet[key] = true
		}
	}
	
	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}
	
	return keys
}

// Generic conversion functions to replace legacy conversion functions

// ConvertToDisplayFormat converts raw service data to a standardized display format
func ConvertToDisplayFormat[T any](services []types.Service, granularity string, converter func(types.Service, string) T) []T {
	return Transform(services, func(service types.Service) T {
		return converter(service, granularity)
	})
}

// ConvertServicesToTableRows converts services to table rows using generic transformation
func ConvertServicesToTableRows(services []types.Service, granularity string) [][]string {
	return Transform(services, func(service types.Service) []string {
		rows := ConvertServiceToSlice(service, granularity)
		if len(rows) > 0 {
			return rows[0] // Return first row for each service
		}
		return []string{}
	})
}

// ConvertMapToSlice converts a service map to a slice while preserving order
func ConvertMapToSlice[T any](serviceMap map[int]T) []T {
	result := make([]T, 0, len(serviceMap))
	for i := 0; i < len(serviceMap); i++ {
		if service, exists := serviceMap[i]; exists {
			result = append(result, service)
		}
	}
	return result
}

// ConvertWithContext converts data with additional context using generic transformation
func ConvertWithContext[TInput, TOutput, TContext any](
	input []TInput,
	context TContext,
	converter func(TInput, TContext) TOutput,
) []TOutput {
	return Transform(input, func(item TInput) TOutput {
		return converter(item, context)
	})
}